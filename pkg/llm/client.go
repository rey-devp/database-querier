package llm

import (
	"context"
	"fmt"
	"strings"
)

type LLMResponse struct {
	Collection string                   `json:"collection"`
	Operation  string                   `json:"operation"`
	Filter     map[string]interface{}   `json:"filter"`
	Pipeline   []map[string]interface{} `json:"pipeline"`
	Sort       map[string]interface{}   `json:"sort"`
	Limit      int64                    `json:"limit"`
	Projection map[string]interface{}   `json:"projection"`
}

type LLMClient interface {
	GenerateQuery(ctx context.Context, prompt string) (*LLMResponse, error)
}

func NewClient(provider, apiKey, model string) (LLMClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required for LLM provider %s", provider)
	}

	switch strings.ToLower(provider) {
	case "gemini":
		return NewGeminiClient(apiKey, model), nil
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", provider)
	}
}
