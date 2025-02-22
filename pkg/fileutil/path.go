package fileutil

import (
	"io/fs"
	"path/filepath"
)

// Path represents a filesystem path with helper methods for common operations.
// It provides a type-safe way to handle paths and perform path-related operations.
type Path string

// String returns the path as a string.
func (p Path) String() string {
	return string(p)
}

// Subpath creates a new Path by joining the current path with additional path elements.
// It uses filepath.Join internally to handle path separators correctly for the current OS.
func (p Path) Subpath(elements ...string) Path {
	all := append([]string{p.String()}, elements...)
	return Path(filepath.Join(all...))
}

// IsDir returns true if the path exists and is a directory.
// The path is relative to the current working directory.
func (p Path) IsDir() bool {
	return IsDir(p.String())
}

// IsFile returns true if the path exists and is a regular file.
// The path is relative to the current working directory.
func (p Path) IsFile() bool {
	return IsFile(p.String())
}

// EnsureDir ensures that the directory at this path exists,
// creating it and any parent directories if necessary.
func (p Path) EnsureDir() error {
	return EnsureDir(p.String())
}

// FileInfo returns the fs.FileInfo for this path.
// Returns nil if the path does not exist or cannot be accessed.
// The path is relative to the current working directory.
func (p Path) FileInfo() fs.FileInfo {
	return FileInfo(p.String())
}

// Exists returns true if the path exists.
// The path is relative to the current working directory.
func (p Path) Exists() bool {
	return Exists(p.String())
}
