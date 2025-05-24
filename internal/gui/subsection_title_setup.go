package gui

import (
	"log"
	"strconv"

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

	updateFrequencyDigitEntry := widget.NewEntry()
	updateFrequencyTimePeriodSelector := widget.NewSelect([]string{"Seconds", "Minutes"}, nil)

	// TODO - make a function that handles the calculation to seconds
	// call the function when either of these widgets is unchanged
	// if it fails to parse - show a warning label or something
	updateFrequencyDigitEntry.OnChanged = func(s string) {
		digitAsInt, err := strconv.Atoi(s)
		if err != nil {
			saveBtn.Disable()
			return
		}
		saveBtn.Enable()
		log.Println(intVal)
	}

	updateFrequencyTimePeriodSelector.OnChanged = func(s string) {
		if s != "" {
			saveBtn.Enable()
		} else {
			saveBtn.Disable()
		}
	}

	updateFrequencyContainer := container.New(
		layout.NewHBoxLayout(),
		updateFrequencyDigitEntry, updateFrequencyTimePeriodSelector,
	)

	return container.New(
		layout.NewFormLayout(),
		widget.NewLabel("Title Template"),
		titleTemplateEntry,
		layout.NewSpacer(),
		tipLabel,
		widget.NewLabel("Update every "),
		updateFrequencyContainer,
		horizontalSpacer(3),
		layout.NewSpacer(),
		layout.NewSpacer(),
		saveBtn,
	)
}
