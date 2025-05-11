package twitch

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/finahdinner/tidal/internal/preferences"
)

func GetUsers(client *http.Client, preferences preferences.PreferencesFormat) (getUsersApiResponseT, error) {
	params := url.Values{}
	params.Add("id", preferences.TwitchConfig.UserId)
	queryUrl := fmt.Sprintf("%s?%s", twitchApiUsersUrl, params.Encode())
	log.Printf("queryUrl: %v", queryUrl)
	usersApiResponse, err := makeGetRequest[getUsersApiResponseT](client, queryUrl, "application/json", preferences)
	if err != nil {
		return usersApiResponse, err
	}
	return usersApiResponse, nil
}

func GetSubscribers(client *http.Client, preferences preferences.PreferencesFormat) (getChannelSubscribersResponseT, error) {
	params := url.Values{}
	params.Add("broadcaster_id", preferences.TwitchConfig.UserId)
	queryUrl := fmt.Sprintf("%s?%s", twitchApiSubscriptionsUrl, params.Encode())
	log.Printf("queryUrl: %v", queryUrl)
	subscribersApiResponse, err := makeGetRequest[getChannelSubscribersResponseT](client, queryUrl, "application/json", preferences)
	if err != nil {
		return subscribersApiResponse, err
	}
	return subscribersApiResponse, nil
}

func GetFollowers(client *http.Client, preferences preferences.PreferencesFormat) error {
	return nil
}

func makeGetRequest[T any](client *http.Client, queryUrl string, mimeType string, preferences preferences.PreferencesFormat) (T, error) {
	var result T

	req, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		return result, fmt.Errorf("unable to construct request for %v - err: %v", queryUrl, err)
	}

	req.Header.Set("Client-Id", preferences.TwitchConfig.ClientId)
	req.Header.Set("Authorization", "Bearer "+preferences.TwitchConfig.Credentials.UserAccessToken)
	req.Header.Set("Accept", mimeType)

	resp, err := client.Do(req)
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
