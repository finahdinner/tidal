package gui

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/finahdinner/tidal/internal/config"
	"github.com/finahdinner/tidal/internal/helpers"
)

var titleSetupWindowSize fyne.Size = fyne.NewSize(700, 1) // height 1 lets the layout determine the height

func (g *GuiWrapper) getTitleSetupSubsection() *fyne.Container {

	titleConfig := config.Preferences.Title

	saveBtn := widget.NewButton("Save", nil)
	titleTemplateEntry := getMultilineEntry("", saveBtn, 6)
	titleTemplateEntry.Scroll = fyne.ScrollVerticalOnly
	titleTemplateEntry.Wrapping = fyne.TextWrapWord
	titleTemplateEntry.SetText(titleConfig.TitleTemplate)
	tipLabel := widget.NewLabelWithStyle("You can use any Variables in your title template\nAccess them using {{VariableName}}", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})

	updateIntervalEntry := widget.NewEntry()
	if config.Preferences.Title.TitleUpdateIntervalMinutes > 0 {
		updateIntervalEntry.SetText(strconv.Itoa(titleConfig.TitleUpdateIntervalMinutes))
	}
	intervalEntryErrorText := canvas.NewText("", color.RGBA{255, 0, 0, 255})

	saveBtn.OnTapped = func() {
		if !titleConfigValid(titleConfig) {
			config.Logger.LogErrorf("unable to save title config - titleConfig is invalid")
			return
		}
		config.Preferences.Title = titleConfig
		config.SavePreferences() // TODO - do I need to check for the error?
		g.closeSecondaryWindow()
	}

	titleTemplateEntry.OnChanged = func(s string) {
		saveBtn.Disable()
		s = strings.TrimSpace(s)
		if s == "" {
			return
		}
		titleConfig.TitleTemplate = s
		if titleConfigValid(titleConfig) {
			saveBtn.Enable()
		}
	}

	updateIntervalEntry.OnChanged = func(s string) {
		saveBtn.Disable()
		intervalEntryErrorText.Text = ""
		s = strings.TrimSpace(s)
		if s == "" {
			return
		}
		if titleConfig.TitleTemplate == "" {
			return
		}
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
		titleConfig.ThrowErrorIfEmptyValue = b
	})
	throwErrorIfEmptyVariable.SetChecked(titleConfig.ThrowErrorIfEmptyValue)

	return container.New(
		layout.NewFormLayout(),
		widget.NewLabel("Title Template"),
		titleTemplateEntry,
		layout.NewSpacer(),
		tipLabel,
		widget.NewLabel("Update every "),
		updateFrequencyContainer,
		layout.NewSpacer(),
		updateImmediatelyCheck,
		layout.NewSpacer(),
		throwErrorIfEmptyVariable,
		layout.NewSpacer(),
		saveBtn,
	)
}

func titleConfigValid(titleConfig config.TitleT) bool {
	return titleConfig.TitleTemplate != "" && titleConfig.TitleUpdateIntervalMinutes <= helpers.MaxTitleUpdateIntervalMinutes && titleConfig.TitleUpdateIntervalMinutes >= helpers.MinTitleUpdateIntervalMinutes
}
