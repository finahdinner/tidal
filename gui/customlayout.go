package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/finahdinner/tidal/config"
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

	var customDialog dialog.Dialog

	showLogsBtn := widget.NewButton("Show Logs", nil)
	dismissBtn := widget.NewButton("Dismiss", func() {
		customDialog.Dismiss()
	})
	btnRow := container.New(
		layout.NewHBoxLayout(),
		layout.NewSpacer(),
		showLogsBtn, dismissBtn,
		layout.NewSpacer(),
	)

	customContent := container.New(
		layout.NewVBoxLayout(),
		widget.NewLabel(dialogText),
		btnRow,
	)
	customDialog = dialog.NewCustomWithoutButtons("Error", customContent, window)
	customDialog.Show()
}

func showInfoDialog(title string, message string, window fyne.Window) {
	config.Logger.LogInfof("%s: %s", title, message)
	dialog.ShowInformation(title, message, window)
}

func (g *GuiWrapper) openSecondaryWindow(title string, canvasObj fyne.CanvasObject, promptWindowSize *fyne.Size) {
	if g.SecondaryWindow == nil {
		g.SecondaryWindow = g.App.NewWindow(title)
		g.SecondaryWindow.SetOnClosed(func() {
			g.SecondaryWindow = nil
		})
		if promptWindowSize != nil {
			g.SecondaryWindow.Resize(*promptWindowSize)
		}
		g.SecondaryWindow.SetContent(canvasObj)
	}
	g.SecondaryWindow.Show()
	g.SecondaryWindow.RequestFocus()
}

func (g *GuiWrapper) closeSecondaryWindow() {
	if g.SecondaryWindow != nil {
		g.SecondaryWindow.Close()
		g.SecondaryWindow = nil
	}
}

func newVariablesDetectedWidget() *widget.RichText {
	variablesDetectedWidget := widget.NewRichText()
	variablesDetectedWidget.Scroll = fyne.ScrollHorizontalOnly
	return variablesDetectedWidget
}
