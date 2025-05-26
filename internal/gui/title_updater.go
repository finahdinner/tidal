package gui

import (
	"fmt"
	"strings"
	"time"

	"github.com/finahdinner/tidal/internal/config"
	"github.com/finahdinner/tidal/internal/helpers"
)

const updateTitleTimeout = 10 * time.Second

var titleVariableUpdaterTicker *time.Ticker
var titleVariableUpdaterTickerDone chan struct{}
var updateDashboardSectionSignal = make(chan struct{}, 1)

func startUpdatingTitle() error {
	updateIntervalMinutes := config.Preferences.TitleConfig.TitleUpdateIntervalMinutes

	if updateIntervalMinutes < helpers.MinTitleUpdateIntervalMinutes ||
		updateIntervalMinutes > helpers.MaxTitleUpdateIntervalMinutes {
		return fmt.Errorf(
			"update interval (%v minutes) is not in the valid range between %v and %v",
			updateIntervalMinutes, helpers.MinTitleUpdateIntervalMinutes, helpers.MaxTitleUpdateIntervalMinutes,
		)
	}

	// updateIntervalSeconds := config.Preferences.TitleConfig.TitleUpdateIntervalMinutes * 60

	titleTemplate := config.Preferences.TitleConfig.TitleTemplate
	aiGeneratedVariableUsedMap := map[string]config.LlmVariableT{}

	for _, v := range config.Preferences.AiGeneratedVariables {
		placeholderName := helpers.GenerateVarPlaceholderString(v.Name)
		if strings.Contains(titleTemplate, placeholderName) {
			aiGeneratedVariableUsedMap[placeholderName] = v
		}
	}

	return nil

	// TODO - list out all AI-generated variables
	// then determine which exist in this specific title
	// then only generate those variables

}
