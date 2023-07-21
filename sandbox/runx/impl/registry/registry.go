package registry

import (
	"context"

	"go.jetpack.io/runx/impl/download"
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

func (r *Registry) GetReleaseMetadata(ctx context.Context, ref types.PkgRef) (types.ReleaseMetadata, error) {
	resolvedRef, err := r.resolveVersion(ref)
	if err != nil {
		return types.ReleaseMetadata{}, err
	}

	path := r.rootPath.Subpath(resolvedRef.Owner, resolvedRef.Repo, resolvedRef.Version, "release.json")

	return fetchCachedJSON(path.String(), func() (types.ReleaseMetadata, error) {
		return r.gh.GetRelease(ctx, ref)
	})
}

func (r *Registry) GetArtifactMetadata(ctx context.Context, ref types.PkgRef, platform types.Platform) (types.ArtifactMetadata, error) {
	resolvedRef, err := r.resolveVersion(ref)
	if err != nil {
		return types.ArtifactMetadata{}, err
	}

	release, err := r.GetReleaseMetadata(ctx, resolvedRef)
	if err != nil {
		return types.ArtifactMetadata{}, err
	}

	artifact := findArtifactForPlatform(release.Artifacts, platform)
	if artifact == nil {
		return types.ArtifactMetadata{}, types.ErrPlatformNotSupported
	}
	return *artifact, nil
}

func (r *Registry) GetArtifact(ctx context.Context, ref types.PkgRef, platform types.Platform) (string, error) {
	resolvedRef, err := r.resolveVersion(ref)
	if err != nil {
		return "", err
	}

	metadata, err := r.GetArtifactMetadata(ctx, ref, platform)
	if err != nil {
		return "", err
	}

	path := r.rootPath.Subpath(resolvedRef.Owner, resolvedRef.Repo, resolvedRef.Version, metadata.Name)
	err = download.DownloadOnce(metadata.DownloadURL, path.String())
	if err != nil {
		return "", err
	}
	return path.String(), nil
}

func (r *Registry) GetPackage(ctx context.Context, ref types.PkgRef, platform types.Platform) (string, error) {
	resolvedRef, err := r.resolveVersion(ref)
	if err != nil {
		return "", err
	}
	installPath := r.rootPath.Subpath(
		resolvedRef.Owner,
		resolvedRef.Repo,
		resolvedRef.Version,
		platform.OS(),
		platform.Arch(),
	)
	// If the installation path already exists, we assume the artifact is already installed.
	// We'll want some way of validating via a checksum or digest in the future.
	if installPath.IsDir() {
		return installPath.String(), nil
	}

	artifactPath, err := r.GetArtifact(ctx, ref, platform)
	if err != nil {
		return "", err
	}

	err = Extract(ctx, artifactPath, installPath.String())
	if err != nil {
		return "", err
	}
	return installPath.String(), nil
}

func (r *Registry) resolveVersion(ref types.PkgRef) (types.PkgRef, error) {
	if ref.Version != "" && ref.Version != "latest" {
		return ref, nil
	}

	releases, err := r.ListReleases(context.Background(), ref.Owner, ref.Repo)
	if err != nil {
		return types.PkgRef{}, err
	}

	found := false
	latestVersion := "latest"
	for _, release := range releases {
		if !release.Draft && !release.Prerelease {
			found = true
			latestVersion = release.TagName
			break
		}
	}

	if !found {
		return types.PkgRef{}, types.ErrReleaseNotFound
	}

	return types.PkgRef{
		Owner:   ref.Owner,
		Repo:    ref.Repo,
		Version: latestVersion,
	}, nil
}
