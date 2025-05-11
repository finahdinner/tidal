package updater

import "log"

func UpdateVariables() {
	for {
		select {
		case <-TickerDone:
			return
		case <-UpdaterTicker.C:
			log.Println("ticker!")
		}
	}
}
