// Package unmarshal provides functionality for unmarshaling data from various file formats
// into Go structures. It supports JSON, JSONC (JSON with comments), YAML, and TOML formats.
//
// The package offers several ways to unmarshal data:
//   - Reader: Unmarshal from an io.Reader with a specified format
//   - File: Unmarshal from a single file
//   - Paths: Unmarshal multiple files from given filesystem paths
//
// Example usage:
//
//	// Unmarshal from a reader
//	var config MyConfig
//	err := unmarshal.Reader(strings.NewReader(data), &config, ".json")
//
//	// Unmarshal from a file
//	var config MyConfig
//	err := unmarshal.File("config.yaml", &config)
//
//	// Unmarshal multiple files
//	configs, err := unmarshal.Paths[MyConfig](fsys, []string{"configs", "extra.yaml"})
//
// The package automatically detects the format based on file extensions:
//   - .json: Standard JSON
//   - .jsonc: JSON with comments
//   - .yml, .yaml: YAML
//   - .toml: TOML
//
// When using the Paths function with directories, it will recursively find and process
// all files with supported extensions within those directories.
package unmarshal
