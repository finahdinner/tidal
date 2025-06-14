package twitch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/finahdinner/tidal/config"
)

var Err401Unauthorised error = errors.New("unauthorised")

func GetStreamInfo(ctx context.Context, prefs config.PreferencesFormat) (*streamInfoT, error) {
	params := url.Values{}
	params.Add("user_id", prefs.TwitchConfig.UserId)
	queryUrl := fmt.Sprintf("%s?%s", twitchApiStreamsUrl, params.Encode())
	config.Logger.LogInfof("queryUrl: %v", queryUrl)
	streamsApiResponse, err := makeGetRequest[getStreamInfoApiResponseT](ctx, queryUrl, "application/json", prefs)
	if err != nil {
		return nil, err
	}
	switch len(streamsApiResponse.Data) {
	case 0:
		return nil, fmt.Errorf("api response returned no stream info for user_id %v", prefs.TwitchConfig.UserId)
	case 1:
		// valid
	default:
		return nil, fmt.Errorf("api response somehow returned more than one stream for user_id %v", prefs.TwitchConfig.UserId)
	}
	return &streamsApiResponse.Data[0], nil
}

func GetSubscribers(ctx context.Context, prefs config.PreferencesFormat) (*getChannelSubscribersResponseT, error) {
	params := url.Values{}
	params.Add("broadcaster_id", prefs.TwitchConfig.UserId)
	queryUrl := fmt.Sprintf("%s?%s", twitchApiSubscriptionsUrl, params.Encode())
	config.Logger.LogInfof("queryUrl: %v", queryUrl)
	subscribersApiResponse, err := makeGetRequest[getChannelSubscribersResponseT](ctx, queryUrl, "application/json", prefs)
	if err != nil {
		return nil, err
	}
	return &subscribersApiResponse, nil
}

func GetFollowers(ctx context.Context, prefs config.PreferencesFormat) (*getChannelFollowersResponseT, error) {
	params := url.Values{}
	params.Add("broadcaster_id", prefs.TwitchConfig.UserId)
	queryUrl := fmt.Sprintf("%s?%s", twitchApiFollowersUrl, params.Encode())
	config.Logger.LogInfof("queryUrl: %v", queryUrl)
	followersApiResponse, err := makeGetRequest[getChannelFollowersResponseT](ctx, queryUrl, "application/json", prefs)
	if err != nil {
		return &followersApiResponse, err
	}
	return &followersApiResponse, nil
}

// PATCH request to /channels endpoint
func UpdateStreamTitle(ctx context.Context, prefs config.PreferencesFormat) error {
	params := url.Values{}
	params.Add("broadcaster_id", config.Preferences.TwitchConfig.UserId)
	queryUrl := fmt.Sprintf("%s?%s", twitchApiChannelsUrl, params.Encode())
	config.Logger.LogInfof("queryUrl: %v", queryUrl)

	reqBody := map[string]string{
		"title": prefs.Title.Value,
	}
	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("unable to parse reqBody - err: %w", err)
	}
	config.Logger.LogInfof("reqBodyJson: %v", string(reqBodyJson))

	// make a PATCH request
	req, err := http.NewRequestWithContext(ctx, "PATCH", queryUrl, bytes.NewBuffer(reqBodyJson))
	if err != nil {
		return fmt.Errorf("unable to construct request using url %q and body %v", queryUrl, reqBodyJson)
	}

	req.Header.Set("Client-Id", prefs.TwitchConfig.ClientId)
	req.Header.Set("Authorization", "Bearer "+prefs.TwitchConfig.Credentials.UserAccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request for %v failed - err: %w", req.URL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unable to update content - http code %v", resp.Status)
	}
	return nil
}

func makeGetRequest[T any](ctx context.Context, queryUrl string, mimeType string, prefs config.PreferencesFormat) (T, error) {
	var result T

	req, err := http.NewRequestWithContext(ctx, "GET", queryUrl, nil)
	if err != nil {
		return result, fmt.Errorf("unable to construct request for %v - err: %w", queryUrl, err)
	}

	req.Header.Set("Client-Id", prefs.TwitchConfig.ClientId)
	req.Header.Set("Authorization", "Bearer "+prefs.TwitchConfig.Credentials.UserAccessToken)
	req.Header.Set("Accept", mimeType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, fmt.Errorf("request for %v failed - err: %w", req.URL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return result, Err401Unauthorised
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, fmt.Errorf("unable to decode response from request to %v - err: %w", req.URL, err)
	}

	return result, nil
}
