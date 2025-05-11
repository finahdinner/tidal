package config

import (
	"log"
	"os"
	"path"
)

const appConfigDirName = "finahdinner-tidal"

func init() {
	globalConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	appConfigDir := path.Join(globalConfigDir, appConfigDirName)
	if !dirExists(appConfigDir) {
		os.Mkdir(appConfigDir, 0755) // 0755 - owner can rwx, others can r-x
	}

	// config file
	appPreferencesPath = path.Join(appConfigDir, preferencesFileName)
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

	// create logger
	appLogPath = path.Join(appConfigDir, logFileName)
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
