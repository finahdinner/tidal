package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

const (
	fontSize   = 15
	headerSize = 18
)

type tidalTheme struct{}

var _ fyne.Theme = (*tidalTheme)(nil) // assert that it implements the fyne.Theme interface

func (t tidalTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameScrollBar {
		return color.RGBA{25, 175, 160, 255} // blue
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (t tidalTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t tidalTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t tidalTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return fontSize
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameInputRadius:
		return 10
	case theme.SizeNameInputBorder:
		return 2
	case theme.SizeNameSelectionRadius:
		return 5
	case theme.SizeNameScrollBarSmall:
		return 10
	default:
		return theme.DefaultTheme().Size(name)
	}
}
