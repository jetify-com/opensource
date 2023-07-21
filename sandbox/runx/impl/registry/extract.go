package registry

import (
	"context"
	"os"
	"path/filepath"

	"github.com/codeclysm/extract"
	"go.jetpack.io/runx/impl/fileutil"
)

func Extract(ctx context.Context, src string, dest string) error {
	tmpDest := src + ".contents"
	defer os.RemoveAll(tmpDest)

	reader, err := os.Open(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	err = extract.Archive(context.Background(), reader, tmpDest, nil /* no renaming of files */)
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
