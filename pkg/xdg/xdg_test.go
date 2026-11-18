package xdg

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDataSubpath(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		subpath  string
		want     string
	}{
		{
			name:     "with environment variable",
			envValue: "/custom/data",
			subpath:  "test",
			want:     "/custom/data/test",
		},
		{
			name:     "without environment variable",
			envValue: "",
			subpath:  "test",
			want:     filepath.Join(os.Getenv("HOME"), ".local/share/test"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("XDG_DATA_HOME", tt.envValue)

			got := DataSubpath(tt.subpath)
			if got != tt.want {
				t.Errorf("DataSubpath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigSubpath(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		subpath  string
		want     string
	}{
		{
			name:     "with environment variable",
			envValue: "/custom/config",
			subpath:  "test",
			want:     "/custom/config/test",
		},
		{
			name:     "without environment variable",
			envValue: "",
			subpath:  "test",
			want:     filepath.Join(os.Getenv("HOME"), ".config/test"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("XDG_CONFIG_HOME", tt.envValue)

			got := ConfigSubpath(tt.subpath)
			if got != tt.want {
				t.Errorf("ConfigSubpath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCacheSubpath(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		subpath  string
		want     string
	}{
		{
			name:     "with environment variable",
			envValue: "/custom/cache",
			subpath:  "test",
			want:     "/custom/cache/test",
		},
		{
			name:     "without environment variable",
			envValue: "",
			subpath:  "test",
			want:     filepath.Join(os.Getenv("HOME"), ".cache/test"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("XDG_CACHE_HOME", tt.envValue)

			got := CacheSubpath(tt.subpath)
			if got != tt.want {
				t.Errorf("CacheSubpath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateSubpath(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		subpath  string
		want     string
	}{
		{
			name:     "with environment variable",
			envValue: "/custom/state",
			subpath:  "test",
			want:     "/custom/state/test",
		},
		{
			name:     "without environment variable",
			envValue: "",
			subpath:  "test",
			want:     filepath.Join(os.Getenv("HOME"), ".local/state/test"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("XDG_STATE_HOME", tt.envValue)

			got := StateSubpath(tt.subpath)
			if got != tt.want {
				t.Errorf("StateSubpath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolveDirWithNoHomeDir(t *testing.T) {
	// Save original HOME environment variable
	originalHome := os.Getenv("HOME")
	defer t.Setenv("HOME", originalHome)

	// Unset HOME to test fallback behavior
	t.Setenv("HOME", "")

	// Test that we get the expected fallback path
	got := dataDir()
	want := filepath.Join("/tmp", ".local/share")
	if got != want {
		t.Errorf("dataDir() with no HOME = %v, want %v", got, want)
	}
}
