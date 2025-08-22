package fileutil

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/tailscale/hujson"
	"github.com/goccy/go-yaml"
)

var ErrUnsupportedFormat = errors.New("unsupported file format")

// supportedExtensions is the list of file extensions that can be processed
var supportedExtensions = []string{".json", ".jsonc", ".yml", ".yaml", ".toml"}

// UnmarshalFile reads and parses a single file into v.
func UnmarshalFile(path string, v any) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	return UnmarshalReader(f, v, filepath.Ext(path))
}

// UnmarshalPaths reads and parses files from the given paths in the filesystem
// into a slice of type T.
func UnmarshalPaths[T any](fsys fs.FS, paths []string) ([]T, error) {
	results := []T{}

	filePaths, err := FindFiles(fsys, paths, supportedExtensions)
	if err != nil {
		return nil, err
	}

	// Process all files
	for _, filePath := range filePaths {
		result, err := unmarshalPath[T](fsys, filePath)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

// unmarshalPath reads and parses a single file into type T.
func unmarshalPath[T any](fsys fs.FS, path string) (T, error) {
	var result T

	f, err := fsys.Open(path)
	if err != nil {
		return result, err
	}
	defer func() { _ = f.Close() }()

	if err := UnmarshalReader(f, &result, filepath.Ext(path)); err != nil {
		return result, err
	}

	return result, nil
}

// UnmarshalReader reads and parses data from an io.Reader into v based on the format specified.
// Supported formats: json, jsonc, yml, yaml, and toml.
// The format string can be a file extension (e.g. ".json") or a format name (e.g. "json").
func UnmarshalReader(r io.Reader, v any, format string) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	format = strings.TrimPrefix(strings.ToLower(format), ".")

	switch format {
	case "json":
		err = json.Unmarshal(data, v)
	case "jsonc":
		err = hujsonUnmarshal(data, v)
	case "yml", "yaml":
		err = yaml.Unmarshal(data, v)
	case "toml":
		err = toml.Unmarshal(data, v)
	default:
		return ErrUnsupportedFormat
	}

	return err
}

// hujsonUnmarshal parses JSONC data using hujson and unmarshals it into v.
// It follows the same pattern as json.Unmarshal.
func hujsonUnmarshal(data []byte, v any) error {
	ast, err := hujson.Parse(data)
	if err != nil {
		return err
	}
	ast.Standardize()
	return json.Unmarshal(ast.Pack(), v)
}

// FindFiles returns a list of files with the given extensions from the paths in the filesystem.
func FindFiles(fsys fs.FS, paths, exts []string) ([]string, error) {
	filePaths := []string{}

	for _, path := range paths {
		info, err := fs.Stat(fsys, path)
		if err != nil {
			return nil, err
		}

		if info.IsDir() {
			// Find files for each extension
			for _, ext := range exts {
				entries, err := fs.Glob(fsys, path+"/*"+ext)
				if err != nil {
					return nil, err
				}
				filePaths = append(filePaths, entries...)
			}
		} else if hasMatchingExtension(path, exts) {
			filePaths = append(filePaths, path)
		}
	}

	return filePaths, nil
}

// hasMatchingExtension checks if a file path has one of the given extensions
func hasMatchingExtension(path string, exts []string) bool {
	ext := filepath.Ext(path)
	for _, supportedExt := range exts {
		if ext == supportedExt {
			return true
		}
	}
	return false
}
