package llm

import (
	"fmt"
	"time"
)

var LlmProviders = []string{"Google Gemini"}

type LLMHandler interface {
	// BuildPrompt([]string) string
	GetResponseText(string, time.Duration) (string, error)
}

func NewLlmHandler(provider string, apiKey string) (LLMHandler, error) {
	var handler LLMHandler
	var err error

	switch provider {
	case "Google Gemini":
		handler, err = newGoogleGeminiHandler(apiKey)
		if err != nil {
			return nil, fmt.Errorf("unable to create GoogleGeminiHandler - err: %w", err)
		}
	default:
		return nil, fmt.Errorf("%v is not a valid LLM provider", provider)
	}
	return handler, nil
}
