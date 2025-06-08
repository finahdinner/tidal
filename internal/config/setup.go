package config

import (
	"log"
	"os"
	"path"
)

const appConfigDirName = "finahdinner-tidal"

var AppConfigDir string

func init() {
	var err error
	AppConfigDir, err = getAppConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	// config file
	appPreferencesPath = path.Join(AppConfigDir, preferencesFileName)
	if fileExists(appPreferencesPath) {
		Preferences, err = GetPreferences()
		if err != nil {
			log.Fatalf("unable to load preferences from disk: %v", err)
		}
	} else {
		err = SavePreferences()
		if err != nil {
			log.Fatalf("unable to save/load default preferences: %v", err)
		}
	}

	// create general logger
	appLogPath = path.Join(AppConfigDir, logFileName)
	Logger, err = newTidalLogger(appLogPath)
	if err != nil {
		log.Fatalf("unable to create logger: %v", err)
	}
}

func dirExists(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func getAppConfigDir() (string, error) {
	globalConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appConfigDir := path.Join(globalConfigDir, appConfigDirName)
	if !dirExists(appConfigDir) {
		os.Mkdir(appConfigDir, 0755) // 0755 - owner can rwx, others can r-x
	}
	return appConfigDir, nil
}
