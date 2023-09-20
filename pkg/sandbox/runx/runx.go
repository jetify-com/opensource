package runx

import (
	"go.jetpack.io/pkg/sandbox/runx/impl"
)

// Install installs the given packages and returns the paths to the directories
// where they were installed.
func Install(pkgs ...string) ([]string, error) {
	return impl.Install(pkgs...)
}

func Run(args ...string) error {
	return impl.Run(args...)
}
