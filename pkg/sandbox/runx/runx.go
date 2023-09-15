package runx

import (
	"go.jetpack.io/pkg/sandbox/runx/impl"
)

func Install(pkgs ...string) error {
	_, err := impl.Install(pkgs...)
	return err
}

func Run(args ...string) error {
	return impl.Run(args...)
}
