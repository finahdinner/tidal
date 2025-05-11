package updater

import "time"

var UpdaterTicker *time.Ticker
var TickerDone chan bool

func UpdateUpdateTicker(interval int) {
	if UpdaterTicker != nil {
		UpdaterTicker.Stop()
	}
	UpdaterTicker = time.NewTicker(time.Duration(interval) * time.Second)
}

func RemoveUpdateTicker() {
	if UpdaterTicker != nil {
		UpdaterTicker.Stop()
	}
	UpdaterTicker = nil
}
