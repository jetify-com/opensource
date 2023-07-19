package registry

import (
	"context"

	"github.com/google/go-github/v53/github"
	"go.jetpack.io/runx/impl/httpcacher"
)

type GithubRegistry struct {
	gh *github.Client
}

func NewGithubRegistry() *GithubRegistry {
	return &GithubRegistry{
		gh: github.NewClient(httpcacher.DefaultClient),
	}
}

func (r *GithubRegistry) ListReleases(ctx context.Context, owner, repo string) ([]ReleaseMetadata, error) {
	releases, resp, err := r.gh.Repositories.ListReleases(ctx, owner, repo, nil /* opts */)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 404 {
		return nil, ErrPackageNotFound
	}
	return convertReleases(releases), nil
}

func convertReleases(releases []*github.RepositoryRelease) []ReleaseMetadata {
	result := []ReleaseMetadata{}
	for _, release := range releases {
		result = append(result, ReleaseMetadata{
			Name:        release.GetName(),
			CreatedAt:   release.GetCreatedAt().Time,
			PublishedAt: release.GetPublishedAt().Time,
			Artifacts:   convertAssets(release.Assets),
		})
	}
	return result
}

func convertAssets(assets []*github.ReleaseAsset) []ArtifactMetadata {
	result := []ArtifactMetadata{}
	for _, asset := range assets {
		result = append(result, ArtifactMetadata{
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
