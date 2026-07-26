package ai

import (
	"context"
	"fmt"

	"go.jetify.com/ai/api"
)

// CountTokens counts the number of tokens in the given messages using the specified model.
//
// This function is useful for estimating costs and checking whether a prompt
// fits within the model's context window before making a generation request.
//
// Example usage:
//
//	count, err := ai.CountTokens(ctx, messages, ai.WithModel(model))
//	if err != nil {
//		// Handle error - model may not support token counting
//	}
//	fmt.Printf("Token count: %d\n", count.InputTokens)
//
// The function accepts [GenerateOption] arguments:
//
//	ai.CountTokens(ctx, messages, ai.WithModel(model), ai.WithTools(tools...))
func CountTokens(ctx context.Context, prompt []api.Message, opts ...GenerateOption) (*api.TokenCount, error) {
	config := buildGenerateConfig(opts)
	return countTokens(ctx, prompt, config)
}

// CountTokensStr is a convenience function for counting tokens in a simple string.
//
// Example usage:
//
//	count, err := ai.CountTokensStr(ctx, "Hello, world!", ai.WithModel(model))
//	if err != nil {
//		// Handle error
//	}
//	fmt.Printf("Token count: %d\n", count.InputTokens)
//
// The string is automatically converted to a [api.UserMessage] before counting.
func CountTokensStr(ctx context.Context, text string, opts ...GenerateOption) (*api.TokenCount, error) {
	msg := &api.UserMessage{
		Content: []api.ContentBlock{&api.TextBlock{Text: text}},
	}
	return CountTokens(ctx, []api.Message{msg}, opts...)
}

func countTokens(ctx context.Context, prompt []api.Message, opts GenerateOptions) (*api.TokenCount, error) {
	// Check if the model implements TokenCounter
	counter, ok := opts.Model.(api.TokenCounter)
	if !ok {
		return nil, api.NewUnsupportedFunctionalityError(
			"token counting",
			fmt.Sprintf("model %q does not support token counting", opts.Model.ModelID()),
		)
	}
	return counter.CountTokens(ctx, prompt, opts.CallOptions)
}

func getEncodingForModel(modelID string) string {
	if len(modelID) == 0 {
		return "cl100k_base"
	}

	if modelID == "gpt-4o" || modelID == "gpt-4o-mini" || modelID == "gpt-5" ||
		modelID == "o1-preview" || modelID == "o3-mini" ||
		(len(modelID) > 0 && modelID[0] == 'o') {
		return "o200k_base"
	}

	return "cl100k_base"
}

