package updater

import (
	"log"
	"net/http"
	"sync"

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

	httpClient := &http.Client{}
	preferences := preferences.Preferences

	var wg sync.WaitGroup
	var mu sync.Mutex

	apiResponses := map[string]any{}

	numApiRequests := 3
	wg.Add(numApiRequests)

	// stream info
	go func() {
		defer wg.Done()
		streamInfo, err := twitch.GetStreamInfo(httpClient, preferences)
		if err != nil {
			log.Printf("unable to get stream info - err: %v", err)
			return
		}
		mu.Lock()
		apiResponses["streamInfo"] = streamInfo
		mu.Unlock()
	}()

	// subscribers
	go func() {
		defer wg.Done()
		subscribersInfo, err := twitch.GetSubscribers(httpClient, preferences)
		if err != nil {
			log.Printf("unable to get subscribers - err: %v", err)
			return
		}
		mu.Lock()
		apiResponses["subscribersInfo"] = subscribersInfo
		mu.Unlock()
	}()

	// followers
	go func() {
		defer wg.Done()
		followersInfo, err := twitch.GetFollowers(httpClient, preferences)
		if err != nil {
			log.Printf("unable to get followers - err: %v", err)
			return
		}
		mu.Lock()
		apiResponses["followersInfo"] = followersInfo
		mu.Unlock()
	}()

	wg.Wait()

	log.Printf("all api responses: %v", apiResponses)

	return nil
}
