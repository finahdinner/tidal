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
	startTidalButton := widget.NewButton("Start Tidal", func() {
		// TODO - implement something proper here
		log.Println("clicked a button")
	})

	bottomRow := container.New(layout.NewBorderLayout(nil, nil, uptimeLabel, startTidalButton), uptimeLabel, startTidalButton)

	dashboardSection = container.NewPadded(container.New(layout.NewBorderLayout(headerContainer, bottomRow, nil, nil), headerContainer, bottomRow, console))
	return dashboardSection
}
