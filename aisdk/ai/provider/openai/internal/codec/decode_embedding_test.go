package codec

import (
	"net/http"
	"testing"

	"github.com/openai/openai-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/api"
)

func TestDecodeEmbedding(t *testing.T) {
	type tc struct {
		name       string
		in         *openai.CreateEmbeddingResponse
		want       api.EmbeddingResponse
		wantErrSub string
	}

	tests := []tc{
		{
			name:       "nil response -> error",
			in:         nil,
			wantErrSub: "response from OpenAI embeddings API is nil",
		},
		{
			name: "maps data and usage; copies vectors; empty headers",
			in: &openai.CreateEmbeddingResponse{
				Data: []openai.Embedding{
					{Embedding: []float64{0.1, 0.2, 0.3}},
					{Embedding: []float64{0.4, 0.5}},
				},
				Usage: openai.CreateEmbeddingResponseUsage{
					PromptTokens: 27,
					TotalTokens:  27,
				},
			},
			want: api.EmbeddingResponse{
				Embeddings: []api.Embedding{
					[]float64{0.1, 0.2, 0.3},
					[]float64{0.4, 0.5},
				},
				Usage: &api.EmbeddingUsage{
					PromptTokens: 27,
					TotalTokens:  27,
				},
				RawResponse: &api.EmbeddingRawResponse{
					Headers: http.Header{},
				},
			},
		},
		{
			name: "empty data yields empty embeddings and zero usage",
			in: &openai.CreateEmbeddingResponse{
				Data: []openai.Embedding{},
				Usage: openai.CreateEmbeddingResponseUsage{
					PromptTokens: 0,
					TotalTokens:  0,
				},
			},
			want: api.EmbeddingResponse{
				Embeddings: []api.Embedding{},
				Usage: &api.EmbeddingUsage{
					PromptTokens: 0,
					TotalTokens:  0,
				},
				RawResponse: &api.EmbeddingRawResponse{
					Headers: http.Header{},
				},
			},
		},
		{
			name: "single long vector",
			in: &openai.CreateEmbeddingResponse{
				Data: []openai.Embedding{
					{Embedding: []float64{1, 2, 3, 4, 5, 6}},
				},
				Usage: openai.CreateEmbeddingResponseUsage{
					PromptTokens: 12,
					TotalTokens:  12,
				},
			},
			want: api.EmbeddingResponse{
				Embeddings: []api.Embedding{
					[]float64{1, 2, 3, 4, 5, 6},
				},
				Usage: &api.EmbeddingUsage{
					PromptTokens: 12,
					TotalTokens:  12,
				},
				RawResponse: &api.EmbeddingRawResponse{
					Headers: http.Header{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeEmbedding(tt.in)

			if tt.wantErrSub != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrSub)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
