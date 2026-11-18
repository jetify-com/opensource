package filecache

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func newTestCache[T any](t *testing.T) *Cache[T] {
	t.Helper()
	return New("test-domain", WithCacheDir[T](t.TempDir()))
}

func TestNewDefaultsToUserCacheDir(t *testing.T) {
	c := New[string]("domain")
	if c.domain != "domain" {
		t.Errorf("domain = %q, want %q", c.domain, "domain")
	}
	if c.cacheDir == "" {
		t.Error("cacheDir is empty, want a default")
	}
}

func TestWithCacheDir(t *testing.T) {
	dir := t.TempDir()
	c := New("domain", WithCacheDir[int](dir))
	if c.cacheDir != dir {
		t.Errorf("cacheDir = %q, want %q", c.cacheDir, dir)
	}
}

func TestGetMissing(t *testing.T) {
	c := newTestCache[string](t)
	val, err := c.Get("missing")
	if !errors.Is(err, NotFound) {
		t.Errorf("Get() error = %v, want NotFound", err)
	}
	if val != "" {
		t.Errorf("Get() = %q, want zero value", val)
	}
	if !IsCacheMiss(err) {
		t.Error("IsCacheMiss(NotFound) = false")
	}
}

func TestSetGet(t *testing.T) {
	type payload struct {
		Name  string
		Count int
	}
	c := newTestCache[payload](t)

	want := payload{Name: "x", Count: 3}
	if err := c.Set("key", want, time.Hour); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	got, err := c.Get("key")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got != want {
		t.Errorf("Get() = %+v, want %+v", got, want)
	}
}

func TestSetOverwrites(t *testing.T) {
	c := newTestCache[int](t)
	if err := c.Set("k", 1, time.Hour); err != nil {
		t.Fatal(err)
	}
	if err := c.Set("k", 2, time.Hour); err != nil {
		t.Fatal(err)
	}
	if got, err := c.Get("k"); err != nil || got != 2 {
		t.Errorf("Get() = %v, %v; want 2, nil", got, err)
	}
}

func TestKeysAreIndependent(t *testing.T) {
	c := newTestCache[string](t)
	if err := c.Set("a", "A", time.Hour); err != nil {
		t.Fatal(err)
	}
	if err := c.Set("b", "B", time.Hour); err != nil {
		t.Fatal(err)
	}
	for key, want := range map[string]string{"a": "A", "b": "B"} {
		if got, err := c.Get(key); err != nil || got != want {
			t.Errorf("Get(%q) = %q, %v; want %q, nil", key, got, err, want)
		}
	}
}

func TestDomainsAreIsolated(t *testing.T) {
	dir := t.TempDir()
	a := New("a", WithCacheDir[string](dir))
	b := New("b", WithCacheDir[string](dir))
	if err := a.Set("key", "from-a", time.Hour); err != nil {
		t.Fatal(err)
	}
	if _, err := b.Get("key"); !errors.Is(err, NotFound) {
		t.Errorf("b.Get() error = %v, want NotFound", err)
	}
}

func TestExpired(t *testing.T) {
	c := newTestCache[string](t)
	if err := c.Set("key", "stale", -time.Second); err != nil {
		t.Fatal(err)
	}
	val, err := c.Get("key")
	if !errors.Is(err, Expired) {
		t.Errorf("Get() error = %v, want Expired", err)
	}
	// The stale value is still returned alongside the Expired error.
	if val != "stale" {
		t.Errorf("Get() = %q, want %q", val, "stale")
	}
	if !IsCacheMiss(err) {
		t.Error("IsCacheMiss(Expired) = false")
	}
}

