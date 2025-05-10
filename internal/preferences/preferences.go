package preferences

import (
	"encoding/json"
	"fmt"
	"os"
)

func SavePreferences() error {
	if err := writeJsonIfSuccessful(appConfigPath, Preferences); err != nil {
		return err
	}
	return nil
}

func GetPreferences() (PreferencesFormat, error) {
	preferences := PreferencesFormat{}
	data, err := os.ReadFile(appConfigPath)
	if err != nil {
		return preferences, err
	}
	if err := json.Unmarshal(data, &preferences); err != nil {
		return preferences, err
	}
	return preferences, nil
}

func writeJsonIfSuccessful(path string, data any) error {

	tmpFile, err := os.CreateTemp("", "tmpconfig_*.json")
	if err != nil {
		return fmt.Errorf("unable to create temporary config file - err: %v", err)
	}
	defer tmpFile.Close()

	encoder := json.NewEncoder(tmpFile)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("unable to write json data to temporary file - err: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("unable to close temporary file - err: %v", err)
	}

	if err := os.Rename(tmpFile.Name(), path); err != nil {
		return fmt.Errorf("unable to rename temporary file to %s - err: %v", path, err)
	}

	return nil
}
