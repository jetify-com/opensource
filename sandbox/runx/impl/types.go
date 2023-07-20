package impl

import "time"

type ArtifactMetadata struct {
	DownloadURL string `json:"download_url"`

	Name          string    `json:"name"`
	DownloadCount int       `json:"download_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	ContentType   string    `json:"content_type"`
	Size          int       `json:"size"`
}
