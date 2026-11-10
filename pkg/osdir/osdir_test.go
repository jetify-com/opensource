package osdir

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

// skipIfSystemUser skips tests that exercise the user or search directories,
// which aren't used when running as the system user (such as root).
func skipIfSystemUser(t *testing.T) {
	t.Helper()
	if isSystemUser() {
		t.Skip("skipping test of user directories when running as the system user")
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

// writeTestFile creates the file at path, including any parent directories.
func writeTestFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestDirTypeSubSearch(t *testing.T) {
	// testSub checks that a file existing only under the first search
	// directory's subdirectory is found with both Sub and a joined path.
	testSub := func(t *testing.T, dt DirType, first, second string) {
		t.Helper()

		data := []byte("hello world")
		writeTestFile(t, filepath.Join(first, "app", "config.json"), data)

		sub := dt.Sub("app")
		got, err := sub.ReadFile("config.json")
		if err != nil {
			t.Fatalf("Sub(%q).ReadFile(%q) error = %v", "app", "config.json", err)
		}
		if string(got) != string(data) {
			t.Errorf("Sub(%q).ReadFile(%q) = %q, want %q", "app", "config.json", got, data)
		}

		got, err = dt.ReadFile("app/config.json")
		if err != nil {
			t.Fatalf("ReadFile(%q) error = %v", "app/config.json", err)
		}
		if string(got) != string(data) {
			t.Errorf("ReadFile(%q) = %q, want %q", "app/config.json", got, data)
		}

		paths, err := sub.JoinSearchPath("config.json")
		if err != nil {
			t.Fatalf("Sub(%q).JoinSearchPath(%q) error = %v", "app", "config.json", err)
		}
		want := []string{
			filepath.Join(dt.User, "app", "config.json"),
			filepath.Join(first, "app", "config.json"),
			filepath.Join(second, "app", "config.json"),
		}
		if !slices.Equal(paths, want) {
			t.Errorf("Sub(%q).JoinSearchPath(%q) = %q, want %q", "app", "config.json", paths, want)
		}
	}

	t.Run("SearchDefault", func(t *testing.T) {
		skipIfSystemUser(t)

		first, second := t.TempDir(), t.TempDir()
		dt := DirType{
			System:        t.TempDir(),
			User:          t.TempDir(),
			SearchDefault: first + string(filepath.ListSeparator) + second,
		}
		testSub(t, dt, first, second)
	})

	t.Run("SearchEnvList", func(t *testing.T) {
		skipIfSystemUser(t)

		first, second := t.TempDir(), t.TempDir()
		t.Setenv("TEST_XDG_CONFIG_DIRS", first+string(filepath.ListSeparator)+second)
		dt := DirType{
			System:        t.TempDir(),
			User:          t.TempDir(),
			Search:        "$TEST_XDG_CONFIG_DIRS",
			SearchDefault: t.TempDir(),
		}
		testSub(t, dt, first, second)
	})
}
