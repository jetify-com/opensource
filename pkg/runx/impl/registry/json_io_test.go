package registry

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type cachedValue struct {
	Name string `json:"name"`
}

func TestFetchCachedJSONMissFetchesAndCaches(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "cache.json")
	calls := 0
	fetch := func() (cachedValue, error) {
		calls++
		return cachedValue{Name: "fresh"}, nil
	}

	got, err := fetchCachedJSON(path, fetch)
	if err != nil || got.Name != "fresh" {
		t.Fatalf("fetchCachedJSON() = %+v, %v; want fresh, nil", got, err)
	}

	var onDisk cachedValue
	if err := readJSON(path, &onDisk); err != nil || onDisk.Name != "fresh" {
		t.Fatalf("cache file = %+v, %v; want fresh", onDisk, err)
	}

	// A second call is served from the fresh cache.
	got, err = fetchCachedJSON(path, fetch)
	if err != nil || got.Name != "fresh" {
		t.Fatalf("second fetchCachedJSON() = %+v, %v", got, err)
	}
	if calls != 1 {
		t.Errorf("fetch called %d times, want 1", calls)
	}
}

func TestFetchCachedJSONStaleRefetches(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cache.json")
	if err := writeJSON(path, cachedValue{Name: "stale"}); err != nil {
		t.Fatal(err)
	}
	old := time.Now().Add(-2 * ttl)
	if err := os.Chtimes(path, old, old); err != nil {
		t.Fatal(err)
	}

	got, err := fetchCachedJSON(path, func() (cachedValue, error) {
		return cachedValue{Name: "fresh"}, nil
	})
	if err != nil || got.Name != "fresh" {
		t.Fatalf("fetchCachedJSON() = %+v, %v; want fresh, nil", got, err)
	}
	var onDisk cachedValue
	if err := readJSON(path, &onDisk); err != nil || onDisk.Name != "fresh" {
		t.Errorf("cache file = %+v, %v; want it refreshed", onDisk, err)
	}
}

func TestFetchCachedJSONFetchErrorFallsBackToStaleCache(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cache.json")
	if err := writeJSON(path, cachedValue{Name: "stale"}); err != nil {
		t.Fatal(err)
	}
	old := time.Now().Add(-2 * ttl)
	if err := os.Chtimes(path, old, old); err != nil {
		t.Fatal(err)
	}

	got, err := fetchCachedJSON(path, func() (cachedValue, error) {
		return cachedValue{}, errors.New("network down")
	})
	if err != nil {
		t.Fatalf("fetchCachedJSON() error = %v, want stale fallback", err)
	}
	if got.Name != "stale" {
		t.Errorf("fetchCachedJSON() = %+v, want the stale cached value", got)
	}
}

func TestFetchCachedJSONFetchErrorWithoutCache(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cache.json")
	wantErr := errors.New("network down")
	got, err := fetchCachedJSON(path, func() (cachedValue, error) {
		return cachedValue{Name: "partial"}, wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Errorf("error = %v, want %v", err, wantErr)
	}
	if got != (cachedValue{}) {
		t.Errorf("fetchCachedJSON() = %+v, want zero value on error", got)
	}
}

func TestFetchCachedJSONCorruptCacheRefetches(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cache.json")
	if err := os.WriteFile(path, []byte("{not json"), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := fetchCachedJSON(path, func() (cachedValue, error) {
		return cachedValue{Name: "fresh"}, nil
	})
	if err != nil || got.Name != "fresh" {
		t.Fatalf("fetchCachedJSON() = %+v, %v; want fresh, nil", got, err)
	}
}

func TestWriteJSONUnmarshalable(t *testing.T) {
	if err := writeJSON(filepath.Join(t.TempDir(), "x.json"), make(chan int)); err == nil {
		t.Error("writeJSON(chan) error = nil, want error")
	}
}

func TestContentDir(t *testing.T) {
	t.Run("single subdirectory is unwrapped", func(t *testing.T) {
		root := t.TempDir()
		inner := filepath.Join(root, "tool-1.0")
		if err := os.Mkdir(inner, 0o755); err != nil {
			t.Fatal(err)
		}
		if got := contentDir(root); got != inner {
			t.Errorf("contentDir() = %q, want %q", got, inner)
		}
	})

	t.Run("single file", func(t *testing.T) {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, "bin"), nil, 0o644); err != nil {
			t.Fatal(err)
		}
		if got := contentDir(root); got != root {
			t.Errorf("contentDir() = %q, want %q", got, root)
		}
	})

	t.Run("multiple entries", func(t *testing.T) {
		root := t.TempDir()
		for _, d := range []string{"a", "b"} {
			if err := os.Mkdir(filepath.Join(root, d), 0o755); err != nil {
				t.Fatal(err)
			}
		}
		if got := contentDir(root); got != root {
			t.Errorf("contentDir() = %q, want %q", got, root)
		}
	})

	t.Run("missing", func(t *testing.T) {
		missing := filepath.Join(t.TempDir(), "missing")
		if got := contentDir(missing); got != missing {
			t.Errorf("contentDir() = %q, want %q", got, missing)
		}
	})
}

func TestCreateSymbolicLink(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "pkg", "mytool_linux_amd64")
	if err := os.MkdirAll(filepath.Dir(src), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(src, []byte("#!/bin/sh\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(root, "bin")

	if err := createSymbolicLink(src, dst, "mytool"); err != nil {
		t.Fatalf("createSymbolicLink() error = %v", err)
	}

	// The link is named after the repo when the binary name contains it.
	link := filepath.Join(dst, "mytool")
	target, err := os.Readlink(link)
	if err != nil {
		t.Fatalf("Readlink(%q) error = %v", link, err)
	}
	if target != src {
		t.Errorf("link target = %q, want %q", target, src)
	}
	info, err := os.Stat(src)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm()&0o111 == 0 {
		t.Errorf("source perms = %v, want executable", info.Mode().Perm())
	}

	// Creating the link again is a no-op.
	if err := createSymbolicLink(src, dst, "mytool"); err != nil {
		t.Errorf("second createSymbolicLink() error = %v", err)
	}
}

func TestCreateSymbolicLinkKeepsBinaryName(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "other-binary")
	if err := os.WriteFile(src, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(root, "bin")
	if err := createSymbolicLink(src, dst, "mytool"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(filepath.Join(dst, "other-binary")); err != nil {
		t.Errorf("expected link named after the binary: %v", err)
	}
}

func TestCreateSymbolicLinkMissingSource(t *testing.T) {
	root := t.TempDir()
	if err := createSymbolicLink(filepath.Join(root, "missing"), filepath.Join(root, "bin"), "x"); err == nil {
		t.Error("createSymbolicLink() with missing source error = nil")
	}
}
