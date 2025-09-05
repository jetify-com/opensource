package openai

import (
	"context"
	"fmt"

	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/openai/internal/codec"
)

// LanguageModel represents an OpenAI language model.
type LanguageModel struct {
	modelID string
	pc      ProviderConfig
}

var _ api.LanguageModel = &LanguageModel{}

// NewLanguageModel creates a new OpenAI language model.
func (p *Provider) NewLanguageModel(modelID string) *LanguageModel {
	// Create model with provider's client
	model := &LanguageModel{
		modelID: modelID,
		pc: ProviderConfig{
			providerName: fmt.Sprintf("%s.responses", p.name),
			client:       p.client,
		},
	}

	return model
}

func (m *LanguageModel) ProviderName() string {
	return m.pc.providerName
}

func (m *LanguageModel) ModelID() string {
	return m.modelID
}

func (m *LanguageModel) SupportedUrls() []api.SupportedURL {
	// TODO: Make configurable via the constructor.
	return []api.SupportedURL{
		{
			MediaType: "image/*",
			URLPatterns: []string{
				"^https?://.*",
			},
		},
	}
}

func (m *LanguageModel) Generate(
	ctx context.Context, prompt []api.Message, opts api.CallOptions,
) (*api.Response, error) {
	params, warnings, err := codec.Encode(m.modelID, prompt, opts)
	if err != nil {
		return nil, err
	}

	openaiResponse, err := m.pc.client.Responses.New(ctx, params)
	if err != nil {
		return nil, err
	}

	response, err := codec.DecodeResponse(openaiResponse)
	if err != nil {
		return nil, err
	}

	response.Warnings = append(response.Warnings, warnings...)
	return response, nil
}

func (m *LanguageModel) Stream(
	ctx context.Context, prompt []api.Message, opts api.CallOptions,
) (*api.StreamResponse, error) {
	// TODO: add warnings to the stream response by adding an initial StreamStart event
	// (it could happen inside of codec.Encode)
	params, _, err := codec.Encode(m.modelID, prompt, opts)
	if err != nil {
		return nil, err
	}

	stream := m.pc.client.Responses.NewStreaming(ctx, params)
	response, err := codec.DecodeStream(stream)
	if err != nil {
		return nil, err
	}

	return response, nil
}
