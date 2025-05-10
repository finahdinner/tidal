package gui

import (
	"fmt"
	"log"
	"reflect"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/finahdinner/tidal/internal/preferences"
)

var variablesSection *fyne.Container

func (g *GuiWrapper) getVariablesSection() *fyne.Container {

	if variablesSection != nil {
		log.Println("variablesSection already exists")
		return variablesSection
	}

	header := canvas.NewText("Variables", theme.Color(theme.ColorNameForeground))
	header.TextSize = headerSize

	twitchVariableCopyColumn := container.New(layout.NewVBoxLayout(), widget.NewLabel("Copy"))
	twitchVariableNameColumn := container.New(layout.NewVBoxLayout(), widget.NewLabel("Name"))
	twitchVariableValueColumn := container.New(layout.NewVBoxLayout(), widget.NewLabel("Value"))
	twitchVariableDescriptionColumn := container.New(layout.NewVBoxLayout(), widget.NewLabel("Description"))

	streamVariables := &preferences.Preferences.StreamVariables

	fields := reflect.TypeOf(*streamVariables)
	vals := reflect.ValueOf(*streamVariables)

	for idx := range reflect.ValueOf(*streamVariables).NumField() {

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
		asInt, err := strconv.Atoi(s)
		if err != nil {
			log.Printf("unable to parse update rate as int - err: %v", err)
			return
		}
		preferences.Preferences.VariableUpdateInterval = asInt
		preferences.SavePreferences()
		log.Printf("saved update frequency as %v seconds", asInt)
	})

	if preferences.Preferences.VariableUpdateInterval > 0 {
		updateRateSelect.SetSelected(strconv.Itoa(preferences.Preferences.VariableUpdateInterval))
	}

	updateRateForm := container.New(
		layout.NewFormLayout(),
		widget.NewLabel("Update every"),
		container.New(layout.NewHBoxLayout(),
			updateRateSelect,
			widget.NewLabel("seconds"),
		),
	)
	updateRateRow := container.New(layout.NewBorderLayout(nil, nil, updateRateForm, nil), updateRateForm)

	variablesSection = container.NewPadded(container.New(
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

	return variablesSection
}

func generatePlaceholderString(varName string) string {
	return fmt.Sprintf("{{%v}}", varName)
}
