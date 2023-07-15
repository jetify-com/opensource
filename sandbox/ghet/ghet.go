package ghet

import (
	"go.jetpack.io/ghet/impl"
)

func Download(pkgs ...string) error {
	return impl.Download(pkgs...)
}
