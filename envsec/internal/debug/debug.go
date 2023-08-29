// Copyright 2023 Jetpack Technologies Inc and contributors. All rights reserved.
// Use of this source code is governed by the license in the LICENSE file.

package debug

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/pkg/errors"
)

var enabled bool

func IsEnvsecDebugEnabled() bool {
	enabled, _ := strconv.ParseBool(os.Getenv("ENVSEC_DEBUG"))
	return enabled
}
func init() {
	enabled = IsEnvsecDebugEnabled()
}

func IsEnabled() bool { return enabled }

func Enable() {
	enabled = true
	log.SetPrefix("[DEBUG] ")
	log.SetFlags(log.Llongfile | log.Ldate | log.Ltime)
	_ = log.Output(2, "Debug mode enabled.")
}

func SetOutput(w io.Writer) {
	log.SetOutput(w)
}

func Log(format string, v ...any) {
	if !enabled {
		return
	}
	_ = log.Output(2, fmt.Sprintf(format, v...))
}

func EarliestStackTrace(err error) error {
	type pkgErrorsStackTracer interface{ StackTrace() errors.StackTrace }
	type redactStackTracer interface{ StackTrace() []runtime.Frame }

	var stErr error
	for err != nil {
		//nolint:errorlint
		switch err.(type) {
		case redactStackTracer, pkgErrorsStackTracer:
			stErr = err
		}
		err = errors.Unwrap(err)
	}
	return stErr
}
