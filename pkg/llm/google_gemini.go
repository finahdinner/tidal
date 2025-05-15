package llm

import (
	"context"
	"strings"
	"time"

	"google.golang.org/genai"
)

const geminiModel = "gemini-2.0-flash"

type GoogleGeminiHandler struct {
	client *genai.Client
}

func NewGoogleGeminiHandler(apiKey string) (*GoogleGeminiHandler, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}
	return &GoogleGeminiHandler{client}, nil
}

func (h *GoogleGeminiHandler) BuildPrompt(promptParts []string) string {
	return strings.Join(promptParts, "\n")
}

func (h *GoogleGeminiHandler) GetResponseText(prompt string, timeoutDuration time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()
	result, err := h.client.Models.GenerateContent(
		ctx,
		geminiModel,
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", err
	}
	return result.Text(), nil
}