func TestSetWithTime(t *testing.T) {
	c := newTestCache[string](t)

	if err := c.SetWithTime("future", "v", time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	if got, err := c.Get("future"); err != nil || got != "v" {
		t.Errorf("Get(future) = %q, %v; want v, nil", got, err)
	}

	if err := c.SetWithTime("past", "v", time.Now().Add(-time.Hour)); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Get("past"); !errors.Is(err, Expired) {
		t.Errorf("Get(past) error = %v, want Expired", err)
	}
}

func TestGetCorruptFile(t *testing.T) {
	c := newTestCache[string](t)
	if err := os.WriteFile(c.filename("key"), []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := c.Get("key")
	if err == nil {
		t.Fatal("Get() error = nil, want JSON error")
	}
	if IsCacheMiss(err) {
		t.Error("a corrupt file should not be reported as a cache miss")
	}
}

func TestSetUnmarshalableValue(t *testing.T) {
	c := newTestCache[chan int](t)
	if err := c.Set("key", make(chan int), time.Hour); err == nil {
		t.Error("Set() error = nil, want marshal error")
	}
	if err := c.SetWithTime("key", make(chan int), time.Now()); err == nil {
		t.Error("SetWithTime() error = nil, want marshal error")
	}
}

func TestGetOrSet(t *testing.T) {
	c := newTestCache[int](t)
	calls := 0
	f := func() (int, time.Duration, error) {
		calls++
		return 42, time.Hour, nil
	}

	for i := 0; i < 3; i++ {
		got, err := c.GetOrSet("key", f)
		if err != nil || got != 42 {
			t.Fatalf("GetOrSet() = %v, %v; want 42, nil", got, err)
		}
	}
	if calls != 1 {
		t.Errorf("f called %d times, want 1", calls)
	}
}

func TestGetOrSetRefreshesExpired(t *testing.T) {
	c := newTestCache[int](t)
	if err := c.Set("key", 1, -time.Second); err != nil {
		t.Fatal(err)
	}
	got, err := c.GetOrSet("key", func() (int, time.Duration, error) {
		return 2, time.Hour, nil
	})
	if err != nil || got != 2 {
		t.Errorf("GetOrSet() = %v, %v; want 2, nil", got, err)
	}
	if got, err := c.Get("key"); err != nil || got != 2 {
		t.Errorf("Get() after refresh = %v, %v; want 2, nil", got, err)
	}
}

func TestGetOrSetError(t *testing.T) {
	c := newTestCache[int](t)
	wantErr := errors.New("boom")
	_, err := c.GetOrSet("key", func() (int, time.Duration, error) {
		return 0, time.Hour, wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Errorf("GetOrSet() error = %v, want %v", err, wantErr)
	}
	if _, err := c.Get("key"); !errors.Is(err, NotFound) {
		t.Errorf("value should not be cached on error, Get() error = %v", err)
	}
}

func TestGetOrSetPropagatesNonMissErrors(t *testing.T) {
	c := newTestCache[int](t)
	if err := os.WriteFile(c.filename("key"), []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	called := false
	_, err := c.GetOrSet("key", func() (int, time.Duration, error) {
		called = true
		return 1, time.Hour, nil
	})
	if err == nil {
		t.Error("GetOrSet() error = nil, want error from corrupt file")
	}
	if called {
		t.Error("f should not be called when Get fails with a non-miss error")
	}
}

func TestGetOrSetWithTime(t *testing.T) {
	c := newTestCache[string](t)
	calls := 0
	f := func() (string, time.Time, error) {
		calls++
		return "v", time.Now().Add(time.Hour), nil
	}
	for i := 0; i < 2; i++ {
		got, err := c.GetOrSetWithTime("key", f)
		if err != nil || got != "v" {
			t.Fatalf("GetOrSetWithTime() = %q, %v; want v, nil", got, err)
		}
	}
	if calls != 1 {
		t.Errorf("f called %d times, want 1", calls)
	}

	wantErr := errors.New("boom")
	_, err := c.GetOrSetWithTime("other", func() (string, time.Time, error) {
		return "", time.Time{}, wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Errorf("GetOrSetWithTime() error = %v, want %v", err, wantErr)
	}
}

func TestClear(t *testing.T) {
	dir := t.TempDir()
	c := New("domain", WithCacheDir[string](dir))
	other := New("other", WithCacheDir[string](dir))
	if err := c.Set("key", "v", time.Hour); err != nil {
		t.Fatal(err)
	}
	if err := other.Set("key", "v", time.Hour); err != nil {
		t.Fatal(err)
	}

	if err := c.Clear(); err != nil {
		t.Fatalf("Clear() error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "domain")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("domain dir still exists after Clear(): %v", err)
	}
	if _, err := c.Get("key"); !errors.Is(err, NotFound) {
		t.Errorf("Get() after Clear() error = %v, want NotFound", err)
	}
	if _, err := other.Get("key"); err != nil {
		t.Errorf("Clear() affected another domain: %v", err)
	}

	// Clearing an already-empty cache is not an error.
	if err := c.Clear(); err != nil {
		t.Errorf("second Clear() error = %v", err)
	}
}

func TestIsCacheMiss(t *testing.T) {
	if IsCacheMiss(nil) {
		t.Error("IsCacheMiss(nil) = true")
	}
	if IsCacheMiss(errors.New("other")) {
		t.Error("IsCacheMiss(other) = true")
	}
}

func TestSetWritesPrivateFiles(t *testing.T) {
	c := newTestCache[string](t)
	path := c.filename("secret")

	// Simulate a cache file left behind by an older version with loose perms.
	if err := os.WriteFile(path, []byte(`{}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := c.Set("secret", "value", time.Hour); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("cache file perm = %o, want 600", perm)
	}

	dirInfo, err := os.Stat(filepath.Dir(path))
	if err != nil {
		t.Fatal(err)
	}
	if perm := dirInfo.Mode().Perm(); perm != 0o700 {
		t.Errorf("cache dir perm = %o, want 700", perm)
	}

	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Errorf("cache dir has %d entries, want 1 (temp file leaked?)", len(entries))
	}
}

func TestGetTightensLegacyPermissions(t *testing.T) {
	c := newTestCache[string](t)
	if err := c.Set("creds", "secret", time.Hour); err != nil {
		t.Fatal(err)
	}
	path := c.filename("creds")
	dir := filepath.Dir(path)

	// Simulate a cache written by an older version: 0755 dir, 0644 file.
	if err := os.Chmod(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, 0o644); err != nil {
		t.Fatal(err)
	}

	// A valid cache hit must not leave the old permissions in place.
	val, err := c.Get("creds")
	if err != nil || val != "secret" {
		t.Fatalf("Get() = %q, %v; want %q, nil", val, err, "secret")
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("cache file perm = %o, want 600", perm)
	}
	dirInfo, err := os.Stat(dir)
	if err != nil {
		t.Fatal(err)
	}
	if perm := dirInfo.Mode().Perm(); perm != 0o700 {
		t.Errorf("cache dir perm = %o, want 700", perm)
	}
}
