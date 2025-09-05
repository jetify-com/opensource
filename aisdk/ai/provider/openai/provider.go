package openai

import (
	"github.com/openai/openai-go/v2"
)

type Provider struct {
	// client is the OpenAI client used to make API calls.
	client openai.Client
	// name is the name of the provider, overrides the default "openai".
	name string
}

type ProviderOption func(*Provider)

func WithClient(c openai.Client) ProviderOption {
	return func(p *Provider) { p.client = c }
}

func WithName(name string) ProviderOption {
	return func(p *Provider) { p.name = name }
}

func NewProvider(opts ...ProviderOption) *Provider {
	p := &Provider{client: openai.NewClient()}
	for _, opt := range opts {
		opt(p)
	}
	if p.name == "" {
		p.name = "openai"
	}

	return p
}
