package runx

import (
	"context"

	"go.jetify.com/pkg/runx/impl/registry"
	"go.jetify.com/pkg/runx/impl/types"
)

func (r *RunX) Install(ctx context.Context, pkgs ...string) ([]string, error) {
	refs := []types.PkgRef{}

	for _, pkg := range pkgs {
		ref, err := types.NewPkgRef(pkg)
		if err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}

	return r.install(ctx, refs...)
}

func (r *RunX) install(ctx context.Context, pkgs ...types.PkgRef) ([]string, error) {
	paths := []string{}
	for _, pkg := range pkgs {
		path, err := r.installOne(ctx, pkg)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	return paths, nil
}

func (r *RunX) installOne(ctx context.Context, ref types.PkgRef) (string, error) {
	reg, err := registry.NewLocalRegistry(ctx, r.GithubAPIToken)
	if err != nil {
		return "", err
	}

	pkgPath, err := reg.GetPackage(context.Background(), ref, types.CurrentPlatform())
	if err != nil {
		return "", err
	}
	return pkgPath, nil
}
