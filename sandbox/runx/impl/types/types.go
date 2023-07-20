package types

import "time"

type Platform struct {
	OS   string
	Arch string
}

type ReleaseMetadata struct {
	// Name here is the release name, which often is the same as the tag name (usually the version).
	// Two things:
	// 1. There's also TagName, consider using that instead
	// 2. Should this be called Tag or Version, so we can reserve Name for tha package name if
	//    we want to include that?
	Name        string             `json:"name"`
	CreatedAt   time.Time          `json:"created_at"`
	PublishedAt time.Time          `json:"published_at"`
	Artifacts   []ArtifactMetadata `json:"artifacts"`
}

type ArtifactMetadata struct {
	// TODO: decide which fields are actually required. We are getting a bunch of them from the
	// github api. But if we want to get releases from other sources, or allow publishes to embed
	// this metadata with their releases, some of these won't apply (i.e. DownloadCount).

	DownloadURL string `json:"download_url"`

	Name          string    `json:"name"`
	DownloadCount int       `json:"download_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	ContentType   string    `json:"content_type"`
	Size          int       `json:"size"`
}
