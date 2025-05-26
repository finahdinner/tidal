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
const titleUpdateTimeout = 10 * time.Second

// const updateTitleTimeout = 10 * time.Second

var updaterTicker *time.Ticker
var updaterTickerDone chan struct{}

var updateVariablesSectionSignal = make(chan struct{}, 1)
var updateDashboardSectionSignal = make(chan struct{}, 1)

// Begins a ticker to update the twitch title
func startUpdater(immediatelyStart bool) error {

	ctx, cancel := context.WithTimeout(context.Background(), titleUpdateTimeout)
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

	if immediatelyStart {
		if err := twitch.UpdateTwitchVariables(ctx); err != nil {
			return fmt.Errorf("unable to update twitch variables - err: %w", err)
		}
		if err := updateTitle(ctx); err != nil {
			return fmt.Errorf("unable to update title - err: %w", err)
		}
	}

	updaterTicker = time.NewTicker(time.Duration(updateIntervalSeconds) * time.Second)
	updaterTickerDone = make(chan struct{})

	err401Chan := make(chan error, 1)
	doneChan := make(chan struct{})

	go func() {
		defer close(doneChan)
		if immediatelyStart {
			if err := twitch.UpdateTwitchVariables(ctx); err != nil {
				if errors.Is(err, twitch.Err401Unauthorised) {
					err401Chan <- fmt.Errorf("unable to update twitch variables - err: %w", err)
				}
				return
			}
			if err := updateTitle(ctx); err != nil {
				errChan <- fmt.Errorf("unable to update title - err: %w", err)
				return
			}
		}
		// run with the ticker

	}()

	// reacher when the ticker stops
	select {
	case err401 := <-err401Chan:
		return fmt.Errorf("unauthorised (401) - err: %w", err401)
	case <-doneChan:
		return nil
	}
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

	var wg sync.WaitGroup
	wg.Add(len(aiGeneratedVariableUsedMap))

	promptsMap := map[string]string{}
	twitchVariableStringReplacer := getTwitchVariablesStringReplacer(config.Preferences.TwitchVariables)

	for placeholderStr, v := range aiGeneratedVariableUsedMap {
		prompt := v.PromptMain
		if v.PromptSuffix != "" {
			prompt += "\n" + v.PromptSuffix
		}
		prompt = twitchVariableStringReplacer.Replace(prompt)
		promptsMap[placeholderStr] = prompt
		wg.Add(1)
	}

	llmProvider := config.Preferences.LlmConfig.Provider
	apiKey := config.Preferences.LlmConfig.ApiKey

	llmHandler, err := llm.NewLlmHandler(llmProvider, apiKey)
	if err != nil {
		return fmt.Errorf("unable to create new llm handler - err: %w", err)
	}

	responsesMap := map[string]string{}
	var responsesMapMutex sync.Mutex

	doneChan := make(chan struct{})
	errChan := make(chan error)

	for placeholderStr, prompt := range promptsMap {
		go func(placeholderStr, prompt string) {
			defer wg.Done()
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

	aiGeneratedVariableResponseReplacer, err := getAiGeneratedVariableResponseReplacer(responsesMap)
	if err != nil {
		return errors.New("unable to create a string replacer for AI-generated variable responses")
	}

	// update preferences with llm variable values AND the new title
	for placeholderStr, response := range responsesMap {
		for idx, v := range config.Preferences.AiGeneratedVariables {
			if v.Name == helpers.GetVarNameFromPlaceholderString(placeholderStr) {
				config.Preferences.AiGeneratedVariables[idx].Value = response
			}
		}
	}
	config.Preferences.Title.Value = aiGeneratedVariableResponseReplacer.Replace(titleTemplate)
	config.SavePreferences()

	return nil
}

func getAiGeneratedVariableResponseReplacer(responsesMap map[string]string) (*strings.Replacer, error) {
	replacementList := []string{}
	for placeholderStr, response := range responsesMap {
		if placeholderStr == "" || response == "" {
			return nil, fmt.Errorf("mapping of %q to %q is not valid - must both be non-empty", placeholderStr, response)
		}
		replacementList = append(replacementList, placeholderStr, response)
	}
	return strings.NewReplacer(replacementList...), nil
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
