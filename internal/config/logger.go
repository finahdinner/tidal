package config

import (
	"fmt"
	"io"
	"log"
	"os"
)

const logFileName = "tidal.log"

var appLogPath string
var Logger *TidalLoggerT

type TidalLoggerT struct {
	logger *log.Logger
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

func (tl *TidalLoggerT) LogInfof(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	tl.logger.Println("INFO: " + msg)
}

func (tl *TidalLoggerT) LogError(msg string) {
	tl.logger.Println("ERROR: " + msg)
}

func (tl *TidalLoggerT) LogErrorf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	tl.logger.Println("ERROR: " + msg)
}
