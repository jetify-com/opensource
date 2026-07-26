package ai

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.jetify.com/ai/api"
)

func TestCountTokens_UnsupportedModel(t *testing.T) {
	ctx := context.Background()
	messages := []api.Message{
		&api.UserMessage{
			Content: []api.ContentBlock{&api.TextBlock{Text: "Hello, world!"}},
		},
	}

	model := &mockLanguageModel{name: "unsupported-model"}
	count, err := CountTokens(ctx, messages, WithModel(model))

	assert.Nil(t, count)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token counting")
	assert.Contains(t, err.Error(), "unsupported-model")
}

func TestCountTokensStr_UnsupportedModel(t *testing.T) {
	ctx := context.Background()
	model := &mockLanguageModel{name: "unsupported-model"}
	count, err := CountTokensStr(ctx, "Hello, world!", WithModel(model))

	assert.Nil(t, count)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token counting")
}

func TestGetEncodingForModel(t *testing.T) {
	tests := []struct {
		modelID  string
		expected string
	}{
		{"gpt-4o", "o200k_base"},
		{"gpt-4o-mini", "o200k_base"},
		{"gpt-5", "o200k_base"},
		{"o1-preview", "o200k_base"},
		{"o3-mini", "o200k_base"},
		{"gpt-4", "cl100k_base"},
		{"gpt-4-turbo", "cl100k_base"},
		{"gpt-3.5-turbo", "cl100k_base"},
		{"unknown-model", "cl100k_base"},
	}

	for _, tt := range tests {
		t.Run(tt.modelID, func(t *testing.T) {
			result := getEncodingForModel(tt.modelID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

