// filecache is a simple local file-based cache
package filecache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

var NotFound = errors.New("not found")
var Expired = errors.New("expired")

type Cache[T any] struct {
	domain   string
	cacheDir string
}

type data[T any] struct {
	Val T
	Exp time.Time
}

type Option[T any] func(*Cache[T])

func New[T any](domain string, opts ...Option[T]) *Cache[T] {
	result := &Cache[T]{domain: domain}

	var err error
	result.cacheDir, err = os.UserCacheDir()
	if err != nil {
		result.cacheDir = "~/.cache"
	}

	for _, opt := range opts {
		opt(result)
	}

	return result
}

func WithCacheDir[T any](dir string) Option[T] {
	return func(c *Cache[T]) {
		c.cacheDir = dir
	}
}

func (c *Cache[T]) Set(key string, val T, dur time.Duration) error {
	d, err := json.Marshal(data[T]{Val: val, Exp: time.Now().Add(dur)})
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(os.WriteFile(c.filename(key), d, 0644))
}

func (c *Cache[T]) SetT(key string, val T, t time.Time) error {
	d, err := json.Marshal(data[T]{Val: val, Exp: t})
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(os.WriteFile(c.filename(key), d, 0644))
}

func (c *Cache[T]) Get(key string) (T, error) {
	path := c.filename(key)
	resultData := data[T]{}

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return resultData.Val, NotFound
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return resultData.Val, errors.WithStack(err)
	}

	if err := json.Unmarshal(content, &resultData); err != nil {
		return resultData.Val, errors.WithStack(err)
	}
	if time.Now().After(resultData.Exp) {
		return resultData.Val, Expired
	}
	return resultData.Val, nil
}

func (c *Cache[T]) GetOrSet(
	key string,
	f func() (T, time.Duration, error),
) (T, error) {
	if val, err := c.Get(key); err == nil || !IsCacheMiss(err) {
		return val, err
	}

	val, dur, err := f()
	if err != nil {
		return val, err
	}

	return val, c.Set(key, val, dur)
}

func (c *Cache[T]) GetOrSetT(
	key string,
	f func() (T, time.Time, error),
) (T, error) {
	if val, err := c.Get(key); err == nil || !IsCacheMiss(err) {
		return val, err
	}

	val, t, err := f()
	if err != nil {
		return val, err
	}

	return val, c.SetT(key, val, t)
}

func IsCacheMiss(err error) bool {
	return errors.Is(err, NotFound) || errors.Is(err, Expired)
}

func (c *Cache[T]) filename(key string) string {
	dir := filepath.Join(c.cacheDir, c.domain)
	_ = os.MkdirAll(dir, 0755)
	return filepath.Join(dir, key)
}
