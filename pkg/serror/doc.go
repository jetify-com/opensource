// Package serror provides structured error handling for Go applications.
//
// serror makes it easy to associate arbitrary structured data with errors in a similar way
// to what slog does for logging. When handling errors, developers often need to include
// contextual information to help diagnose what went wrong. The traditional approach using
// fmt.Errorf interpolates this data into a string:
//
//	return fmt.Errorf("failed to process user %d: %w", userID, err)
//
// While simple, this approach has limitations:
//   - The contextual data becomes just part of the error message
//   - There's no way to programmatically access these values later
//   - The data loses its type information
//   - The format becomes rigid and harder to parse
//
// serror solves these problems by allowing you to attach structured data to errors:
//
//	return serror.Wrap(err, "failed to process user",
//	    "user_id", userID,
//	    "attempt", attempt,
//	    "status", status,
//	)
//
// Key Features:
//   - Structured Attributes: Attach key-value pairs to errors using slog-like syntax
//   - Error Wrapping: Fully compatible with Go's error wrapping conventions
//   - slog Integration: Implements LogValuer for automatic structured logging
//   - JSON Support: Full JSON marshaling/unmarshaling capabilities
//   - Attribute Access: Get attributes using dot notation (e.g., "user.id")
//   - Thread Safety: Immutable design ensures safe concurrent usage
//   - Call Site Capture: Automatically records where errors are created
//
// Basic Usage:
//
//	// Create a new error with attributes
//	err := serror.New("validation failed",
//	    "code", 400,
//	    "field", "email",
//	)
//
//	// Wrap an existing error
//	if err != nil {
//	    return serror.Wrap(err, "failed to validate user",
//	        "user_id", id,
//	        "email", email,
//	    )
//	}
//
//	// Add attributes to an existing error
//	// Note: With() returns a new error - serror errors are immutable
//	newErr := err.With("retry", true, "attempt", 3)
//
//	// Access attributes using dot notation
//	code := err.Get("code")
//	email := err.Get("user.email")
//
//	// Log the error with slog
//	logger.Error("operation failed", "err", err)
//	// Output: ERROR operation failed err.code=400 err.field=email err.msg="validation failed"
//
// The API is designed to be familiar to Go developers who use slog:
//
//	// slog style
//	logger.Error("request failed",
//	    "user_id", userID,
//	    "status", status,
//	)
//
//	// serror follows the same pattern
//	return serror.New("request failed",
//	    "user_id", userID,
//	    "status", status,
//	)
//
// Attribute constructors are also available for type safety:
//
//	err := serror.New("validation failed",
//	    serror.Int("code", 400),
//	    serror.String("field", "email"),
//	    serror.Bool("valid", false),
//	    serror.Group("user",
//	        serror.Int("id", 123),
//	        serror.String("name", "alice"),
//	    ),
//	)
//
// All errors are immutable - methods like With() return new instances rather than
// modifying the original error. This ensures thread safety and prevents accidental
// modifications.
package serror
