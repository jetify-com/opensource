package serror

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"runtime"
	"strings"
	"time"

	"github.com/lmittmann/tint"
	"go.jetify.com/pkg/serror/internal/record"
)

// Error represents a structured error with attributes that conform to slog conventions.
// Error is immutable - methods like With() return new instances rather than
// modifying the original error.
type Error struct {
	record record.Record
	cause  error
}

// Ensure Error implements the error interface
var _ error = Error{}

// Ensure Error implements slog.LogValuer
var _ slog.LogValuer = Error{}

// For testing
var timeNow = time.Now

// New creates a new Error with the given message and attributes.
// The attributes should be alternating key-value pairs, similar to slog conventions.
func New(msg string, args ...any) Error {
	// Capture caller's program counter
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip runtime.Callers + New

	e := new(timeNow(), msg, nil, pcs[0])
	e.add(args...)

	return e
}

func new(t time.Time, msg string, cause error, pc uintptr) Error {
	return Error{
		cause:  cause,
		record: record.NewRecord(t, msg, pc),
	}
}

// Wrap creates a new Error that wraps an existing error with additional context.
// If cause is nil, it behaves like New.
func Wrap(cause error, msg string, args ...any) Error {
	// Capture caller's program counter
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip runtime.Callers + Wrap

	e := new(timeNow(), msg, cause, pcs[0])
	e.add(args...)

	return e
}

func (e *Error) add(args ...any) {
	if e == nil {
		return
	}
	e.record.Add(args...)
}

func (e Error) Error() string {
	var buf bytes.Buffer
	e.format(&buf)
	return strings.TrimSuffix(buf.String(), "\n")
}

func (e Error) format(w io.Writer) {
	// Create an slog.Record
	record := e.record.ToSlog()
	record.Time = time.Time{} // Set time to zero so we don't print it

	// If we have a cause, add it as an error attribute
	if e.cause != nil {
		record.AddAttrs(slog.String("cause", e.cause.Error()))
	}

	// Create a tint handler that writes to our writer
	handler := tint.NewHandler(w, &tint.Options{
		Level:   slog.LevelError,
		NoColor: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove the level attribute
			if a.Key == slog.LevelKey {
				return slog.Attr{}
			}
			return a
		},
	})

	// Handle the complete record
	_ = handler.Handle(context.Background(), record)
}

// Unwrap returns the underlying cause of the error, if any.
// This method implements Go's error unwrapping interface, allowing Error
// to work with errors.Is, errors.As, and errors.Unwrap.
func (e Error) Unwrap() error {
	return e.cause
}

// With returns a new Error with additional attributes added to the existing ones.
// The attributes should be alternating key-value pairs, similar to slog conventions.
func (e Error) With(args ...any) Error {
	newErr := e.clone()
	newErr.add(args...)

	return newErr
}

// clone returns a copy of the error with all its attributes and cause.
// Modifications to the cloned error will not affect the original.
func (e Error) clone() Error {
	return Error{
		cause:  e.cause,
		record: e.record.Clone(),
	}
}

// LogValue implements slog.LogValuer, allowing Error to be used as a logging value.
// It returns a structured representation of the error including its message, attributes,
// and cause if present.
func (e Error) LogValue() slog.Value {
	attrs := make([]slog.Attr, e.record.NumAttrs()+2)

	// Add the message
	attrs = append(attrs, slog.String("msg", e.record.Message))

	e.record.Attrs(func(a record.Attr) bool {
		attrs = append(attrs, a.ToSlog())
		return true
	})
	// TODO: Change to use.Error as the message?

	// Add the cause
	if e.cause != nil {
		attrs = append(attrs, slog.Any("cause", e.cause))
	}

	return slog.GroupValue(attrs...)
}
