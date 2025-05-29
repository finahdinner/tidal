package gui

import (
	"errors"
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/finahdinner/tidal/internal/config"
	"github.com/finahdinner/tidal/internal/helpers"
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

	var uptimeTicker *time.Ticker
	uptimeTickerDone := make(chan bool, 1)

	uptimeLabel := widget.NewLabel("")

	startTidalButton.OnTapped = func() {
		config.Logger.LogInfo("starting the ticker")

		if !config.Preferences.HasPopulatedTwitchCredentials() {
			showErrorDialog(
				errors.New("twitch configuration is not populated"),
				"You must first configure your Twitch credentials before starting Tidal.",
				g.PrimaryWindow,
			)
			return
		}

		if !config.Preferences.HasPopulatedTitleConfig() {
			showErrorDialog(
				errors.New("title setup is not populated"),
				"You must first configure your Title Setup before starting Tidal.",
				g.PrimaryWindow,
			)
			return
		}

		go func() {
			go func() {
				uptimeSeconds := 0
				fyne.Do(func() {
					uptimeLabel.SetText(fmt.Sprintf("Uptime: %s", helpers.GetTimeStringFromSeconds(uptimeSeconds)))
				})
				uptimeTicker = time.NewTicker(1 * time.Second)
				defer uptimeTicker.Stop()
				for {
					select {
					case <-uptimeTickerDone:
						return
					case <-uptimeTicker.C:
						uptimeSeconds += 1
						// TODO - proper time format
						fyne.Do(func() {
							uptimeLabel.SetText(fmt.Sprintf("Uptime: %s", helpers.GetTimeStringFromSeconds(uptimeSeconds)))
						})
					}
				}
			}()
			if err := startUpdater(); err != nil {
				stopUpdater()
				uptimeTickerDone <- true
				uptimeTicker = nil
				fyne.Do(func() {
					startTidalButton.Enable()
					stopTidalButton.Disable()
					uptimeLabel.SetText("")
				})
				if errors.Is(err, twitch.Err401Unauthorised) {
					showErrorDialog(err, "Twitch API returned 401 Unauthorised.\nEnsure you have set up your Twitch credentials correctly.", g.PrimaryWindow)
				} else {
					showErrorDialog(err, fmt.Sprintf("Error encountered during title update process - err: %v", err), g.PrimaryWindow)
				}
				g.App.SendNotification(fyne.NewNotification("Tidal stopped", "Please check the app."))
			}
		}()
		startTidalButton.Disable()
		stopTidalButton.Enable()
	}

	stopTidalButton.OnTapped = func() {
		stopUpdater()
		uptimeTickerDone <- true
		uptimeTicker = nil
		uptimeLabel.SetText("")
		config.Logger.LogInfo("tidal stopped")
		stopTidalButton.Disable()
		startTidalButton.Enable()
	}

	buttonContainer := container.New(layout.NewFormLayout(), startTidalButton, stopTidalButton)

	titleSetupButton := widget.NewButtonWithIcon("Title Setup", theme.SettingsIcon(), func() {
		g.openSecondaryWindow("Title Setup", g.getTitleSetupSubsection(), &titleSetupWindowSize)
	})
	bottomLeftContainer := container.New(
		layout.NewHBoxLayout(),
		titleSetupButton,
		uptimeLabel,
	)

	bottomRow := container.New(
		layout.NewBorderLayout(nil, nil, bottomLeftContainer, buttonContainer),
		bottomLeftContainer,
		buttonContainer,
	)

	dashboardSection = container.NewPadded(container.New(layout.NewBorderLayout(nil, bottomRow, nil, nil), bottomRow, console))
	return dashboardSection
}
