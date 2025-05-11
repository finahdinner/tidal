package updater

import (
	"log"
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

	// httpClient := &http.Client{}
	// preferences := preferences.Preferences
	// var wg sync.WaitGroup
	// var mu sync.Mutex
	// responses := map[string]any{
	// 	"streamResponse":      nil,
	// 	"subscribersResponse": nil,
	// 	"followerssResponse":  nil,
	// }
	// wg.Wait()

	// var wg sync.WaitGroup
	// httpClient := &http.Client{}
	// preferences := preferences.Preferences
	// usersApiResponse, err := twitch.GetUsers(httpClient, preferences)
	// if err != nil {
	// 	return err
	// }
	// log.Printf("usersApiResponse data: %v", usersApiResponse.Data)
	return nil
}
