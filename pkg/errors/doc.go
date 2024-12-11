package errors

// Package errors provides a wrapper around github.com/pkg/errors that adds
// errors.Join from the standard library.
//
// It is a drop-in replacement for github.com/pkg/errors.
//
// If we do ever want to drop the dependency on github.com/pkg/errors, we can
// do so by removing the wrapper and implementing the missing functions.
