package anthropic

import (
	"context"

	"github.com/anthropics/anthropic-sdk-go"
	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/anthropic/codec"
)

// ModelOption is a function type that modifies a LanguageModel.
type ModelOption func(*LanguageModel)

// WithClient returns a ModelOption that sets the client.
func WithClient(client anthropic.Client) ModelOption {
	// TODO: Instead of only supporting an anthropic.Client, we can "flatten"
	// the options supported by the Anthropic SDK.
	return func(m *LanguageModel) {
		m.client = client
	}
}

// LanguageModel represents an Anthropic language model.
type LanguageModel struct {
	modelID string
	client  anthropic.Client
}

var _ api.LanguageModel = &LanguageModel{}
var _ api.TokenCounter = &LanguageModel{}

// NewLanguageModel creates a new Anthropic language model.
func NewLanguageModel(modelID string, opts ...ModelOption) *LanguageModel {
	// Create model with default settings
	model := &LanguageModel{
		modelID: modelID,
		client:  anthropic.NewClient(), // Default client
	}

	// Apply options
	for _, opt := range opts {
		opt(model)
	}

	return model
}

func (m *LanguageModel) ProviderName() string {
	return ProviderName
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
	params, warnings, err := codec.EncodeParams(m.modelID, prompt, opts)
	if err != nil {
		return nil, err
	}

	message, err := m.client.Beta.Messages.New(ctx, params)
	if err != nil {
		return nil, err
	}

	response, err := codec.DecodeResponse(message)
	if err != nil {
		return nil, err
	}

	response.Warnings = append(response.Warnings, warnings...)
	return response, nil
}

func (m *LanguageModel) Stream(
	ctx context.Context, prompt []api.Message, opts api.CallOptions,
) (*api.StreamResponse, error) {
	return nil, api.NewUnsupportedFunctionalityError("streaming generation", "")
}

func (m *LanguageModel) CountTokens(
	ctx context.Context, prompt []api.Message, opts api.CallOptions,
) (*api.TokenCount, error) {
	anthropicPrompt, err := codec.EncodePrompt(prompt)
	if err != nil {
		return nil, err
	}

	params := anthropic.BetaMessageCountTokensParams{
		Model: anthropic.Model(m.modelID),
	}

	if len(anthropicPrompt.System) > 0 {
		params.System = anthropic.BetaMessageCountTokensParamsSystemUnion{
			OfBetaTextBlockArray: anthropicPrompt.System,
		}
	}
	if len(anthropicPrompt.Messages) > 0 {
		params.Messages = anthropicPrompt.Messages
	}

	if len(opts.Tools) > 0 {
		tools, err := codec.EncodeTools(opts.Tools, opts.ToolChoice)
		if err != nil {
			return nil, err
		}
		if len(tools.Tools) > 0 {
			countTokensTools := make([]anthropic.BetaMessageCountTokensParamsToolUnion, len(tools.Tools))
			for i, tool := range tools.Tools {
				countTokensTool := anthropic.BetaMessageCountTokensParamsToolUnion{}
				if tool.OfTool != nil {
					countTokensTool.OfTool = tool.OfTool
				} else if tool.OfBashTool20241022 != nil {
					countTokensTool.OfBashTool20241022 = tool.OfBashTool20241022
				} else if tool.OfBashTool20250124 != nil {
					countTokensTool.OfBashTool20250124 = tool.OfBashTool20250124
				} else if tool.OfCodeExecutionTool20250522 != nil {
					countTokensTool.OfCodeExecutionTool20250522 = tool.OfCodeExecutionTool20250522
				} else if tool.OfCodeExecutionTool20250825 != nil {
					countTokensTool.OfCodeExecutionTool20250825 = tool.OfCodeExecutionTool20250825
				} else if tool.OfComputerUseTool20241022 != nil {
					countTokensTool.OfComputerUseTool20241022 = tool.OfComputerUseTool20241022
				} else if tool.OfMemoryTool20250818 != nil {
					countTokensTool.OfMemoryTool20250818 = tool.OfMemoryTool20250818
				} else if tool.OfComputerUseTool20250124 != nil {
					countTokensTool.OfComputerUseTool20250124 = tool.OfComputerUseTool20250124
				} else if tool.OfTextEditor20241022 != nil {
					countTokensTool.OfTextEditor20241022 = tool.OfTextEditor20241022
				} else if tool.OfComputerUseTool20251124 != nil {
					countTokensTool.OfComputerUseTool20251124 = tool.OfComputerUseTool20251124
				} else if tool.OfTextEditor20250124 != nil {
					countTokensTool.OfTextEditor20250124 = tool.OfTextEditor20250124
				} else if tool.OfTextEditor20250429 != nil {
					countTokensTool.OfTextEditor20250429 = tool.OfTextEditor20250429
				} else if tool.OfTextEditor20250728 != nil {
					countTokensTool.OfTextEditor20250728 = tool.OfTextEditor20250728
				} else if tool.OfWebSearchTool20250305 != nil {
					countTokensTool.OfWebSearchTool20250305 = tool.OfWebSearchTool20250305
				} else if tool.OfWebFetchTool20250910 != nil {
					countTokensTool.OfWebFetchTool20250910 = tool.OfWebFetchTool20250910
				} else if tool.OfToolSearchToolBm25_20251119 != nil {
					countTokensTool.OfToolSearchToolBm25_20251119 = tool.OfToolSearchToolBm25_20251119
				} else if tool.OfToolSearchToolRegex20251119 != nil {
					countTokensTool.OfToolSearchToolRegex20251119 = tool.OfToolSearchToolRegex20251119
				} else if tool.OfMCPToolset != nil {
					countTokensTool.OfMCPToolset = tool.OfMCPToolset
				}
				countTokensTools[i] = countTokensTool
			}
			params.Tools = countTokensTools
		}
	}

	result, err := m.client.Beta.Messages.CountTokens(ctx, params)
	if err != nil {
		return nil, err
	}

	return &api.TokenCount{
		InputTokens: int(result.InputTokens),
	}, nil
}
