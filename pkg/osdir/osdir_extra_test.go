//go:build unix

package osdir

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func forceSystem(t *testing.T) {
	t.Helper()
	forceSystemUser = true
	t.Cleanup(func() { forceSystemUser = false })
}

func TestSub(t *testing.T) {
	dt := DirType{
		System:        "/sys",
		User:          "/user",
		UserDefault:   "/userdef",
		Search:        "/search",
		SearchDefault: "/searchdef",
	}
	got := dt.Sub("app")
	want := DirType{
		System:      "/sys/app",
		User:        "/user/app",
		UserDefault: "/userdef/app",
		// Search and SearchDefault are lists that may contain environment
		// variables, so the subdirectory is joined onto each search
		// directory later instead of onto the raw values.
		Search:        "/search",
		SearchDefault: "/searchdef",
		sub:           "app",
	}
	if got != want {
		t.Errorf("Sub() = %+v, want %+v", got, want)
	}
	// Sub has a value receiver and must not modify the original.
	if dt.System != "/sys" {
		t.Errorf("Sub() modified the receiver: %+v", dt)
	}
}

func TestJoinPathInvalidPath(t *testing.T) {
	dt := DirType{System: "/sys", User: "/user"}
	for _, p := range []string{"", "../escape", "a/../../escape", "/abs/path"} {
		if _, err := dt.JoinPath(p); err == nil {
			t.Errorf("JoinPath(%q) error = nil, want error", p)
		}
		if _, err := dt.JoinSearchPath(p); err == nil {
			t.Errorf("JoinSearchPath(%q) error = nil, want error", p)
		}
	}
}

func TestJoinPathNested(t *testing.T) {
	forceSystem(t)
	dt := DirType{System: "/sys"}
	got, err := dt.JoinPath("a/b/c.json")
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join("/sys", "a", "b", "c.json"); got != want {
		t.Errorf("JoinPath() = %q, want %q", got, want)
	}
}

func TestJoinPathNoSuitableDirectory(t *testing.T) {
	t.Run("System", func(t *testing.T) {
		forceSystem(t)
		t.Setenv("OSDIR_TEST_UNSET", "")
		for _, sys := range []string{"", "relative/dir", "$OSDIR_TEST_UNSET/x"} {
			dt := DirType{System: sys, User: "/user"}
			if _, err := dt.JoinPath("f"); err == nil {
				t.Errorf("JoinPath() with System=%q error = nil, want error", sys)
			}
		}
	})

	t.Run("User", func(t *testing.T) {
		skipIfSystemUser(t)
		t.Setenv("OSDIR_TEST_UNSET", "")
		dt := DirType{System: "/sys", User: "$OSDIR_TEST_UNSET", UserDefault: "relative"}
		if _, err := dt.JoinPath("f"); err == nil {
			t.Error("JoinPath() error = nil, want error")
		}
	})
}

func TestJoinPathExpandsEnv(t *testing.T) {
	skipIfSystemUser(t)
	t.Setenv("OSDIR_TEST_DIR", "/from/env")
	dt := DirType{User: "${OSDIR_TEST_DIR}/sub", UserDefault: "/default"}
	got, err := dt.JoinPath("f")
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join("/from/env/sub", "f"); got != want {
		t.Errorf("JoinPath() = %q, want %q", got, want)
	}
}

func TestJoinSearchPathSystem(t *testing.T) {
	forceSystem(t)
	dt := DirType{System: "/sys", Search: "/a:/b"}
	got, err := dt.JoinSearchPath("f")
	if err != nil {
		t.Fatal(err)
	}
	// System users never search additional directories.
	if want := []string{filepath.Join("/sys", "f")}; !slices.Equal(got, want) {
		t.Errorf("JoinSearchPath() = %q, want %q", got, want)
	}

	if _, err := (DirType{Search: "/a"}).JoinSearchPath("f"); err == nil {
		t.Error("JoinSearchPath() with no System dir error = nil, want error")
	}
}

