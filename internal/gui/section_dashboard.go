package gui

import (
	"errors"
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/finahdinner/tidal/internal/config"
	"github.com/finahdinner/tidal/internal/twitch"
)

var dashboardSection *fyne.Container

func (g *GuiWrapper) getDashboardSection() *fyne.Container {

	if dashboardSection != nil {
		config.Logger.LogDebug("dashboardSection already exists")
		return dashboardSection
	}

	consoleTextGrid := widget.NewTextGrid()
	consoleTextGrid.SetText(config.Preferences.ActivityConsoleOutput)

	consoleTextGrid.SetStyleRange(0, 0, 100, 100, &widget.CustomTextGridStyle{FGColor: color.White, BGColor: color.Black})
	consoleBg := canvas.NewRectangle(color.Black)
	consoleBg.Resize(consoleTextGrid.MinSize())
	console := container.NewScroll(container.NewStack(consoleBg, consoleTextGrid))

	startTidalButton := widget.NewButton("Start Tidal", nil) // TODO - disable this if no title is set up
	stopTidalButton := widget.NewButton("Stop Tidal", nil)
	stopTidalButton.Disable()

	startTidalButton.OnTapped = func() {
		config.Logger.LogInfo("starting the ticker")

		if !config.Preferences.IsValidForUpdatingTwitchVariables() {
			showErrorDialog(
				errors.New("twitch configuration is not populated"),
				"You must first configure your Twitch credentials before starting Tidal.",
				g.PrimaryWindow,
			)
			return
		}

		go func() {
			// TODO - have a preference for choosing whether to immediately start or not
			if err := startUpdater(true); err != nil {
				fyne.Do(func() {
					startTidalButton.Enable()
					stopTidalButton.Disable()
				})
				if errors.Is(err, twitch.Err401Unauthorised) {
					showErrorDialog(err, "Twitch API returned 401 Unauthorised.\nEnsure you have set up your Twitch credentials correctly.", g.PrimaryWindow)
				} else {
					showErrorDialog(err, fmt.Sprintf("Error encountered during title update process - err: %w", err), g.PrimaryWindow)
				}
				return
			}
		}()
		startTidalButton.Disable()
		stopTidalButton.Enable()
	}

	stopTidalButton.OnTapped = func() {
		stopUpdater()
		config.Logger.LogInfo("tidal stopped")
		stopTidalButton.Disable()
		startTidalButton.Enable()
	}

	buttonContainer := container.New(layout.NewFormLayout(), startTidalButton, stopTidalButton)

	// uptimeLabel := widget.NewLabel("Uptime: <placeholder>")
	titleSetupButton := widget.NewButtonWithIcon("Title setup", theme.SettingsIcon(), func() {
		g.openSecondaryWindow("Title Setup", g.getTitleSetupSubsection(), &titleSetupWindowSize)
	})
	bottomRow := container.New(
		layout.NewBorderLayout(nil, nil, titleSetupButton, buttonContainer),
		titleSetupButton,
		buttonContainer,
	)

	dashboardSection = container.NewPadded(container.New(layout.NewBorderLayout(nil, bottomRow, nil, nil), bottomRow, console))
	return dashboardSection
}
