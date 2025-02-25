package unmarshal

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
	"gopkg.in/yaml.v3"
)

var ErrUnsupportedFormat = errors.New("unsupported file format")

// supportedExtensions is the list of file extensions that can be processed
var supportedExtensions = []string{".json", ".jsonc", ".yml", ".yaml", ".toml"}

// Reader reads and parses data from an io.Reader into v based on the format specified.
// Supported formats: json, jsonc, yml, yaml, and toml.
// The format string can be a file extension (e.g. ".json") or a format name (e.g. "json").
func Reader(r io.Reader, v any, format string) error {
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

// File reads and parses a single file into v.
func File(path string, v any) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return Reader(f, v, filepath.Ext(path))
}

// Paths reads and parses files from the given paths in the filesystem into a slice of type T.
func Paths[T any](fsys fs.FS, paths []string) ([]T, error) {
	results := []T{}

	filePaths, err := findFiles(fsys, paths, supportedExtensions)
	if err != nil {
		return nil, err
	}

	// Process all files
	for _, filePath := range filePaths {
		f, err := fsys.Open(filePath)
		if err != nil {
			return nil, err
		}

		var result T
		if err := Reader(f, &result, filepath.Ext(filePath)); err != nil {
			f.Close()
			return nil, err
		}
		f.Close()

		results = append(results, result)
	}

	return results, nil
}

// findFiles returns a list of files with the given extensions from the paths in the filesystem.
func findFiles(fsys fs.FS, paths []string, exts []string) ([]string, error) {
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
