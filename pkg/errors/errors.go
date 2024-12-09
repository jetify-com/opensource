package errors

import (
	stderrors "errors"

	"github.com/pkg/errors"
)

// New returns a new error with the given message.
// It is a wrapper around github.com/pkg/errors.New
func New(message string) error {
	return errors.New(message)
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
// It is a wrapper around github.com/pkg/errors.Errorf
func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

// WithMessage annotates err with a new message.
// If err is nil, WithMessage returns nil.
func WithMessage(err error, message string) error {
	return errors.WithMessage(err, message)
}

// WithMessagef annotates err with the format specifier.
// If err is nil, WithMessagef returns nil.
func WithMessagef(err error, format string, args ...interface{}) error {
	return errors.WithMessagef(err, format, args...)
}

// WithStack annotates err with a stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
func WithStack(err error) error {
	return errors.WithStack(err)
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//	type causer interface {
//	       Cause() error
//	}
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	return errors.Cause(err)
}

// Is reports whether any error in err's chain matches target.
// It is a wrapper around errors.Is from the standard library
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true. Otherwise, it returns false.
// It is a wrapper around errors.As from the standard library
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
// It is a wrapper around errors.Unwrap from the standard library
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// Join returns an error that wraps the given errors.
// Any nil error values are discarded.
// Join returns nil if errs contains no non-nil values.
// The error formats as the concatenation of the strings obtained
// by calling the Error method of each element of errs, with a newline
// between each string.
// It is a wrapper around errors.Join from the standard library
func Join(errs ...error) error {
	return stderrors.Join(errs...)
}
