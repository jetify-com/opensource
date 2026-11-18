package fileutil

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestPathString(t *testing.T) {
	if got := Path("/a/b").String(); got != "/a/b" {
		t.Errorf("String() = %q, want /a/b", got)
	}
}

func TestPathSubpath(t *testing.T) {
	tests := []struct {
		base Path
		elem []string
		want string
	}{
		{"/a", []string{"b", "c"}, filepath.Join("/a", "b", "c")},
		{"/a", nil, "/a"},
		{"/a/", []string{"b/"}, filepath.Join("/a", "b")},
		{"/a/b", []string{".."}, "/a"},
		{"rel", []string{"x"}, filepath.Join("rel", "x")},
	}
	for _, tt := range tests {
		if got := tt.base.Subpath(tt.elem...); got.String() != tt.want {
			t.Errorf("Path(%q).Subpath(%q) = %q, want %q", tt.base, tt.elem, got, tt.want)
		}
	}
}

func TestPathFilesystemHelpers(t *testing.T) {
	root := Path(t.TempDir())
	dir := root.Subpath("dir")
	file := root.Subpath("file.txt")
	missing := root.Subpath("missing")

	if err := os.Mkdir(dir.String(), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file.String(), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name                  string
		p                     Path
		isDir, isFile, exists bool
	}{
		{"dir", dir, true, false, true},
		{"file", file, false, true, true},
		{"missing", missing, false, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.IsDir(); got != tt.isDir {
				t.Errorf("IsDir() = %v, want %v", got, tt.isDir)
			}
			if got := tt.p.IsFile(); got != tt.isFile {
				t.Errorf("IsFile() = %v, want %v", got, tt.isFile)
			}
			if got := tt.p.Exists(); got != tt.exists {
				t.Errorf("Exists() = %v, want %v", got, tt.exists)
			}
			if got := tt.p.FileInfo() != nil; got != tt.exists {
				t.Errorf("FileInfo() != nil is %v, want %v", got, tt.exists)
			}
		})
	}
}

func TestEnsureDir(t *testing.T) {
	root := Path(t.TempDir())
	nested := root.Subpath("a", "b", "c")

	if err := nested.EnsureDir(); err != nil {
		t.Fatalf("EnsureDir() error = %v", err)
	}
	if !nested.IsDir() {
		t.Fatal("EnsureDir() did not create the directory")
	}
	// Idempotent.
	if err := nested.EnsureDir(); err != nil {
		t.Errorf("second EnsureDir() error = %v", err)
	}

	// Fails when a file is in the way.
	file := root.Subpath("file")
	if err := os.WriteFile(file.String(), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := file.Subpath("child").EnsureDir(); err == nil {
		t.Error("EnsureDir() under a regular file succeeded, want error")
	}
}

func TestWriteFile(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "nested", "dir", "out.txt")

	if err := WriteFile(path, []byte("first")); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := WriteFile(path, []byte("second")); err != nil {
		t.Fatalf("WriteFile() overwrite error = %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "second" {
		t.Errorf("file contents = %q, want %q", got, "second")
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("file perms = %v, want 0600", perm)
	}

	// No temp files are left behind by the atomic rename.
	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Errorf("directory has %d entries, want 1", len(entries))
	}
}

func TestWriteFileParentIsFile(t *testing.T) {
	root := t.TempDir()
	blocker := filepath.Join(root, "blocker")
	if err := os.WriteFile(blocker, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := WriteFile(filepath.Join(blocker, "out.txt"), []byte("x")); err == nil {
		t.Error("WriteFile() under a regular file succeeded, want error")
	}
}

func TestGlob(t *testing.T) {
	// doublestar.Glob is called with osFS, which opens paths relative to
	// the working directory.
	root := t.TempDir()
	t.Chdir(root)
	for _, p := range []string{"a.json", "sub/b.json", "sub/deeper/c.json", "sub/d.yaml"} {
		if err := WriteFile(p, []byte("{}")); err != nil {
			t.Fatal(err)
		}
	}

	got, err := Glob("**/*.json")
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]bool{"a.json": true, "sub/b.json": true, "sub/deeper/c.json": true}
	if len(got) != len(want) {
		t.Fatalf("Glob() = %v, want %v", got, want)
	}
	for _, g := range got {
		if !want[filepath.ToSlash(g)] {
			t.Errorf("Glob() returned unexpected %q", g)
		}
	}
}

func TestUnmarshalFile(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "rec.yaml")
	if err := os.WriteFile(path, []byte("id: abc\nexpected: 3\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var rec testRecord
	if err := UnmarshalFile(path, &rec); err != nil {
		t.Fatal(err)
	}
	if rec.ID != "abc" || rec.Expected != 3 {
		t.Errorf("UnmarshalFile() = %+v", rec)
	}

	if err := UnmarshalFile(filepath.Join(root, "missing.json"), &rec); err == nil {
		t.Error("UnmarshalFile(missing) error = nil")
	}

	txt := filepath.Join(root, "rec.txt")
	if err := os.WriteFile(txt, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := UnmarshalFile(txt, &rec); !errors.Is(err, ErrUnsupportedFormat) {
		t.Errorf("UnmarshalFile(.txt) error = %v, want ErrUnsupportedFormat", err)
	}
}

func TestUnmarshalPaths(t *testing.T) {
	fsys := fstest.MapFS{
		"data/a.json":  {Data: []byte(`{"id":"a"}`)},
		"data/b.yaml":  {Data: []byte("id: b\n")},
		"data/c.toml":  {Data: []byte(`id = "c"`)},
		"data/d.jsonc": {Data: []byte("{\"id\": \"d\", // comment\n}")},
		"data/skip.md": {Data: []byte("# ignored")},
		"single.yml":   {Data: []byte("id: s\n")},
	}

	got, err := UnmarshalPaths[testRecord](fsys, []string{"data", "single.yml"})
	if err != nil {
		t.Fatal(err)
	}
	ids := map[string]bool{}
	for _, r := range got {
		ids[r.ID] = true
	}
	for _, want := range []string{"a", "b", "c", "d", "s"} {
		if !ids[want] {
			t.Errorf("UnmarshalPaths() missing record %q; got %+v", want, got)
		}
	}
	if len(got) != 5 {
		t.Errorf("UnmarshalPaths() returned %d records, want 5", len(got))
	}
}

func TestUnmarshalPathsErrors(t *testing.T) {
	fsys := fstest.MapFS{
		"bad.json": {Data: []byte("{")},
	}
	if _, err := UnmarshalPaths[testRecord](fsys, []string{"missing"}); err == nil {
		t.Error("UnmarshalPaths(missing) error = nil")
	}
	if _, err := UnmarshalPaths[testRecord](fsys, []string{"bad.json"}); err == nil {
		t.Error("UnmarshalPaths(bad.json) error = nil")
	}
}
