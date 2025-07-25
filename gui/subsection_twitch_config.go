package gui

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/finahdinner/tidal/config"
	"github.com/finahdinner/tidal/helpers"
	"github.com/finahdinner/tidal/twitch"
)

var twitchConfigWindowSize fyne.Size = fyne.NewSize(800, 1) // height 1 lets the layout determine the height

func (g *GuiWrapper) getTwitchConfigSubsection() *fyne.Container {

	channelUsernameEntry := widget.NewEntry()
	channelUserIdEntry := widget.NewPasswordEntry()
	channelUserIdEntry.Disable()
	appClientIdEntry := widget.NewPasswordEntry()
	appClientSecretEntry := widget.NewPasswordEntry()
	appClientRedirectUri := widget.NewEntry()
	channelAccessTokenEntry := widget.NewPasswordEntry()
	channelAccessTokenEntry.Disable()

	twitchConfig := config.Preferences.TwitchConfig

	channelUsernameEntry.SetText(twitchConfig.UserName)
	channelUserIdEntry.SetText(twitchConfig.UserId)
	appClientIdEntry.SetText(twitchConfig.ClientId)
	appClientSecretEntry.SetText(twitchConfig.ClientSecret)
	appClientRedirectUri.SetText(twitchConfig.ClientRedirectUri)
	channelAccessTokenEntry.SetText(twitchConfig.Credentials.UserAccessToken)

	configForm := container.New(
		layout.NewGridLayoutWithColumns(2),
		widget.NewLabel("Twitch Username"), channelUsernameEntry,
		widget.NewLabel("Client ID"), appClientIdEntry,
		widget.NewLabel("Client Secret"), appClientSecretEntry,
		widget.NewLabel("Redirect URI"), appClientRedirectUri,
		horizontalSpacer(10), layout.NewSpacer(),
		widget.NewLabel("Twitch User ID"), channelUserIdEntry,
		widget.NewLabel("Access Token"), channelAccessTokenEntry,
	)

	channelHeader := canvas.NewText("Twitch Channel", theme.Color(theme.ColorNameForeground))
	channelHeader.TextSize = headerSize

	innerContainer := container.New(
		layout.NewVBoxLayout(),
		horizontalSpacer(10),
		configForm,
		horizontalSpacer(4),
	)

	saveConfigButton := widget.NewButton("Save", nil)
	saveConfigButton.Disable()

	authenticateButton := widget.NewButton("Authenticate", nil)

	buttonContainer := container.New(layout.NewHBoxLayout(), saveConfigButton, authenticateButton)

	bottomButtonRow := container.New(
		layout.NewBorderLayout(nil, nil, nil, buttonContainer),
		buttonContainer,
	)

	rightSpacer := verticalSpacer(30)
	outerContainer := container.New(
		layout.NewBorderLayout(nil, bottomButtonRow, nil, rightSpacer),
		bottomButtonRow,
		rightSpacer,
		innerContainer,
	)

	for _, entry := range []*widget.Entry{
		channelUsernameEntry, appClientIdEntry, appClientSecretEntry, appClientRedirectUri,
	} {
		entry.OnChanged = func(_ string) {
			saveConfigButton.Enable()
			authenticateButton.Disable()
		}
	}

	saveConfigButton.OnTapped = func() {
		prevPreferences := config.Preferences
		if err := handleSaveTwitchConfig(
			channelUsernameEntry, appClientIdEntry, appClientSecretEntry,
			appClientRedirectUri, channelUserIdEntry, channelAccessTokenEntry,
		); err != nil {
			// restore old preferences
			config.Preferences = prevPreferences
			if err2 := config.SavePreferences(); err2 != nil {
				err = err2
			}
			showErrorDialog(
				err,
				"Unable to save twitch configuration - see logs for details.",
				g.SecondaryWindow,
			)
			return
		}
		saveConfigButton.Disable()
		authenticateButton.Enable()
	}

	authenticateButton.OnTapped = func() {
		go func() {
			prevPreferences := config.Preferences
			if err := handleAuthenticate(
				channelUserIdEntry,
				channelAccessTokenEntry,
			); err != nil {
				// restore old preferences
				config.Preferences = prevPreferences
				if err2 := config.SavePreferences(); err2 != nil {
					err = err2
				}
				showErrorDialog(
					err,
					"Unable to authenticate using Twitch credentials - see logs for details.",
					g.SecondaryWindow,
				)
				fyne.Do(func() { saveConfigButton.Enable() }) // to encourage to change settings + save again
			}
			fyne.Do(func() { authenticateButton.Disable() }) // to encourage to authenticate again
		}()
	}

	return container.NewPadded(outerContainer)
}

