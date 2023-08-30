// Copyright 2022 Jetpack Technologies Inc and contributors. All rights reserved.
// Use of this source code is governed by the license in the LICENSE file.

package envsec

import (
	"path"
	"strings"

	"go.jetpack.io/envsec/internal/debug"
)

const PATH_PREFIX = "/jetpack-data/env"

func varPath(envId EnvId, varName string) string {
	path := path.Join(projectPath(envId), envId.EnvName, varName)
	debug.Log("varpath: %s \n\n", path)
	return path
}

func projectPath(envId EnvId) string {
	return path.Join(PATH_PREFIX, envId.ProjectId)
}

func nameFromPath(path string) string {
	subpaths := strings.Split(path, "/")
	if len(subpaths) == 0 {
		return ""
	}
	return subpaths[len(subpaths)-1]
}
