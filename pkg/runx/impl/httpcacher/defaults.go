package httpcacher

import (
	"os"
	"path/filepath"
)

const xdgSubdir = "jetify.com/http"

// It's important to note that with the current implementation, the cache
// must be a private cache: we're doing nothing to filter out requests/responses
// with cookies and other sensitive information, so if it is shared, a user could
// use that to gain access to resources they shouldn't have access to.
//
// TODO: consider supporting a shared cache. Consider changing the default caching
// directory structure, to separate the private cache from the shared cache (which
// could be copied between machines).

func defaultCacheDir() string {
	cacheHome, err := os.UserCacheDir()
	if err != nil {
		cacheHome = "~/.cache"
	}
	return filepath.Join(cacheHome, xdgSubdir)
}
