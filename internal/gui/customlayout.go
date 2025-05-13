package gui

import (
	"errors"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"github.com/finahdinner/tidal/internal/config"
)

func horizontalSpacer(height float32) *canvas.Rectangle {
	padding := canvas.NewRectangle(color.Transparent)
	padding.SetMinSize(fyne.NewSize(0, height))
	return padding
}

func verticalSpacer(width float32) *canvas.Rectangle {
	padding := canvas.NewRectangle(color.Transparent)
	padding.SetMinSize(fyne.NewSize(width, 0))
	return padding
}

func showErrorDialog(err error, dialogText string, window fyne.Window) {
	config.Logger.LogError(err.Error())
	dialog.ShowError(errors.New(dialogText), window)
}

// func showInfoDialog(title string, message string, window fyne.Window) {
// 	config.Logger.LogInfof("%s: %s", title, message)
// 	dialog.ShowInformation(title, message, window)
// }

func (g *GuiWrapper) openSecondaryWindow(title string, canvasObj fyne.CanvasObject, promptWindowSize fyne.Size) {
	if g.SecondaryWindow == nil {
		g.SecondaryWindow = g.App.NewWindow(title)
		g.SecondaryWindow.SetOnClosed(func() {
			g.SecondaryWindow = nil
		})
		g.SecondaryWindow.Resize(promptWindowSize)
		g.SecondaryWindow.SetContent(canvasObj)
	}
	g.SecondaryWindow.Show()
	g.SecondaryWindow.RequestFocus()
}
