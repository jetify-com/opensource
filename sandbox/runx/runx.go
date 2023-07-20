package runx

import (
	"go.jetpack.io/runx/impl"
)

func Install(pkgs ...string) error {
	return impl.Install(pkgs...)
}
