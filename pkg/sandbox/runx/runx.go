package runx

import (
	"context"

	"go.jetpack.io/pkg/sandbox/runx/impl"
	"go.jetpack.io/pkg/sandbox/runx/impl/runopt"
)

type RunX interface {
	// Install installs the given packages and returns the paths to the directories
	// where they were installed.
	Install(ctx context.Context, pkgs ...string) ([]string, error)
	Run(ctx context.Context, args ...string) error
}

func New(opts runopt.Opts) RunX {
	return impl.New(opts)
}
