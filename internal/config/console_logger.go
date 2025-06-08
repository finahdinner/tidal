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

func (cl *ConsoleLoggerT) updateLogPath() error {
	appConfigDir, err := getAppConfigDir()
	if err != nil {
		return err
	}
	timeNowStr := time.Now().Format("20060102150405")
	cl.logPath = path.Join(appConfigDir, timeNowStr+".log")
	return nil
}

func (cl *ConsoleLoggerT) deleteLogPath() {
	cl.logPath = ""
}

func (cl *ConsoleLoggerT) pushToLog(text string) error {
	if cl.logPath == "" {
		if err := cl.updateLogPath(); err != nil {
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

type ConsoleLoggerI interface {
	pushToLog(string) error
	deleteLogPath()
	updateLogPath() error
}

var ConsoleLogger *ConsoleLoggerI
