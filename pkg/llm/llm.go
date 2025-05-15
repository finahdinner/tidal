package llm

import "time"

type LLMHandler interface {
	BuildPrompt(string) string
	GetResponseText(string, time.Duration) string
}
