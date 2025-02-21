package unmarshal

import (
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

type testRecord struct {
	ID       string         `json:"id,omitempty" yaml:"id" toml:"id"`
	Input    string         `json:"input,omitempty" yaml:"input" toml:"input"`
	Expected int            `json:"expected,omitempty" yaml:"expected" toml:"expected"`
	Tags     []string       `json:"tags,omitempty" yaml:"tags" toml:"tags"`
	Metadata map[string]any `json:"metadata,omitempty" yaml:"metadata" toml:"metadata"`
}

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		format  string
		want    *testRecord
		wantErr bool
	}{
		{
			name: "valid json",
			input: `{
				"id": "test1",
				"input": "hello",
				"expected": 42,
				"tags": ["test", "example"]
			}`,
			format: ".json",
			want: &testRecord{
				ID:       "test1",
				Input:    "hello",
				Expected: 42,
				Tags:     []string{"test", "example"},
			},
		},
		{
			name: "valid jsonc with comments",
			input: `{
				// This is a comment
				"id": "test1",
				"input": "hello", /* inline comment */
				"expected": 42,
				"tags": ["test", "example"]
			}`,
			format: ".jsonc",
			want: &testRecord{
				ID:       "test1",
				Input:    "hello",
				Expected: 42,
				Tags:     []string{"test", "example"},
			},
		},
		{
			name: "valid yaml",
			input: `
id: test1
input: hello
expected: 42
tags: 
  - test
  - example`,
			format: ".yaml",
			want: &testRecord{
				ID:       "test1",
				Input:    "hello",
				Expected: 42,
				Tags:     []string{"test", "example"},
			},
		},
		{
			name:    "invalid format",
			input:   "{}",
			format:  ".invalid",
			wantErr: true,
		},
		{
			name:    "invalid json",
			input:   "{invalid json}",
			format:  ".json",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			var got testRecord
			err := Reader(r, &got, tt.format)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, &got)
		})
	}
}

func TestFindFiles(t *testing.T) {
	// Create test filesystem
	fsys := fstest.MapFS{
		"config1.json":    {Data: []byte(`{}`)},
		"config2.yaml":    {Data: []byte(`{}`)},
		"dir/nested.json": {Data: []byte(`{}`)},
		"dir/other.yaml":  {Data: []byte(`{}`)},
		"empty":           {Data: []byte(`{}`), Mode: fs.ModeDir},
		"unsupported.txt": {Data: []byte(`{}`)},
	}

	tests := []struct {
		name    string
		paths   []string
		exts    []string
		want    []string
		wantErr bool
		errPath string
	}{
		{
			name:  "single file",
			paths: []string{"config1.json"},
			exts:  []string{".json"},
			want:  []string{"config1.json"},
		},
		{
			name:  "multiple extensions",
			paths: []string{"config1.json", "config2.yaml"},
			exts:  []string{".json", ".yaml"},
			want:  []string{"config1.json", "config2.yaml"},
		},
		{
			name:  "directory with multiple files",
			paths: []string{"dir"},
			exts:  []string{".json", ".yaml"},
			want:  []string{"dir/nested.json", "dir/other.yaml"},
		},
		{
			name:  "mixed paths",
			paths: []string{"config1.json", "dir"},
			exts:  []string{".json"},
			want:  []string{"config1.json", "dir/nested.json"},
		},
		{
			name:  "empty directory",
			paths: []string{"empty"},
			exts:  []string{".json"},
			want:  []string{},
		},
		{
			name:  "unsupported extension",
			paths: []string{"unsupported.txt"},
			exts:  []string{".json"},
			want:  []string{},
		},
		{
			name:    "nonexistent file",
			paths:   []string{"nonexistent.json"},
			exts:    []string{".json"},
			wantErr: true,
			errPath: "nonexistent.json",
		},
		{
			name:    "nonexistent directory",
			paths:   []string{"nonexistent"},
			exts:    []string{".json"},
			wantErr: true,
			errPath: "nonexistent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findFiles(fsys, tt.paths, tt.exts)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errPath != "" {
					assert.Contains(t, err.Error(), tt.errPath)
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
