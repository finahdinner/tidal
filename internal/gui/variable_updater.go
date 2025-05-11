package gui

import (
	"errors"
	"log"
	"time"

	"github.com/finahdinner/tidal/internal/twitch"
)

var updaterTicker *time.Ticker
var updaterTickerDone chan struct{}
var updateVariablesSectionSignal = make(chan struct{}, 1)

func startUpdatingVariables(interval int) error {
	if interval < 0 {
		return errors.New("interval value must be a positive integer")
	}
	if updaterTicker != nil {
		return errors.New("ticker already running - stop it first")
	}

	updaterTicker = time.NewTicker(time.Duration(interval) * time.Second)
	updaterTickerDone = make(chan struct{})

	go func() {
		// // initial update, before the ticker
		// if err := twitch.UpdateVariables(); err != nil {
		// 	log.Printf("failed - err: %v", err)
		// }
		for {
			select {
			case <-updaterTickerDone:
				log.Println("updaterTickerDone closed")
				return
			case <-updaterTicker.C:
				if err := twitch.UpdateVariables(); err != nil {
					log.Printf("failed - err: %v", err)
					continue
				}
				select {
				case updateVariablesSectionSignal <- struct{}{}:
					// signal to update widgets in variables section
				default:
					// reached if updateSignalChan is full
				}
			}
		}
	}()
	return nil
}

func stopUpdaterTicker() {
	if updaterTicker != nil {
		log.Println("ticker 'updaterTicker' stopped")
		updaterTicker.Stop()
		updaterTicker = nil
	}
	if updaterTickerDone != nil {
		log.Println("chan 'updaterTickerDone' closed")
		close(updaterTickerDone)
		updaterTickerDone = nil
	}
}
