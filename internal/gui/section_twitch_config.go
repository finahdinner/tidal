package gui

import (
	"errors"
	"fmt"
	"regexp"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/finahdinner/tidal/internal/preferences"
)

func (g *GuiWrapper) getTwitchConfigSection() *fyne.Container {

	channelUsernameEntry := widget.NewEntry()
	channelUserIdEntry := widget.NewPasswordEntry()
	channelUserIdEntry.Disable()
	appClientIdEntry := widget.NewPasswordEntry()
	appClientSecretEntry := widget.NewPasswordEntry()
	appClientRedirectUri := widget.NewEntry()
	channelAccessTokenEntry := widget.NewPasswordEntry()
	channelAccessTokenEntry.Disable()

	twitchConfig := preferences.Preferences.TwitchConfig

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
		horizontalSpacer(20), layout.NewSpacer(),
		widget.NewLabel("Twitch User ID"), channelUserIdEntry,
		widget.NewLabel("Access Token"), channelAccessTokenEntry,
	)

	applicationHeader := canvas.NewText("Application", theme.Color(theme.ColorNameForeground))
	applicationHeader.TextSize = headerSize

	channelHeader := canvas.NewText("Twitch Channel", theme.Color(theme.ColorNameForeground))
	channelHeader.TextSize = headerSize

	innerContainer := container.New(
		layout.NewVBoxLayout(),
		applicationHeader,
		horizontalSpacer(10),
		configForm,
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

	// functionality

	for _, entry := range []*widget.Entry{
		channelUsernameEntry, appClientIdEntry, appClientSecretEntry, appClientRedirectUri,
	} {
		entry.OnChanged = func(_ string) {
			saveConfigButton.Enable()
			authenticateButton.Disable()
		}
	}

	saveConfigButton.OnTapped = func() {
		if err := handleSaveTwitchConfig(
			channelUsernameEntry, appClientIdEntry, appClientSecretEntry, appClientRedirectUri,
		); err != nil {
			showErrorDialog(
				err,
				fmt.Sprintf("unable to save twitch config:\n%v", err),
				g.PrimaryWindow,
			)
		}
		saveConfigButton.Disable()
		authenticateButton.Enable()
	}

	return container.NewPadded(outerContainer)
}

func handleSaveTwitchConfig(
	channelUsernameEntry *widget.Entry,
	appClientIdEntry *widget.Entry,
	appClientSecretEntry *widget.Entry,
	appClientRedirectUri *widget.Entry,
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

	preferences.Preferences.TwitchConfig = preferences.TwitchConfigT{
		UserName:          twitchUsername,
		UserId:            "",
		ClientId:          clientId,
		ClientSecret:      clientSecret,
		ClientRedirectUri: clientRedirectUri,
		Credentials:       preferences.CredentialsT{},
	}

	if err := preferences.SavePreferences(); err != nil {
		return fmt.Errorf("unable to save preferences - err: %v", err)
	}

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
