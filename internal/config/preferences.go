package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/finahdinner/tidal/internal/helpers"
)

const preferencesFileName = "preferences.json"

var appPreferencesPath string
var Preferences PreferencesFormat = defaultPreferences

func SavePreferences() error {
	if err := writeJsonIfSuccessful(appPreferencesPath, Preferences); err != nil {
		return err
	}
	return nil
}

func GetPreferences() (PreferencesFormat, error) {
	prefs := PreferencesFormat{}
	data, err := os.ReadFile(appPreferencesPath)
	if err != nil {
		return prefs, err
	}
	if err := json.Unmarshal(data, &prefs); err != nil {
		return prefs, err
	}
	return prefs, nil
}

func writeJsonIfSuccessful(path string, data any) error {

	tmpFile, err := os.CreateTemp("", "tmpconfig_*.json")
	if err != nil {
		return fmt.Errorf("unable to create temporary config file - err: %w", err)
	}
	defer tmpFile.Close()

	encoder := json.NewEncoder(tmpFile)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("unable to write json data to temporary file - err: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("unable to close temporary file - err: %w", err)
	}

	if err := os.Rename(tmpFile.Name(), path); err != nil {
		return fmt.Errorf("unable to rename temporary file to %s - err: %w", path, err)
	}

	return nil
}

// Returns a slice of the variable names and a map of the variable names to config.TwitchVariableT objects
func GetAllTwitchVariables() ([]string, map[string]TwitchVariableT) {
	varMap := helpers.GenerateMapFromHomogenousStruct[TwitchVariablesT, TwitchVariableT](Preferences.TwitchVariables)
	varNameSlice := make([]string, 0, len(varMap))
	for v, _ := range varMap {
		varNameSlice = append(varNameSlice, v)
	}
	return varNameSlice, varMap
}

func GetAllAiGeneratedVariables() ([]string, map[string]LlmVariableT) {
	varSlice := make([]string, 0, len(Preferences.AiGeneratedVariables))
	varMap := make(map[string]LlmVariableT)
	for _, v := range Preferences.AiGeneratedVariables {
		varSlice = append(varSlice, v.Name)
		varMap[v.Name] = v
	}
	return varSlice, varMap
}
