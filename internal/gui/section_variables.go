package gui

import (
	"fmt"
	"reflect"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/finahdinner/tidal/internal/config"
)

const (
	streamVarNamePlaceholderPrefix = "{{"
	streamVarNamePlaceholderSuffix = "}}"
	streamVarPlaceholderValue      = "-"
)

func (g *GuiWrapper) getVariablesSection() *fyne.Container {

	header := canvas.NewText("Variables", theme.Color(theme.ColorNameForeground))
	header.TextSize = headerSize

	twitchVariableCopyColumn := container.New(layout.NewVBoxLayout(), widget.NewLabel("Copy"))
	twitchVariableNameColumn := container.New(layout.NewVBoxLayout(), widget.NewLabel("Name"))
	twitchVariableValueColumn := container.New(layout.NewVBoxLayout(), widget.NewLabel("Value"))
	twitchVariableDescriptionColumn := container.New(layout.NewVBoxLayout(), widget.NewLabel("Description"))

	streamVariables := &config.Preferences.StreamVariables

	fields := reflect.TypeOf(*streamVariables)
	vals := reflect.ValueOf(*streamVariables)

	for idx := range reflect.ValueOf(*streamVariables).NumField() {

		varName := fields.Field(idx).Name
		varPlaceholderName := generatePlaceholderString(varName)

		nameLabel := widget.NewLabel(varPlaceholderName)
		nameLabelCopyButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
			labelObj := twitchVariableNameColumn.Objects[idx+1]
			if entry, ok := labelObj.(*widget.Label); ok {
				// TODO - also add a brief popup to confirm copying to clipboard
				g.App.Clipboard().SetContent(entry.Text)
			}
		})

		twitchVariableNameColumn.Objects = append(
			twitchVariableNameColumn.Objects, nameLabel,
		)
		twitchVariableCopyColumn.Objects = append(
			twitchVariableCopyColumn.Objects, nameLabelCopyButton,
		)

		streamVariable := vals.Field(idx).Interface().(config.StreamVariableT)

		twitchVariableValueColumn.Objects = append(
			twitchVariableValueColumn.Objects, widget.NewLabel(valueOrPlaceholderValue(streamVariable.Value)),
		)
		twitchVariableDescriptionColumn.Objects = append(
			twitchVariableDescriptionColumn.Objects, widget.NewLabel(streamVariable.Description),
		)
	}

	variablesSection := container.NewPadded(container.New(
		layout.NewVBoxLayout(),
		header,
		container.New(
			layout.NewHBoxLayout(),
			twitchVariableCopyColumn,
			twitchVariableNameColumn,
			twitchVariableValueColumn,
			twitchVariableDescriptionColumn,
		),
	))

	// set up a listener to update widgets whenever the ticker updates stream variables
	go func() {
		for range updateVariablesSectionSignal {
			config.Logger.LogInfo("updating stream variable widgets")

			for rowIdx := 1; rowIdx < len(twitchVariableValueColumn.Objects); rowIdx++ {

				varPlaceholderName := twitchVariableNameColumn.Objects[rowIdx].(*widget.Label).Text
				varName := getVarNameFromPlaceholderString(varPlaceholderName)
				streamVariablesV := reflect.ValueOf(config.Preferences.StreamVariables)
				streamVariablesT := streamVariablesV.Type()
				streamVariableObjType := reflect.TypeOf(config.StreamVariableT{})

				for fieldIdx := range streamVariablesT.NumField() {
					fieldName := streamVariablesT.Field(fieldIdx).Name
					if fieldName == varName {
						// populate the value on this row with StreamVariableT.Value
						valueField := streamVariablesV.Field(fieldIdx)
						if valueField.Kind() == reflect.Struct && valueField.Type() == streamVariableObjType {
							valueField := valueField.FieldByName("Value")
							if valueField.IsValid() {
								newValue := valueField.String()
								fyne.Do(func() {
									twitchVariableValueColumn.Objects[rowIdx].(*widget.Label).SetText(valueOrPlaceholderValue(newValue))
									twitchVariableValueColumn.Objects[rowIdx].Refresh()
								})
								config.Logger.LogInfof("updated field name %v to value %v", fieldName, newValue)
							}
						}
					}
				}
			}
		}
	}()

	return variablesSection
}

func generatePlaceholderString(varName string) string {
	return fmt.Sprintf("%v%v%v", streamVarNamePlaceholderPrefix, varName, streamVarNamePlaceholderSuffix)
}

func getVarNameFromPlaceholderString(placeholderString string) string {
	return strings.Replace(strings.Replace(placeholderString, streamVarNamePlaceholderPrefix, "", 1), streamVarNamePlaceholderSuffix, "", 1)
}

func valueOrPlaceholderValue(txt string) string {
	if txt == "" {
		return streamVarPlaceholderValue
	}
	return txt
}
