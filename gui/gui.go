package gui

import (
	"errors"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/finahdinner/tidal/config"
	"github.com/finahdinner/tidal/helpers"
	"github.com/finahdinner/tidal/twitch"
	"github.com/skratchdot/open-golang/open"
)

type GuiWrapper struct {
	App             fyne.App
	PrimaryWindow   fyne.Window
	SecondaryWindow fyne.Window
}

var Gui *GuiWrapper

func init() {

	a := app.NewWithID(config.AppName)

	icon, err := fyne.LoadResourceFromPath("assets/icon.png")
	if err != nil {
		config.Logger.LogErrorf("unable to load logo - err: %v", err)
	}
	a.SetIcon(icon)

	a.Settings().SetTheme(&tidalTheme{})

	primaryWindow := a.NewWindow("Tidal")
	primaryWindow.Resize(fyne.NewSize(900, 600))
	primaryWindow.SetMaster()

	Gui = &GuiWrapper{
		App:           a,
		PrimaryWindow: primaryWindow,
	}

	menuMap := map[string]func() fyne.CanvasObject{
		"Console":                Gui.getConsoleSection,
		"Stream Variables":       Gui.getStreamVariablesSection,
		"AI-Generated Variables": Gui.getAiGeneratedVariablesSection,
		"Help":                   Gui.getHelpSection,
	}
	menuItemNames := []string{"Console", "Stream Variables", "AI-Generated Variables", "Help"}

	contentContainer := container.New(layout.NewStackLayout())
	menuButtons := container.New(layout.NewGridLayoutWithColumns(len(menuItemNames)))

	for btnIdx, menuItemName := range menuItemNames {
		newBtn := widget.NewButton(menuItemName, nil)
		newBtn.OnTapped = func() {
			for otherBtnIdx, otherBtn := range menuButtons.Objects {
				otherBtn := otherBtn.(*widget.Button)
				if btnIdx != otherBtnIdx {
					otherBtn.Enable()
				}
			}
			contentContainer.Objects = []fyne.CanvasObject{menuMap[menuItemName]()}
			contentContainer.Refresh()
			newBtn.Disable()
		}
		menuButtons.Objects = append(menuButtons.Objects, newBtn)
	}

	if len(menuButtons.Objects) == 0 {
		showErrorDialog(
			errors.New("unable to preload subsections"),
			"Unable to preload subsections",
			Gui.SecondaryWindow,
		)
		return
	}

	// by default, open the first section
	if btn, ok := menuButtons.Objects[0].(*widget.Button); ok {
		btn.OnTapped() // trigger its onTap function
	} else {
		showErrorDialog(
			errors.New("unable to load default subsection"),
			"Unable to load default subsection",
			Gui.SecondaryWindow,
		)
		return
	}

	bottomRibbon := Gui.getBottomRibbon()

	mainSplit := container.New(
		layout.NewBorderLayout(menuButtons, bottomRibbon, nil, nil),
		menuButtons,
		bottomRibbon,
		contentContainer,
	)

	Gui.PrimaryWindow.SetContent(mainSplit)
	Gui.PrimaryWindow.Show()
}

func (g *GuiWrapper) getBottomRibbon() *fyne.Container {

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
					showErrorDialog(err, "Unable to update title - see logs for details.", g.PrimaryWindow)
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
		open.Run(config.AppConfigDir)
	})

	bottomLeftContainer := container.New(
		layout.NewHBoxLayout(),
		titleSetupButton,
		openConfigFolderBtn,
		uptimeLabel,
	)

	return container.New(
		layout.NewBorderLayout(nil, nil, bottomLeftContainer, buttonContainer),
		bottomLeftContainer,
		buttonContainer,
	)
}
