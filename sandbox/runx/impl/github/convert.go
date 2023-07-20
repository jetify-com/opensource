package github

import (
	githubimpl "github.com/google/go-github/v53/github"
	"go.jetpack.io/runx/impl/types"
)

func convertGithubReleases(releases []*githubimpl.RepositoryRelease) []types.ReleaseMetadata {
	result := []types.ReleaseMetadata{}
	for _, release := range releases {
		if release == nil {
			continue
		}
		result = append(result, convertGithubRelease(*release))
	}
	return result
}

func convertGithubRelease(release githubimpl.RepositoryRelease) types.ReleaseMetadata {
	return types.ReleaseMetadata{
		Name:        release.GetName(),
		CreatedAt:   release.GetCreatedAt().Time,
		PublishedAt: release.GetPublishedAt().Time,
		Artifacts:   convertGithubAssets(release.Assets),
	}
}

func convertGithubAssets(assets []*githubimpl.ReleaseAsset) []types.ArtifactMetadata {
	result := []types.ArtifactMetadata{}
	for _, asset := range assets {
		if asset == nil {
			continue
		}
		result = append(result, types.ArtifactMetadata{
			DownloadURL:   asset.GetBrowserDownloadURL(),
			Name:          asset.GetName(),
			DownloadCount: asset.GetDownloadCount(),
			CreatedAt:     asset.GetCreatedAt().Time,
			UpdatedAt:     asset.GetUpdatedAt().Time,
			ContentType:   asset.GetContentType(),
			Size:          asset.GetSize(),
		})
	}
	return result
}
