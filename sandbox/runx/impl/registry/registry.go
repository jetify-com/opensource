package registry

import (
	"context"

	"go.jetpack.io/runx/impl/fileutil"
	"go.jetpack.io/runx/impl/github"
	"go.jetpack.io/runx/impl/types"
)

type Registry struct {
	rootPath fileutil.Path
	gh       *github.Client
}

func NewLocalRegistry(rootDir string) (*Registry, error) {
	rootPath := fileutil.Path(rootDir)
	err := rootPath.EnsureDir()
	if err != nil {
		return nil, err
	}

	return &Registry{
		rootPath: rootPath,
		gh:       github.NewClient(),
	}, nil
}

func (r *Registry) ListReleases(ctx context.Context, owner, repo string) ([]types.ReleaseMetadata, error) {
	path := r.rootPath.Subpath(owner, repo, "releases.json")

	return fetchCachedJSON(path.String(), func() ([]types.ReleaseMetadata, error) {
		return r.gh.ListReleases(ctx, owner, repo)
	})
}

func (r *Registry) GetRelease(ctx context.Context, ref types.PkgRef) (types.ReleaseMetadata, error) {
	resolvedRef, err := r.resolveVersion(ref)
	if err != nil {
		return types.ReleaseMetadata{}, err
	}

	path := r.rootPath.Subpath(resolvedRef.Owner, resolvedRef.Repo, resolvedRef.Version, "release.json")

	return fetchCachedJSON(path.String(), func() (types.ReleaseMetadata, error) {
		return r.gh.GetRelease(ctx, ref)
	})
}

func (r *Registry) resolveVersion(ref types.PkgRef) (types.PkgRef, error) {
	if ref.Version != "" && ref.Version != "latest" {
		return ref, nil
	}

	releases, err := r.ListReleases(context.Background(), ref.Owner, ref.Repo)
	if err != nil {
		return types.PkgRef{}, err
	}

	if len(releases) == 0 {
		return types.PkgRef{}, types.ErrReleaseNotFound
	}

	latestVersion := releases[0].Name
	return types.PkgRef{
		Owner:   ref.Owner,
		Repo:    ref.Repo,
		Version: latestVersion,
	}, nil
}
