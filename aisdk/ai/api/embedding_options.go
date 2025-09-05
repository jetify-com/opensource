package api

import "net/http"

// EmbeddingOption represent the options for generating embeddings.
type EmbeddingOption func(*EmbeddingOptions)

// EmbeddingOptions represents the options for generating embeddings.
type EmbeddingOptions struct {
	// Headers are additional HTTP headers to be sent with the request.
	// Only applicable for HTTP-based providers.
	Headers http.Header

	// BaseURL is the base URL for the API endpoint.
	BaseURL *string

	// ProviderMetadata contains additional provider-specific metadata.
	// The metadata is passed through to the provider from the AI SDK and enables
	// provider-specific functionality that can be fully encapsulated in the provider.
	ProviderMetadata *ProviderMetadata
}

func (o EmbeddingOptions) GetProviderMetadata() *ProviderMetadata { return o.ProviderMetadata }
