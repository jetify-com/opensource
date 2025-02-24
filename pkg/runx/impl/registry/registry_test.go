package registry

import (
	"errors"
	"os"
	"slices"
	"testing"

	"go.jetify.com/pkg/runx/impl/types"
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
			file, err := os.CreateTemp(t.TempDir(), "testfile")
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
		if got != test.want {
			t.Errorf("isKnownArchive(%s) = %v, want %v", test.name, got, test.want)
		}
	}
}

func TestIsArtifactForPlatform(t *testing.T) {
	tests := []struct {
		name     string
		platform types.Platform
		want     bool
	}{
		{"linux-x86_64", types.NewPlatform("linux", "x86_64"), true},
		{"linux_x86_64_1.1.0_SNAPSHOT-abcde12345", types.NewPlatform("linux", "x86_64"), true},
		{"linux_x86_64-1.2.3-rc2.tar.gz", types.NewPlatform("linux", "amd64"), true},
		{"no-os-no-arch", types.NewPlatform("", ""), false},
		{"no-os-no-arch", types.NewPlatform("linux", "amd64"), false},
		{"no_os_arm64", types.NewPlatform("linux", "arm64"), false},
		{"linux_no-arch", types.NewPlatform("linux", "arm64"), false},
		{"macos-universal-1.2.3", types.NewPlatform("darwin", "arm64"), true},
		{"mac_386-1.2.3-snapshot_abcdef123456", types.NewPlatform("darwin", "386"), true},
	}

	for _, test := range tests {
		got := isArtifactForPlatform(test.name, test.platform)
		if got != test.want {
			t.Errorf("isArtifactForPlatform(%s) = %v, want %v", test.name, got, test.want)
		}
	}
}

func TestFindArtifactForPlatform(t *testing.T) {
	tests := []struct {
		artifacts   []types.ArtifactMetadata
		platform    types.Platform
		want        bool
		errTypeWant error
	}{
		{[]types.ArtifactMetadata{{Name: "mac-amd64.tar.gz"}}, types.NewPlatform("darwin", "amd64"), true, nil},
		{[]types.ArtifactMetadata{{Name: "linux-amd64.deb"}}, types.NewPlatform("linux", "amd64"), true, nil},
		{[]types.ArtifactMetadata{{Name: "linux-amd64"}}, types.NewPlatform("linux", "amd64"), true, nil},
		{[]types.ArtifactMetadata{{Name: "darwin_arm64.tar"}}, types.NewPlatform("windows", "amd64"), false, types.ErrPlatformNotSupported},
		{[]types.ArtifactMetadata{{Name: "darwin_arm64"}}, types.NewPlatform("windows", "amd64"), false, types.ErrPlatformNotSupported},
		{[]types.ArtifactMetadata{{Name: "mac-amd64.tar.gz"}}, types.NewPlatform("", ""), false, types.ErrPlatformNotSupported},
	}

	for _, test := range tests {
		got, err := findArtifactForPlatform(test.artifacts, test.platform)
		if slices.Contains(test.artifacts, got) != test.want {
			t.Errorf("findArtifactForPlatform() didn't return %v", got)
		}
		if !errors.Is(err, test.errTypeWant) {
			t.Errorf("findArtifactForPlatform() returned %v error instead of %v", err, test.errTypeWant)
		}
	}
}
