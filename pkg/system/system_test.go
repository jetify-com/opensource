package system

import (
	"bytes"
	"io"
	"io/fs"
	"os"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func TestSystem_GetFS(t *testing.T) {
	tests := []struct {
		name     string
		system   *System
		wantType fs.FS
	}{
		{
			name:     "nil system returns os filesystem",
			system:   nil,
			wantType: os.DirFS("."),
		},
		{
			name:     "system with nil fs returns os filesystem",
			system:   &System{},
			wantType: os.DirFS("."),
		},
		{
			name:     "system with custom fs returns that fs",
			system:   &System{fs: fstest.MapFS{}},
			wantType: fstest.MapFS{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.system.GetFS()
			assert.IsType(t, tt.wantType, got)
		})
	}
}

func TestSystem_SetFS(t *testing.T) {
	s := &System{}
	mockFS := fstest.MapFS{}

	s.SetFS(mockFS)

	assert.Equal(t, mockFS, s.GetFS())
}

func TestSystem_GetStdin(t *testing.T) {
	tests := []struct {
		name     string
		system   *System
		wantType io.Reader
	}{
		{
			name:     "nil system returns os.Stdin",
			system:   nil,
			wantType: os.Stdin,
		},
		{
			name:     "system with nil reader returns os.Stdin",
			system:   &System{},
			wantType: os.Stdin,
		},
		{
			name:     "system with custom reader returns that reader",
			system:   &System{inReader: bytes.NewBuffer(nil)},
			wantType: &bytes.Buffer{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.system.GetStdin()
			assert.IsType(t, tt.wantType, got)
		})
	}
}

func TestSystem_SetStdin(t *testing.T) {
	s := &System{}
	mockReader := bytes.NewBuffer(nil)

	s.SetStdin(mockReader)

	assert.Equal(t, mockReader, s.GetStdin())
}

func TestSystem_GetStdout(t *testing.T) {
	tests := []struct {
		name     string
		system   *System
		wantType io.Writer
	}{
		{
			name:     "nil system returns os.Stdout",
			system:   nil,
			wantType: os.Stdout,
		},
		{
			name:     "system with nil writer returns os.Stdout",
			system:   &System{},
			wantType: os.Stdout,
		},
		{
			name:     "system with custom writer returns that writer",
			system:   &System{outWriter: bytes.NewBuffer(nil)},
			wantType: &bytes.Buffer{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.system.GetStdout()
			assert.IsType(t, tt.wantType, got)
		})
	}
}

func TestSystem_SetStdout(t *testing.T) {
	s := &System{}
	mockWriter := bytes.NewBuffer(nil)

	s.SetStdout(mockWriter)

	assert.Equal(t, mockWriter, s.GetStdout())
}

func TestSystem_GetStderr(t *testing.T) {
	tests := []struct {
		name     string
		system   *System
		wantType io.Writer
	}{
		{
			name:     "nil system returns os.Stderr",
			system:   nil,
			wantType: os.Stderr,
		},
		{
			name:     "system with nil writer returns os.Stderr",
			system:   &System{},
			wantType: os.Stderr,
		},
		{
			name:     "system with custom writer returns that writer",
			system:   &System{errWriter: bytes.NewBuffer(nil)},
			wantType: &bytes.Buffer{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.system.GetStderr()
			assert.IsType(t, tt.wantType, got)
		})
	}
}

func TestSystem_SetStderr(t *testing.T) {
	s := &System{}
	mockWriter := bytes.NewBuffer(nil)

	s.SetStderr(mockWriter)

	assert.Equal(t, mockWriter, s.GetStderr())
}
