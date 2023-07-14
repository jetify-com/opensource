package impl

import (
	"context"
	"strings"

	"github.com/google/go-github/v53/github"
)

func Download(pkgs ...string) error {
	for _, pkg := range pkgs {
		owner, repo, version := parsePkg(pkg)
		err := DownloadVersion(owner, repo, version)
		if err != nil {
			return err
		}
	}
	return nil
}

func parsePkg(pkg string) (owner, repo, version string) {
	parts := strings.SplitN(pkg, "@", 2)
	if len(parts) == 1 {
		return parts[0], "", ""
	}
	return parts[0], parts[1], ""
}

func DownloadLatest(owner string, repo string) error {
	gh := github.NewClient(nil)
	_, _, err := gh.Repositories.GetLatestRelease(context.Background(), owner, repo)
	if err != nil {
		return err
	}
	return nil
}

func DownloadVersion(owner, repo, version string) error {
	if version == "latest" {
		version = ""
	}
	return nil
}
