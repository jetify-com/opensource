package httpcacher

import (
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
)

func newTransport(cacheDir string) *httpcache.Transport {
	cache := diskcache.New(cacheDir)
	return httpcache.NewTransport(cache)
}
