package updater

import (
	"errors"
	"time"
)

var UpdaterTicker *time.Ticker
var TickerDone chan bool

func UpdateUpdateTicker(interval int) error {
	if interval < 0 {
		return errors.New("interval value must be a positive interger")
	}
	if UpdaterTicker != nil {
		UpdaterTicker.Stop()
	}
	UpdaterTicker = time.NewTicker(time.Duration(interval) * time.Second)
	return nil
}

func RemoveUpdateTicker() {
	if UpdaterTicker != nil {
		UpdaterTicker.Stop()
	}
	UpdaterTicker = nil
}
