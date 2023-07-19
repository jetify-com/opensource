package registry

import (
	"context"

	"go.jetpack.io/runx/impl/fileutil"
	"go.jetpack.io/runx/impl/github"
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

func (r *Registry) ListReleases(ctx context.Context, owner, repo string) ([]ReleaseMetadata, error) {
	return nil, nil
}
