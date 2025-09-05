package codec

import (
	"net/http"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"go.jetify.com/ai/api"
)

// EncodeEmbedding builds OpenAI params + request options from the unified API options.
func EncodeEmbedding(
	modelID string,
	values []string,
	opts api.EmbeddingOptions,
) (openai.EmbeddingNewParams, []option.RequestOption, []api.CallWarning, error) {
	var reqOpts []option.RequestOption
	if opts.Headers != nil {
		reqOpts = append(reqOpts, applyHeaders(opts.Headers)...)
	}

	if opts.BaseURL != nil {
		reqOpts = append(reqOpts, option.WithBaseURL(*opts.BaseURL))
	}

	params := openai.EmbeddingNewParams{
		Model: openai.EmbeddingModel(modelID),
		Input: openai.EmbeddingNewParamsInputUnion{
			OfArrayOfStrings: values,
		},
		EncodingFormat: openai.EmbeddingNewParamsEncodingFormatFloat,
	}

	var warnings []api.CallWarning

	return params, reqOpts, warnings, nil
}

// applyHeaders applies the provided HTTP headers to the request options.
func applyHeaders(headers http.Header) []option.RequestOption {
	var reqOpts []option.RequestOption
	for k, vs := range headers {
		for _, v := range vs {
			reqOpts = append(reqOpts, option.WithHeaderAdd(k, v))
		}
	}
	return reqOpts
}
