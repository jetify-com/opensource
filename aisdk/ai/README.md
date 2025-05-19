# AI SDK for Go

### Build powerful AI applications and agents using a unified API.

[![Version](https://img.shields.io/github/v/release/jetify-com/ai?color=green&label=version&sort=semver)](https://github.com/jetify-com/ai/releases)
[![Coverage](https://img.shields.io/badge/coverage-90%25-success)]()
[![Go Reference](https://pkg.go.dev/badge/go.jetify.com/ai)](https://pkg.go.dev/go.jetify.com/ai)
[![License](https://img.shields.io/github/license/jetify-com/ai)]()
[![Join Discord](https://img.shields.io/discord/903306922852245526?color=7389D8&label=discord&logo=discord&logoColor=ffffff&cacheSeconds=1800)](https://discord.gg/jetify)

## Introduction

**AI SDK for Go** is a unified interface for interacting with multiple AI providers including OpenAI, Anthropic, and more.
Inspired by [Vercel's AI SDK](https://github.com/vercel/ai) for TypeScript, we bring a similar developer experience to the Go ecosystem.

### The Problem

Building AI applications go today means dealing with:
- **Fragmented ecosystems** - Each provider has different APIs, authentication, and patterns
- **Vendor lock-in** - Switching providers requires rewriting significant application code
- **Poor Go developer experience** - Official Go SDKs are often auto-generated from OpenAPI specs, resulting in unidiomatic Go code
- **Complex multi-modal handling** - Different providers handle images, files, and tools differently

### Our Solution

The AI SDK provides a **unified interface** across multiple AI providers, with key advantages:

1. **Provider abstraction** - Common interfaces for language models, embeddings, and image generation
2. **Go-first design** - Built specifically for Go developers with idiomatic patterns and strong typing
3. **Production-ready** - Comprehensive error handling, automatic retries, rate limiting, and robust provider failover
4. **Multi-modal by default** - First-class support for text, images, files, and structured outputs across all providers
5. **Extensible architecture** - Clean interfaces make it easy to add new providers while maintaining backward compatibility

### When to Choose AI SDK

**Choose AI SDK when you want to:**
- Build applications that might need to switch AI providers
- Avoid vendor lock-in from the start
- Use the same API patterns across different model types (LLM, embeddings, image generation)
- Leverage Go's performance and type safety for AI applications
- Have a consistent developer experience across your team

**Stick with provider SDKs when you:**
- Are building a proof of concept with a single provider
- Need cutting-edge provider-specific features immediately upon release
- Have team expertise with specific provider APIs

## Features

* [x] **Multi-Provider Support** – [OpenAI](#), [Anthropic](#), with more coming
* [x] **Multi-Modal Inputs** – Text, images, and files in conversations
* [x] **Tool Calling** – Function calling with parallel execution
* [x] **Language Models** – Text generation with streaming support
* [ ] **Embedding Models** – Text embeddings for semantic search
* [ ] **Image Models** – Generate images from text prompts
* [ ] **Structured Outputs** – JSON generation with schema validation

### Language Models

* [x] Text generation (streaming & non-streaming)
* [x] Multi-modal conversations (text + images + files)
* [x] System messages and conversation history
* [x] Tool/function calling with structured schemas
* [ ] JSON output with schema validation

### Provider-Specific Features

* [x] **OpenAI** - Web search, computer use, file search tools
* [x] **Anthropic** - Claude's advanced reasoning and tool use

## Status

- [x] Private Alpha: We are testing the SDK with a select group of developers.
- [x] Public Alpha: Open to all developers, but breaking changes still expected.
- [ ] Public Beta: Stable enough for most non-enterprise use cases.
- [ ] v1 Release: Ready for production use at scale with guranteed API stability.

We are currently in **Public Alpha**. The SDK functionality is stable but it's API may have breaking changes. While in alpha, minor version bumps indicate breaking changes (`0.1.0` -> `0.2.0` would indicate a breaking change). Watch "releases" of this repo to get notified of major updates.

## Installation

```bash
go get go.jetify.com/ai
```

## Quickstart

Get started with a simple text generation example:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "go.jetify.com/ai"
    "go.jetify.com/ai/provider/openai"
)

func main() {
    // Set up your model (defaults to Claude 3.7 Sonnet)
    model := openai.NewLanguageModel("gpt-4o")

    // Generate text
    response, err := aisdk.GenerateText(
        context.Background(),
        "Explain quantum computing in simple terms",
        aisdk.WithModel(model),
        aisdk.WithMaxTokens(200),
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(response.Text)
}
```

## Usage

### Basic Text Generation

```go
response, err := aisdk.GenerateText(
    ctx,
    "Write a haiku about Go programming",
    aisdk.WithTemperature(0.8),
    aisdk.WithMaxTokens(100),
)
```

### Multi-Modal Conversations

```go
response, err := aisdk.GenerateText(
    ctx,
    aisdk.UserMessage(
        "What's in this image?",
        &api.ImageBlock{URL: "https://example.com/chart.png"},
    ),
    aisdk.AssistantMessage("I can see a bar chart showing sales data."),
    aisdk.UserMessage("What insights can you draw?"),
)
```

### Provider Switching

```go
// Use OpenAI
openaiModel := openai.NewLanguageModel("gpt-4o")

// Use Anthropic
claudeModel := anthropic.NewLanguageModel("claude-3-7-sonnet-20250219")

// Same code works with both
response1, _ := aisdk.GenerateText(ctx, prompt, aisdk.WithModel(openaiModel))
response2, _ := aisdk.GenerateText(ctx, prompt, aisdk.WithModel(claudeModel))
```

## Advanced Features

### Tool Calling

Enable models to call functions with structured schemas:

```go
weatherTool := &api.FunctionTool{
    Name: "get_weather",
    Description: "Get current weather for a location",
    InputSchema: &jsonschema.Definition{
        Type: jsonschema.Object,
        Properties: map[string]jsonschema.Definition{
            "location": {Type: jsonschema.String},
        },
        Required: []string{"location"},
    },
}

response, err := aisdk.GenerateText(
    ctx,
    "What's the weather like in San Francisco?",
    aisdk.WithTools(weatherTool),
)

// Handle tool calls in response.ToolCalls
```

### Structured JSON Output

Generate data matching your Go structs:

```go
type Analysis struct {
    Summary   string   `json:"summary"`
    Sentiment string   `json:"sentiment"`
    Topics    []string `json:"topics"`
}

response, err := aisdk.GenerateText(
    ctx,
    "Analyze this text...",
    aisdk.WithResponseFormat(&api.ResponseFormat{
        Type:   "json",
        Schema: schemaForStruct(Analysis{}), // helper function
    }),
)
```

### Provider-Specific Features

Access unique capabilities like OpenAI's web search:

```go
response, err := aisdk.GenerateText(
    ctx,
    "What happened in San Francisco this week?",
    aisdk.WithTools(&openai.WebSearchTool{
        SearchContextSize: "high",
    }),
    aisdk.WithModel(openaiModel),
)
```

For detailed examples, see our [examples directory](examples/) and [documentation](https://docs.jetify.com/ai-sdk).

## Configuration

### Environment Variables

```bash
export OPENAI_API_KEY="your-openai-key"
export ANTHROPIC_API_KEY="your-anthropic-key"
```

### Global Defaults

```go
// Set default model for all operations
aisdk.SetDefaultLanguageModel(openai.NewLanguageModel("gpt-4o"))

// Now you can omit WithModel() in calls
response, err := aisdk.GenerateText(ctx, "Hello world")
```

### Provider-Specific Options

```go
// OpenAI-specific configuration
metadata := &openai.Metadata{
    User: "user-123",
    Store: true,
    ParallelToolCalls: false,
}

response, err := aisdk.GenerateText(
    ctx,
    prompt,
    aisdk.WithProviderMetadata("openai", metadata),
)
```

## Documentation

Comprehensive documentation is available:

* **[API Reference](https://pkg.go.dev/go.jetify.com/ai)** - Complete Go package documentation
* **[Provider Guides](https://docs.jetify.com/ai-sdk/providers)** - OpenAI, Anthropic, and more
* **[Examples](examples/)** - Real-world usage patterns
* **[Migration Guide](https://docs.jetify.com/ai-sdk/migration)** - Moving from other SDKs

## Community & Support

Join our community and get help:

* **Discord** – [https://discord.gg/jetify](https://discord.gg/jetify) (best for quick questions & showcase)
* **GitHub Discussions** – [Discussions](https://github.com/jetify/ai-sdk/discussions) (best for ideas & design questions)
* **Issues** – [Bug reports & feature requests](https://github.com/jetify/ai-sdk/issues)

## Contributing

We 💖 contributions! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

```
Fork  ➜  Clone ➜  Create branch ➜  Commit ➜  Push ➜  PR 🙌
```

## Roadmap

See [open issues](https://github.com/jetify/ai-sdk/issues) for the full roadmap. High-level priorities:

* [ ] **Custom Provider Interface** - Build your own provider implementations
* [ ] **Azure OpenAI Support** - Enterprise Azure deployments
* [ ] **Google AI Integration** - Gemini models and Vertex AI
* [ ] **Streaming Improvements** - Server-sent events and better error handling
* [ ] **Fine-tuning APIs** - Train custom models across providers

## Related Work

Similar projects and alternatives:

* **[LangChain Go](https://github.com/tmc/langchaingo)** - Comprehensive AI framework for Go
* **[OpenAI Go](https://github.com/openai/openai-go)** - Official OpenAI SDK for Go
* **[Anthropic Go](https://github.com/anthropics/anthropic-sdk-go)** - Official Anthropic SDK for Go
* **[LangSmith](https://docs.smith.langchain.com/)** - Observability and testing for LLM applications

Our SDK focuses specifically on **provider abstraction** and **production reliability** with a smaller, more focused API surface.

## License

Licensed under the **MIT License** – see [LICENSE](LICENSE) for details.

## Acknowledgments

* Thanks to all [contributors](https://github.com/jetify/ai-sdk/graphs/contributors) ❤️
* Inspired by the [Vercel AI SDK](https://github.com/vercel/ai) and Go community best practices
* Made possible by our [community](#community--support) and the amazing AI provider APIs

---

*Happy building! 🚀*
