package config

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

const logFileName = "tidal.log"

var appLogPath string
var Logger *TidalLoggerT

type TidalLoggerT struct {
	fileLogger   *log.Logger
	stdoutLogger *log.Logger
	bufferLogger *log.Logger
	buffer       *bytes.Buffer
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
	fileLogger := log.New(multiWriter, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
	stdoutLogger := log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)

	buffer := &bytes.Buffer{}
	bufferLogger := log.New(buffer, "", log.Ldate|log.Ltime)

	return &TidalLoggerT{
		fileLogger:   fileLogger,
		stdoutLogger: stdoutLogger,
		bufferLogger: bufferLogger,
		buffer:       buffer,
	}, nil
}

func (tl *TidalLoggerT) LogDebug(msg string) {
	tl.stdoutLogger.Println("DEBUG: " + msg)
}

func (tl *TidalLoggerT) LogDebugf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	tl.stdoutLogger.Println("DEBUG: " + msg)
}

func (tl *TidalLoggerT) LogInfo(msg string) {
	tl.fileLogger.Println("INFO: " + msg)
}

func (tl *TidalLoggerT) LogInfof(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	tl.fileLogger.Println("INFO: " + msg)
}

func (tl *TidalLoggerT) LogError(msg string) {
	tl.fileLogger.Println("ERROR: " + msg)
}

func (tl *TidalLoggerT) LogErrorf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	tl.fileLogger.Println("ERROR: " + msg)
}

func (tl *TidalLoggerT) LogToBuffer(msg string) string {
	tl.buffer.Reset()
	tl.bufferLogger.Print(msg)
	return tl.buffer.String()
}

func (tl *TidalLoggerT) LogToBufferf(format string, args ...any) string {
	tl.buffer.Reset()
	msg := fmt.Sprintf(format, args...)
	tl.bufferLogger.Print(msg)
	return tl.buffer.String()
}
