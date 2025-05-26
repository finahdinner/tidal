package gui

import (
	"errors"
	"fmt"
	"image/color"
	"math"
	"reflect"
	"strings"
	"time"

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
	multilineEntryHeight              = 5
	promptEmptyStreamValuePlaceholder = "<<N/A>>"
)

var promptWindowSize fyne.Size = fyne.NewSize(600, 1) // height 1 lets the layout determine the height

func (g *GuiWrapper) getVariablesSection() *fyne.Container {

	twitchVariablesHeader := canvas.NewText("Twitch Variables", theme.Color(theme.ColorNameForeground))
	twitchVariablesHeader.TextSize = headerSize

	twitchVariablesSettingsButton := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		g.openSecondaryWindow("Twitch Configuration", g.getTwitchConfigSubsection(), &twitchConfigWindowSize)
	})
	twitchVariablesHeaderRow := container.New(
		layout.NewHBoxLayout(),
		twitchVariablesSettingsButton,
		verticalSpacer(1),
		twitchVariablesHeader,
	)

	twitchVariableCopyColumn := container.New(layout.NewVBoxLayout(), widget.NewLabel("Copy"))
	twitchVariableNameColumn := container.New(layout.NewVBoxLayout(), widget.NewLabel("Name"))
	twitchVariableValueColumn := container.New(layout.NewVBoxLayout(), widget.NewLabel("Value"))
	twitchVariableDescriptionColumn := container.New(layout.NewVBoxLayout(), widget.NewLabel("Description"))

	twitchVariables := &config.Preferences.TwitchVariables
	twitchVariablesStringReplacer := getTwitchVariablesStringReplacer(*twitchVariables)

	g.populateRowsWithExistingTwitchVariables(
		twitchVariables,
		twitchVariableCopyColumn,
		twitchVariableNameColumn,
		twitchVariableValueColumn,
		twitchVariableDescriptionColumn,
	)

	// set up a listener to update widgets whenever the ticker updates twitch variables
	go func() {
		for range updateTwitchVariablesSectionSignal {
			config.Logger.LogInfo("updating stream variable widgets")

			for rowIdx := 1; rowIdx < len(twitchVariableValueColumn.Objects); rowIdx++ {

				varPlaceholderName := twitchVariableNameColumn.Objects[rowIdx].(*widget.Label).Text
				varName := helpers.GetVarNameFromPlaceholderString(varPlaceholderName)
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

	aiGeneratedVariablesSettingsButton := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		g.openSecondaryWindow("LLM Configuration", g.getLlmConfigSubsection(), &llmConfigWindowSize)
	})
	aiGeneratedVariablesHeaderRow := container.New(
		layout.NewHBoxLayout(),
		aiGeneratedVariablesSettingsButton,
		verticalSpacer(1),
		aiGeneratedVariablesHeader,
	)

	aiGeneratedVariableCopyColumn := container.New(layout.NewVBoxLayout(), widget.NewLabel("Copy"))
	aiGeneratedVariableNameColumn := container.New(layout.NewVBoxLayout(), widget.NewLabel("Name"))
	aiGeneratedEditColumn := container.New(layout.NewVBoxLayout(), layout.NewSpacer())
	aiGeneratedVariableRemoveColumn := container.New(layout.NewVBoxLayout(), layout.NewSpacer())

	aiGeneratedVariables := config.Preferences.AiGeneratedVariables

	g.populateRowsWithExistingAiGeneratedVariables(
		aiGeneratedVariables,
		twitchVariablesStringReplacer,
		aiGeneratedVariableCopyColumn,
		aiGeneratedVariableNameColumn,
		aiGeneratedEditColumn,
		aiGeneratedVariableRemoveColumn,
	)

	addAiGeneratedVariableBtn := widget.NewButton("Add Variable", func() {
		g.openSecondaryWindow(
			"Add AI-Generated Variable",
			g.getAiGeneratedVariableSection(
				false,
				twitchVariablesStringReplacer,
				"",
				"Add your main prompt here",
				"Add a suffix to your prompt here",
				aiGeneratedVariableCopyColumn,
				aiGeneratedVariableNameColumn,
				aiGeneratedEditColumn,
				aiGeneratedVariableRemoveColumn,
			),
			&promptWindowSize,
		)
	})
	addAiGeneratedVariableBtnRow := container.New(layout.NewBorderLayout(nil, nil, addAiGeneratedVariableBtn, nil), addAiGeneratedVariableBtn)

	return container.NewPadded(container.NewScroll(container.New(
		layout.NewVBoxLayout(),
		twitchVariablesHeaderRow,
		container.New(
			layout.NewHBoxLayout(),
			twitchVariableCopyColumn,
			twitchVariableNameColumn,
			twitchVariableValueColumn,
			twitchVariableDescriptionColumn,
		),
		horizontalSpacer(8),
		aiGeneratedVariablesHeaderRow,
		container.New(
			layout.NewVBoxLayout(),
			container.New(
				layout.NewHBoxLayout(),
				aiGeneratedVariableCopyColumn,
				aiGeneratedVariableNameColumn,
				// aiGeneratedVariableValueColumn,
				aiGeneratedEditColumn,
				aiGeneratedVariableRemoveColumn,
			),
			horizontalSpacer(3),
			addAiGeneratedVariableBtnRow,
		),
	)))
}

