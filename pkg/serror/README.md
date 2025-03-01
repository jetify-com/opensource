# serror

[![Go Reference](https://pkg.go.dev/badge/github.com/jetify-com/serror.svg)](https://pkg.go.dev/github.com/jetify-com/serror)
[![Go Report Card](https://goreportcard.com/badge/github.com/jetify-com/serror)](https://goreportcard.com/report/github.com/jetify-com/serror)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

An error handling library for Go that makes it easy to associate arbitrary
structured data with errors in a similar way to what `slog` does for logging.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Key Concepts](#key-concepts)
- [Design Philosophy](#design-philosophy)
- [Contributing](#contributing)
- [License](#license)
- [Related Resources](#related-resources)

## Overview

`serror` provides a clean API for creating, wrapping, and logging errors with
associated structured data.

When handling errors in Go, developers often need to include contextual
information to help diagnose what went wrong. The traditional approach using
`fmt.Errorf` interpolates this data into a string:

```go
return fmt.Errorf("failed to process user %d: %w", userID, err)
```

While simple, this approach has limitations:

- The contextual data (like `userID`) becomes just part of the error message
- There's no way to programmatically access these values later
- The data loses its type information
- The format becomes rigid and harder to parse

`serror` solves these problems by allowing you to attach structured data to
errors:

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

Inspired by Go's `log/slog` package, `serror` uses the same familiar key-value
pair pattern for attaching attributes. If you've used `slog`, you'll feel right
at home:

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

- **Structured Attributes**: Attach key-value pairs to errors using `slog`-like
  syntax
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
// Output:
// ERROR request failed err.userID=123 err.path=/api/v1/users err.msg="failed to process request"
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
// Note: With() returns a new error because serror errors are immutable
oldErr := serror.New("operation failed", "code", 500)
newErr := oldErr.With("retry", true, "attempt", 3)

// oldErr is unchanged, newErr has all attributes
```

The immutable design ensures thread safety and prevents accidental
modifications. Each call to `With()` returns a new error instance with the
combined attributes from the original error plus the new ones.

### Accessing Attributes

```go
// Get simple attribute values
value := err.Get("code").Int64()

// When getting attributes from groups, you can use a dot notation:
value := err.Get("user.id").Int64()
name := err.Get("user.name").String()
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

When logged, all structured attributes are automatically included in the log
output, making it easy to correlate errors with their context.

## Design Philosophy

- **Standard Library Aligned**: Works with `errors.Is`, `errors.As`,
  `errors.Unwrap`
- **slog Compatible**: Follows `slog` patterns for attribute handling
- **Thread Safe**: Immutable design prevents data races
- **Simple API**: Familiar interface for Go developers

## Documentation

For detailed documentation and API references, please visit our
[Go Package Documentation](https://pkg.go.dev/github.com/jetify-com/serror).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major
changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## Community

- **Issues:** [GitHub Issues](https://github.com/jetify-com/serror/issues)
- **Discussions:**
  [GitHub Discussions](https://github.com/jetify-com/serror/discussions)

## License

This project is licensed under the Apache License, Version 2.0 - see the
[LICENSE](LICENSE) file for details.

## Related Resources

### Articles

A list of helpful articles and resources about error handling in Go.

- [Errors are values](https://go.dev/blog/errors-are-values) by Rob Pike\
  A foundational article explaining Go's philosophy that errors are just values
  that can be programmed with, not exceptional conditions.
- [Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors) by Damien
  Neil and Jonathan Amsterdam\
  Introduces the error wrapping features added in Go 1.13 including unwrapping,
  Is and As functions.
- [Experiment, Simplify, Ship](https://go.dev/blog/experiment) by Russ Cox\
  Discusses the Go development process and philosophy, including how error
  handling evolved and continues to improve in Go.
- [Error handling in Upspin](https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html)
  by Rob Pike and Andrew Gerrand\
  A detailed look at Upspin's custom error package that prioritizes both
  user-friendly messages and diagnostic details for developers.
- [Failure is your Domain](https://www.gobeyond.dev/failure-is-your-domain/) by
  Ben Johnson\
  Explores the concept of domain-specific errors and why building a custom error
  package can improve application design.
- [Working with errors.As()](https://blog.carlana.net/post/2020/working-with-errors-as/)
  by Carl Scharenberg\
  A practical guide to using Go 1.13's errors.As() function to safely
  type-assert wrapped errors.

### Talks

A list of helpful talks about error handling in Go.

- [Don't Just Check Errors Handle Them Gracefully](https://www.youtube.com/watch?v=lsBF58Q-DnY)
  by Dave Cheney\
  A talk about patterns for handling errors in Go.
  ([Slides](https://github.com/gophercon/2016-talks/blob/master/DaveCheney-DontCheckErrorsHandleThemGracefully/GopherCon%202016.pdf))
- [Error Handling in Go](https://www.youtube.com/watch?v=YkOa5ZrNR_s) by Marwan
  Sulaiman\
  A talk about programmable errors and how to design an architecture to manage
  system failures.
  ([Slides](https://github.com/gophercon/2019-talks/blob/master/MarwanSulaiman-HandlingGoErrors/handling-go-errors.pdf))

### Libraries

A list of other error handling libraries you might find useful.

- [pkg/errors](https://github.com/pkg/errors) by Dave Cheney\
  Provides error wrapping with stack traces and was the inspiration for Go
  1.13's error wrapping. Still very popular, but archived.
- [cockroachdb/errors](https://github.com/cockroachdb/errors) by Cockroach Labs\
  A comprehensive error handling library with rich features including stack
  traces, error wrapping, and cause hierarchies. Fairly complex.
- [emperror](https://github.com/emperror/emperror) by Márk Sági-Kazár\
  A collection of error handling tools with adapters for various logging
  libraries.
- [eris](https://github.com/rotisserie/eris) by Rotisserie\
  Features stack traces, error wrapping, and formatted error output.
- [juju/errors](https://github.com/juju/errors) by Canonical\
  Provides an easy way to annotate errors without losing the original error
  context.
- [upspin/errors](https://github.com/upspin/upspin/tree/master/errors) by Rob
  Pike and the Upspin authors\
  An interesting approach to structured errors with error codes and operation
  metadata that has inspired others.
