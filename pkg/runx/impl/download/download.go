package download

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/cavaliergopher/grab/v3"
	"go.jetify.com/pkg/fileutil"
)

type Client struct {
	githubAPIToken string
}

func NewClient(accessToken string) *Client {
	return &Client{
		githubAPIToken: accessToken,
	}
}

func (c *Client) DownloadOnce(url, dest string) error {
	dir := filepath.Dir(dest)
	if err := fileutil.EnsureDir(dir); err != nil {
		return err
	}

	info := fileutil.FileInfo(dest)
	if info != nil && info.Mode().IsRegular() && info.Size() > 0 {
		// We've already downloaded it
		return nil
	}
	return c.Download(url, dest)
}

func (c *Client) Download(url, dest string) error {
	if fileutil.IsDir(dest) {
		return errors.New("destination is a directory")
	}

	client := grab.DefaultClient
	tmpDest := dest + ".crdownload"

	// Grab supports partial and resumable downloads, so we don't automatically
	// delete tmpDest until we're done with the download. That way, if we retry
	// we can continue the download where we left off. That said, if the remote file
	// changes, this can result in a corrupted file. For now we assume remote files
	// don't change, since we're are downloading "immutable" releases, but to be safe
	// we'll want to add checksum validation.

	req, err := grab.NewRequest(tmpDest, url)
	if err != nil {
		return err
	}

	req.HTTPRequest.Header.Add("Accept", "application/octet-stream")
	if c.githubAPIToken != "" {
		req.HTTPRequest.Header.Add("Authorization", "Bearer "+c.githubAPIToken)
	}

	resp := client.Do(req)
	if err := resp.Err(); err != nil {
		return err
	}

	err = os.Rename(tmpDest, dest)
	if err != nil {
		return err
	}

	return nil
}
