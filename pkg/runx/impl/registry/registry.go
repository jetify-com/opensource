package registry

import (
	"bytes"
	"context"
	"os"
	"path/filepath"

	"go.jetpack.io/pkg/runx/impl/download"
	"go.jetpack.io/pkg/runx/impl/fileutil"
	"go.jetpack.io/pkg/runx/impl/github"
	"go.jetpack.io/pkg/runx/impl/types"
)

var xdgInstallationSubdir = "runx/pkgs"

type Registry struct {
	rootPath fileutil.Path
	gh       *github.Client
}

func NewLocalRegistry(ctx context.Context, githubAPIToken string) (*Registry, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	rootDir := filepath.Join(cacheDir, xdgInstallationSubdir)
	rootPath := fileutil.Path(rootDir)

	if err := rootPath.EnsureDir(); err != nil {
		return nil, err
	}

	return &Registry{
		rootPath: rootPath,
		gh:       github.NewClient(ctx, githubAPIToken),
	}, nil
}

func (r *Registry) ListReleases(ctx context.Context, owner, repo string) ([]types.ReleaseMetadata, error) {
	path := r.rootPath.Subpath(owner, repo, "releases.json")

	return fetchCachedJSON(path.String(), func() ([]types.ReleaseMetadata, error) {
		return r.gh.ListReleases(ctx, owner, repo)
	})
}

func (r *Registry) GetReleaseMetadata(ctx context.Context, ref types.PkgRef) (types.ReleaseMetadata, error) {
	resolvedRef, err := r.ResolveVersion(ref)
	if err != nil {
		return types.ReleaseMetadata{}, err
	}

	path := r.rootPath.Subpath(resolvedRef.Owner, resolvedRef.Repo, resolvedRef.Version, "release.json")

	return fetchCachedJSON(path.String(), func() (types.ReleaseMetadata, error) {
		return r.gh.GetRelease(ctx, ref)
	})
}

func (r *Registry) GetArtifactMetadata(ctx context.Context, ref types.PkgRef, platform types.Platform) (types.ArtifactMetadata, error) {
	resolvedRef, err := r.ResolveVersion(ref)
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
	resolvedRef, err := r.ResolveVersion(ref)
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
	resolvedRef, err := r.ResolveVersion(ref)
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

	if isKnownArchive(filepath.Base(artifactPath)) {
		err = Extract(ctx, artifactPath, installPath.String())
	} else if isExecutableBinary(artifactPath) {
		err = createSymbolicLink(artifactPath, installPath.String(), resolvedRef.Repo)
	}
	if err != nil {
		return "", err
	}
	return installPath.String(), nil
}

func (r *Registry) ResolveVersion(ref types.PkgRef) (types.PkgRef, error) {
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

	found := false
	latestVersion := "latest"
	for _, release := range releases {
		if !release.Draft && !release.Prerelease {
			found = true
			latestVersion = release.TagName
			break
		}
	}

	// Return the first draft or prerelease if we couldn't find a stable release:
	if !found {
		latestVersion = releases[0].TagName
	}

	return types.PkgRef{
		Owner:   ref.Owner,
		Repo:    ref.Repo,
		Version: latestVersion,
	}, nil
}

// Best effort heuristic to determine if the artifact is an executable binary.
func isExecutableBinary(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	header := make([]byte, 4)
	_, err = file.Read(header)
	if err != nil {
		return false
	}

	switch {
	case bytes.HasPrefix(header, []byte("#!")): // Shebang
		return true
	case bytes.HasPrefix(header, []byte{0x7f, 0x45}): // ELF
		return true
	case bytes.Equal(header, []byte{0xfe, 0xed, 0xfa, 0xce}): // MachO32 BE
		return true
	case bytes.Equal(header, []byte{0xfe, 0xed, 0xfa, 0xcf}): // MachO64 BE
		return true
	case bytes.Equal(header, []byte{0xca, 0xfe, 0xba, 0xbe}): // Java class
		return true
	case bytes.Equal(header, []byte{0xCF, 0xFA, 0xED, 0xFE}): // Little-endian mac 64-bit
		return true
	case bytes.Equal(header, []byte{0xCE, 0xFA, 0xED, 0xFE}): // Little-endian mac 32-bit
		return true
	default:
		return false
	}
}