func (g *GuiWrapper) getNewCopyButton(rowIdx int, variableNameColumn *fyne.Container) *fyne.Container {
	nameLabelCopyButtonBg := canvas.NewRectangle(theme.Color(theme.ColorNameButton))
	nameLabelCopyButtonBtn := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), nil)
	nameLabelCopyButtonBtn.OnTapped = func() {
		labelObj := variableNameColumn.Objects[rowIdx+1]
		if entry, ok := labelObj.(*widget.Label); ok {
			g.App.Clipboard().SetContent(entry.Text)
			nameLabelCopyButtonBg.FillColor = color.RGBA{0, 255, 0, 255} // green
		} else {
			nameLabelCopyButtonBg.FillColor = color.RGBA{255, 0, 0, 255} // red
		}
		time.AfterFunc(500*time.Millisecond, func() {
			fyne.Do(func() {
				nameLabelCopyButtonBg.FillColor = theme.Color(theme.ColorNameButton)
				nameLabelCopyButtonBtn.Refresh()
			})
		})
	}
	return container.New(layout.NewStackLayout(), nameLabelCopyButtonBg, nameLabelCopyButtonBtn)
}

func getMultilineEntry(text string, saveBtn *widget.Button, lineHeight int) *widget.Entry {
	e := widget.NewMultiLineEntry()
	e.SetText(text)
	e.SetMinRowsVisible(lineHeight)
	e.OnChanged = func(_ string) {
		saveBtn.Enable()
	}
	return e
}

func getMultilinePreview(parentEntryWidgets []*widget.Entry, variableReplacer *strings.Replacer, saveBtn *widget.Button, lineHeight int) *widget.Entry {
	e := getMultilineEntry(
		buildStringFromEntryWidgets(parentEntryWidgets, variableReplacer),
		saveBtn,
		lineHeight,
	)
	for _, entry := range parentEntryWidgets {
		entry.OnChanged = func(text string) {
			e.SetText(buildStringFromEntryWidgets(parentEntryWidgets, variableReplacer))
		}
	}
	return e
}

