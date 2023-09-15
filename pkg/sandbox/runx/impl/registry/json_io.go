package registry

import (
	"encoding/json"
	"os"
	"time"

	"go.jetpack.io/pkg/sandbox/runx/impl/fileutil"
)

// TODO: Generalize package by:
// - Making the cache TTL configurable
// - Making the backing cache configurable (e.g. in-memory, filesystem, etc.)
// - Allowing for passthrough cases (in-memory + filesystem, etc.)

const ttl = 24 * time.Hour

func fetchCachedJSON[T any](abspath string, fetchFunc func() (T, error)) (T, error) {
	path := fileutil.Path(abspath)
	var results T

	// First try to load the cached copy:
	loadErr := readJSON(path.String(), &results)
	fresh := path.FileInfo() != nil && time.Since(path.FileInfo().ModTime()) < ttl

	// If we loaded it without errors and it's fresh, return it:
	if loadErr == nil && fresh {
		return results, nil
	}

	// Otherwise, fetch it:
	results, fetchErr := fetchFunc()
	// TODO: we could use something like `singleflight` to ensure that if there are
	// concurrent calls to `fetchFunc` in the same process that only one of them runs.

	// We failed to fetch from the web, but have a stale local copy. Return that
	// as a best effort:
	if fetchErr != nil && loadErr == nil {
		return results, nil
	}

	// We failed to fetch it and don't have a local copy:
	if fetchErr != nil {
		var zero T
		// return the fetch error:
		return zero, fetchErr
	}

	// We retrieved from the web, we're gonna return that.

	// First, update the cache.
	// Saving to the cache is best effort, even if we fail, we'll still want to return
	// the fresh data we just fetched
	_ = writeJSON(path.String(), results)

	// Now return the fresh data:
	return results, nil
}

func readJSON(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func writeJSON(path string, v any) error {
	// We're writing with indentation for debugging purposes, but if we want to
	// optimize storage, we could save without indentation (or in a binary format like
	// CBOR or BSON)
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	// Use renameio to make the write atomic
	err = fileutil.WriteFile(path, data)
	if err != nil {
		return err
	}
	return nil
}
