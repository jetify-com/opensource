package filecache

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// writeFile atomically writes data to path with 0600 permissions. Cached
// values may contain credentials, so they must not be readable by other
// users. Writing to a temp file and renaming also replaces any existing file
// (and its possibly looser permissions) and avoids leaving a truncated file
// behind if the process dies mid-write.
func writeFile(path string, data []byte) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".tmp-*")
	if err != nil {
		return errors.WithStack(err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName) // no-op after a successful rename

	// CreateTemp already uses 0600, but be explicit.
	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()
		return errors.WithStack(err)
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return errors.WithStack(err)
	}
	if err := tmp.Close(); err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(os.Rename(tmpName, path))
}
