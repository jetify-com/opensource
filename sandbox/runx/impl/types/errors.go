package types

import "errors"

var ErrPackageNotFound = errors.New("package not found")
var ErrReleaseNotFound = errors.New("release not found")
var ErrPlatformNotSupported = errors.New("package doesn't support platform")
