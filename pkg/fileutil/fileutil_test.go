package fileutil

import (
	"io/fs"
	"testing"
	"testing/fstest"
)

func TestIsDir(t *testing.T) {
	fsys := fstest.MapFS{
		"dir/":   &fstest.MapFile{Mode: fs.ModeDir},
		"file":   &fstest.MapFile{Data: []byte("content")},
		"empty/": &fstest.MapFile{Mode: fs.ModeDir},
	}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"existing directory", "dir", true},
		{"regular file", "file", false},
		{"empty directory", "empty", true},
		{"non-existent path", "nonexistent", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isDir(fsys, tt.path); got != tt.expected {
				t.Errorf("isDir(%q) = %v, want %v", tt.path, got, tt.expected)
			}
		})
	}
}

func TestIsFile(t *testing.T) {
	fsys := fstest.MapFS{
		"dir/":  &fstest.MapFile{Mode: fs.ModeDir},
		"file":  &fstest.MapFile{Data: []byte("content")},
		"empty": &fstest.MapFile{Data: []byte{}},
	}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"directory", "dir", false},
		{"regular file", "file", true},
		{"empty file", "empty", true},
		{"non-existent path", "nonexistent", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isFile(fsys, tt.path); got != tt.expected {
				t.Errorf("isFile(%q) = %v, want %v", tt.path, got, tt.expected)
			}
		})
	}
}

func TestExists(t *testing.T) {
	fsys := fstest.MapFS{
		"dir/":  &fstest.MapFile{Mode: fs.ModeDir},
		"file":  &fstest.MapFile{Data: []byte("content")},
		"empty": &fstest.MapFile{Data: []byte{}},
	}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"existing directory", "dir", true},
		{"existing file", "file", true},
		{"empty file", "empty", true},
		{"non-existent path", "nonexistent", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := exists(fsys, tt.path); got != tt.expected {
				t.Errorf("exists(%q) = %v, want %v", tt.path, got, tt.expected)
			}
		})
	}
}

func TestFileInfo(t *testing.T) {
	fsys := fstest.MapFS{
		"dir/":  &fstest.MapFile{Mode: fs.ModeDir},
		"file":  &fstest.MapFile{Data: []byte("content")},
		"empty": &fstest.MapFile{Data: []byte{}},
	}

	tests := []struct {
		name          string
		path          string
		expectNil     bool
		expectIsDir   bool
		expectRegular bool
	}{
		{"directory", "dir", false, true, false},
		{"regular file", "file", false, false, true},
		{"empty file", "empty", false, false, true},
		{"non-existent path", "nonexistent", true, false, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			info := fileInfo(fsys, test.path)
			if test.expectNil {
				if info != nil {
					t.Errorf("fileInfo(%q) = %v, want nil", test.path, info)
				}
				return
			}
			if info == nil {
				t.Fatalf("fileInfo(%q) = nil, want non-nil", test.path)
			}
			if got := info.IsDir(); got != test.expectIsDir {
				t.Errorf("fileInfo(%q).IsDir() = %v, want %v", test.path, got, test.expectIsDir)
			}
			if got := info.Mode().IsRegular(); got != test.expectRegular {
				t.Errorf("fileInfo(%q).Mode().IsRegular() = %v, want %v", test.path, got, test.expectRegular)
			}
		})
	}
}
