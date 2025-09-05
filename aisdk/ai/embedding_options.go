package ai

import (
	"net/http"

	"go.jetify.com/ai/api"
)

// EmbeddingOptions bundles the model + per-call embedding options.
type EmbeddingOptions[T any] struct {
	EmbeddingOptions api.EmbeddingOptions
	Model            api.EmbeddingModel[T]
}

// EmbeddingOption mutates EmbeddingOptions.
type EmbeddingOption[T any] func(*EmbeddingOptions[T])

// WithEmbeddingHeaders sets extra HTTP headers for this embedding call.
// Only applies to HTTP-backed providers.
func WithEmbeddingHeaders[T any](headers http.Header) EmbeddingOption[T] {
	return func(o *EmbeddingOptions[T]) {
		o.EmbeddingOptions.Headers = headers
	}
}

// WithEmbeddingProviderMetadata sets provider-specific metadata for the embedding call.
func WithEmbeddingProviderMetadata[T any](provider string, metadata any) EmbeddingOption[T] {
	return func(o *EmbeddingOptions[T]) {
		if o.EmbeddingOptions.ProviderMetadata == nil {
			o.EmbeddingOptions.ProviderMetadata = api.NewProviderMetadata(map[string]any{})
		}
		o.EmbeddingOptions.ProviderMetadata.Set(provider, metadata)
	}
}

// WithEmbeddingBaseURL sets the base URL for the embedding API endpoint.
func WithEmbeddingBaseURL[T any](baseURL string) EmbeddingOption[T] {
	url := baseURL
	return func(o *EmbeddingOptions[T]) {
		o.EmbeddingOptions.BaseURL = &url
	}
}

// WithEmbeddingEmbeddingOptions sets the entire api.EmbeddingOptions struct.
func WithEmbeddingEmbeddingOptions[T any](embeddingOptions api.EmbeddingOptions) EmbeddingOption[T] {
	return func(o *EmbeddingOptions[T]) {
		o.EmbeddingOptions = embeddingOptions
	}
}

// buildEmbeddingConfig combines multiple options into a single EmbeddingOptions.
func buildEmbeddingConfig[T any](
	model api.EmbeddingModel[T], opts []EmbeddingOption[T],
) EmbeddingOptions[T] {
	config := EmbeddingOptions[T]{
		EmbeddingOptions: api.EmbeddingOptions{},
		Model:            model,
	}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}
