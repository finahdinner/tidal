package twitch

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/finahdinner/tidal/internal/helpers"
	"github.com/finahdinner/tidal/internal/preferences"
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
		log.Printf("ctx val: %v\n", ctx.Value(ctxServerKey{}))
		log.Println("handleAuthCodeReceived before")
		code, err := handleAuthCodeReceived(w, r.WithContext(ctx), csrfToken)
		log.Printf("handleAuthCodeReceived after - code: %v\n", code)
		if err != nil {
			log.Printf("not valid: %v\n", err)
			return
		}
		fmt.Fprintln(w, "You may now close this browser.")
		codeChan <- code
	})

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("server error: %v\n", err)
		}
	}()

	log.Printf("listener set up at %s\n", hostAndPort)

	return nil
}

func SendGetRequestForAuthCode(csrfToken string) {
	params := url.Values{}
	params.Add("client_id", preferences.Preferences.TwitchConfig.ClientId)
	params.Add("force_verify", "true") // re-authorise each time
	params.Add("redirect_uri", preferences.Preferences.TwitchConfig.ClientRedirectUri)
	params.Add("response_type", "code")
	params.Add("scope", "channel:read:subscriptions") // add to the scopes if required
	params.Add("state", csrfToken)

	fullAuthUrl := fmt.Sprintf("%s?%s", twitchApiAuthoriseUrl, params.Encode())

	log.Printf("fullAuthUrl: %v\n", fullAuthUrl)

	helpers.OpenUrlInBrowser(fullAuthUrl)
	log.Printf("Please complete the authentication in your browser.")
}

func GetUserAccessTokenFromAuthCode(authCode string) (*userAccessTokenInfoT, error) {
	userAccessTokenInfo := &userAccessTokenInfoT{}

	params := url.Values{}
	params.Add("client_id", preferences.Preferences.TwitchConfig.ClientId)
	params.Add("client_secret", preferences.Preferences.TwitchConfig.ClientSecret)
	params.Add("code", authCode)
	params.Add("grant_type", "authorization_code")
	params.Add("redirect_uri", preferences.Preferences.TwitchConfig.ClientRedirectUri)

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
	if preferences.Preferences.TwitchConfig.UserName == "" {
		return "", fmt.Errorf("username must be populated")
	}

	params := url.Values{}
	params.Add("login", preferences.Preferences.TwitchConfig.UserName)

	queryUrl := fmt.Sprintf("%s?%s", twitchApiUsersUrl, params.Encode())

	req, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		return "", fmt.Errorf("unable to construct request - err: %w", err)
	}

	req.Header.Set("Client-Id", preferences.Preferences.TwitchConfig.ClientId)
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
	params.Add("client_id", preferences.Preferences.TwitchConfig.ClientId)
	params.Add("client_secret", preferences.Preferences.TwitchConfig.ClientSecret)
	params.Add("grant_type", "refresh_token")
	params.Add("refresh_token", preferences.Preferences.TwitchConfig.Credentials.UserAccessRefreshToken)

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
	log.Printf("first request from: %v - shutting down in 2 seconds...\n", r.URL)

	if server, ok := r.Context().Value(ctxServerKey{}).(*http.Server); ok {
		go shutDownServerGracefully(server, 2*time.Second)
	}

	queryParams := r.URL.Query()
	log.Printf("queryParams: %v\n", queryParams)
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
		log.Printf("error during shutdown: %v\n", err)
	} else {
		log.Println("server shut down gracefully")
	}
}
