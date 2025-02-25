package types

import "errors"

var (
	ErrPackageNotFound      = errors.New("package not found")
	ErrReleaseNotFound      = errors.New("release not found")
	ErrPlatformNotSupported = errors.New("package doesn't support platform")
)
