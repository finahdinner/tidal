package gui

import (
	"errors"
	"log"
	"time"

	"github.com/finahdinner/tidal/internal/twitch"
)

var updaterTicker *time.Ticker
var updateVariablesSectionSignal = make(chan struct{}, 1)

func startUpdatingVariables(interval int) error {
	if interval < 0 {
		return errors.New("interval value must be a positive integer")
	}
	if updaterTicker != nil {
		return errors.New("ticker already exists - call stopUpdaterTicker() first")
	}

	updaterTicker = time.NewTicker(time.Duration(interval) * time.Second)

	go func() {
		for range updaterTicker.C {
			if err := twitch.UpdateVariables(); err != nil {
				log.Printf("failed - err: %v", err)
				continue
			}
			log.Println("tick!")
			select {
			case updateVariablesSectionSignal <- struct{}{}:
				// signal to update widgets in variables section
			default:
				// reached if updateSignalChan is full
			}
		}
	}()
	return nil
}

func stopUpdaterTicker() {
	updaterTicker.Stop()
	updaterTicker = nil
}
