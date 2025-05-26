package gui

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/finahdinner/tidal/internal/config"
	"github.com/finahdinner/tidal/internal/helpers"
	"github.com/finahdinner/tidal/internal/twitch"
	"github.com/finahdinner/tidal/pkg/llm"
)

const llmResponseTimeout = 5 * time.Second
const singleCycleTimeout = 10 * time.Second

var updaterTicker *time.Ticker
var updaterTickerDone chan struct{}

var updateVariablesSectionSignal = make(chan struct{}, 1)
var updateDashboardSectionSignal = make(chan struct{}, 1)

// Begins a ticker to update the twitch title
func startUpdater() error {

	ctx, cancel := context.WithTimeout(context.Background(), singleCycleTimeout)
	defer cancel()

	if updaterTicker != nil {
		return errors.New("ticker already running - stop it first")
	}

	updateIntervalMinutes := config.Preferences.Title.TitleUpdateIntervalMinutes
	if updateIntervalMinutes < helpers.MinTitleUpdateIntervalMinutes ||
		updateIntervalMinutes > helpers.MaxTitleUpdateIntervalMinutes {
		return fmt.Errorf(
			"update interval (%v minutes) is not in the valid range between %v and %v",
			updateIntervalMinutes, helpers.MinTitleUpdateIntervalMinutes, helpers.MaxTitleUpdateIntervalMinutes,
		)
	}
	updateIntervalSeconds := updateIntervalMinutes * 60

	updaterTicker = time.NewTicker(time.Duration(updateIntervalSeconds) * time.Second)
	updaterTickerDone = make(chan struct{})

	errChan := make(chan error, 1)
	doneChan := make(chan struct{})

	go func() {
		defer close(doneChan)
		if config.Preferences.Title.UpdateImmediatelyOnStart {
			if err := updateCycle(ctx); err != nil {
				errChan <- fmt.Errorf("unable to complete update cycle - err: %w", err)
			}
		}

		for {
			select {
			case <-updaterTickerDone:
				config.Logger.LogInfo("updateTicker finished")
				return
			case <-updaterTicker.C:
				if err := updateCycle(ctx); err != nil {
					errChan <- fmt.Errorf("unable to complete update cycle - err: %w", err)
					continue
				}

				select {
				case updateVariablesSectionSignal <- struct{}{}:
					// signal to update widgets in variables sections
				default:
					// reached if updateVariablesSectionSignal is full
					config.Logger.LogDebug("updateVariablesSectionSignal chan is full - skipping")
				}

				select {
				case updateDashboardSectionSignal <- struct{}{}:
					// signal to update dashboard section
				default:
					// reached if updateDashboardSectionSignal is full
					config.Logger.LogDebug("updateDashboardSectionSignal chan is full - skipping")
				}
			}
		}
	}()

	// reached when the ticker stops
	select {
	case err := <-errChan:
		stopUpdater()
		return fmt.Errorf("ticker stopped due to error - err: %w", err)
	case <-doneChan:
		return nil
	}
}

// One single update cycle - updates Twitch variables then updates the title
func updateCycle(ctx context.Context) error {
	if err := twitch.UpdateTwitchVariables(ctx); err != nil {
		if errors.Is(err, twitch.Err401Unauthorised) {
			return fmt.Errorf("unable to update twitch variables - err: %w", err)
		}
	}
	if err := updateTitle(ctx); err != nil {
		return fmt.Errorf("unable to update title - err: %w", err)
	}
	return nil
}

