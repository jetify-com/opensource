// TODO: publish as it's own shared package that other binaries can use.
// Right now we have other copies in other binaries. For example, devbox
// has its own copy.

package fileutil

import (
	"os"
	"path/filepath"
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
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func (p Path) IsDir() bool {
	return IsDir(p.String())
}

func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

func (p Path) IsFile() bool {
	return IsFile(p.String())
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
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
