package codec

import (
	"net/http"
	"testing"

	"github.com/openai/openai-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/api"
)

func TestEncodeEmbedding(t *testing.T) {
	tests := []struct {
		name            string
		modelID         string
		values          []string
		headers         http.Header
		wantReqOptsLen  int
		wantWarningsLen int
		expectedParams  openai.EmbeddingNewParams
	}{
		{
			name:            "no headers, two values",
			modelID:         "text-embedding-3-small",
			values:          []string{"hello", "world"},
			headers:         nil,
			wantReqOptsLen:  0,
			wantWarningsLen: 0,
			expectedParams: openai.EmbeddingNewParams{
				Model: openai.EmbeddingModel("text-embedding-3-small"),
				Input: openai.EmbeddingNewParamsInputUnion{
					OfArrayOfStrings: []string{"hello", "world"},
				},
				EncodingFormat: openai.EmbeddingNewParamsEncodingFormatFloat,
			},
		},
		{
			name:    "with single and multi-value headers",
			modelID: "text-embedding-3-small",
			values:  []string{"a", "b", "c"},
			headers: func() http.Header {
				h := http.Header{}
				h.Add("X-One", "1")
				h.Add("X-Multi", "A")
				h.Add("X-Multi", "B")
				return h
			}(),
			// 1 option for X-One + 2 options for X-Multi
			wantReqOptsLen:  3,
			wantWarningsLen: 0,
			expectedParams: openai.EmbeddingNewParams{
				Model: openai.EmbeddingModel("text-embedding-3-small"),
				Input: openai.EmbeddingNewParamsInputUnion{
					OfArrayOfStrings: []string{"a", "b", "c"},
				},
				EncodingFormat: openai.EmbeddingNewParamsEncodingFormatFloat,
			},
		},
		{
			name:            "empty input slice",
			modelID:         "text-embedding-3-large",
			values:          []string{},
			headers:         nil,
			wantReqOptsLen:  0,
			wantWarningsLen: 0,
			expectedParams: openai.EmbeddingNewParams{
				Model: openai.EmbeddingModel("text-embedding-3-large"),
				Input: openai.EmbeddingNewParamsInputUnion{
					OfArrayOfStrings: []string{},
				},
				EncodingFormat: openai.EmbeddingNewParamsEncodingFormatFloat,
			},
		},
		{
			name:            "different model id",
			modelID:         "text-embedding-3-small",
			values:          []string{"only one"},
			headers:         http.Header{},
			wantReqOptsLen:  0,
			wantWarningsLen: 0,
			expectedParams: openai.EmbeddingNewParams{
				Model: openai.EmbeddingModel("text-embedding-3-small"),
				Input: openai.EmbeddingNewParamsInputUnion{
					OfArrayOfStrings: []string{"only one"},
				},
				EncodingFormat: openai.EmbeddingNewParamsEncodingFormatFloat,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := api.EmbeddingOptions{Headers: tt.headers}

			params, reqOpts, warnings, err := EncodeEmbedding(tt.modelID, tt.values, opts)
			require.NoError(t, err)

			// Request options (opaque type): assert count derived from headers
			assert.Len(t, reqOpts, tt.wantReqOptsLen)

			// Warnings (currently none expected)
			assert.Len(t, warnings, tt.wantWarningsLen)

			// Params: model id
			assert.Equal(t, openai.EmbeddingModel(tt.modelID), params.Model)

			// Params: input union mirrors provided values
			assert.Equal(t, tt.values, params.Input.OfArrayOfStrings)
		})
	}
}
