package preferences

import (
	"fmt"
	"log"
	"os"
	"path"
)

const appConfigDirName = "finahdinner-tidal"

var Preferences PreferencesFormat = defaultPreferences

var appConfigPath string

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

func init() {
	globalConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	appConfigDir := path.Join(globalConfigDir, appConfigDirName)
	if !dirExists(appConfigDir) {
		os.Mkdir(appConfigDir, 0755) // 0755 - owner can rwx, others can r-x
	}
	appConfigPath = path.Join(appConfigDir, "config.json")
	fmt.Printf("appConfigPath: %s\n", appConfigPath)
	if fileExists(appConfigPath) {
		Preferences, err = GetPreferences()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = SavePreferences()
		if err != nil {
			log.Fatal(err)
		}
	}
}
