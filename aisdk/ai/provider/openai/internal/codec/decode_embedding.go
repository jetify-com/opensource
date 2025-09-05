package codec

import (
	"net/http"

	"github.com/openai/openai-go/v2"
	"go.jetify.com/ai/api"
)

// DecodeEmbedding maps the OpenAI embedding API response to the unified api.EmbeddingResponse.
func DecodeEmbedding(resp *openai.CreateEmbeddingResponse) (api.EmbeddingResponse, error) {
	if resp == nil {
		return api.EmbeddingResponse{}, api.NewEmptyResponseBodyError("response from OpenAI embeddings API is nil")
	}

	embs := make([]api.Embedding, len(resp.Data))
	for i, d := range resp.Data {
		vec := make([]float64, len(d.Embedding))
		copy(vec, d.Embedding)
		embs[i] = vec
	}

	usage := &api.EmbeddingUsage{
		PromptTokens: resp.Usage.PromptTokens,
		TotalTokens:  resp.Usage.TotalTokens,
	}

	return api.EmbeddingResponse{
		Embeddings: embs,
		Usage:      usage,
		RawResponse: &api.EmbeddingRawResponse{
			Headers: http.Header{},
		},
	}, nil
}
