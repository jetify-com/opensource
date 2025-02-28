# serror

[![Go Reference](https://pkg.go.dev/badge/github.com/jetify-com/serror.svg)](https://pkg.go.dev/github.com/jetify-com/serror)
[![Go Report Card](https://goreportcard.com/badge/github.com/jetify-com/serror)](https://goreportcard.com/report/github.com/jetify-com/serror)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

An error handling library for Go that makes it easy to associate arbitrary structured data with errors
in a similar way to what `slog` does for logging.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Key Concepts](#key-concepts)
- [Design Philosophy](#design-philosophy)
- [Contributing](#contributing)
- [License](#license)
- [Credits](#credits)

## Overview
`serror` provides a clean API for creating, wrapping, and logging errors with associated structured data.

When handling errors in Go, developers often need to include contextual information to help diagnose what went wrong. The traditional approach using `fmt.Errorf` interpolates this data into a string:

```go
return fmt.Errorf("failed to process user %d: %w", userID, err)
```

While simple, this approach has limitations:
- The contextual data (like `userID`) becomes just part of the error message
- There's no way to programmatically access these values later
- The data loses its type information
- The format becomes rigid and harder to parse

`serror` solves these problems by allowing you to attach structured data to errors:

```go
return serror.Wrap(err, "failed to process user",
    "user_id", userID,
    "attempt", attempt,
    "status", status,
)
```

This approach:
- Preserves the original data types
- Makes values accessible programmatically via `err.Get("user_id")`
- Integrates seamlessly with structured logging
- Maintains the flexibility to add or modify context

Inspired by Go's `log/slog` package, `serror` uses the same familiar key-value pair pattern for attaching attributes. If you've used `slog`, you'll feel right at home:

```go
// slog style
logger.Error("request failed", 
    "user_id", userID,
    "status", status,
)

// serror follows the same pattern
return serror.New("request failed",
    "user_id", userID,
    "status", status,
)
```

## Features

- **Structured Attributes**: Attach key-value pairs to errors using `slog`-like syntax
- **Error Wrapping**: Fully compatible with Go's error wrapping conventions
- **slog Integration**: Implements `LogValuer` for automatic structured logging
- **JSON Support**: Full JSON marshaling/unmarshaling capabilities
- **Attribute Access**: Get attributes using dot notation (e.g., "user.id")
- **Thread Safety**: Immutable design ensures safe concurrent usage
- **Call Site Capture**: Automatically records where errors are created

## Installation

### Using go get
```bash
go get github.com/jetify-com/serror
```

### From Source
```bash
git clone https://github.com/jetify-com/serror.git
cd serror
go install
```

## Quick Start

```go
// Create a new error with attributes
err := serror.New("failed to process request",
    "userID", 123,
    "path", "/api/v1/users",
)

// Wrap an existing error with additional context
if err != nil {
    return serror.Wrap(err, "user operation failed",
        "operation", "create",
        "retry", false,
    )
}

// Log the error with slog
slog.Error("request failed", "err", err)
```

## Key Concepts

### Creating Errors

```go
// Basic error with attributes
err := serror.New("operation failed",
    "code", 500,
    "component", "database",
)

// Using attribute constructors
err := serror.New("validation failed",
    serror.Int("code", 400),
    serror.String("field", "email"),
    serror.Bool("valid", false),
)

// Group related attributes
err := serror.New("validation failed",
    serror.Group("user",
        serror.Int("id", 123),
        serror.String("name", "alice"),
    ),
)
```

### Wrapping Errors

```go
baseErr := serror.New("database error", 
    "table", "users",
)

err := serror.Wrap(baseErr, "query failed",
    "query_id", "abc123",
)
```

### Adding Attributes

```go
// Add more attributes to an existing error
// Note: With() returns a new error - serror errors are immutable
oldErr := serror.New("operation failed", "code", 500)
newErr := oldErr.With("retry", true, "attempt", 3)

// oldErr is unchanged, newErr has all attributes
```

The immutable design ensures thread safety and prevents accidental modifications. Each call to `With()` returns a new error instance with the combined attributes from the original error plus the new ones.

### Accessing Attributes

```go
// Get attribute values using dot notation
value := err.Get("user.id")
name := err.Get("user.name")
```

### JSON Support

```go
// Marshal to JSON
data, err := json.Marshal(serror)

// Unmarshal from JSON
var newErr serror.Error
err := json.Unmarshal(data, &newErr)
```

### Integration with slog

```go
err := serror.New("operation failed",
    "user_id", 123,
    "status", 500,
)

// serror.Error implements slog.LogValuer
logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
logger.Error("request failed", "err", err)

// Output:
// ERROR request failed err.user_id=123 err.status=500 err.msg="operation failed"
```

When logged, all structured attributes are automatically included in the log output, making it easy to correlate errors with their context.

## Design Philosophy

- **Standard Library Aligned**: Works with `errors.Is`, `errors.As`, `errors.Unwrap`
- **slog Compatible**: Follows `slog` patterns for attribute handling
- **Thread Safe**: Immutable design prevents data races
- **Simple API**: Familiar interface for Go developers

## Documentation

For detailed documentation and API references, please visit our [Go Package Documentation](https://pkg.go.dev/github.com/jetify-com/serror).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## Community

- **Issues:** [GitHub Issues](https://github.com/jetify-com/serror/issues)
- **Discussions:** [GitHub Discussions](https://github.com/jetify-com/serror/discussions)

## License

This project is licensed under the Apache License, Version 2.0 - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Thanks to the Go team for the excellent `log/slog` package design
- Inspired by various error handling patterns in the Go community