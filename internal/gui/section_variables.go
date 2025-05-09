package gui

import (
	"fmt"
	"reflect"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/finahdinner/tidal/internal/preferences"
)

func (g *GuiWrapper) getVariablesSection() *fyne.Container {

	header := canvas.NewText("Variables", theme.Color(theme.ColorNameForeground))
	header.TextSize = headerSize

	twitchVariableCopyColumn := container.New(layout.NewVBoxLayout(), widget.NewLabel("Copy"))
	twitchVariableNameColumn := container.New(layout.NewVBoxLayout(), widget.NewLabel("Name"))
	twitchVariableValueColumn := container.New(layout.NewVBoxLayout(), widget.NewLabel("Value"))
	twitchVariableDescriptionColumn := container.New(layout.NewVBoxLayout(), widget.NewLabel("Description"))

	streamVariables := g.Preferences.StreamVariables

	fields := reflect.TypeOf(streamVariables)
	vals := reflect.ValueOf(streamVariables)

	for idx := range reflect.ValueOf(streamVariables).NumField() {

		varName := fields.Field(idx).Name
		varPlaceholderName := generatePlaceholderString(varName)

		nameLabel := widget.NewLabel(varPlaceholderName)
		nameLabelCopyButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
			labelObj := twitchVariableNameColumn.Objects[idx+1]
			if entry, ok := labelObj.(*widget.Entry); ok {
				// TODO - also add a brief popup to confirm copying to clipboard
				g.Clipboard.SetContent(entry.Text)
			}
		})

		twitchVariableNameColumn.Objects = append(
			twitchVariableNameColumn.Objects, nameLabel,
		)
		twitchVariableCopyColumn.Objects = append(
			twitchVariableCopyColumn.Objects, nameLabelCopyButton,
		)

		streamVariable := vals.Field(idx).Interface().(preferences.StreamVariableT)

		streamVariableValue := streamVariable.Value
		if streamVariable.Value == "" {
			streamVariableValue = "N/A"
		}

		twitchVariableValueColumn.Objects = append(
			twitchVariableValueColumn.Objects, widget.NewLabel(streamVariableValue),
		)
		twitchVariableDescriptionColumn.Objects = append(
			twitchVariableDescriptionColumn.Objects, widget.NewLabel(streamVariable.Description),
		)
	}

	updateRateSelect := widget.NewSelect([]string{"10", "20", "30", "60"}, func(s string) {
		// TODO - save updateRateSeconds in preferences
		fmt.Printf("%v seconds\n", s)
	})
	updateRateSelect.SetSelected("10") // TODO set to whatever is saved in preferences, else set a default
	updateRateForm := container.New(
		layout.NewFormLayout(),
		widget.NewLabel("Update every"),
		container.New(layout.NewHBoxLayout(),
			updateRateSelect,
			widget.NewLabel("seconds"),
		),
	)
	updateRateRow := container.New(layout.NewBorderLayout(nil, nil, updateRateForm, nil), updateRateForm)

	return container.NewPadded(container.New(
		layout.NewVBoxLayout(),
		header,
		container.New(
			layout.NewHBoxLayout(),
			twitchVariableCopyColumn,
			twitchVariableNameColumn,
			twitchVariableValueColumn,
			twitchVariableDescriptionColumn,
		),
		horizontalSpacer(20),
		updateRateRow,
	))
}

func generatePlaceholderString(varName string) string {
	return fmt.Sprintf("{{%v}}", varName)
}