func (g *GuiWrapper) populateRowsWithExistingTwitchVariables(
	twitchVariables *config.TwitchVariablesT,
	twitchVariableCopyColumn *fyne.Container,
	twitchVariableNameColumn *fyne.Container,
	twitchVariableValueColumn *fyne.Container,
	twitchVariableDescriptionColumn *fyne.Container,
) {

	twitchVariableCopyColumn.Objects = twitchVariableCopyColumn.Objects[:1]
	twitchVariableNameColumn.Objects = twitchVariableNameColumn.Objects[:1]
	twitchVariableValueColumn.Objects = twitchVariableValueColumn.Objects[:1]
	twitchVariableDescriptionColumn.Objects = twitchVariableDescriptionColumn.Objects[:1]

	fields := reflect.TypeOf(*twitchVariables)
	vals := reflect.ValueOf(*twitchVariables)

	for idx := range vals.NumField() {

		varName := fields.Field(idx).Name
		varPlaceholderName := helpers.GenerateVarPlaceholderString(varName)

		nameLabel := widget.NewLabel(varPlaceholderName)
		nameLabelCopyButton := g.getNewCopyButton(idx, twitchVariableNameColumn)

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
}

func (g *GuiWrapper) populateRowsWithExistingAiGeneratedVariables(
	aiGeneratedVariables []config.LlmVariableT,
	twitchVariablesStringReplacer *strings.Replacer,
	aiGeneratedVariableCopyColumn *fyne.Container,
	aiGeneratedVariableNameColumn *fyne.Container,
	aiGeneratedEditColumn *fyne.Container,
	aiGeneratedVariableRemoveColumn *fyne.Container,
) {

	aiGeneratedVariableCopyColumn.Objects = aiGeneratedVariableCopyColumn.Objects[:1]
	aiGeneratedVariableNameColumn.Objects = aiGeneratedVariableNameColumn.Objects[:1]
	aiGeneratedEditColumn.Objects = aiGeneratedEditColumn.Objects[:1]
	aiGeneratedVariableRemoveColumn.Objects = aiGeneratedVariableRemoveColumn.Objects[:1]

	for idx, aiGenVar := range aiGeneratedVariables {
		name := aiGenVar.Name

		nameLabel := widget.NewLabel(helpers.GenerateVarPlaceholderString(name))
		nameLabelCopyButton := g.getNewCopyButton(idx, aiGeneratedVariableNameColumn)

		aiGeneratedVariableNameColumn.Objects = append(
			aiGeneratedVariableNameColumn.Objects,
			nameLabel,
		)

		aiGeneratedVariableCopyColumn.Objects = append(
			aiGeneratedVariableCopyColumn.Objects, nameLabelCopyButton,
		)

		aiGeneratedEditColumn.Objects = append(
			aiGeneratedEditColumn.Objects,
			widget.NewButton("Edit", func() {
				g.openSecondaryWindow(
					"Edit AI-Generated Variable",
					g.getAiGeneratedVariableSection(
						true,
						twitchVariablesStringReplacer,
						name,
						aiGenVar.PromptMain,
						aiGenVar.PromptSuffix,
						aiGeneratedVariableCopyColumn,
						aiGeneratedVariableNameColumn,
						aiGeneratedEditColumn,
						aiGeneratedVariableRemoveColumn,
					),
					&promptWindowSize,
				)
			}),
		)

		aiGeneratedVariableRemoveColumn.Objects = append(
			aiGeneratedVariableRemoveColumn.Objects,
			widget.NewButton("Remove", func() {
				variableIdx := -1
				existingVars := config.Preferences.AiGeneratedVariables
				for idx, val := range existingVars {
					if val.Name == name {
						variableIdx = idx
						break
					}
				}
				if variableIdx == -1 {
					showErrorDialog(
						fmt.Errorf("unable to find variable with name %q - cannot remove", name),
						fmt.Sprintf("Unable to remove variable %q - cannot be found in existing variables", name),
						g.SecondaryWindow,
					)
					return
				}
				// remove the variable at that index
				config.Preferences.AiGeneratedVariables = append(
					existingVars[:variableIdx],
					existingVars[variableIdx+1:]...,
				)
				config.SavePreferences()

				g.populateRowsWithExistingAiGeneratedVariables(
					config.Preferences.AiGeneratedVariables,
					twitchVariablesStringReplacer,
					aiGeneratedVariableCopyColumn,
					aiGeneratedVariableNameColumn,
					aiGeneratedEditColumn,
					aiGeneratedVariableRemoveColumn,
				)

				showInfoDialog(
					"Variable successfully removed",
					fmt.Sprintf("AI-generated variable %q has successfully been removed.", name),
					g.PrimaryWindow,
				)
			}),
		)
	}
}

func (g *GuiWrapper) getAiGeneratedVariableSection(
	editExisting bool,
	twitchVariablesStringReplacer *strings.Replacer,
	variableName string,
	promptMainText string,
	promptSuffixText string,
	aiGeneratedVariableCopyColumn *fyne.Container,
	aiGeneratedVariableNameColumn *fyne.Container,
	aiGeneratedEditColumn *fyne.Container,
	aiGeneratedVariableRemoveColumn *fyne.Container,
) *fyne.Container {
	saveBtn := widget.NewButton("Save", nil)
	variableNameEntry := widget.NewEntry()
	variableNameEntry.SetText(variableName)

	// if editing an existing variable, don't let the user rename it
	if editExisting {
		variableNameEntry.Disable()
	}

	promptEntryMain := getMultilineEntry(promptMainText, saveBtn, multilineEntryHeight)
	promptEntrySuffix := getMultilineEntry(promptSuffixText, saveBtn, multilineEntryHeight)
	promptPreviewLineHeight := int(math.Trunc((1.5 * multilineEntryHeight)))
	promptPreview := getMultilinePreview(
		[]*widget.Entry{promptEntryMain, promptEntrySuffix},
		twitchVariablesStringReplacer,
		saveBtn,
		promptPreviewLineHeight,
	)

	saveBtn.OnTapped = func() {
		varName := strings.TrimSpace(variableNameEntry.Text)

		// if creating a variable and the name is empty
		if varName == "" {
			showErrorDialog(
				errors.New("variable name is empty - cannot save"),
				"Unable to save - variable name must not be empty",
				g.SecondaryWindow,
			)
			return
		}

		existingVariableNamesLower := make(map[string]struct{})
		for _, variable := range config.Preferences.AiGeneratedVariables {
			existingVariableNamesLower[strings.ToLower(variable.Name)] = struct{}{}
		}
		twitchVariablesMap := helpers.GenerateMapFromHomogenousStruct[
			config.TwitchVariablesT, config.TwitchVariableT,
		](config.Preferences.TwitchVariables)
		for name := range twitchVariablesMap {
			existingVariableNamesLower[strings.ToLower(name)] = struct{}{}
		}

		// if creating a variable and the name is taken
		if !editExisting {
			if _, exists := existingVariableNamesLower[strings.ToLower(varName)]; exists {
				showErrorDialog(
					fmt.Errorf("variable name %q already exists - choose a new name", varName),
					fmt.Sprintf("Unable to save - variable name %q already exists", varName),
					g.SecondaryWindow,
				)
				return
			}
		}

		promptMainText := strings.TrimSpace(promptEntryMain.Text)
		promptSuffixText := strings.TrimSpace(promptEntrySuffix.Text)

		if promptMainText == "" {
			showErrorDialog(
				errors.New("main prompt must not be empty - cannot save"),
				"Unable to save - main prompt must not be empty",
				g.SecondaryWindow,
			)
			return
		}

		if editExisting {
			existingVarIdx := -1
			for idx, val := range config.Preferences.AiGeneratedVariables {
				if val.Name == varName {
					existingVarIdx = idx
					break
				}
			}
			if existingVarIdx == -1 {
				showErrorDialog(
					fmt.Errorf("unable to find existing variable with name %q", varName),
					fmt.Sprintf("Unable to save variable - existing variable with name %q not found", varName),
					g.SecondaryWindow,
				)
				return
			}
			config.Preferences.AiGeneratedVariables[existingVarIdx] = config.LlmVariableT{
				Name:         varName,
				Value:        "", // reset the value
				PromptMain:   promptMainText,
				PromptSuffix: promptSuffixText,
			}
		} else {
			config.Preferences.AiGeneratedVariables = append(
				config.Preferences.AiGeneratedVariables,
				config.LlmVariableT{
					Name:         varName,
					Value:        "",
					PromptMain:   promptMainText,
					PromptSuffix: promptSuffixText,
				},
			)
		}
		config.SavePreferences()

		g.populateRowsWithExistingAiGeneratedVariables(
			config.Preferences.AiGeneratedVariables,
			twitchVariablesStringReplacer,
			aiGeneratedVariableCopyColumn,
			aiGeneratedVariableNameColumn,
			aiGeneratedEditColumn,
			aiGeneratedVariableRemoveColumn,
		)

		g.closeSecondaryWindow()
		showInfoDialog(
			"Variable successfully saved",
			fmt.Sprintf("AI-generated variable %q has successfully been saved.", varName),
			g.PrimaryWindow,
		)
	}

	form := container.New(
		layout.NewFormLayout(),
		widget.NewLabel("Name"),
		variableNameEntry,
		widget.NewLabel("Main Prompt"),
		promptEntryMain,
		widget.NewLabel("Prompt suffix"),
		promptEntrySuffix,
		widget.NewLabel("Preview"),
		promptPreview,
	)
	return container.New(layout.NewVBoxLayout(), form, container.New(layout.NewBorderLayout(nil, nil, nil, saveBtn), saveBtn))
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
		replacementList = append(replacementList, helpers.GenerateVarPlaceholderString(name), value)
	}
	return strings.NewReplacer(replacementList...)
}

func valueOrPlaceholderValue(txt string) string {
	if txt == "" {
		return helpers.VariablePlaceholderValue
	}
	return txt
}
