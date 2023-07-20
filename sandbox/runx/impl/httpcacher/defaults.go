package httpcacher

import (
	"path/filepath"

	"github.com/adrg/xdg"
)

const xdgSubdir = "jetpack.io/http"

// It's important to note that with the current implementation, the cache
// must be a private cache: we're doing nothing to filter out requests/responses
// with cookies and other sensitive information, so if it is shared, a user could
// use that to gain access to resources they shouldn't have access to.
//
// TODO: consider supporting a shared cache. Consider changing the default caching
// directory structure, to separate the private cache from the shared cache (which
// could be copied between machines).
var defaultCacheDir = filepath.Join(xdg.CacheHome, xdgSubdir)
