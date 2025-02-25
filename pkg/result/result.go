package result

import (
	"errors"
	"fmt"
)

// Result is a type that can hold either a Value (T) or an error.
type Result[T any] struct {
	value T
	err   error
}

// Constructors
// ------------

// Ok returns a Result containing the provided Value.
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value}
}

// Err returns a Result containing the provided error.
func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

// Errf returns a Result containing a formatted error.
func Errf[T any](format string, args ...interface{}) Result[T] {
	return Err[T](fmt.Errorf(format, args...))
}

// From converts a (Value, error) pair into a Result.
func From[T any](value T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return Ok(value)
}

// Predicates
// ----------

// IsOk reports whether the Result holds a valid Value (Err == nil).
func (r Result[T]) IsOk() bool {
	return r.err == nil
}

// IsErr reports whether the Result holds an error (Err != nil).
func (r Result[T]) IsErr() bool {
	return r.err != nil
}

// Methods
// -------

// Unwrap returns the underlying error, if any. This is
// useful for integrating with Go 1.13+ error wrapping.
func (r Result[T]) Err() error {
	return r.err
}

// Get returns (Value, error). If r.IsErr(), Value will be
// the zero Value for T.
func (r Result[T]) Get() (T, error) {
	return r.value, r.err
}

// MustGet returns the Value if r.IsOk(), otherwise it panics.
// Use with caution in production code.
func (r Result[T]) MustGet() T {
	if r.err != nil {
		panic(r.err)
	}
	return r.value
}

// OrElse returns the stored Value if IsOk(), or else returns fallback.
func (r Result[T]) OrElse(fallback T) T {
	if r.err != nil {
		return fallback
	}
	return r.value
}

// Formatting
// --

// String provides a textual representation for debugging/logging.
func (r Result[T]) String() string {
	if r.IsOk() {
		return fmt.Sprintf("Ok(%v)", r.value)
	}
	return fmt.Sprintf("Err(%v)", r.err)
}

// GoString provides a Go-syntax representation (used in fmt %#v, etc.).
func (r Result[T]) GoString() string {
	if r.IsOk() {
		return fmt.Sprintf("Ok[%T](%#v)", r.value, r.value)
	}
	return fmt.Sprintf("Err[%T](%q)", r.value, r.err)
}

// Actions
// -------

// Do executes a function and wraps its result in a Result type.
// If the function panics, the panic is caught and converted to an error Result.
// For panics that are already errors, they are used directly.
// For other panic values, they are converted to error strings.
func Do[T any](fn func() T) (result Result[T]) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				result = Err[T](err)
			} else {
				result = Err[T](errors.New(fmt.Sprint(r)))
			}
		}
	}()
	return Ok(fn())
}

// Go runs a function asynchronously and returns a channel of its Result.
// The function is expected to return a Value and an error.
// The channel is closed after the function completes and its result is sent.
func Go[T any](f func() (T, error)) chan Result[T] {
	out := make(chan Result[T])
	go func() {
		defer close(out)
		out <- From(f())
	}()
	return out
}
