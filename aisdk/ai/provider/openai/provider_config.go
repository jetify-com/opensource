package openai

import "github.com/openai/openai-go/v2"

type ProviderConfig struct {
	providerName string
	client       openai.Client
}
