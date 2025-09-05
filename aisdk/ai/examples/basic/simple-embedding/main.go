package main

import (
	"context"
	"log"

	"github.com/k0kubun/pp/v3"
	"go.jetify.com/ai"
	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/openai"
)

func example() error {
	// Initialize the OpenAI provider
	provider := openai.NewProvider()

	// Create a model
	model := provider.NewEmbeddingModel("text-embedding-3-small")

	// Generate text
	response, err := ai.EmbedMany(
		context.Background(),
		model,
		[]string{
			"Artificial intelligence is the simulation of human intelligence in machines.",
			"Machine learning is a subset of AI that enables systems to learn from data.",
		},
	)
	if err != nil {
		return err
	}

	// Print the response:
	printResponse(response)

	return nil
}

func printResponse(response api.EmbeddingResponse) {
	printer := pp.New()
	printer.SetOmitEmpty(true)
	printer.Print(response)
}

func main() {
	if err := example(); err != nil {
		log.Fatal(err)
	}
}
