package twitch

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/finahdinner/tidal/internal/helpers"
	"github.com/finahdinner/tidal/internal/preferences"
)

func UpdateStreamVariables(ctx context.Context) error {

	// httpClient := &http.Client{}
	prefs := preferences.Preferences

	// if the access token expires in <100 seconds, refresh it
	accessTokenExpiryTimestamp := preferences.Preferences.TwitchConfig.Credentials.ExpiryUnixTimestamp
	if time.Now().Unix()+100 > accessTokenExpiryTimestamp {
		newUserAccessTokenInfo, err := getUserAccessTokenFromRefreshToken(ctx)
		if err != nil {
			return fmt.Errorf("unable to refresh access code - err: %w", err)
		}
		preferences.Preferences.TwitchConfig.Credentials = preferences.CredentialsT{
			UserAccessToken:        newUserAccessTokenInfo.AccessToken,
			UserAccessRefreshToken: newUserAccessTokenInfo.RefreshToken,
			UserAccessScope:        newUserAccessTokenInfo.Scope,
			ExpiryUnixTimestamp:    time.Now().Unix() + int64(newUserAccessTokenInfo.ExpiresIn),
		}
		if err := preferences.SavePreferences(); err != nil {
			return fmt.Errorf("unable to save preferences - error: %v", err)
		}
		prefs = preferences.Preferences
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	rawApiResponses := RawApiResponses{}

	numRawApiResponses, err := helpers.NumFieldsInStruct(rawApiResponses)
	if err != nil {
		return err
	}
	log.Printf("numRawApiResponses: %v", numRawApiResponses)

	wg.Add(numRawApiResponses)

	err401Chan := make(chan error, 1)

	// stream info
	go func() {
		defer wg.Done()
		streamInfo, err := GetStreamInfo(ctx, prefs)
		if err != nil {
			log.Printf("unable to get stream info - err: %v", err)
			if errors.Is(err, Err401Unauthorised) {
				err401Chan <- err
			}
			streamInfo = nil
		}
		mu.Lock()
		rawApiResponses.StreamInfo = streamInfo
		mu.Unlock()
	}()

	// subscribers
	go func() {
		defer wg.Done()
		subscribersInfo, err := GetSubscribers(ctx, prefs)
		if err != nil {
			log.Printf("unable to get subscribers - err: %v", err)
			if errors.Is(err, Err401Unauthorised) {
				err401Chan <- err
			}
			subscribersInfo = nil
		}
		mu.Lock()
		rawApiResponses.SubscribersInfo = subscribersInfo
		mu.Unlock()
	}()

	// followers
	go func() {
		defer wg.Done()
		followersInfo, err := GetFollowers(ctx, prefs)
		if err != nil {
			log.Printf("unable to get followers - err: %v", err)
			if errors.Is(err, Err401Unauthorised) {
				err401Chan <- err
			}
			followersInfo = nil
		}
		mu.Lock()
		rawApiResponses.FollowersInfo = followersInfo
		mu.Unlock()
	}()

	requestsDone := make(chan struct{})

	go func() {
		wg.Wait()
		close(requestsDone)
	}()

	select {
	case err := <-err401Chan:
		log.Printf("401 unauthorised http code - invalid oauth, cancelling early - err: %v", err)
		return err
	case <-ctx.Done():
		log.Println("context timed out - returning early")
		return ctx.Err()
	case <-requestsDone:
		// continue
	}

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
	} else {
		prefs.StreamVariables.NumViewers.Value = ""
		prefs.StreamVariables.StreamCategory.Value = ""
		prefs.StreamVariables.StreamUptime.Value = ""
	}

	if rawApiResponses.SubscribersInfo != nil {
		prefs.StreamVariables.NumSubscribers.Value = strconv.Itoa(rawApiResponses.SubscribersInfo.Total)
	} else {
		prefs.StreamVariables.NumSubscribers.Value = ""
	}

	if rawApiResponses.FollowersInfo != nil {
		prefs.StreamVariables.NumFollowers.Value = strconv.Itoa(rawApiResponses.FollowersInfo.Total)
	} else {
		prefs.StreamVariables.NumFollowers.Value = ""
	}

	preferences.Preferences = prefs

	if err := preferences.SavePreferences(); err != nil {
		// restore old preferences
		preferences.Preferences = prevPrefs
		return fmt.Errorf("unable to save new preferences - err: %w", err)
	}

	log.Println("updated preferences with new values")
	return nil
}
