package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/finahdinner/tidal/internal/config"
)

var llmConfigWindowSize fyne.Size = fyne.NewSize(400, 1) // height 1 lets the layout determine the height

func (g *GuiWrapper) getLlmConfigSubsection() *fyne.Container {

	llmProviders := []string{"Google Gemini"}
	llmProviderSelect := widget.NewSelect(llmProviders, nil)
	llmProviderSelect.SetSelected(config.Preferences.LlmConfig.Provider)

	llmApiKeyEntry := widget.NewPasswordEntry()
	llmApiKeyEntry.SetText(config.Preferences.LlmConfig.ApiKey)

	saveButton := widget.NewButton("Save", nil)
	saveButton.OnTapped = func() {
		config.Preferences.LlmConfig = config.LlmConfigT{
			Provider: llmProviderSelect.Selected,
			ApiKey:   llmApiKeyEntry.Text,
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
		layout.NewSpacer(),
		saveButton,
	)
}
