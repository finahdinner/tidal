package twitch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/finahdinner/tidal/internal/preferences"
)

var Err401Unauthorised error = errors.New("unauthorised")

func GetStreamInfo(ctx context.Context, preferences preferences.PreferencesFormat) (*streamInfoT, error) {
	params := url.Values{}
	params.Add("user_id", preferences.TwitchConfig.UserId)
	queryUrl := fmt.Sprintf("%s?%s", twitchApiStreamsUrl, params.Encode())
	log.Printf("queryUrl: %v", queryUrl)
	streamsApiResponse, err := makeGetRequest[getStreamInfoApiResponseT](ctx, queryUrl, "application/json", preferences)
	if err != nil {
		return nil, err
	}
	switch len(streamsApiResponse.Data) {
	case 0:
		return nil, fmt.Errorf("api response returned no stream info for user_id %v", preferences.TwitchConfig.UserId)
	case 1:
		// valid
	default:
		return nil, fmt.Errorf("api response somehow returned more than one stream for user_id %v", preferences.TwitchConfig.UserId)
	}
	return &streamsApiResponse.Data[0], nil
}

func GetSubscribers(ctx context.Context, preferences preferences.PreferencesFormat) (*getChannelSubscribersResponseT, error) {
	params := url.Values{}
	params.Add("broadcaster_id", preferences.TwitchConfig.UserId)
	queryUrl := fmt.Sprintf("%s?%s", twitchApiSubscriptionsUrl, params.Encode())
	log.Printf("queryUrl: %v", queryUrl)
	subscribersApiResponse, err := makeGetRequest[getChannelSubscribersResponseT](ctx, queryUrl, "application/json", preferences)
	if err != nil {
		return nil, err
	}
	return &subscribersApiResponse, nil
}

func GetFollowers(ctx context.Context, preferences preferences.PreferencesFormat) (*getChannelFollowersResponseT, error) {
	params := url.Values{}
	params.Add("broadcaster_id", preferences.TwitchConfig.UserId)
	queryUrl := fmt.Sprintf("%s?%s", twitchApiFollowersUrl, params.Encode())
	log.Printf("queryUrl: %v", queryUrl)
	followersApiResponse, err := makeGetRequest[getChannelFollowersResponseT](ctx, queryUrl, "application/json", preferences)
	if err != nil {
		return &followersApiResponse, err
	}
	return &followersApiResponse, nil
}

func makeGetRequest[T any](ctx context.Context, queryUrl string, mimeType string, preferences preferences.PreferencesFormat) (T, error) {
	var result T

	req, err := http.NewRequestWithContext(ctx, "GET", queryUrl, nil)
	if err != nil {
		return result, fmt.Errorf("unable to construct request for %v - err: %w", queryUrl, err)
	}

	req.Header.Set("Client-Id", preferences.TwitchConfig.ClientId)
	req.Header.Set("Authorization", "Bearer "+preferences.TwitchConfig.Credentials.UserAccessToken)
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
