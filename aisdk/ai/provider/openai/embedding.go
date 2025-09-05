package openai

import (
	"context"
	"fmt"

	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/openai/internal/codec"
)

// EmbeddingModel represents an OpenAI embedding model.
type EmbeddingModel struct {
	modelID string
	pc      ProviderConfig
}

var _ api.EmbeddingModel[string] = &EmbeddingModel{}

// NewEmbeddingModel creates a new OpenAI embedding model.
func (p *Provider) NewEmbeddingModel(modelID string) *EmbeddingModel {
	// Create model with provider's client
	model := &EmbeddingModel{
		modelID: modelID,
		pc: ProviderConfig{
			providerName: fmt.Sprintf("%s.embedding", p.name),
			client:       p.client,
		},
	}

	return model
}

func (m *EmbeddingModel) ProviderName() string {
	return m.pc.providerName
}

func (m *EmbeddingModel) SpecificationVersion() string {
	return "v2"
}

func (m *EmbeddingModel) ModelID() string {
	return m.modelID
}

// SupportsParallelCalls implements api.EmbeddingModel.
func (m *EmbeddingModel) SupportsParallelCalls() bool {
	return true
}

// MaxEmbeddingsPerCall implements api.EmbeddingModel.
func (m *EmbeddingModel) MaxEmbeddingsPerCall() *int {
	max := 2048
	return &max
}

// DoEmbed implements api.EmbeddingModel.
func (m *EmbeddingModel) DoEmbed(
	ctx context.Context,
	values []string,
	opts api.EmbeddingOptions,
) (api.EmbeddingResponse, error) {
	embeddingParams, openaiOpts, _, err := codec.EncodeEmbedding(
		m.modelID,
		values,
		opts,
	)
	if err != nil {
		return api.EmbeddingResponse{}, err
	}

	resp, err := m.pc.client.Embeddings.New(ctx, embeddingParams, openaiOpts...)
	if err != nil {
		return api.EmbeddingResponse{}, err
	}

	return codec.DecodeEmbedding(resp)
}
