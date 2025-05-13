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
	"github.com/finahdinner/tidal/internal/helpers"
)

const (
	varNamePlaceholderPrefix = "{{"
	varNamePlaceholderSuffix = "}}"
	varPlaceholderValue      = "-"

	multilineEntryHeight              = 3
	promptEmptyStreamValuePlaceholder = "<<N/A>>"
)

var promptWindowSize fyne.Size = fyne.NewSize(600, 400)

func (g *GuiWrapper) getVariablesSection() *fyne.Container {

	twitchVariablesHeader := canvas.NewText("Twitch Variables", theme.Color(theme.ColorNameForeground))
	twitchVariablesHeader.TextSize = headerSize

	twitchVariableCopyColumn := container.New(layout.NewVBoxLayout(), widget.NewLabel("Copy"))
	twitchVariableNameColumn := container.New(layout.NewVBoxLayout(), widget.NewLabel("Name"))
	twitchVariableValueColumn := container.New(layout.NewVBoxLayout(), widget.NewLabel("Value"))
	twitchVariableDescriptionColumn := container.New(layout.NewVBoxLayout(), widget.NewLabel("Description"))

	twitchVariables := &config.Preferences.TwitchVariables

	fields := reflect.TypeOf(*twitchVariables)
	vals := reflect.ValueOf(*twitchVariables)

	for idx := range vals.NumField() {

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

		twitchVariable := vals.Field(idx).Interface().(config.TwitchVariableT)

		twitchVariableValueColumn.Objects = append(
			twitchVariableValueColumn.Objects, widget.NewLabel(valueOrPlaceholderValue(twitchVariable.Value)),
		)
		twitchVariableDescriptionColumn.Objects = append(
			twitchVariableDescriptionColumn.Objects, widget.NewLabel(twitchVariable.Description),
		)
	}

	// set up a listener to update widgets whenever the ticker updates twitch variables
	go func() {
		for range updateTwitchVariablesSectionSignal {
			config.Logger.LogInfo("updating stream variable widgets")

			for rowIdx := 1; rowIdx < len(twitchVariableValueColumn.Objects); rowIdx++ {

				varPlaceholderName := twitchVariableNameColumn.Objects[rowIdx].(*widget.Label).Text
				varName := getVarNameFromPlaceholderString(varPlaceholderName)
				twitchVariablesV := reflect.ValueOf(config.Preferences.TwitchVariables)
				twitchVariablesT := twitchVariablesV.Type()
				twitchVariableObjType := reflect.TypeOf(config.TwitchVariableT{})

				for fieldIdx := range twitchVariablesT.NumField() {
					fieldName := twitchVariablesT.Field(fieldIdx).Name
					if fieldName == varName {
						// populate the value on this row with TwitchVariableT.Value
						valueField := twitchVariablesV.Field(fieldIdx)
						if valueField.Kind() == reflect.Struct && valueField.Type() == twitchVariableObjType {
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

	aiGeneratedVariablesHeader := canvas.NewText("AI-generated Variables", theme.Color(theme.ColorNameForeground))
	aiGeneratedVariablesHeader.TextSize = headerSize

	aiGeneratedVariableCopyColumn := container.New(layout.NewVBoxLayout(), widget.NewLabel("Copy"))
	aiGeneratedVariableNameColumn := container.New(layout.NewVBoxLayout(), widget.NewLabel("Name"))
	aiGeneratedVariableValueColumn := container.New(layout.NewVBoxLayout(), widget.NewLabel("Value"))
	aiGeneratedVariablePromptColumn := container.New(layout.NewVBoxLayout(), widget.NewLabel("Prompt"))
	aiGeneratedVariableRemoveColumn := container.New(layout.NewVBoxLayout(), layout.NewSpacer())

	aiGeneratedVariables := config.Preferences.AiGeneratedVariables

	for idx, aiGenVar := range aiGeneratedVariables {
		name := aiGenVar.Name

		nameLabel := widget.NewLabel(generatePlaceholderString(name))
		nameLabelCopyButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
			labelObj := aiGeneratedVariableNameColumn.Objects[idx+1]
			if entry, ok := labelObj.(*widget.Label); ok {
				// TODO - also add a brief popup to confirm copying to clipboard
				g.App.Clipboard().SetContent(entry.Text)
			}
		})

		aiGeneratedVariableNameColumn.Objects = append(
			aiGeneratedVariableNameColumn.Objects,
			nameLabel,
		)

		aiGeneratedVariableCopyColumn.Objects = append(
			aiGeneratedVariableCopyColumn.Objects, nameLabelCopyButton,
		)

		aiGeneratedVariableValueColumn.Objects = append(
			aiGeneratedVariableValueColumn.Objects,
			widget.NewLabel(valueOrPlaceholderValue(aiGenVar.Value)),
		)

		aiGeneratedVariablePromptColumn.Objects = append(
			aiGeneratedVariablePromptColumn.Objects,
			widget.NewLabel(valueOrPlaceholderValue(aiGenVar.Prompt)),
		)

	}

	twitchVariablesStringReplacer := getTwitchVariablesStringReplacer(*twitchVariables)
	addAiGeneratedVariableBtn := widget.NewButton("Add Variable", func() {
		saveBtn := widget.NewButton("Save", nil) // TODO - save to config.json
		promptEntryMain := getMultilineEntry("prompt entry main", saveBtn)
		promptEntrySuffix := getMultilineEntry("prompt entry suffix", saveBtn)
		promptPreview := getMultilinePreview([]*widget.Entry{promptEntryMain, promptEntrySuffix}, twitchVariablesStringReplacer, saveBtn)
		c := container.New(
			layout.NewFormLayout(),
			widget.NewLabel("Main Prompt"),
			promptEntryMain,
			widget.NewLabel("Prompt suffix"),
			promptEntrySuffix,
			widget.NewLabel("Preview"),
			promptPreview,
		)
		g.openSecondaryWindow("Add AI-Generated Variable", c, promptWindowSize)
	})
	addAiGeneratedVariableBtnRow := container.New(layout.NewBorderLayout(nil, nil, addAiGeneratedVariableBtn, nil), addAiGeneratedVariableBtn)

	return container.NewPadded(container.NewScroll(container.New(
		layout.NewVBoxLayout(),
		twitchVariablesHeader,
		container.New(
			layout.NewHBoxLayout(),
			twitchVariableCopyColumn,
			twitchVariableNameColumn,
			twitchVariableValueColumn,
			twitchVariableDescriptionColumn,
		),
		horizontalSpacer(8),
		aiGeneratedVariablesHeader,
		container.New(
			layout.NewVBoxLayout(),
			container.New(
				layout.NewHBoxLayout(),
				aiGeneratedVariableCopyColumn,
				aiGeneratedVariableNameColumn,
				aiGeneratedVariableValueColumn,
				aiGeneratedVariablePromptColumn,
				aiGeneratedVariableRemoveColumn,
			),
			addAiGeneratedVariableBtnRow,
		),
	)))
}

func getMultilineEntry(text string, saveBtn *widget.Button) *widget.Entry {
	e := widget.NewMultiLineEntry()
	e.SetText(text)
	e.SetMinRowsVisible(multilineEntryHeight)
	e.OnChanged = func(_ string) {
		saveBtn.Enable()
	}
	return e
}

func getMultilinePreview(parentEntryWidgets []*widget.Entry, variableReplacer *strings.Replacer, saveBtn *widget.Button) *widget.Entry {
	e := widget.NewMultiLineEntry()
	e.SetText(buildStringFromEntryWidgets(parentEntryWidgets, variableReplacer))
	e.SetMinRowsVisible(3)
	e.Disable()
	for _, entry := range parentEntryWidgets {
		entry.OnChanged = func(text string) {
			e.SetText(buildStringFromEntryWidgets(parentEntryWidgets, variableReplacer))
		}
	}
	return e
}

func buildStringFromEntryWidgets(entryWidgets []*widget.Entry, variableReplacer *strings.Replacer) string {
	promptParts := []string{}
	for _, e := range entryWidgets {
		if e.Text != "" {
			promptParts = append(promptParts, e.Text)
		}
	}
	concatenatedPrompt := strings.Join(promptParts, "\n")
	return variableReplacer.Replace(concatenatedPrompt)
}

func getTwitchVariablesStringReplacer(twitchVariables config.TwitchVariablesT) *strings.Replacer {
	replacementList := []string{}
	twitchVariablesMap := helpers.GenerateMapFromHomogenousStruct[config.TwitchVariablesT, config.TwitchVariableT](twitchVariables)
	for name, val := range twitchVariablesMap {
		value := val.Value
		if value == "" {
			value = "<<N/A>>"
		}
		replacementList = append(replacementList, generatePlaceholderString(name), value)
	}
	return strings.NewReplacer(replacementList...)
}

func generatePlaceholderString(varName string) string {
	return fmt.Sprintf("%v%v%v", varNamePlaceholderPrefix, varName, varNamePlaceholderSuffix)
}

func getVarNameFromPlaceholderString(placeholderString string) string {
	return strings.Replace(strings.Replace(placeholderString, varNamePlaceholderPrefix, "", 1), varNamePlaceholderSuffix, "", 1)
}

func valueOrPlaceholderValue(txt string) string {
	if txt == "" {
		return varPlaceholderValue
	}
	return txt
}
