package try

import (
	"errors"
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

// Wrap converts a (Value, error) pair into a Try.
func Wrap[T any](value T, err error) Try[T] {
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

// Unwrap returns (Value, error). If r.IsErr(), Value will be
// the zero Value for T.
func (r Try[T]) Unwrap() (T, error) {
	return r.value, r.err
}

// Get returns the underlying value, if any.
func (r Try[T]) Get() T {
	return r.value
}

// MustGet returns the Value if r.IsOk(), otherwise it panics.
// Use with caution in production code.
func (r Try[T]) MustGet() T {
	if r.err != nil {
		panic(r.err)
	}
	return r.value
}

// OrElse returns the stored Value if IsOk(), or else returns fallback.
func (r Try[T]) OrElse(fallback T) T {
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

// Actions
// -------

// Do executes a function and wraps its result in a Try type.
// If the function panics, the panic is caught and converted to an error Try.
// For panics that are already errors, they are used directly.
// For other panic values, they are converted to error strings.
func Do[T any](fn func() T) (result Try[T]) {
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

// Go runs a function asynchronously and returns a channel of type Try[T].
// The function is expected to return a Value and an error.
// The channel is closed after the function completes and its result is sent.
func Go[T any](f func() (T, error)) chan Try[T] {
	out := make(chan Try[T])
	go func() {
		defer close(out)
		out <- Wrap(f())
	}()
	return out
}
