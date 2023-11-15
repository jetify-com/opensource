package registry

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/codeclysm/extract/v3"
	"go.jetpack.io/pkg/runx/impl/fileutil"
)

func Extract(ctx context.Context, src string, dest string) error {
	tmpDest := src + ".contents"
	defer os.RemoveAll(tmpDest)

	reader, err := os.Open(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	err = extract.Archive(ctx, reader, tmpDest, nil /* no renaming of files */)
	if err != nil {
		return err
	}

	// Automatically flatten contents if they are inside a single directory
	srcDir := contentDir(tmpDest)

	parent := filepath.Dir(dest)
	err = fileutil.EnsureDir(parent)
	if err != nil {
		return err
	}

	err = os.Rename(srcDir, dest)
	if err != nil {
		return err
	}
	return nil
}

func contentDir(path string) string {
	contents, err := os.ReadDir(path)
	if err != nil {
		return path
	}
	if len(contents) != 1 {
		return path
	}
	if !contents[0].IsDir() {
		return path
	}
	return filepath.Join(path, contents[0].Name())
}

func createSymbolicLink(src, dst, repoName string) error {
	if err := os.MkdirAll(dst, 0700); err != nil {
		return err
	}
	if err := os.Chmod(src, 0755); err != nil {
		return err
	}
	binaryName := filepath.Base(src)
	// This is a good guess for the binary name. In the future we could allow
	// user to customize.
	if strings.Contains(binaryName, repoName) {
		binaryName = repoName
	}
	err := os.Symlink(src, filepath.Join(dst, binaryName))
	if errors.Is(err, os.ErrExist) {
		// TODO: verify symlink points to the right place
		return nil
	}
	return err
}
