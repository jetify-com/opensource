package registry

import "go.jetpack.io/runx/impl/fileutil"

type Registry struct {
	rootPath fileutil.Path
}

func NewLocalRegistry(rootDir string) (*Registry, error) {
	rootPath := fileutil.Path(rootDir)
	err := rootPath.EnsureDir()
	if err != nil {
		return nil, err
	}

	return &Registry{
		rootPath: rootPath,
	}, nil
}
