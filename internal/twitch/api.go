package twitch

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/finahdinner/tidal/internal/preferences"
)

func GetStreamInfo(ctx context.Context, preferences preferences.PreferencesFormat) (*streamInfoT, error) {
	params := url.Values{}
	params.Add("user_id", preferences.TwitchConfig.UserId)
	queryUrl := fmt.Sprintf("%s?%s", twitchApiStreamsUrl, params.Encode())
	log.Printf("queryUrl: %v", queryUrl)
	streamsApiResponse, err := makeGetRequest[getStreamInfoApiResponseT](ctx, queryUrl, "application/json", preferences)
	if err != nil {
		return nil, err
	}
	if len(streamsApiResponse.Data) > 1 {
		return nil, fmt.Errorf("api response somehow returned more than one stream for user_id %v", preferences.TwitchConfig.UserId)
	}
	if len(streamsApiResponse.Data) == 0 {
		return nil, fmt.Errorf("api response returned no stream info for user_id %v", preferences.TwitchConfig.UserId)
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
		return &subscribersApiResponse, err
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
		return result, fmt.Errorf("unable to construct request for %v - err: %v", queryUrl, err)
	}

	req.Header.Set("Client-Id", preferences.TwitchConfig.ClientId)
	req.Header.Set("Authorization", "Bearer "+preferences.TwitchConfig.Credentials.UserAccessToken)
	req.Header.Set("Accept", mimeType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, fmt.Errorf("request for %v failed - err: %v", req.URL, err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, fmt.Errorf("unable to decode response from request to %v - err: %v", req.URL, err)
	}

	return result, nil
}

// func buildRequest(queryUrl string, mimeType string, preferences preferences.PreferencesFormat) (*http.Request, error) {

// 	req, err := http.NewRequest("GET", queryUrl, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to construct request for %v - err: %v", queryUrl, err)
// 	}

// 	req.Header.Set("Client-Id", preferences.TwitchConfig.ClientId)
// 	req.Header.Set("Authorization", "Bearer "+preferences.TwitchConfig.Credentials.UserAccessToken)
// 	req.Header.Set("Accept", mimeType)

// 	return req, nil
// }

// func makeGetRequestGetRawResponse(req *http.Request) (io.ReadCloser, error) {
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("request for %v failed - err: %v", req.URL, err)
// 	}
// 	defer resp.Body.Close()
// 	return resp.Body, nil
// }
