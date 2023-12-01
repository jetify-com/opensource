// Copyright 2022 Jetpack Technologies Inc and contributors. All rights reserved.
// Use of this source code is governed by the license in the LICENSE file.

package envsec

import (
	"os"
	"strconv"
	"strings"
)

// isDev determines whether we are running in dev mode, by default.
// Specific settings may still be overridable by specific env-vars.
var isDev bool = false

func nameFromPath(path string) string {
	subpaths := strings.Split(path, "/")
	if len(subpaths) == 0 {
		return ""
	}
	return subpaths[len(subpaths)-1]
}

func SetDevMode(dev bool) {
	isDev = dev
}

func IsDevMode() bool {
	devEnv := os.Getenv("ENVSEC_DEV")
	if devEnv != "" {
		isDev, err := strconv.ParseBool(devEnv)
		return err == nil && isDev
	}
	return isDev
}
