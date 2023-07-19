package impl

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"

	"github.com/adrg/xdg"
	"github.com/cavaliergopher/grab/v3"
	"github.com/codeclysm/extract"
	"github.com/google/go-github/v53/github"
	"go.jetpack.io/runx/impl/httpcacher"
	"go.jetpack.io/runx/impl/pkgref"
)

var xdgInstallationSubdir = "jetpack.io/pkgs"

func Install(pkgs ...string) error {
	refs := []pkgref.PkgRef{}

	for _, pkg := range pkgs {
		ref, err := pkgref.FromString(pkg)
		if err != nil {
			return err
		}
		refs = append(refs, ref)
	}

	for _, ref := range refs {
		err := install(ref)
		if err != nil {
			return err
		}
	}
	return nil
}

func install(ref pkgref.PkgRef) error {
	// Figure out latest release:
	release, err := getReleaseMetadata(ref)
	if err != nil {
		return err
	}

	resolvedRef := pkgref.PkgRef{
		Owner:   ref.Owner,
		Repo:    ref.Repo,
		Version: release.GetTagName(),
	}
	fmt.Printf("Installing %s...\n", resolvedRef)

	// Figure out which asset to download:
	artifact, err := getArtifactMetadata(release)
	if err != nil {
		return err
	}

	installPath := filepath.Join(xdg.CacheHome, xdgInstallationSubdir, resolvedRef.Owner, resolvedRef.Repo, resolvedRef.Version)
	err = os.MkdirAll(installPath, 0700)
	if err != nil {
		return err
	}

	// TODO: Add httpcacher
	grabResp, err := grab.Get(installPath, artifact.DownloadURL)
	if err != nil {
		return err
	}

	reader, err := grabResp.Open()
	if err != nil {
		return err
	}

	// TODO: only extract if we haven't already extracted
	err = extract.Archive(context.Background(), reader, installPath, nil)
	if err != nil {
		return err
	}

	err = os.Remove(grabResp.Filename)
	if err != nil {
		return err
	}

	return nil
}

func getArtifactMetadata(releaseMeta *github.RepositoryRelease) (*ArtifactMetadata, error) {
	// Attempt to figure out the right artifact for the current platform.
	// TODO:
	// - Pass platform as an argument
	// - Support different "templates" for the artifact names

	assetNames := []string{}
	for _, asset := range releaseMeta.Assets {
		assetNames = append(assetNames, *asset.Name)
	}
	fmt.Println(assetNames)

	for _, asset := range releaseMeta.Assets {
		if isAssetForCurrentPlatform(asset) {
			return &ArtifactMetadata{
				DownloadURL:   asset.GetBrowserDownloadURL(),
				Name:          asset.GetName(),
				DownloadCount: asset.GetDownloadCount(),
				CreatedAt:     asset.GetCreatedAt().Time,
				UpdatedAt:     asset.GetUpdatedAt().Time,
				ContentType:   asset.GetContentType(),
				Size:          asset.GetSize(),
			}, nil
		}
	}
	return nil, nil
}

func isAssetForCurrentPlatform(asset *github.ReleaseAsset) bool {
	tokens := strings.FieldsFunc(strings.ToLower(asset.GetName()), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	hasOS := false
	hasArch := false

	for _, token := range tokens {
		if token == runtime.GOOS {
			hasOS = true
		}
		if token == runtime.GOARCH {
			hasArch = true
		}
	}
	return hasOS && hasArch
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
