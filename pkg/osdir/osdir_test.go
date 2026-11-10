package osdir

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// testFileName is the file name used by the ReadFile tests.
const testFileName = "test.txt"

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
		writeTestFile(t, searchDir, want)

		checkReadFile(t, dt, want)
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
		writeTestFile(t, searchDir2, want)

		checkReadFile(t, dt, want)
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
		writeTestFile(t, userDir, want)
		writeTestFile(t, searchDir, []byte("from search dir"))

		checkReadFile(t, dt, want)
	})

	t.Run("NotFound", func(t *testing.T) {
		skipIfSystemUser(t)

		dt := DirType{
			System: t.TempDir(),
			User:   t.TempDir(),
			Search: t.TempDir(),
		}
		checkReadFileNotExist(t, dt)
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
		writeTestFile(t, searchDir, []byte("from search dir"))

		checkReadFileNotExist(t, dt)

		want := []byte("from system dir")
		writeTestFile(t, systemDir, want)
		checkReadFile(t, dt, want)
	})
}

func writeTestFile(t *testing.T, dir string, data []byte) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, testFileName), data, 0o600); err != nil {
		t.Fatal(err)
	}
}

func checkReadFile(t *testing.T, dt DirType, want []byte) {
	t.Helper()
	got, err := dt.ReadFile(testFileName)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", testFileName, err)
	}
	if string(got) != string(want) {
		t.Errorf("ReadFile(%q) = %q, want %q", testFileName, got, want)
	}
}

func checkReadFileNotExist(t *testing.T, dt DirType) {
	t.Helper()
	_, err := dt.ReadFile(testFileName)
	if !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("ReadFile(%q) error = %v, want %v", testFileName, err, fs.ErrNotExist)
	}
}
