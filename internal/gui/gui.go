package gui

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/finahdinner/tidal/internal/config"
)

type GuiWrapper struct {
	App             fyne.App
	PrimaryWindow   fyne.Window
	SecondaryWindow fyne.Window
}

var Gui *GuiWrapper

func init() {

	a := app.NewWithID(config.AppName)

	icon, err := fyne.LoadResourceFromPath("internal/assets/icon.png")
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

	menuMap := map[string]func() *fyne.Container{
		"Dashboard": Gui.getDashboardSection,
		"Variables": Gui.getVariablesSection,
	}
	menuItemNames := []string{"Dashboard", "Variables"}

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

	mainSplit := container.New(
		layout.NewBorderLayout(menuButtons, nil, nil, nil),
		menuButtons,
		contentContainer,
	)

	// mainSplit := container.New(layout.NewBorderLayout(nil, nil, menuList, nil), menuList, contentContainer)

	Gui.PrimaryWindow.SetContent(mainSplit)
	Gui.PrimaryWindow.Show()
}
