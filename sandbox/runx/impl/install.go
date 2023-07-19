package impl

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"

	"github.com/adrg/xdg"
	"github.com/cavaliergopher/grab/v3"
	"github.com/codeclysm/extract"
	"go.jetpack.io/runx/impl/github"
	"go.jetpack.io/runx/impl/types"
)

var xdgInstallationSubdir = "jetpack.io/pkgs"

func Install(pkgs ...string) error {
	refs := []types.PkgRef{}

	for _, pkg := range pkgs {
		ref, err := types.NewPkgRef(pkg)
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

func install(ref types.PkgRef) error {
	gh := github.NewClient()
	// Figure out latest release:
	release, err := gh.GetRelease(context.Background(), ref)
	if err != nil {
		return err
	}

	resolvedRef := types.PkgRef{
		Owner:   ref.Owner,
		Repo:    ref.Repo,
		Version: release.Name,
	}
	fmt.Printf("Installing %s...\n", resolvedRef)

	// Figure out which asset to download:
	artifact, err := getArtifactForCurrentPlatform(release)
	if err != nil {
		return err
	}
	if artifact == nil {
		return errors.New("no artifact found")
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

	fmt.Printf("Installed at %s...\n", installPath)

	return nil
}

func getArtifactForCurrentPlatform(release types.ReleaseMetadata) (*types.ArtifactMetadata, error) {
	// Attempt to figure out the right artifact for the current platform.
	// TODO:
	// - Support different "templates" for the artifact names if our default heuristic doesn't work.

	platform := types.Platform{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	for _, artifact := range release.Artifacts {
		if isArtifactForPlatform(artifact, platform) {
			return &artifact, nil
		}
	}
	return nil, nil
}

func isArtifactForPlatform(artifact types.ArtifactMetadata, platform types.Platform) bool {
	// Invalid platform:
	if platform.Arch == "" || platform.OS == "" {
		return false
	}

	tokens := strings.FieldsFunc(strings.ToLower(artifact.Name), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	hasOS := false
	hasArch := false

	for _, token := range tokens {
		if token == platform.OS {
			hasOS = true
		}
		if token == platform.Arch {
			hasArch = true
		}
	}
	return hasOS && hasArch
}
