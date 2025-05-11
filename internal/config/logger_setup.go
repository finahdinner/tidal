package config

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"path"
// )

// type TidalLoggerT struct {
// 	logger *log.Logger
// }

// const logFileName = "tidal.log"

// var TidalLogger *TidalLoggerT

// func newTidalLogger() (*TidalLoggerT, error) {
// 	globalConfigDir, err := os.UserConfigDir()
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to locate UserConfigDir - err: %w", err)
// 	}
// 	appConfigDir := path.Join(globalConfigDir, appConfigDirName)
// 	// logsFilePath := path.Join(globalConfigDir, logFileName)
// 	if !dirExists(logsFilePath) {
// 		os.Mkdir(logsFilePath, 0755) // 0755 - owner can rwx, others can r-x
// 	}

// 	return nil, nil
// }

// func (tl *TidalLoggerT) LogInfo(msg string) {

// }

// func (tl *TidalLoggerT) LogError(msg string) {

// }

// func init() {
// 	var err error
// 	TidalLogger, err = newTidalLogger()
// 	if err != nil {
// 		log.Fatalf("unable to create Tidal logger - err: %w", err)
// 	}
// }
