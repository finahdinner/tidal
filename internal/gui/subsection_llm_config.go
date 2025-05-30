package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/finahdinner/tidal/internal/config"
	"github.com/finahdinner/tidal/pkg/llm"
)

var llmConfigWindowSize fyne.Size = fyne.NewSize(600, 1) // height 1 lets the layout determine the height

func (g *GuiWrapper) getLlmConfigSubsection() *fyne.Container {

	saveButton := widget.NewButton("Save", nil)

	llmProviderSelect := widget.NewSelect(llm.LlmProviders, nil)
	llmProviderSelect.SetSelected(config.Preferences.LlmConfig.Provider)

	llmApiKeyEntry := widget.NewPasswordEntry()
	llmApiKeyEntry.SetText(config.Preferences.LlmConfig.ApiKey)

	defaultPromptSuffixEntry := getMultilineEntry(
		config.Preferences.LlmConfig.DefaultPromptSuffix,
		saveButton, tallerMultilineEntryHeight, fyne.ScrollVerticalOnly, fyne.TextWrapWord,
	)

	saveButton.OnTapped = func() {
		config.Preferences.LlmConfig = config.LlmConfigT{
			Provider:            llmProviderSelect.Selected,
			ApiKey:              llmApiKeyEntry.Text,
			DefaultPromptSuffix: defaultPromptSuffixEntry.Text,
		}
		if err := config.SavePreferences(); err != nil {
			showErrorDialog(
				fmt.Errorf("unable to save LLM configuration - err: %w", err),
				"Unable to save LLM configuration.",
				g.SecondaryWindow,
			)
			return
		}
		saveButton.Disable()
		g.closeSecondaryWindow()
	}

	return container.New(
		layout.NewFormLayout(),
		widget.NewLabel("Provider"),
		llmProviderSelect,
		widget.NewLabel("API Key"),
		llmApiKeyEntry,
		widget.NewLabel("Default Prompt Suffix"),
		defaultPromptSuffixEntry,
		layout.NewSpacer(),
		saveButton,
	)
}
