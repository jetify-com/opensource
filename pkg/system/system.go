package system

import (
	"io"
	"io/fs"
	"os"
)

// System provides an abstraction layer for filesystem and I/O operations,
// allowing for easier testing and dependency injection. It wraps the standard
// filesystem, stdin, stdout, and stderr functionality.
type System struct {
	fs        fs.StatFS
	inReader  io.Reader
	outWriter io.Writer
	errWriter io.Writer

	// TODO: Add logger and environment variables support.
}

// GetFS returns the filesystem interface used by the System.
// If no filesystem is set or System is nil, it returns the OS filesystem.
func (s *System) GetFS() fs.FS {
	if s == nil || s.fs == nil {
		return os.DirFS(".")
	}
	return s.fs
}

// SetFS sets the filesystem interface to be used by the System.
func (s *System) SetFS(fs fs.StatFS) {
	s.fs = fs
}

// GetStdin returns the current input reader.
// If no reader is set or System is nil, it returns os.Stdin.
func (s *System) GetStdin() io.Reader {
	if s == nil || s.inReader == nil {
		return os.Stdin
	}
	return s.inReader
}

// SetStdin sets the input reader to be used by the System.
func (s *System) SetStdin(inReader io.Reader) {
	s.inReader = inReader
}

// GetStdout returns the current output writer.
// If no writer is set or System is nil, it returns os.Stdout.
func (s *System) GetStdout() io.Writer {
	if s == nil || s.outWriter == nil {
		return os.Stdout
	}
	return s.outWriter
}

// SetStdout sets the output writer to be used by the System.
func (s *System) SetStdout(outWriter io.Writer) {
	s.outWriter = outWriter
}

// GetStderr returns the current error writer.
// If no writer is set or System is nil, it returns os.Stderr.
func (s *System) GetStderr() io.Writer {
	if s == nil || s.errWriter == nil {
		return os.Stderr
	}
	return s.errWriter
}

// SetStderr sets the error writer to be used by the System.
func (s *System) SetStderr(errWriter io.Writer) {
	s.errWriter = errWriter
}
