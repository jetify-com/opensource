package osdir

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDirTypeJoinPath(t *testing.T) {
	t.Run("User", func(t *testing.T) {
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
