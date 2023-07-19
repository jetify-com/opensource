package gh

import (
	"context"

	"github.com/google/go-github/v53/github"
	"go.jetpack.io/runx/impl/httpcacher"
	"go.jetpack.io/runx/impl/pkgref"
	"go.jetpack.io/runx/impl/types"
)

type Client struct {
	gh *github.Client
}

func NewClient() *Client {
	return &Client{
		gh: github.NewClient(httpcacher.DefaultClient),
	}
}

func (c *Client) ListReleases(ctx context.Context, owner, repo string) ([]types.ReleaseMetadata, error) {
	releases, resp, err := c.gh.Repositories.ListReleases(ctx, owner, repo, nil /* opts */)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 404 {
		return nil, types.ErrPackageNotFound
	}
	return convertGithubReleases(releases), nil
}

func (c *Client) GetRelease(ctx context.Context, ref pkgref.PkgRef) (types.ReleaseMetadata, error) {
	var release *github.RepositoryRelease
	var resp *github.Response
	var err error

	if ref.Version == "" || ref.Version == "latest" {
		release, _, err = c.gh.Repositories.GetLatestRelease(context.Background(), ref.Owner, ref.Repo)
	} else {
		release, _, err = c.gh.Repositories.GetReleaseByTag(context.Background(), ref.Owner, ref.Repo, ref.Version)
	}

	if err != nil {
		return types.ReleaseMetadata{}, err
	}
	if resp.StatusCode == 404 || release == nil {
		return types.ReleaseMetadata{}, types.ErrPackageNotFound
	}

	return convertGithubRelease(*release), nil
}
