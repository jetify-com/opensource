package registry

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsBinary(t *testing.T) {
	tests := []struct {
		name   string
		header []byte
		want   bool
	}{
		{"Shebang", []byte("#!/bin/bash\n"), true},
		{"ELF", []byte{0x7f, 0x45, 0x4c, 0x46}, true},
		{"MachO32 BE", []byte{0xfe, 0xed, 0xfa, 0xce}, true},
		{"MachO64 BE", []byte{0xfe, 0xed, 0xfa, 0xcf}, true},
		{"Java Class", []byte{0xca, 0xfe, 0xba, 0xbe}, true},
		{"MachO64 LE", []byte{0xcf, 0xfa, 0xed, 0xfe}, true},
		{"MachO32 LE", []byte{0xce, 0xfa, 0xed, 0xfe}, true},
		{"Unknown", []byte{0xaa, 0xbb, 0xcc, 0xdd}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			file, err := os.CreateTemp("", "testfile")
			if err != nil {
				t.Fatalf("Could not create temp file: %v", err)
			}
			defer os.Remove(file.Name())

			_, err = file.Write(test.header)
			if err != nil {
				t.Fatalf("Could not write to temp file: %v", err)
			}
			file.Close()

			got := isExecutableBinary(file.Name())
			if got != test.want {
				t.Errorf("isBinary() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestIsKnownArchive(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"archive.tar", true},
		{"archive.tar.gz", true},
		{"archive.deb", false},
		{"archive", false},
	}
	for _, test := range tests {
		got := isKnownArchive(test.name)
		t.Logf("%s", filepath.Ext(test.name))
		if got != test.want {
			t.Errorf("isKnownArchive(%s) = %v, want %v", test.name, got, test.want)
		}
	}
}
