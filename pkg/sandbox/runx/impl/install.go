package impl

import (
	"context"
	"os"
	"path/filepath"

	"go.jetpack.io/pkg/sandbox/runx/impl/registry"
	"go.jetpack.io/pkg/sandbox/runx/impl/types"
)

var xdgInstallationSubdir = "jetpack.io/pkgs"

func Install(pkgs ...string) ([]string, error) {
	refs := []types.PkgRef{}

	for _, pkg := range pkgs {
		ref, err := types.NewPkgRef(pkg)
		if err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}

	return install(refs...)
}

func install(pkgs ...types.PkgRef) ([]string, error) {
	paths := []string{}
	for _, pkg := range pkgs {
		path, err := installOne(pkg)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	return paths, nil
}

func installOne(ref types.PkgRef) (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = "~/.cache"
	}
	rootDir := filepath.Join(cacheDir, xdgInstallationSubdir)
	reg, err := registry.NewLocalRegistry(rootDir)
	if err != nil {
		return "", err
	}

	pkgPath, err := reg.GetPackage(context.Background(), ref, types.CurrentPlatform())
	if err != nil {
		return "", err
	}
	return pkgPath, nil
}
