package github

import (
	"context"

	githubimpl "github.com/google/go-github/v53/github"
	"go.jetpack.io/pkg/sandbox/runx/impl/httpcacher"
	"go.jetpack.io/pkg/sandbox/runx/impl/types"
)

type Client struct {
	gh *githubimpl.Client
}

func NewClient() *Client {
	return &Client{
		gh: githubimpl.NewClient(httpcacher.DefaultClient),
	}
}

func (c *Client) ListReleases(ctx context.Context, owner, repo string) ([]types.ReleaseMetadata, error) {
	releases, resp, err := c.gh.Repositories.ListReleases(ctx, owner, repo, &githubimpl.ListOptions{
		PerPage: 100, // Max allowed by GitHub API
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 404 {
		return nil, types.ErrPackageNotFound
	}
	return convertGithubReleases(releases), nil
}

func (c *Client) GetRelease(ctx context.Context, ref types.PkgRef) (types.ReleaseMetadata, error) {
	var release *githubimpl.RepositoryRelease
	var resp *githubimpl.Response
	var err error

	if ref.Version == "" || ref.Version == "latest" {
		release, resp, err = c.gh.Repositories.GetLatestRelease(context.Background(), ref.Owner, ref.Repo)
	} else {
		release, resp, err = c.gh.Repositories.GetReleaseByTag(context.Background(), ref.Owner, ref.Repo, ref.Version)
	}

	if err != nil {
		return types.ReleaseMetadata{}, err
	}

	if resp == nil || release == nil || resp.StatusCode == 404 {
		return types.ReleaseMetadata{}, types.ErrPackageNotFound
	}

	return convertGithubRelease(*release), nil
}
