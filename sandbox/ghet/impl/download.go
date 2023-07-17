package impl

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-github/v53/github"
	"go.jetpack.io/ghet/impl/fetch"
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
	gh := github.NewClient(fetch.HTTPClient())
	// Figure out latest release:
	release, _, err := gh.Repositories.GetLatestRelease(context.Background(), ref.Owner, ref.Repo)
	if err != nil {
		// TODO: handle when repo doesn't exist, when tag doesn't exist, and figure out caching
		return err
	}
	resp, err := json.MarshalIndent(release, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(resp))
	return nil
}
