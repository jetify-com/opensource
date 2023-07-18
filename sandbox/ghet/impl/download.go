package impl

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-github/v53/github"
	"go.jetpack.io/ghet/impl/httpcacher"
	"go.jetpack.io/ghet/impl/pkgref"
)

func Download(pkgs ...string) error {
	refs := []pkgref.PkgRef{}

	for _, pkg := range pkgs {
		ref, err := pkgref.FromString(pkg)
		if err != nil {
			return err
		}
		refs = append(refs, ref)
	}

	for _, ref := range refs {
		err := download(ref)
		if err != nil {
			return err
		}
	}
	return nil
}

func download(ref pkgref.PkgRef) error {
	fmt.Printf("Downloading %s...\n", ref)
	// Figure out latest release:
	release, err := getReleaseMetadata(ref)
	if err != nil {
		return err
	}

	resp, err := json.MarshalIndent(release, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(resp))
	return nil
}

func getReleaseMetadata(ref pkgref.PkgRef) (*github.RepositoryRelease, error) {
	gh := github.NewClient(httpcacher.DefaultClient)
	var release *github.RepositoryRelease
	var err error

	// TODO: handle when repo doesn't exist, when tag doesn't exist, etc.

	if ref.Version == "" || ref.Version == "latest" {
		release, _, err = gh.Repositories.GetLatestRelease(context.Background(), ref.Owner, ref.Repo)
		if err != nil {
			return nil, err
		}
	} else {
		release, _, err = gh.Repositories.GetReleaseByTag(context.Background(), ref.Owner, ref.Repo, ref.Version)
		if err != nil {
			return nil, err
		}
	}
	return release, nil
}
