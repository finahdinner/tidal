package gui

import (
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/finahdinner/tidal/internal/preferences"
)

var dashboardSection *fyne.Container

func (g *GuiWrapper) getDashboardSection() *fyne.Container {

	if dashboardSection != nil {
		log.Println("dashboardSection already exists")
		return dashboardSection
	}

	header := canvas.NewText("Dashboard", theme.Color(theme.ColorNameForeground))
	header.TextSize = headerSize
	headerContainer := container.NewVBox(header, horizontalSpacer(5))

	consoleTextGrid := widget.NewTextGrid()
	consoleTextGrid.SetText(preferences.Preferences.ActivityConsoleOutput)

	consoleTextGrid.SetStyleRange(0, 0, 100, 100, &widget.CustomTextGridStyle{FGColor: color.White, BGColor: color.Black})
	consoleBg := canvas.NewRectangle(color.Black)
	consoleBg.Resize(consoleTextGrid.MinSize())
	console := container.NewScroll(container.NewStack(consoleBg, consoleTextGrid))

	uptimeLabel := widget.NewLabel("Uptime: <placeholder>")

	startTidalButton := widget.NewButton("Start Tidal", nil)
	stopTidalButton := widget.NewButton("Stop Tidal", nil)
	stopTidalButton.Disable()

	startTidalButton.OnTapped = func() {
		// err := startUpdatingVariables(5)
		if err := startUpdatingVariables(5); err != nil {
			log.Println("failed to create ticker")
			return // TODO - have proper error handling and/or a popup
		}
		log.Println("created a ticker")
		// go func() { startUpdatingVariables() }()
		startTidalButton.Disable()
		stopTidalButton.Enable()
	}

	stopTidalButton.OnTapped = func() {
		log.Println("stopped the ticker")
		stopUpdaterTicker()
		stopTidalButton.Disable()
		startTidalButton.Enable()
	}

	buttonContainer := container.New(layout.NewFormLayout(), startTidalButton, stopTidalButton)
	bottomRow := container.New(layout.NewBorderLayout(nil, nil, uptimeLabel, buttonContainer), uptimeLabel, buttonContainer)

	dashboardSection = container.NewPadded(container.New(layout.NewBorderLayout(headerContainer, bottomRow, nil, nil), headerContainer, bottomRow, console))
	return dashboardSection
}