func TestJoinSearchPathUser(t *testing.T) {
	skipIfSystemUser(t)
	t.Setenv("OSDIR_TEST_UNSET", "")
	t.Setenv("OSDIR_TEST_SEARCH", "/env")

	tests := []struct {
		name    string
		dt      DirType
		want    []string
		wantErr bool
	}{
		{
			name: "user then search",
			dt:   DirType{User: "/user", Search: "/a:/b", SearchDefault: "/def"},
			want: []string{"/user/f", "/a/f", "/b/f"},
		},
		{
			name: "search default when search is empty",
			dt:   DirType{User: "/user", SearchDefault: "/def1:/def2"},
			want: []string{"/user/f", "/def1/f", "/def2/f"},
		},
		{
			name: "search from env",
			dt:   DirType{User: "/user", Search: "$OSDIR_TEST_SEARCH"},
			want: []string{"/user/f", "/env/f"},
		},
		{
			name: "skips relative and unexpandable search dirs",
			dt:   DirType{User: "/user", Search: "relative:$OSDIR_TEST_UNSET/x:/ok"},
			want: []string{"/user/f", "/ok/f"},
		},
		{
			name: "search only when user dir is unusable",
			dt:   DirType{User: "$OSDIR_TEST_UNSET", Search: "/a"},
			want: []string{"/a/f"},
		},
		{
			name: "user only",
			dt:   DirType{User: "/user"},
			want: []string{"/user/f"},
		},
		{
			name:    "nothing usable",
			dt:      DirType{User: "$OSDIR_TEST_UNSET", Search: "relative"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.dt.JoinSearchPath("f")
			if tt.wantErr {
				if err == nil {
					t.Errorf("JoinSearchPath() = %q, want error", got)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			want := make([]string, len(tt.want))
			for i, w := range tt.want {
				want[i] = filepath.FromSlash(w)
			}
			if !slices.Equal(got, want) {
				t.Errorf("JoinSearchPath() = %q, want %q", got, want)
			}
		})
	}
}

func TestReadWriteFileSystem(t *testing.T) {
	forceSystem(t)
	dt := DirType{System: t.TempDir()}

	if _, err := dt.ReadFile("missing.txt"); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("ReadFile(missing) error = %v, want ErrNotExist", err)
	}

	if err := dt.WriteFile("nested/dir/file.txt", []byte("one")); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	// Overwrites truncate existing content.
	if err := dt.WriteFile("nested/dir/file.txt", []byte("2")); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	got, err := dt.ReadFile("nested/dir/file.txt")
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(got) != "2" {
		t.Errorf("ReadFile() = %q, want %q", got, "2")
	}

	info, err := os.Stat(filepath.Join(dt.System, "nested/dir/file.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("file perms = %v, want 0600", perm)
	}
}

func TestWriteFileErrors(t *testing.T) {
	forceSystem(t)
	if err := (DirType{System: t.TempDir()}).WriteFile("../escape", nil); err == nil {
		t.Error("WriteFile(../escape) error = nil, want error")
	}
	if err := (DirType{}).WriteFile("f", nil); err == nil {
		t.Error("WriteFile() with no directory error = nil, want error")
	}

	// A regular file where a parent directory should be.
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "blocker"), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := (DirType{System: root}).WriteFile("blocker/f", nil); err == nil {
		t.Error("WriteFile() under a regular file error = nil, want error")
	}
}

func TestExpand(t *testing.T) {
	t.Setenv("HOME", "/home/tester")
	t.Setenv("OSDIR_A", "a")
	t.Setenv("OSDIR_B", "b")
	t.Setenv("OSDIR_EMPTY", "")

	tests := map[string]string{
		"/plain/path":           "/plain/path",
		"$HOME/.cache":          "/home/tester/.cache",
		"${HOME}/x":             "/home/tester/x",
		"/$OSDIR_A/${OSDIR_B}":  "/a/b",
		"$OSDIR_EMPTY/x":        "",
		"/ok/$OSDIR_EMPTY":      "",
		"$OSDIR_A/$OSDIR_EMPTY": "",
		"":                      "",
	}
	for in, want := range tests {
		if got := expand(in); got != want {
			t.Errorf("expand(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestHomeDir(t *testing.T) {
	t.Setenv("HOME", "/home/tester")
	got, err := HomeDir()
	if err != nil || got != "/home/tester" {
		t.Errorf("HomeDir() = %q, %v; want /home/tester, nil", got, err)
	}

	// With HOME unset, HomeDir falls back to the user database.
	t.Setenv("HOME", "")
	got, err = HomeDir()
	if err == nil && got == "" {
		t.Error("HomeDir() returned an empty path without an error")
	}
}

func TestBinDirs(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	xdgBin := filepath.Join(home, ".local", "bin")
	homeBin := filepath.Join(home, "bin")

	withSystem := func(dirs []string) []string {
		if isSystemUser() {
			return append([]string{"/usr/local/bin"}, dirs...)
		}
		return append(dirs, "/usr/local/bin")
	}

	tests := []struct {
		name string
		path string
		want []string
	}{
		{
			name: "neither in PATH",
			path: "/usr/bin:/bin",
			want: []string{xdgBin, homeBin},
		},
		{
			name: "xdg bin first",
			path: "/usr/bin:" + xdgBin + ":/bin",
			want: []string{xdgBin, homeBin},
		},
		{
			name: "home bin first",
			path: homeBin + "/:/usr/bin",
			want: []string{homeBin, xdgBin},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PATH", tt.path)
			for _, system := range []bool{false, true} {
				forceSystemUser = system
				got := BinDirs()
				if want := withSystem(slices.Clone(tt.want)); !slices.Equal(got, want) {
					t.Errorf("BinDirs() (system=%v) = %q, want %q", isSystemUser(), got, want)
				}
			}
			forceSystemUser = false
		})
	}
}
