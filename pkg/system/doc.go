// Package system provides an abstraction layer for filesystem and I/O operations,
// making it easier to write testable code by allowing dependency injection of
// system-level interfaces.
//
// The primary type is System, which wraps standard filesystem operations and
// I/O streams (stdin, stdout, stderr). This abstraction allows for easy mocking
// in tests and provides a clean interface for system operations.
//
// Example usage:
//
//	sys := &system.System{}
//
//	// Use default OS filesystem
//	fs := sys.GetFS()
//
//	// Override with custom filesystem for testing
//	sys.SetFS(fstest.MapFS{})
//
//	// Use custom writers/readers
//	var buf bytes.Buffer
//	sys.SetStdout(&buf)
//	sys.SetStderr(&buf)
//
// The package is designed to be used as a dependency in other packages that need
// to interact with the filesystem or standard I/O streams, making those packages
// more testable and maintainable.
package system
