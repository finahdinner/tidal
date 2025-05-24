package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var titleSetupWindowSize fyne.Size = fyne.NewSize(600, 1) // height 1 lets the layout determine the height

func (g *GuiWrapper) getTitleSetupSubsection() *fyne.Container {

	saveBtn := widget.NewButton("Save", nil)
	titleTemplateEntry := getMultilineEntry("", saveBtn, 6)
	tipLabel := widget.NewLabel(`You can use any Variables in your title template
Access them using {{VariableName}}`)

	return container.New(
		layout.NewFormLayout(),
		widget.NewLabel("Title Template"),
		titleTemplateEntry,
		layout.NewSpacer(),
		tipLabel,
		layout.NewSpacer(),
		saveBtn,
	)
}
