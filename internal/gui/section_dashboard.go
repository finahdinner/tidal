package gui

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"strings"
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

// type Console struct {
// 	box    *fyne.Container
// 	scroll *container.Scroll
// }

// var dashboardSection *fyne.Container

func (g *GuiWrapper) getDashboardSection() *fyne.Container {

	// if dashboardSection != nil {
	// 	config.Logger.LogDebug("dashboardSection already exists")
	// 	return dashboardSection
	// }

	consoleBox := container.New(layout.NewVBoxLayout())

	consoleStrLines := strings.Split(config.Preferences.ActivityConsoleOutput, "\n")
	log.Println(consoleStrLines)
	for _, s := range consoleStrLines {
		line := widget.NewRichTextFromMarkdown(fmt.Sprintf("`%s`", s))
		line.Wrapping = fyne.TextWrapWord
		line.Scroll = fyne.ScrollNone
		consoleBox.Objects = append(consoleBox.Objects, line)
	}
	consoleBox.Refresh()

	// TODO - add saved console output in preferences? or maybe remove that preference
	consoleBoxBg := canvas.NewRectangle(color.Black)
	// consoleBoxBg.SetMinSize(fyne.NewSize(400, 300))

	consoleScroll := container.NewVScroll(consoleBox)
	// consoleScroll.SetMinSize(fyne.NewSize(400, 300))

	consoleStack := container.New(layout.NewStackLayout(), consoleBoxBg, consoleScroll)

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
		consoleBox.Objects = append(consoleBox.Objects, widget.NewRichTextFromMarkdown("`hi!hi! hi!hi! hi!hi! hi!hi! hi!hi! hi!hi! hi!hi! hi!hi! hi!hi! hi!hi! hi!hi! hi!hi! hi!hi! hi!hi! hi!hi! hi!hi! hi!hi! hi!hi! hi!hi! hi!hi! hi!hi! hi!hi! hi!hi! hi!hi!`"))
		consoleScroll.ScrollToBottom()
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

	return container.NewPadded(container.New(layout.NewBorderLayout(nil, bottomRow, nil, nil), bottomRow, consoleStack))
	// return dashboardSection
}
