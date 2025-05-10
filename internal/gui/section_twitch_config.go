package gui

import (
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

	authenticateButton := widget.NewButton("Authenticate", nil)
	saveConfigButton := widget.NewButton("Save", nil)
	// saveConfigButton.Disable()
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

	return container.NewPadded(outerContainer)
}
