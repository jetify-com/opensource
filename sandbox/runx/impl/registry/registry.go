package registry

import (
	"context"

	"go.jetpack.io/runx/impl/fileutil"
	"go.jetpack.io/runx/impl/gh"
)

type Registry struct {
	rootPath fileutil.Path
	gh       *gh.Client
}

func NewLocalRegistry(rootDir string) (*Registry, error) {
	rootPath := fileutil.Path(rootDir)
	err := rootPath.EnsureDir()
	if err != nil {
		return nil, err
	}

	return &Registry{
		rootPath: rootPath,
		gh:       gh.NewClient(),
	}, nil
}

func (r *Registry) ListReleases(ctx context.Context, owner, repo string) ([]ReleaseMetadata, error) {
	return nil, nil
}
