package twitch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/finahdinner/tidal/config"
	"github.com/finahdinner/tidal/helpers"
)

type ctxServerKey struct{}

func CreateAuthCodeListener(hostAndPort string, codeChan chan string, csrfToken string) error {

	if hostAndPort == "" {
		return fmt.Errorf("hostAndPort %q is not valid", hostAndPort)
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    hostAndPort,
		Handler: mux,
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// ctx to allow handler to shutdown the server
		ctx := context.WithValue(r.Context(), ctxServerKey{}, server)
		config.Logger.LogInfof("ctx val: %v", ctx.Value(ctxServerKey{}))
		code, err := handleAuthCodeReceived(w, r.WithContext(ctx), csrfToken)
		if err != nil {
			config.Logger.LogInfof("not valid: %v", err)
			return
		}
		fmt.Fprintln(w, "You may now close this browser.")
		codeChan <- code
	})

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			config.Logger.LogInfof("server error: %v", err)
		}
	}()

	config.Logger.LogInfof("listener set up at %s", hostAndPort)

	return nil
}

func SendGetRequestForAuthCode(csrfToken string) {
	params := url.Values{}
	params.Add("client_id", config.Preferences.TwitchConfig.ClientId)
	params.Add("force_verify", "true") // re-authorise each time
	params.Add("redirect_uri", config.Preferences.TwitchConfig.ClientRedirectUri)
	params.Add("response_type", "code")
	params.Add("scope", "channel:read:subscriptions channel:manage:broadcast") // add to the scopes if required
	params.Add("state", csrfToken)

	fullAuthUrl := fmt.Sprintf("%s?%s", twitchApiAuthoriseUrl, params.Encode())

	config.Logger.LogInfof("fullAuthUrl: %v", fullAuthUrl)

	helpers.OpenUrlInBrowser(fullAuthUrl)
	config.Logger.LogInfof("Please complete the authentication in your browser.")
}

func GetUserAccessTokenFromAuthCode(authCode string) (*userAccessTokenInfoT, error) {
	userAccessTokenInfo := &userAccessTokenInfoT{}

	params := url.Values{}
	params.Add("client_id", config.Preferences.TwitchConfig.ClientId)
	params.Add("client_secret", config.Preferences.TwitchConfig.ClientSecret)
	params.Add("code", authCode)
	params.Add("grant_type", "authorization_code")
	params.Add("redirect_uri", config.Preferences.TwitchConfig.ClientRedirectUri)

	resp, err := http.Post(twitchApiTokenUrl, "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
	if err != nil {
		return userAccessTokenInfo, fmt.Errorf("error requesting userAccess token: %v", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(userAccessTokenInfo); err != nil {
		return userAccessTokenInfo, fmt.Errorf("error decoding response: %v", err)
	}
	return userAccessTokenInfo, nil
}

func GetTwitchUserId(accessToken string) (string, error) {
	if config.Preferences.TwitchConfig.UserName == "" {
		return "", fmt.Errorf("username must be populated")
	}

	params := url.Values{}
	params.Add("login", config.Preferences.TwitchConfig.UserName)

	queryUrl := fmt.Sprintf("%s?%s", twitchApiUsersUrl, params.Encode())

	req, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		return "", fmt.Errorf("unable to construct request - err: %w", err)
	}

	req.Header.Set("Client-Id", config.Preferences.TwitchConfig.ClientId)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed - err: %w", err)
	}
	defer resp.Body.Close()

	usersResponse := getUsersApiResponseT{}
	if err := json.NewDecoder(resp.Body).Decode(&usersResponse); err != nil {
		return "", fmt.Errorf("unable to decode response")
	}

	if len(usersResponse.Data) == 0 {
		return "", fmt.Errorf("no twitch users returned")
	}

	twitchUserId := usersResponse.Data[0].Id
	if twitchUserId == "" {
		return "", fmt.Errorf("no twitch user id returned")
	}

	return twitchUserId, nil
}

func getUserAccessTokenFromRefreshToken(ctx context.Context) (*userAccessTokenInfoT, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	userAccessTokenInfo := &userAccessTokenInfoT{}

	params := url.Values{}
	params.Add("client_id", config.Preferences.TwitchConfig.ClientId)
	params.Add("client_secret", config.Preferences.TwitchConfig.ClientSecret)
	params.Add("grant_type", "refresh_token")
	params.Add("refresh_token", config.Preferences.TwitchConfig.Credentials.UserAccessRefreshToken)

	req, err := http.NewRequestWithContext(ctx, "POST", twitchApiTokenUrl, strings.NewReader(params.Encode()))
	if err != nil {
		return userAccessTokenInfo, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return userAccessTokenInfo, fmt.Errorf("error requesting userAccess token: %v", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(userAccessTokenInfo); err != nil {
		return userAccessTokenInfo, fmt.Errorf("error decoding response: %v", err)
	}
	return userAccessTokenInfo, nil
}

func handleAuthCodeReceived(_ http.ResponseWriter, r *http.Request, csrfToken string) (string, error) {
	config.Logger.LogInfof("first request from: %v - shutting down in 2 seconds...", r.URL)

	if server, ok := r.Context().Value(ctxServerKey{}).(*http.Server); ok {
		go shutDownServerGracefully(server, 2*time.Second)
	}

	queryParams := r.URL.Query()
	config.Logger.LogInfof("queryParams: %v", queryParams)
	authCode := queryParams.Get("code")

	if authCode == "" {
		return "", fmt.Errorf("error %v - missing authorisation code", http.StatusBadRequest)
	} else if state := queryParams.Get("state"); state != csrfToken {
		return "", fmt.Errorf("error %v - invalid state received", http.StatusBadRequest)
	} else {
		return authCode, nil
	}
}

func shutDownServerGracefully(server *http.Server, timeoutDuration time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		config.Logger.LogInfof("error during shutdown: %v", err)
	} else {
		config.Logger.LogInfo("server shut down gracefully")
	}
}
