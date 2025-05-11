package twitch

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/finahdinner/tidal/internal/helpers"
	"github.com/finahdinner/tidal/internal/preferences"
)

func UpdateVariables() error {

	httpClient := &http.Client{}
	prefs := preferences.Preferences

	var wg sync.WaitGroup
	var mu sync.Mutex

	rawApiResponses := RawApiResponses{}

	numRawApiResponses, err := helpers.NumFieldsInStruct(rawApiResponses)
	if err != nil {
		return err
	}
	log.Printf("numRawApiResponses: %v", numRawApiResponses)

	wg.Add(numRawApiResponses)

	// stream info
	go func() {
		defer wg.Done()
		streamInfo, err := GetStreamInfo(httpClient, prefs)
		if err != nil {
			log.Printf("unable to get stream info - err: %v", err)
			streamInfo = nil
		}
		mu.Lock()
		rawApiResponses.StreamInfo = streamInfo
		mu.Unlock()
	}()

	// subscribers
	go func() {
		defer wg.Done()
		subscribersInfo, err := GetSubscribers(httpClient, prefs)
		if err != nil {
			log.Printf("unable to get subscribers - err: %v", err)
			subscribersInfo = nil
		}
		mu.Lock()
		rawApiResponses.SubscribersInfo = subscribersInfo
		mu.Unlock()
	}()

	// followers
	go func() {
		defer wg.Done()
		followersInfo, err := GetFollowers(httpClient, prefs)
		if err != nil {
			log.Printf("unable to get followers - err: %v", err)
			followersInfo = nil
		}
		mu.Lock()
		rawApiResponses.FollowersInfo = followersInfo
		mu.Unlock()
	}()

	wg.Wait()

	log.Printf("all api responses: %v", rawApiResponses)

	prevPrefs := preferences.Preferences

	if rawApiResponses.StreamInfo != nil {
		prefs.StreamVariables.NumViewers.Value = strconv.Itoa(rawApiResponses.StreamInfo.ViewerCount)
		prefs.StreamVariables.StreamCategory.Value = rawApiResponses.StreamInfo.GameName
		streamStartedAt := rawApiResponses.StreamInfo.StartedAt
		t, err := time.Parse(time.RFC3339, streamStartedAt)
		if err == nil {
			secondsSinceStreamStart := int(time.Since(t).Seconds())
			prefs.StreamVariables.StreamUptime.Value = strconv.Itoa(secondsSinceStreamStart)
		}
	}

	if rawApiResponses.SubscribersInfo != nil {
		prefs.StreamVariables.NumSubscribers.Value = strconv.Itoa(rawApiResponses.SubscribersInfo.Total)
	}

	if rawApiResponses.FollowersInfo != nil {
		prefs.StreamVariables.NumFollowers.Value = strconv.Itoa(rawApiResponses.FollowersInfo.Total)
	}

	preferences.Preferences = prefs

	if err := preferences.SavePreferences(); err != nil {
		// restore old preferences
		preferences.Preferences = prevPrefs
		return fmt.Errorf("unable to save new preferences - err: %v", err)
	}

	log.Println("successfully updated preferences with new values")
	return nil
}
