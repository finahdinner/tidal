package gui

import (
	"errors"
	"fmt"
	"image/color"
	"strconv"
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

var titleSetupWindowSize fyne.Size = fyne.NewSize(700, 1) // height 1 lets the layout determine the height

func (g *GuiWrapper) getTitleSetupSubsection() *fyne.Container {

	titleConfig := config.Preferences.Title

	saveBtn := widget.NewButton("Save", nil)
	if !titleConfigValid(titleConfig) {
		saveBtn.Disable()
	}

	titleTemplateEntry := getMultilineEntry(titleConfig.TitleTemplate, saveBtn, 6, fyne.ScrollVerticalOnly, fyne.TextWrapWord)
	tipLabel := widget.NewLabelWithStyle(
		fmt.Sprintf("You can use any Variables in your title template\nAccess them using %sVariableName%s", helpers.VarNamePlaceholderPrefix, helpers.VarNamePlaceholderSuffix),
		fyne.TextAlignLeading, fyne.TextStyle{Italic: true},
	)

	variablesDetectedWidget := widget.NewRichText()
	variablesDetectedWidget.Scroll = fyne.ScrollHorizontalOnly
	variablesDetected := []string{}
	variablesDetectedIndices := map[string]int{} // index position in the slice above

	allVariablesNamesMap := map[string]struct{}{}
	twitchVarNamesSlice, _ := config.GetAllTwitchVariables()
	aiGeneratedVarNamesSlice, _ := config.GetAllAiGeneratedVariables()
	for _, v := range append(twitchVarNamesSlice, aiGeneratedVarNamesSlice...) {
		allVariablesNamesMap[v] = struct{}{}
	}

	updateVariablesDetectedLabel(
		titleConfig.TitleTemplate,
		allVariablesNamesMap,
		&variablesDetected,
		variablesDetectedIndices,
		variablesDetectedWidget,
	)

	updateIntervalEntry := widget.NewEntry()
	if config.Preferences.Title.TitleUpdateIntervalMinutes > 0 {
		updateIntervalEntry.SetText(strconv.Itoa(titleConfig.TitleUpdateIntervalMinutes))
	}
	intervalEntryErrorText := canvas.NewText("", color.RGBA{255, 0, 0, 255})

	saveBtn.OnTapped = func() {
		if !titleConfigValid(titleConfig) {
			config.Logger.LogErrorf("unable to save title config - titleConfig is invalid")
			showErrorDialog(
				errors.New("title config is not valid"),
				"Title configuration fields are not all valid/populated",
				g.SecondaryWindow,
			)
			return
		}
		config.Preferences.Title = titleConfig
		config.SavePreferences() // TODO - do I need to check for the error?
		g.closeSecondaryWindow()
	}

	titleTemplateEntry.OnChanged = func(s string) {
		saveBtn.Disable()
		s = strings.TrimSpace(s)
		titleConfig.TitleTemplate = s

		hasUndefinedVariables := updateVariablesDetectedLabel(
			titleConfig.TitleTemplate,
			allVariablesNamesMap,
			&variablesDetected,
			variablesDetectedIndices,
			variablesDetectedWidget,
		)

		if titleConfigValid(titleConfig) && !hasUndefinedVariables {
			saveBtn.Enable()
		}
	}

	updateIntervalEntry.OnChanged = func(s string) {
		saveBtn.Disable()
		titleConfig.TitleUpdateIntervalMinutes = -1 // will be updated if s is valid
		intervalEntryErrorText.Text = ""
		s = strings.TrimSpace(s)
		updateIntervalMinutes, err := strconv.Atoi(s)
		if err != nil {
			config.Logger.LogErrorf("unable to convert text %q to an int - err: %v", s, err)
			return
		}
		if updateIntervalMinutes < helpers.MinTitleUpdateIntervalMinutes ||
			updateIntervalMinutes > helpers.MaxTitleUpdateIntervalMinutes {
			config.Logger.LogDebugf("updateIntervalMinutes must be %v<=x<=%v", helpers.MinTitleUpdateIntervalMinutes, helpers.MaxTitleUpdateIntervalMinutes)
			intervalEntryErrorText.Text = fmt.Sprintf("Interval must be between %v and %v, inclusive.", helpers.MinTitleUpdateIntervalMinutes, helpers.MaxTitleUpdateIntervalMinutes)
			return
		}
		titleConfig.TitleUpdateIntervalMinutes = updateIntervalMinutes
		if titleConfigValid(titleConfig) {
			saveBtn.Enable()
		}
	}

	updateFrequencyContainer := container.New(
		layout.NewFormLayout(),
		container.New(
			layout.NewGridLayoutWithColumns(2),
			updateIntervalEntry,
			widget.NewLabel("Minutes"),
		),
		intervalEntryErrorText,
	)

	updateImmediatelyCheck := widget.NewCheck("Update immediately on start", func(b bool) {
		titleConfig.UpdateImmediatelyOnStart = b
	})
	updateImmediatelyCheck.SetChecked(titleConfig.UpdateImmediatelyOnStart)

	throwErrorIfEmptyVariable := widget.NewCheck("Throw error if using an empty variable", func(b bool) {
		titleConfig.ThrowErrorIfEmptyVariable = b
	})
	throwErrorIfEmptyVariable.SetChecked(titleConfig.ThrowErrorIfEmptyVariable)

	throwErrorIfNonExistentVariable := widget.NewCheck("Throw error if using a non-existent variable", func(b bool) {
		titleConfig.ThrowErrorIfNonExistentVariable = b
	})
	throwErrorIfNonExistentVariable.SetChecked(titleConfig.ThrowErrorIfNonExistentVariable)

	return container.New(
		layout.NewFormLayout(),
		widget.NewLabel("Title Template"),
		titleTemplateEntry,
		layout.NewSpacer(),
		tipLabel,
		widget.NewLabel("Variables Detected"),
		variablesDetectedWidget,
		widget.NewLabel("Update Every "),
		updateFrequencyContainer,
		layout.NewSpacer(),
		updateImmediatelyCheck,
		layout.NewSpacer(),
		throwErrorIfEmptyVariable,
		layout.NewSpacer(),
		throwErrorIfNonExistentVariable,
		layout.NewSpacer(),
		saveBtn,
	)
}

func titleConfigValid(titleConfig config.TitleT) bool {
	return titleConfig.TitleTemplate != "" && titleConfig.TitleUpdateIntervalMinutes <= helpers.MaxTitleUpdateIntervalMinutes && titleConfig.TitleUpdateIntervalMinutes >= helpers.MinTitleUpdateIntervalMinutes
}

func removeFromStringSlicePreserveOrder(slice *[]string, removalIdx int) error {
	if removalIdx >= len(*slice) {
		return errors.New("provided index is out of range")
	}
	*slice = append((*slice)[:removalIdx], (*slice)[removalIdx+1:]...)
	return nil
}

func updateVariablesDetectedLabel(
	titleTemplate string,
	allVariablesNamesMap map[string]struct{},
	variablesDetectedPtr *[]string,
	variablesDetectedIndices map[string]int,
	variablesDetectedWidget *widget.RichText,
) bool {
	tmpVariablesDetected := helpers.ExtractVariableNamesFromText(titleTemplate)
	tmpVariablesDetectedSet := map[string]struct{}{}
	for _, v := range tmpVariablesDetected {
		tmpVariablesDetectedSet[v] = struct{}{}
	}

	variablesDetected := *variablesDetectedPtr

	// remove any variables that haven't been detected
	for _, v := range variablesDetected {
		if _, exists := tmpVariablesDetectedSet[v]; !exists {
			removeFromStringSlicePreserveOrder(&variablesDetected, variablesDetectedIndices[v])
			delete(variablesDetectedIndices, v)
		}
	}

	// insert any new variables that weren't previously being tracked
	for _, v := range tmpVariablesDetected {
		if _, exists := variablesDetectedIndices[v]; !exists {
			variablesDetected = append(variablesDetected, v)
			variablesDetectedIndices[v] = len(variablesDetected) - 1
		}
	}

	// rebuild segments
	numUndefinedVars := 0
	variablesDetectedWidget.Segments = []widget.RichTextSegment{}
	for _, v := range variablesDetected {
		segment := &widget.TextSegment{
			Text:  "",
			Style: widget.RichTextStyleInline,
		}
		if _, exists := allVariablesNamesMap[v]; exists {
			segment.Text = fmt.Sprintf("✅ %s  ", v)
			segment.Style.ColorName = theme.ColorGreen
		} else {
			segment.Text = fmt.Sprintf("❌ %s  ", v)
			segment.Style.ColorName = theme.ColorRed
			numUndefinedVars++
		}
		variablesDetectedWidget.Segments = append(
			variablesDetectedWidget.Segments,
			segment,
		)
		variablesDetectedWidget.Refresh()
	}

	// modify the actual slice being passed in
	*variablesDetectedPtr = variablesDetected

	return numUndefinedVars > 0
}
