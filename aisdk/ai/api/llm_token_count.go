package api

import "context"

// TokenCount represents the result of counting tokens in a prompt.
type TokenCount struct {
	// InputTokens is the total number of tokens in the input prompt.
	InputTokens int `json:"input_tokens"`
}

// TokenCounter is an optional interface that LanguageModels can implement
// to support token counting functionality.
//
// Token counting allows estimating the cost of a request and checking
// whether a prompt fits within the model's context window before making
// an API call.
type TokenCounter interface {
	// CountTokens counts the number of tokens in the given messages.
	// Returns the token count or an error if counting fails.
	CountTokens(ctx context.Context, prompt []Message, opts CallOptions) (*TokenCount, error)
}

