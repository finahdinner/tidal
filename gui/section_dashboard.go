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
	"github.com/finahdinner/tidal/config"
	"github.com/finahdinner/tidal/helpers"
	"github.com/finahdinner/tidal/twitch"

	"github.com/skratchdot/open-golang/open"
)

type ActivityConsoleT struct {
	box    *fyne.Container
	scroll *container.Scroll
	stack  *fyne.Container
}

var ActivityConsole *ActivityConsoleT

var dashboardSection *fyne.Container

func init() {
	if ActivityConsole == nil {
		ActivityConsole = NewActivityConsole()
	}
}

func NewActivityConsole() *ActivityConsoleT {
	consoleBox := container.New(layout.NewVBoxLayout())
	consoleBoxBg := canvas.NewRectangle(color.Black)
	consoleScroll := container.NewVScroll(consoleBox)
	consoleStack := container.New(layout.NewStackLayout(), consoleBoxBg, consoleScroll)
	return &ActivityConsoleT{consoleBox, consoleScroll, consoleStack}
}

// Append a new line to the activity console
func (ac *ActivityConsoleT) pushToConsole(text string) error {
	if err := config.ConsoleLogger.PushToLog(text); err != nil {
		return err
	}
	line := widget.NewRichTextFromMarkdown(fmt.Sprintf("`%s`", text))
	line.Wrapping = fyne.TextWrapWord
	line.Scroll = fyne.ScrollNone
	fyne.Do(func() {
		ac.box.Objects = append(ac.box.Objects, line)
		ac.scroll.ScrollToBottom()
		ac.box.Refresh()
	})
	return nil
}

// Clears the console and closes the console log file
func (ac *ActivityConsoleT) clearConsole() {
	config.ConsoleLogger.DeleteInstance()
	ac.box.Objects = []fyne.CanvasObject{}
	fyne.Do(func() {
		ac.box.Refresh()
	})
}

func (g *GuiWrapper) getDashboardSection() *fyne.Container {

	if ActivityConsole == nil {
		ActivityConsole = NewActivityConsole()
	}

	if dashboardSection != nil {
		return dashboardSection
	}

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

		config.ConsoleLogger.NewInstance()
		startTidalButton.Disable()
		stopTidalButton.Enable()

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
	}

	stopTidalButton.OnTapped = func() {
		stopUpdater()
		uptimeTickerDone <- true
		uptimeTicker = nil
		uptimeLabel.SetText("")
		config.Logger.LogInfo("tidal stopped")
		config.ConsoleLogger.DeleteInstance()
		ActivityConsole.clearConsole()
		stopTidalButton.Disable()
		startTidalButton.Enable()
	}

	buttonContainer := container.New(layout.NewFormLayout(), startTidalButton, stopTidalButton)

	titleSetupButton := widget.NewButtonWithIcon("Title Setup", theme.SettingsIcon(), func() {
		g.openSecondaryWindow("Title Setup", g.getTitleSetupSubsection(), &titleSetupWindowSize)
	})

	openConfigFolderBtn := widget.NewButtonWithIcon("Config Folder", theme.FolderIcon(), func() {
		fmt.Println(config.AppConfigDir)
		open.Run(config.AppConfigDir)
	})

	bottomLeftContainer := container.New(
		layout.NewHBoxLayout(),
		titleSetupButton,
		openConfigFolderBtn,
		uptimeLabel,
	)

	bottomRow := container.New(
		layout.NewBorderLayout(nil, nil, bottomLeftContainer, buttonContainer),
		bottomLeftContainer,
		buttonContainer,
	)

	dashboardSection = container.NewPadded(container.New(layout.NewBorderLayout(nil, bottomRow, nil, nil), bottomRow, ActivityConsole.stack))
	return dashboardSection
}
