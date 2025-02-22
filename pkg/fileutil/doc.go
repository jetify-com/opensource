// Package fileutil provides utilities for working with files and paths in a
// safe and convenient way.
//
// The package offers two main sets of functionality:
//
//  1. File operations through both direct functions and a Path type:
//     - Checking file/directory existence and type
//     - Creating directories
//     - Writing files atomically
//     - Getting file information
//     - Globbing files with pattern matching
//
//  2. File unmarshaling utilities:
//     - Support for JSON, JSONC (JSON with comments), YAML, and TOML formats
//     - Batch processing of multiple files
//     - Type-safe unmarshaling into Go structs
//
// # Path Type
//
// The Path type provides a type-safe way to work with filesystem paths:
//
//	path := fileutil.Path("base/dir")
//	subpath := path.Subpath("nested", "path")
//	if subpath.IsDir() {
//	    // Handle directory...
//	}
//
// # File Operations
//
// Basic file operations are available both as methods on Path and as standalone
// functions:
//
//	// Using Path type
//	path := fileutil.Path("config")
//	if path.Exists() {
//	    info := path.FileInfo()
//	    // ...
//	}
//
//	// Using standalone functions
//	if fileutil.Exists("config") {
//	    info := fileutil.FileInfo("config")
//	    // ...
//	}
//
// # File Unmarshaling
//
// The package provides utilities for unmarshaling structured data from files:
//
//	var config MyConfig
//	err := fileutil.UnmarshalFile("config.yaml", &config)
//
//	// Or process multiple files
//	configs, err := fileutil.UnmarshalPaths[MyConfig](fs, []string{"configs"})
//
// All paths in this package are relative to the current working directory unless
// specified otherwise.
package fileutil
