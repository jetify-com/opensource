package fileutil

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/google/renameio/v2"
)

type osFS struct{}

func (osFS) Open(name string) (fs.File, error) {
	return os.Open(name)
}

// IsDir returns true if the path exists and is a directory.
func IsDir(path string) bool {
	return isDir(osFS{}, path)
}

// IsFile returns true if the path exists and is a regular file.
func IsFile(path string) bool {
	return isFile(osFS{}, path)
}

// Exists returns true if the path exists.
func Exists(path string) bool {
	return exists(osFS{}, path)
}

// FileInfo returns the fs.FileInfo for the given path.
func FileInfo(path string) fs.FileInfo {
	return fileInfo(osFS{}, path)
}

// Glob returns all files that match the given pattern.
func Glob(pattern string) ([]string, error) {
	return doublestar.Glob(osFS{}, pattern)
}

// EnsureDir ensures that the directory at the given path exists,
// creating it and any parent directories if necessary.
func EnsureDir(path string) error {
	if IsDir(path) {
		return nil
	}
	return os.MkdirAll(path, 0o700 /* as suggested by xdg spec */)
}

// WriteFile writes data to the named file, creating it if necessary.
// If the file already exists, it is replaced.
// The file is written atomically by writing to a temporary file and renaming it.
func WriteFile(path string, data []byte) error {
	// First ensure the directory exists:
	dir := filepath.Dir(path)
	err := EnsureDir(dir)
	if err != nil {
		return err
	}

	return renameio.WriteFile(path, data, 0o600)
}

// isDir returns true if the path exists and is a directory in the given filesystem.
func isDir(fsys fs.FS, path string) bool {
	info := fileInfo(fsys, path)
	if info == nil {
		return false
	}
	return info.IsDir()
}

// isFile returns true if the path exists and is a regular file in the given filesystem.
func isFile(fsys fs.FS, path string) bool {
	info := fileInfo(fsys, path)
	if info == nil {
		return false
	}
	return info.Mode().IsRegular()
}

// exists returns true if the path exists in the given filesystem.
func exists(fsys fs.FS, path string) bool {
	return fileInfo(fsys, path) != nil
}

// fileInfo returns the fs.FileInfo for the given path in the given filesystem.
// Returns nil if the path does not exist or cannot be accessed.
func fileInfo(fsys fs.FS, path string) fs.FileInfo {
	info, err := fs.Stat(fsys, path)
	if err != nil {
		return nil
	}
	return info
}
