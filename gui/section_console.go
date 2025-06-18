package gui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/finahdinner/tidal/config"
)

type ActivityConsoleT struct {
	box    *fyne.Container
	scroll *container.Scroll
	stack  *fyne.Container
}

var ActivityConsole *ActivityConsoleT

var consoleSection fyne.CanvasObject

func init() {
	if ActivityConsole == nil {
		ActivityConsole = NewActivityConsole()
	}
}

func NewActivityConsole() *ActivityConsoleT {
	consoleBox := container.New(layout.NewVBoxLayout())
	consoleBoxBg := canvas.NewRectangle(color.Black)
	consoleScroll := container.NewVScroll(consoleBox)
	consoleStack := container.New(layout.NewStackLayout(), consoleBoxBg, consoleScroll)
	return &ActivityConsoleT{consoleBox, consoleScroll, consoleStack}
}

// Append a new line to the activity console
func (ac *ActivityConsoleT) pushToConsole(text string) error {
	if err := config.ConsoleLogger.PushToLog(text); err != nil {
		return err
	}
	line := widget.NewRichTextFromMarkdown(fmt.Sprintf("`%s`", text))
	line.Wrapping = fyne.TextWrapWord
	line.Scroll = fyne.ScrollNone
	fyne.Do(func() {
		ac.box.Objects = append(ac.box.Objects, line)
		ac.scroll.ScrollToBottom()
		ac.box.Refresh()
	})
	return nil
}

// Clears the console and closes the console log file
func (ac *ActivityConsoleT) clearConsole() {
	config.ConsoleLogger.DeleteInstance()
	ac.box.Objects = []fyne.CanvasObject{}
	fyne.Do(func() {
		ac.box.Refresh()
	})
}

func (g *GuiWrapper) getConsoleSection() fyne.CanvasObject {

	if ActivityConsole == nil {
		ActivityConsole = NewActivityConsole()
	}

	if consoleSection != nil {
		return consoleSection
	}

	openHelpFunc := func() {
		g.openSecondaryWindow("Console Help", getConsoleHelpSection(), &helpWindowSize)
	}

	consoleSection = sectionWrapper(
		"Console",
		nil,
		openHelpFunc,
		ActivityConsole.stack,
		false,
		false,
		true,
	)

	return consoleSection
}

func getConsoleHelpSection() fyne.CanvasObject {

	markdownLines := []string{
		"- When the **Start Tidal** button is pressed, each title update will be logged to the Console.",
		"- When the **Stop Tidal** button is pressed, or Tidal stops updating the title for any other reason, the Console will be cleared.",
		"- However, all logs will exist in a log file found in the logs folder - use the **Config Folder** button to open this folder.",
	}

	return helpSectionWrapper("Console Help", markdownLines)
}
