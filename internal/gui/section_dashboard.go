package gui

import (
	"errors"
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
		config.Logger.LogInfo("dashboardSection already exists")
		return dashboardSection
	}

	header := canvas.NewText("Dashboard", theme.Color(theme.ColorNameForeground))
	header.TextSize = headerSize
	headerContainer := container.NewVBox(header, horizontalSpacer(5))

	consoleTextGrid := widget.NewTextGrid()
	consoleTextGrid.SetText(config.Preferences.ActivityConsoleOutput)

	consoleTextGrid.SetStyleRange(0, 0, 100, 100, &widget.CustomTextGridStyle{FGColor: color.White, BGColor: color.Black})
	consoleBg := canvas.NewRectangle(color.Black)
	consoleBg.Resize(consoleTextGrid.MinSize())
	console := container.NewScroll(container.NewStack(consoleBg, consoleTextGrid))

	uptimeLabel := widget.NewLabel("Uptime: <placeholder>")

	startTidalButton := widget.NewButton("Start Tidal", nil)
	stopTidalButton := widget.NewButton("Stop Tidal", nil)
	stopTidalButton.Disable()

	startTidalButton.OnTapped = func() {
		config.Logger.LogInfo("starting the ticker")

		if !config.Preferences.IsValidForUpdatingStreamVariables() {
			showErrorDialog(
				errors.New("twitch configuration is not populated"),
				"You must first configure your Twitch credentials before starting Tidal.",
				Gui.PrimaryWindow,
			)
			return
		}

		updateInterval := config.Preferences.VariableUpdateInterval
		if updateInterval <= 0 {
			showErrorDialog(
				errors.New("updateInterval is not a positive integer"),
				"Unexpected Error\nCorrupted configuration - updateInterval is not a positive integer.",
				Gui.PrimaryWindow,
			)
			return
		}

		go func() {
			if err := startUpdatingVariables(updateInterval); err != nil {
				fyne.Do(func() {
					startTidalButton.Enable()
					stopTidalButton.Disable()
				})
				if errors.Is(err, twitch.Err401Unauthorised) {
					showErrorDialog(err, "Twitch API returned 401 Unauthorised.\nEnsure you have set up your Twitch credentials correctly.", Gui.PrimaryWindow)
				}
			}
		}()
		startTidalButton.Disable()
		stopTidalButton.Enable()
	}

	stopTidalButton.OnTapped = func() {
		config.Logger.LogInfo("attempting to stop the ticker")
		stopUpdaterTicker()
		stopTidalButton.Disable()
		startTidalButton.Enable()
	}

	buttonContainer := container.New(layout.NewFormLayout(), startTidalButton, stopTidalButton)
	bottomRow := container.New(layout.NewBorderLayout(nil, nil, uptimeLabel, buttonContainer), uptimeLabel, buttonContainer)

	dashboardSection = container.NewPadded(container.New(layout.NewBorderLayout(headerContainer, bottomRow, nil, nil), headerContainer, bottomRow, console))
	return dashboardSection
}
