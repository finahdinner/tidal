package updater

import (
	"log"
	"net/http"

	"github.com/finahdinner/tidal/internal/preferences"
	"github.com/finahdinner/tidal/internal/twitch"
)

func StartUpdatingVariables() {
	for {
		select {
		case <-TickerDone:
			return
		case <-UpdaterTicker.C:
			if err := updateVariables(); err != nil {
				log.Printf("failed - err: %v", err)
				return // TODO - consider whether I need to return if an error
			}
			log.Println("done - exiting")
			return
		}
	}
}

func updateVariables() error {
	// var wg sync.WaitGroup
	httpClient := &http.Client{}
	preferences := preferences.Preferences
	usersApiResponse, err := twitch.GetUsers(httpClient, preferences)
	if err != nil {
		return err
	}
	log.Printf("usersApiResponse data: %v", usersApiResponse.Data)
	// wg.Wait()
	return nil
}
