package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

const (
	appConfigDirName = "finahdinner-tidal"
	configFileName   = "config.json"
	logFileName      = "tidal.log"
)

var appConfigPath string
var Preferences PreferencesFormat = defaultPreferences

var appLogPath string
var TidalLogger *TidalLoggerT

type TidalLoggerT struct {
	logger *log.Logger
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

	// config file
	appConfigPath = path.Join(appConfigDir, configFileName)
	log.Printf("appConfigPath: %s\n", appConfigPath)
	if fileExists(appConfigPath) {
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
	log.Printf("appLogPath: %s\n", logFileName)
	TidalLogger, err = newTidalLogger(appLogPath)
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

func newTidalLogger(logPath string) (*TidalLoggerT, error) {
	if !fileExists(logPath) {
		_, err := os.Create(logPath)
		if err != nil {
			return nil, fmt.Errorf("unable to create log file: %w", err)
		}
	}
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0644) // r-w for owner, r-- for other
	if err != nil {
		return nil, fmt.Errorf("unable to open log file: %w", err)
	}
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger := log.New(multiWriter, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
	return &TidalLoggerT{logger: logger}, nil
}

func (tl *TidalLoggerT) LogInfo(msg string) {
	tl.logger.Println("INFO: " + msg)
}

func (tl *TidalLoggerT) LogError(msg string) {
	tl.logger.Println("ERROR: " + msg)
}
