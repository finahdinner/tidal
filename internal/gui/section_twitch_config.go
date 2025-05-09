package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func (g *GuiWrapper) getTwitchConfigSection() *fyne.Container {
	return container.New(layout.NewVBoxLayout(), widget.NewLabel("getTwitchConfigSection"))
}
