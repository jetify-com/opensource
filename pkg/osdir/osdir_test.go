package osdir

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// skipIfSystemUser skips tests that exercise the non-system user directories,
// which are never used when the test process is running as the system user
// (for example, as root).
func skipIfSystemUser(t *testing.T) {
	t.Helper()
	if isSystemUser() {
		t.Skip("skipping user directory test when running as the system user")
	}
}

func TestDirTypeJoinPath(t *testing.T) {
	t.Run("User", func(t *testing.T) {
		skipIfSystemUser(t)

		dt := DirType{
			System:      "/var/cache",
			User:        "/home/user/custom",
			UserDefault: "/home/user/.cache",
		}

		path := "."
		got, err := dt.JoinPath(path)
		if err != nil {
			t.Errorf("JoinPath(%q) error = %v", path, err)
		}

		want := filepath.Join(dt.User, path)
		if got != want {
			t.Errorf("JoinPath(%q) = %q, want %q", path, got, want)
		}
	})

	t.Run("UserDefault", func(t *testing.T) {
		skipIfSystemUser(t)

		t.Setenv("TEST_XDG_CACHE_HOME", "")

		dt := DirType{
			System:      "/var/cache",
			User:        "$TEST_XDG_CACHE_HOME",
			UserDefault: "/home/user/.cache",
		}

		path := "."
		got, err := dt.JoinPath(path)
		if err != nil {
			t.Errorf("JoinPath(%q) error = %v", path, err)
		}

		want := filepath.Join(dt.UserDefault, path)
		if got != want {
			t.Errorf("JoinPath(%q) = %q, want %q", path, got, want)
		}
	})

	t.Run("System", func(t *testing.T) {
		forceSystemUser = true
		t.Cleanup(func() { forceSystemUser = false })

		dt := DirType{
			System:      "/var/cache",
			User:        "/home/user/custom",
			UserDefault: "/home/user/.cache",
		}

		path := "."
		got, err := dt.JoinPath(path)
		if err != nil {
			t.Errorf("JoinPath(%q) error = %v", path, err)
		}

		want := filepath.Join(dt.System, path)
		if got != want {
			t.Errorf("JoinPath(%q) = %q, want %q", path, got, want)
		}
	})
}

func TestDirTypeWriteFile(t *testing.T) {
	t.Run("User", func(t *testing.T) {
		skipIfSystemUser(t)

		systemDir := t.TempDir()
		userDir := t.TempDir()
		userDefaultDir := t.TempDir()

		dt := DirType{
			System:      systemDir,
			User:        userDir,
			UserDefault: userDefaultDir,
		}

		path := "test.txt"
		data := []byte("hello world")

		err := dt.WriteFile(path, data)
		if err != nil {
			t.Errorf("WriteFile(%q, %q) error = %v", path, data, err)
		}

		fullPath := filepath.Join(userDir, path)
		got, err := os.ReadFile(fullPath)
		if err != nil {
			t.Errorf("ReadFile(%q) error = %v", fullPath, err)
		}

		if string(got) != string(data) {
			t.Errorf("ReadFile(%q) = %q, want %q", fullPath, got, data)
		}
	})

	t.Run("UserDefault", func(t *testing.T) {
		skipIfSystemUser(t)

		t.Setenv("TEST_XDG_CACHE_HOME", "")

		systemDir := t.TempDir()
		userDefaultDir := t.TempDir()

		dt := DirType{
			System:      systemDir,
			User:        "$TEST_XDG_CACHE_HOME",
			UserDefault: userDefaultDir,
		}

		path := "test.txt"
		data := []byte("hello world")

		err := dt.WriteFile(path, data)
		if err != nil {
			t.Errorf("WriteFile(%q, %q) error = %v", path, data, err)
		}

		fullPath := filepath.Join(userDefaultDir, path)
		got, err := os.ReadFile(fullPath)
		if err != nil {
			t.Errorf("ReadFile(%q) error = %v", fullPath, err)
		}

		if string(got) != string(data) {
			t.Errorf("ReadFile(%q) = %q, want %q", fullPath, got, data)
		}
	})

	t.Run("System", func(t *testing.T) {
		forceSystemUser = true
		t.Cleanup(func() { forceSystemUser = false })

		systemDir := t.TempDir()
		userDir := t.TempDir()
		userDefaultDir := t.TempDir()

		dt := DirType{
			System:      systemDir,
			User:        userDir,
			UserDefault: userDefaultDir,
		}

		path := "test.txt"
		data := []byte("hello world")

		err := dt.WriteFile(path, data)
		if err != nil {
			t.Errorf("WriteFile(%q, %q) error = %v", path, data, err)
		}

		fullPath := filepath.Join(systemDir, path)
		got, err := os.ReadFile(fullPath)
		if err != nil {
			t.Errorf("ReadFile(%q) error = %v", fullPath, err)
		}

		if string(got) != string(data) {
			t.Errorf("ReadFile(%q) = %q, want %q", fullPath, got, data)
		}
	})
}

