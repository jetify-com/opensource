// Package osdir provides a cross-platform interface for well-known system and
// user directories.
//
// The pre-defined [DirType] package variables use default paths that are
// specific to the current operating system. The standard XDG environment
// variables (XDG_CACHE_HOME, XDG_CONFIG_HOME, XDG_DATA_HOME, XDG_STATE_HOME)
// override these paths when set to a non-empty value, even on macOS and
// Windows.
//
// For bin directories, [BinDirs] returns common binary directories that are
// sorted according to their order in PATH.
package osdir

import (
	"cmp"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"slices"
)

// HomeDir returns the current user's home directory by first calling
// [os.UserHomeDir]. If that fails, it attempts to get the home directory from
// [os/user.Current].
func HomeDir() (string, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		u, err := user.Current()
		if err != nil {
			return "", err
		}
		if u.HomeDir == "" {
			return "", errors.New("current user has no home directory")
		}
		return u.HomeDir, nil
	}
	return path, nil
}

// BinDirs returns directories for storing user binaries.
func BinDirs() []string {
	// The two most common directories for user binaries are ~/.local/bin
	// and ~/bin. Prioritize whichever comes first in $PATH. If neither are
	// in $PATH, use the XDG standard ~/.local/bin.
	dirs := make([]string, 0, 3)
	xdgBin := expand("$HOME/.local/bin")
	homeBin := expand("$HOME/bin")
	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		switch filepath.Clean(dir) {
		case xdgBin:
			dirs = append(dirs, xdgBin, homeBin)
		case homeBin:
			dirs = append(dirs, homeBin, xdgBin)
		}
	}
	if len(dirs) == 0 {
		dirs = append(dirs, xdgBin, homeBin)
	}

	// If we're root, put /usr/local/bin first, otherwise try it last.
	if isSystemUser() {
		dirs = slices.Insert(dirs, 0, "/usr/local/bin")
	} else {
		dirs = append(dirs, "/usr/local/bin")
	}
	return dirs
}

// DirType specifies the directory paths for a certain category of files. All
// fields must be absolute paths. Paths may reference environment variables with
// ${var} or $var. Paths that reference undefined or empty environment variables
// expand to an empty string.
type DirType struct {
	// System is the directory to use when running as the system user. On
	// Unix, a process running with an euid of 0 is a system user. On
	// Windows, it's a process running with an acess token that has elevated
	// UAC privileges.
	System string

	// User is the directory to use when running as a non-system user.
	User string

	// UserDefault is the default directory to use when User is empty or
	// contains environment variables that resolve to an empty string.
	UserDefault string

	// Search specifies additional directories to search when reading files.
	// Methods that write data such as WriteFile or MkdirAll do not consult
	// these directories.
	Search string

	// SearchDefault is the default list of directories to use when Search
	// is empty or contains environment variables that resolve to an empty
	// string.
	SearchDefault string
}

// Sub returns a DirType whose paths are subdirectories of d. For example, the
// following calls read the same config file:
//
//	d.ReadFile("app/config.json")
//	d.Sub("app").ReadFile("config.json")
func (d DirType) Sub(dir string) DirType {
	d.System = filepath.Join(d.System, dir)
	d.User = filepath.Join(d.User, dir)
	d.UserDefault = filepath.Join(d.UserDefault, dir)
	d.Search = filepath.Join(d.Search, dir)
	d.SearchDefault = filepath.Join(d.SearchDefault, dir)
	return d
}

// JoinPath joins path to the system or user directory of d and returns the
// resulting absolute path. The path must be relative and cannot contain any
// ".." elements.
func (d DirType) JoinPath(path string) (string, error) {
	err := validPath(path)
	if err != nil {
		return "", err
	}

	basepath := ""
	if isSystemUser() {
		basepath = expand(d.System)
	} else {
		basepath = cmp.Or(expand(d.User), expand(d.UserDefault))
	}

	if !filepath.IsAbs(basepath) {
		return "", fmt.Errorf("no suitable directory")
	}
	return filepath.Join(basepath, path), nil
}

