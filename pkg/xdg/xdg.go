// Copyright 2025 Jetify Inc. and contributors. All rights reserved.
// Use of this source code is governed by the license in the LICENSE file.

// Package xdg provides access to XDG Base Directory Specification paths.
package xdg

import (
	"os"
	"path/filepath"
)

func DataSubpath(subpath string) string {
	return filepath.Join(dataDir(), subpath)
}

func ConfigSubpath(subpath string) string {
	return filepath.Join(configDir(), subpath)
}

func CacheSubpath(subpath string) string {
	return filepath.Join(cacheDir(), subpath)
}

func StateSubpath(subpath string) string {
	return filepath.Join(stateDir(), subpath)
}

func dataDir() string   { return resolveDir("XDG_DATA_HOME", ".local/share") }
func configDir() string { return resolveDir("XDG_CONFIG_HOME", ".config") }
func cacheDir() string  { return resolveDir("XDG_CACHE_HOME", ".cache") }
func stateDir() string  { return resolveDir("XDG_STATE_HOME", ".local/state") }

func resolveDir(envvar, defaultPath string) string {
	dir := os.Getenv(envvar)
	if dir != "" {
		return dir
	}

	home, err := os.UserHomeDir()
	if err != nil {
		home = "/tmp"
	}

	return filepath.Join(home, defaultPath)
}
