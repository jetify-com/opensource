package try

import (
	"fmt"
)

// Try is a type that can hold either a Value (T) or an error.
type Try[T any] struct {
	value T
	err   error
}

// Constructors
// ------------

// Ok returns a Try containing the provided Value.
func Ok[T any](value T) Try[T] {
	return Try[T]{value: value}
}

// Err returns a Try containing the provided error.
func Err[T any](err error) Try[T] {
	return Try[T]{err: err}
}

// Errf returns a Try containing a formatted error.
func Errf[T any](format string, args ...interface{}) Try[T] {
	return Err[T](fmt.Errorf(format, args...))
}

// From converts a (Value, error) pair into a Try.
func From[T any](value T, err error) Try[T] {
	if err != nil {
		return Err[T](err)
	}
	return Ok(value)
}

// Predicates
// ----------

// IsOk reports whether the Try holds a valid Value (Err == nil).
func (r Try[T]) IsOk() bool {
	return r.err == nil
}

// IsErr reports whether the Try holds an error (Err != nil).
func (r Try[T]) IsErr() bool {
	return r.err != nil
}

// Methods
// -------

// Err returns the underlying error, if any. This is
// useful for integrating with Go 1.13+ error wrapping.
func (r Try[T]) Err() error {
	return r.err
}

// Get returns (Value, error). If r.IsErr(), Value will be
// the zero Value for T.
func (r Try[T]) Get() (T, error) {
	return r.value, r.err
}

// MustGet returns the Value if r.IsOk(), otherwise it panics.
// Use with caution in production code.
func (r Try[T]) MustGet() T {
	if r.err != nil {
		panic(r.err)
	}
	return r.value
}

// GetOrElse returns the stored Value if IsOk(), or else returns fallback.
func (r Try[T]) GetOrElse(fallback T) T {
	if r.err != nil {
		return fallback
	}
	return r.value
}

// Formatting
// --

// String provides a textual representation for debugging/logging.
func (r Try[T]) String() string {
	if r.IsOk() {
		return fmt.Sprintf("Ok(%v)", r.value)
	}
	return fmt.Sprintf("Err(%v)", r.err)
}

// GoString provides a Go-syntax representation (used in fmt %#v, etc.).
func (r Try[T]) GoString() string {
	if r.IsOk() {
		return fmt.Sprintf("Ok[%T](%#v)", r.value, r.value)
	}
	return fmt.Sprintf("Err[%T](%q)", r.value, r.err)
}