// JoinSearchPath returns a slice of absolute paths by joining path with the
// base directory and search directories of d. The path must be relative and
// cannot contain any ".." elements.
func (d DirType) JoinSearchPath(path string) ([]string, error) {
	firstPath, err := d.JoinPath(path)
	if isSystemUser() {
		// No searching for system paths.
		if err != nil {
			return nil, err
		}
		return []string{firstPath}, nil
	}

	var joined []string
	if firstPath != "" {
		joined = append(joined, firstPath)
	}

	search := filepath.SplitList(cmp.Or(d.Search, d.SearchDefault))
	for _, basepath := range search {
		basepath = expand(basepath)
		if filepath.IsAbs(basepath) {
			joined = append(joined, filepath.Join(basepath, path))
		}
	}
	if len(joined) == 0 {
		return nil, fmt.Errorf("no suitable directory")
	}
	return joined, nil
}

// OpenFile calls [DirType.JoinPath] with the given filename and opens the
// resulting path. When flag has the [os.O_RDONLY] bit set, OpenFile calls
// [DirType.JoinSearchPath] and tries opening each path until it succeeds. When
// flag has the [os.O_CREATE] bit set, OpenFile will create any necessary
// directories if they don't already exist.
func (d DirType) OpenFile(name string, flag int, perm fs.FileMode) (*os.File, error) {
	// Don't use search paths when writing.
	if flag&os.O_RDONLY == 0 {
		absPath, err := d.JoinPath(name)
		if err != nil {
			return nil, err
		}
		return mkdirOpenFile(absPath, flag, perm)
	}

	// If we're opening read-only, then also try the search path list until
	// we find the file.
	search, err := d.JoinSearchPath(name)
	if err != nil {
		return nil, err
	}

	var errs []error
	for _, dir := range search {
		f, err := mkdirOpenFile(filepath.Join(dir, name), flag, perm)
		if err == nil {
			return f, nil
		}
		errs = append(errs, err)
	}
	return nil, errors.Join(errs...)
}

// ReadFile calls [DirType.OpenFile] on the named file and returns its contents.
func (d DirType) ReadFile(name string) ([]byte, error) {
	f, err := d.OpenFile(name, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	return b, err
}

// WriteFile writes data to the named file, creating it and any parent
// directories if necessary. It truncates the file if it already exists.
func (d DirType) WriteFile(name string, data []byte) error {
	f, err := d.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}

// expand is like [os.ExpandEnv], but differs in two ways to better support
// file paths:
//
//  1. If any environment variables fail to expand (or expand to an empty
//     string) it returns an empty string.
//  2. HOME expands using [HomeDir], which will lookup the user's home directory
//     from other locations when the environment variable is not set.
func expand(path string) string {
	missing := false
	expanded := os.Expand(path, func(k string) string {
		v := ""
		if k == "HOME" {
			v, _ = HomeDir()
		} else {
			v = os.Getenv(k)
		}
		missing = missing || v == ""
		return v
	})
	if missing {
		return ""
	}
	return expanded
}

// mkdirOpenFile opens a file and creates any missing parent directories if flag
// has [os.O_CREATE] set.
func mkdirOpenFile(path string, flag int, perm fs.FileMode) (*os.File, error) {
	if flag&os.O_CREATE != 0 {
		err := os.MkdirAll(filepath.Dir(path), 0o700)
		if err != nil {
			return nil, err
		}
	}
	return os.OpenFile(path, flag, perm)
}

// validPath returns an error if path is not local as defined by
// [filepath.IsLocal].
func validPath(path string) error {
	if path == "" {
		return errors.New("path must not be empty")
	}
	if !filepath.IsLocal(path) {
		return fmt.Errorf("path must be relative without any \"..\" elements: %s", path)
	}
	return nil
}
