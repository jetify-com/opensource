// TODO: publish as it's own shared package that other binaries can use.
// Right now we have other copies in other binaries. For example, devbox
// has its own copy.

package fileutil

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/rogpeppe/go-internal/renameio"
)

type Path string

func (p Path) String() string {
	return string(p)
}

func (p Path) Subpath(elements ...string) Path {
	all := append([]string{p.String()}, elements...)
	return Path(filepath.Join(all...))
}

func IsDir(path string) bool {
	info := FileInfo(path)
	if info == nil {
		return false
	}
	return info.IsDir()
}

func (p Path) IsDir() bool {
	return IsDir(p.String())
}

func IsFile(path string) bool {
	info := FileInfo(path)
	if info == nil {
		return false
	}
	return info.Mode().IsRegular()
}

func (p Path) IsFile() bool {
	return IsFile(p.String())
}

func Exists(path string) bool {
	return FileInfo(path) != nil
}

func (p Path) Exists() bool {
	return Exists(p.String())
}

func EnsureDir(path string) error {
	if IsDir(path) {
		return nil
	}
	return os.MkdirAll(path, 0700 /* as suggested by xdg spec */)
}

func (p Path) EnsureDir() error {
	return EnsureDir(p.String())
}

func FileInfo(path string) fs.FileInfo {
	info, err := os.Stat(path)
	if err != nil {
		return nil
	}
	return info
}

func (p Path) FileInfo() fs.FileInfo {
	return FileInfo(p.String())
}

func WriteFile(path string, data []byte) error {
	// First ensure the directory exists:
	dir := filepath.Dir(path)
	err := EnsureDir(dir)
	if err != nil {
		return err
	}
	// Write using `renameio` to ensure an atomic write:
	return renameio.WriteFile(path, data)
}