func handleSaveTwitchConfig(
	channelUsernameEntry *widget.Entry,
	appClientIdEntry *widget.Entry,
	appClientSecretEntry *widget.Entry,
	appClientRedirectUri *widget.Entry,
	channelUserIdEntry *widget.Entry,
	channelAccessTokenEntry *widget.Entry,
) error {
	twitchUsername := channelUsernameEntry.Text
	clientId := appClientIdEntry.Text
	clientSecret := appClientSecretEntry.Text
	clientRedirectUri := appClientRedirectUri.Text

	if twitchUsername == "" || clientId == "" || clientSecret == "" || clientRedirectUri == "" {
		return errors.New("not all twitch application fields were populated")
	}

	err := validateRedirectUri(clientRedirectUri)
	if err != nil {
		return errors.New("redirect URI is not valid")
	}

	config.Preferences.TwitchConfig = config.TwitchConfigT{
		UserName:          twitchUsername,
		UserId:            "",
		ClientId:          clientId,
		ClientSecret:      clientSecret,
		ClientRedirectUri: clientRedirectUri,
		Credentials:       config.CredentialsT{},
	}

	fyne.Do(func() {
		channelUserIdEntry.SetText(config.Preferences.TwitchConfig.UserId)
		channelAccessTokenEntry.SetText(config.Preferences.TwitchConfig.Credentials.UserAccessToken)
	})

	if err := config.SavePreferences(); err != nil {
		return fmt.Errorf("unable to save preferences - err: %w", err)
	}

	return nil
}

func handleAuthenticate(channelUserIdEntry *widget.Entry, channelAccessTokenEntry *widget.Entry) error {
	codeChan := make(chan string)
	csrfToken := helpers.GenerateCsrfToken(32)
	hostAndPort := strings.Replace(strings.Replace(config.Preferences.TwitchConfig.ClientRedirectUri, "https://", "", 1), "http://", "", 1)

	if helpers.PortInUse(hostAndPort) {
		config.Logger.LogInfof("%s is already in use - not creating a new one", hostAndPort)
	} else {
		config.Logger.LogInfo("creating authcode listener")
		if err := twitch.CreateAuthCodeListener(hostAndPort, codeChan, csrfToken); err != nil {
			return fmt.Errorf("unable to open listener port - error: %v", err)
		}
	}

	twitch.SendGetRequestForAuthCode(csrfToken)
	authCode := <-codeChan
	config.Logger.LogInfof("auth code: %v", authCode)

	userAccessTokenInfo, err := twitch.GetUserAccessTokenFromAuthCode(authCode)
	if err != nil {
		return fmt.Errorf("unable to retrieve user access token information - error: %v", err)
	}
	config.Logger.LogInfof("userAccessTokenInfo: %v", userAccessTokenInfo)

	twitchUserId, err := twitch.GetTwitchUserId(userAccessTokenInfo.AccessToken)
	if err != nil {
		return fmt.Errorf("unable to retrieve twitch user id - error: %v", err)
	}
	config.Logger.LogInfof("twitchUserId: %v", twitchUserId)

	config.Preferences.TwitchConfig.Credentials = config.CredentialsT{
		UserAccessToken:        userAccessTokenInfo.AccessToken,
		UserAccessRefreshToken: userAccessTokenInfo.RefreshToken,
		UserAccessScope:        userAccessTokenInfo.Scope,
		ExpiryUnixTimestamp:    time.Now().Unix() + int64(userAccessTokenInfo.ExpiresIn),
	}
	config.Preferences.TwitchConfig.UserId = twitchUserId

	if err := config.SavePreferences(); err != nil {
		return fmt.Errorf("unable to save preferences - error: %v", err)
	}

	fyne.Do(func() {
		channelUserIdEntry.SetText(config.Preferences.TwitchConfig.UserId)
		channelAccessTokenEntry.SetText(config.Preferences.TwitchConfig.Credentials.UserAccessToken)
	})

	config.Logger.LogInfo("successfully authenticated (got access token + twitch user id)")
	return nil
}

func validateRedirectUri(redirectUri string) error {
	regexPattern := `^https?://localhost:\d+$`
	compiledPattern, err := regexp.Compile(regexPattern)
	if err != nil {
		return fmt.Errorf("unable to parse regexPattern %s - %w", regexPattern, err)
	}
	redirectUriBytes := []byte(redirectUri)
	isValid, err := regexp.Match(compiledPattern.String(), redirectUriBytes)
	if err != nil {
		return fmt.Errorf("unable to conduct comparison between pattern %s and redirectUri %s", regexPattern, redirectUri)
	}
	if !isValid {
		return fmt.Errorf("redirectUri %s is not valid", redirectUri)
	}
	return nil
}
