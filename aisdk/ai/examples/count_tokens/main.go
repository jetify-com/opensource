package main

import (
	"context"
	"fmt"
	"log"
	"os"

	anthropicsdk "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/joho/godotenv"
	"go.jetify.com/ai"
	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/anthropic"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		log.Fatal("ANTHROPIC_API_KEY environment variable is required")
	}

	ctx := context.Background()

	client := anthropicsdk.NewClient(option.WithAPIKey(apiKey))
	model := anthropic.NewLanguageModel(
		"claude-sonnet-4-20250514",
		anthropic.WithClient(client),
	)

	// Test with a simple string
	fmt.Println("=== CountTokensStr ===")
	text := "Hello, world! How are you today?"
	count, err := ai.CountTokensStr(ctx, text, ai.WithModel(model))
	if err != nil {
		log.Fatalf("CountTokensStr error: %v", err)
	}
	fmt.Printf("Model: %s\n", model.ModelID())
	fmt.Printf("Message: %q\n", text)
	fmt.Printf("Input tokens: %d\n\n", count.InputTokens)

	// Test with multiple messages
	fmt.Println("=== CountTokens with multiple messages ===")
	messages := []api.Message{
		&api.SystemMessage{Content: "You are a helpful assistant."},
		&api.UserMessage{
			Content: []api.ContentBlock{
				&api.TextBlock{Text: "What is the capital of France?"},
			},
		},
		&api.AssistantMessage{
			Content: []api.ContentBlock{
				&api.TextBlock{Text: "The capital of France is Paris."},
			},
		},
		&api.UserMessage{
			Content: []api.ContentBlock{
				&api.TextBlock{Text: "What about Germany?"},
			},
		},
	}

	count, err = ai.CountTokens(ctx, messages, ai.WithModel(model))
	if err != nil {
		log.Fatalf("CountTokens error: %v", err)
	}
	fmt.Printf("Model: %s\n", model.ModelID())
	fmt.Println("Messages:")
	fmt.Println("  [system] You are a helpful assistant.")
	fmt.Println("  [user] What is the capital of France?")
	fmt.Println("  [assistant] The capital of France is Paris.")
	fmt.Println("  [user] What about Germany?")
	fmt.Printf("Input tokens: %d\n", count.InputTokens)
}

