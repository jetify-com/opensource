package download

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/cavaliergopher/grab/v3"
	"go.jetpack.io/pkg/runx/impl/fileutil"
)

func DownloadOnce(url string, dest string) error {
	dir := filepath.Dir(dest)
	if err := fileutil.EnsureDir(dir); err != nil {
		return err
	}

	info := fileutil.FileInfo(dest)
	if info != nil && info.Mode().IsRegular() && info.Size() > 0 {
		// We've already downloaded it
		return nil
	}
	return Download(url, dest)
}

func Download(url string, dest string) error {
	if fileutil.IsDir(dest) {
		return errors.New("destination is a directory")
	}

	client := grab.DefaultClient
	tmpDest := dest + ".crdownload"

	// Grab supports partial and resumable downloads, so we don't automatically
	// delete tmpDest until we're done with the download. That way, if we retry
	// we can continue the download where we left off. That said, if the remote file
	// changes, this can result in a corrupted file. For now we assume remote files
	// don't change, since we're are downloading "immutable" releases, but to be safe
	// we'll want to add checksum validation.

	req, err := grab.NewRequest(tmpDest, url)
	if err != nil {
		return err
	}

	resp := client.Do(req)
	if err := resp.Err(); err != nil {
		return err
	}

	err = os.Rename(tmpDest, dest)
	if err != nil {
		return err
	}

	return nil
}