func TestDirTypeReadFile(t *testing.T) {
	const name = "test.txt"

	t.Run("SearchFallback", func(t *testing.T) {
		skipIfSystemUser(t)

		userDir := t.TempDir()
		searchDir := t.TempDir()
		dt := DirType{
			System: t.TempDir(),
			User:   userDir,
			Search: searchDir,
		}

		want := []byte("from search dir")
		if err := os.WriteFile(filepath.Join(searchDir, name), want, 0o600); err != nil {
			t.Fatal(err)
		}

		got, err := dt.ReadFile(name)
		if err != nil {
			t.Fatalf("ReadFile(%q) error = %v", name, err)
		}
		if string(got) != string(want) {
			t.Errorf("ReadFile(%q) = %q, want %q", name, got, want)
		}
	})

	t.Run("SearchDefaultFallback", func(t *testing.T) {
		skipIfSystemUser(t)
		t.Setenv("TEST_XDG_CONFIG_DIRS", "")

		userDir := t.TempDir()
		searchDir1 := t.TempDir()
		searchDir2 := t.TempDir()
		dt := DirType{
			System:        t.TempDir(),
			User:          userDir,
			Search:        "$TEST_XDG_CONFIG_DIRS",
			SearchDefault: searchDir1 + string(filepath.ListSeparator) + searchDir2,
		}

		want := []byte("from second search dir")
		if err := os.WriteFile(filepath.Join(searchDir2, name), want, 0o600); err != nil {
			t.Fatal(err)
		}

		got, err := dt.ReadFile(name)
		if err != nil {
			t.Fatalf("ReadFile(%q) error = %v", name, err)
		}
		if string(got) != string(want) {
			t.Errorf("ReadFile(%q) = %q, want %q", name, got, want)
		}
	})

	t.Run("UserWins", func(t *testing.T) {
		skipIfSystemUser(t)

		userDir := t.TempDir()
		searchDir := t.TempDir()
		dt := DirType{
			System: t.TempDir(),
			User:   userDir,
			Search: searchDir,
		}

		want := []byte("from user dir")
		if err := os.WriteFile(filepath.Join(userDir, name), want, 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(searchDir, name), []byte("from search dir"), 0o600); err != nil {
			t.Fatal(err)
		}

		got, err := dt.ReadFile(name)
		if err != nil {
			t.Fatalf("ReadFile(%q) error = %v", name, err)
		}
		if string(got) != string(want) {
			t.Errorf("ReadFile(%q) = %q, want %q", name, got, want)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		skipIfSystemUser(t)

		dt := DirType{
			System: t.TempDir(),
			User:   t.TempDir(),
			Search: t.TempDir(),
		}
		_, err := dt.ReadFile(name)
		if !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("ReadFile(%q) error = %v, want %v", name, err, fs.ErrNotExist)
		}
	})

	t.Run("SystemIgnoresSearch", func(t *testing.T) {
		forceSystemUser = true
		t.Cleanup(func() { forceSystemUser = false })

		systemDir := t.TempDir()
		searchDir := t.TempDir()
		dt := DirType{
			System: systemDir,
			User:   t.TempDir(),
			Search: searchDir,
		}
		if err := os.WriteFile(filepath.Join(searchDir, name), []byte("from search dir"), 0o600); err != nil {
			t.Fatal(err)
		}

		_, err := dt.ReadFile(name)
		if !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("ReadFile(%q) error = %v, want %v", name, err, fs.ErrNotExist)
		}

		want := []byte("from system dir")
		if err := os.WriteFile(filepath.Join(systemDir, name), want, 0o600); err != nil {
			t.Fatal(err)
		}
		got, err := dt.ReadFile(name)
		if err != nil {
			t.Fatalf("ReadFile(%q) error = %v", name, err)
		}
		if string(got) != string(want) {
			t.Errorf("ReadFile(%q) = %q, want %q", name, got, want)
		}
	})
}