// Assumes Twitch variables have been updated already
func updateTitle(ctx context.Context) error {
	titleTemplate := config.Preferences.Title.TitleTemplate
	aiGeneratedVariableUsedMap := map[string]config.LlmVariableT{}

	for _, v := range config.Preferences.AiGeneratedVariables {
		placeholderName := helpers.GenerateVarPlaceholderString(v.Name)
		if strings.Contains(titleTemplate, placeholderName) {
			aiGeneratedVariableUsedMap[placeholderName] = v
		}
	}

	promptsMap := map[string]string{}
	twitchVariableStringReplacer, err := getTwitchVariablesStringReplacer(config.Preferences.TwitchVariables)
	if err != nil {
		return fmt.Errorf("unable to get twitch variables string replacer - err: %v", err)
	}

	for placeholderStr, v := range aiGeneratedVariableUsedMap {
		prompt := v.PromptMain
		if v.PromptSuffix != "" {
			prompt += "\n" + v.PromptSuffix
		}
		prompt = twitchVariableStringReplacer.Replace(prompt)
		promptsMap[placeholderStr] = prompt
	}

	llmProvider := config.Preferences.LlmConfig.Provider
	apiKey := config.Preferences.LlmConfig.ApiKey

	llmHandler, err := llm.NewLlmHandler(llmProvider, apiKey)
	if err != nil {
		return fmt.Errorf("unable to create new llm handler - err: %w", err)
	}

	responsesMap := map[string]string{}

	var wg sync.WaitGroup
	var responsesMapMutex sync.Mutex
	doneChan := make(chan struct{})
	errChan := make(chan error, 1)

	for placeholderStr, prompt := range promptsMap {
		wg.Add(1)
		go func(placeholderStr, prompt string) {
			defer wg.Done()
			config.Logger.LogDebugf("sending prompt: %q", prompt)
			response, err := llmHandler.GetResponseText(prompt, llmResponseTimeout)
			if err != nil {
				errChan <- fmt.Errorf("unable to get response text for %v - err: %w", prompt, err)
				return
			}
			responsesMapMutex.Lock()
			responsesMap[placeholderStr] = response
			responsesMapMutex.Unlock()
		}(placeholderStr, prompt)
	}

	go func() {
		wg.Wait()
		close(doneChan)
	}()

	select {
	case err := <-errChan:
		return fmt.Errorf("unable to retrieve all LLM responses - err: %w", err)
	case <-doneChan:
		//
	}

	newPreferences := config.Preferences

	// update preferences with llm variable values AND the new title
	for placeholderStr, response := range responsesMap {
		for idx, v := range newPreferences.AiGeneratedVariables {
			if v.Name == helpers.GetVarNameFromPlaceholderString(placeholderStr) {
				newPreferences.AiGeneratedVariables[idx].Value = response
			}
		}
	}

	// used to replace ALL mentioned variables with their respective value
	fullVariableReplacementMap := responsesMap

	allTwitchVariablesMap := helpers.GenerateMapFromHomogenousStruct[
		config.TwitchVariablesT, config.TwitchVariableT,
	](config.Preferences.TwitchVariables)

	twitchVariablesUsedInTitleMap := map[string]config.TwitchVariableT{}
	for varName, twitchVar := range allTwitchVariablesMap {
		varNamePlaceholder := helpers.GenerateVarPlaceholderString(varName)
		if strings.Contains(titleTemplate, varNamePlaceholder) {
			twitchVariablesUsedInTitleMap[varName] = twitchVar
		}
	}
	config.Logger.LogDebugf("twitchVariablesUsedInTitleMap: %v", twitchVariablesUsedInTitleMap)

	for varName, twitchVar := range twitchVariablesUsedInTitleMap {
		replaceFrom := helpers.GenerateVarPlaceholderString(varName)
		if _, exists := fullVariableReplacementMap[replaceFrom]; exists {
			return fmt.Errorf("conflicting variable name: %q", replaceFrom)
		}
		fullVariableReplacementMap[replaceFrom] = twitchVar.Value
	}
	config.Logger.LogDebugf("fullVariableReplacementMap: %v", fullVariableReplacementMap)

	allVariablesReplacer, err := helpers.GetStringReplacerFromMap(
		fullVariableReplacementMap, !config.Preferences.Title.ThrowErrorIfEmptyValue, false,
	)
	if err != nil {
		return fmt.Errorf("unable to construct allVariablesReplacer - err: %w", err)
	}

	// new title
	newPreferences.Title.Value = allVariablesReplacer.Replace(titleTemplate)

	config.Logger.LogDebugf("attempting to update stream title to %q", newPreferences.Title.Value)
	if err := twitch.UpdateStreamTitle(ctx, newPreferences); err != nil {
		return fmt.Errorf("unable to update stream title - err: %w", err)
	}
	config.Logger.LogInfof("successfully updated title to %q", newPreferences.Title.Value)

	config.Preferences = newPreferences
	config.SavePreferences()

	return nil
}

func stopUpdater() {
	if updaterTicker != nil {
		config.Logger.LogInfo("updaterTicker stopped")
		updaterTicker.Stop()
		updaterTicker = nil
	}
	if updaterTickerDone != nil {
		config.Logger.LogInfo("updaterTickerDone closed")
		close(updaterTickerDone)
		updaterTickerDone = nil
	}
}
