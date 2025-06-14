package config

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

const consoleLogsFolderName = "console_logs"

type ConsoleLoggerT struct {
	logPath string
}

func (cl *ConsoleLoggerT) NewInstance() error {
	appConfigDir, err := getAppConfigDir()
	if err != nil {
		return err
	}
	consoleLogFolder := path.Join(appConfigDir, consoleLogsFolderName)
	if !dirExists(consoleLogFolder) {
		os.Mkdir(consoleLogFolder, 0755)
	}
	timeNowStr := time.Now().Format("20060102150405")
	cl.logPath = path.Join(consoleLogFolder, timeNowStr+".log")
	return nil
}

func (cl *ConsoleLoggerT) DeleteInstance() {
	cl.logPath = ""
}

func (cl *ConsoleLoggerT) PushToLog(text string) error {
	if cl.logPath == "" {
		if err := cl.NewInstance(); err != nil {
			return fmt.Errorf("unable to update log path: %w", err)
		}
	}
	file, err := os.OpenFile(cl.logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.WriteString(strings.TrimSpace(text) + "\n"); err != nil {
		return err
	}
	return nil
}

var ConsoleLogger *ConsoleLoggerT = &ConsoleLoggerT{}
