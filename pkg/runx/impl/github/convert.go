package github

import (
	githubimpl "github.com/google/go-github/v53/github"
	"go.jetify.com/pkg/runx/impl/types"
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
		TagName:     release.GetTagName(),
		CreatedAt:   release.GetCreatedAt().Time,
		PublishedAt: release.GetPublishedAt().Time,
		Artifacts:   convertGithubAssets(release.Assets),
		Draft:       release.GetDraft(),
		Prerelease:  release.GetPrerelease(),
	}
}

func convertGithubAssets(assets []*githubimpl.ReleaseAsset) []types.ArtifactMetadata {
	result := []types.ArtifactMetadata{}
	for _, asset := range assets {
		if asset == nil {
			continue
		}
		result = append(result, types.ArtifactMetadata{
			URL:                asset.GetURL(),
			BrowserDownloadURL: asset.GetBrowserDownloadURL(),
			Name:               asset.GetName(),
			DownloadCount:      asset.GetDownloadCount(),
			CreatedAt:          asset.GetCreatedAt().Time,
			UpdatedAt:          asset.GetUpdatedAt().Time,
			ContentType:        asset.GetContentType(),
			Size:               asset.GetSize(),
		})
	}
	return result
}
