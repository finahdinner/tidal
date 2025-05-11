package gui

import (
	"errors"
	"log"
	"time"

	"github.com/finahdinner/tidal/internal/twitch"
)

var updaterTicker *time.Ticker
var tickerDone chan bool

func updateUpdateTicker(interval int) error {
	if interval < 0 {
		return errors.New("interval value must be a positive interger")
	}
	if updaterTicker != nil {
		updaterTicker.Stop()
	}
	updaterTicker = time.NewTicker(time.Duration(interval) * time.Second)
	return nil
}

func removeUpdateTicker() {
	if updaterTicker != nil {
		updaterTicker.Stop()
	}
	updaterTicker = nil
}

func startUpdatingVariables() {
	for {
		select {
		case <-tickerDone:
			return
		case <-updaterTicker.C:
			if err := twitch.UpdateVariables(); err != nil {
				log.Printf("failed - err: %v", err)
				continue
			}
			log.Println("done")
		}
	}
}
