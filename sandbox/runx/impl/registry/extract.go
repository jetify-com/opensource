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

	parent := filepath.Dir(dest)
	err = fileutil.EnsureDir(parent)
	if err != nil {
		return err
	}

	err = os.Rename(tmpDest, dest)
	if err != nil {
		return err
	}
	return nil
}
