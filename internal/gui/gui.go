package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type GuiWrapper struct {
	App           fyne.App
	PrimaryWindow fyne.Window
	Clipboard     fyne.Clipboard
}

var Gui *GuiWrapper

func init() {

	a := app.NewWithID("tidal.preferences")
	a.Settings().SetTheme(&tidalTheme{})

	primaryWindow := a.NewWindow("Tidal")
	primaryWindow.Resize(fyne.NewSize(1000, 600))
	primaryWindow.SetMaster()

	Gui = &GuiWrapper{
		App:           a,
		PrimaryWindow: primaryWindow,
	}

	menuMap := map[string]func() *fyne.Container{
		"Twitch Config": Gui.getTwitchConfigSection,
		"Variables":     Gui.getVariablesSection,
		"Dashboard":     Gui.getDashboardSection,
	}
	menuItemNames := []string{"Dashboard", "Twitch Config", "Variables"}
	menuList := widget.NewList(
		func() int {
			return len(menuItemNames)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Example Label Spacer")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(menuItemNames[i])
		},
	)

	contentContainer := container.New(layout.NewStackLayout())

	menuList.OnSelected = func(id widget.ListItemID) {
		selectedItem := menuItemNames[id]
		contentContainer.Objects = []fyne.CanvasObject{menuMap[selectedItem]()}
		contentContainer.Refresh()
	}
	menuList.Select(0)

	mainSplit := container.New(layout.NewBorderLayout(nil, nil, menuList, nil), menuList, contentContainer)

	Gui.PrimaryWindow.SetContent(mainSplit)
	Gui.PrimaryWindow.Show()
}
