package twitch

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/finahdinner/tidal/internal/config"
	"github.com/finahdinner/tidal/internal/helpers"
)

func UpdateTwitchVariables(ctx context.Context) error {

	prefs := config.Preferences

	// if the access token expires in <100 seconds, refresh it
	accessTokenExpiryTimestamp := config.Preferences.TwitchConfig.Credentials.ExpiryUnixTimestamp
	if time.Now().Unix()+100 > accessTokenExpiryTimestamp {
		newUserAccessTokenInfo, err := getUserAccessTokenFromRefreshToken(ctx)
		if err != nil {
			return fmt.Errorf("unable to refresh access code - err: %w", err)
		}
		config.Preferences.TwitchConfig.Credentials = config.CredentialsT{
			UserAccessToken:        newUserAccessTokenInfo.AccessToken,
			UserAccessRefreshToken: newUserAccessTokenInfo.RefreshToken,
			UserAccessScope:        newUserAccessTokenInfo.Scope,
			ExpiryUnixTimestamp:    time.Now().Unix() + int64(newUserAccessTokenInfo.ExpiresIn),
		}
		if err := config.SavePreferences(); err != nil {
			return fmt.Errorf("unable to save preferences - error: %v", err)
		}
		prefs = config.Preferences
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	rawApiResponses := RawApiResponses{}

	numRawApiResponses, err := helpers.NumFieldsInStruct(rawApiResponses)
	if err != nil {
		return err
	}
	config.Logger.LogInfof("numRawApiResponses: %v", numRawApiResponses)

	wg.Add(numRawApiResponses)

	err401Chan := make(chan error, 1)

	// stream info
	go func() {
		defer wg.Done()
		streamInfo, err := GetStreamInfo(ctx, prefs)
		if err != nil {
			config.Logger.LogInfof("unable to get stream info - err: %v", err)
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
			config.Logger.LogInfof("unable to get subscribers - err: %v", err)
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
			config.Logger.LogInfof("unable to get followers - err: %v", err)
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
		config.Logger.LogInfof("401 unauthorised http code - invalid oauth, cancelling early - err: %v", err)
		return err
	case <-ctx.Done():
		config.Logger.LogInfo("context timed out - returning early")
		return ctx.Err()
	case <-requestsDone:
		// continue
	}

	config.Logger.LogInfof("all api responses: %v", rawApiResponses)

	prevPrefs := config.Preferences

	if rawApiResponses.StreamInfo != nil {
		prefs.TwitchVariables.NumViewers.Value = strconv.Itoa(rawApiResponses.StreamInfo.ViewerCount)
		prefs.TwitchVariables.StreamCategory.Value = rawApiResponses.StreamInfo.GameName
		streamStartedAt := rawApiResponses.StreamInfo.StartedAt
		t, err := time.Parse(time.RFC3339, streamStartedAt)
		if err == nil {
			secondsSinceStreamStart := int(time.Since(t).Seconds())
			prefs.TwitchVariables.StreamUptime.Value = strconv.Itoa(secondsSinceStreamStart)
		}
	} else {
		prefs.TwitchVariables.NumViewers.Value = ""
		prefs.TwitchVariables.StreamCategory.Value = ""
		prefs.TwitchVariables.StreamUptime.Value = ""
	}

	if rawApiResponses.SubscribersInfo != nil {
		prefs.TwitchVariables.NumSubscribers.Value = strconv.Itoa(rawApiResponses.SubscribersInfo.Total)
	} else {
		prefs.TwitchVariables.NumSubscribers.Value = ""
	}

	if rawApiResponses.FollowersInfo != nil {
		prefs.TwitchVariables.NumFollowers.Value = strconv.Itoa(rawApiResponses.FollowersInfo.Total)
	} else {
		prefs.TwitchVariables.NumFollowers.Value = ""
	}

	config.Preferences = prefs

	if err := config.SavePreferences(); err != nil {
		// restore old preferences
		config.Preferences = prevPrefs
		return fmt.Errorf("unable to save new preferences - err: %w", err)
	}

	config.Logger.LogInfo("updated preferences with new values")
	return nil
}
