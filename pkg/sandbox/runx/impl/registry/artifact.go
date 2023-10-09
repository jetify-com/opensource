package registry

import (
	"path/filepath"
	"strings"
	"unicode"

	"go.jetpack.io/pkg/sandbox/runx/impl/types"
)

func findArtifactForPlatform(artifacts []types.ArtifactMetadata, platform types.Platform) *types.ArtifactMetadata {
	var artifactForPlatform types.ArtifactMetadata
	for _, artifact := range artifacts {
		if isArtifactForPlatform(artifact, platform) {
			artifactForPlatform = artifact
			if isKnownArchive(artifact.Name) {
				// We only consider known archives because sometimes releases contain multiple files
				// for the same platform. Some times those files are alternative installation methods
				// like `.dmg`, `.msi`, or `.deb`, and sometimes they are metadata files like `.sha256`
				// or a `.sig` file. We don't want to install those.
				return &artifact
			}
		}
	}
	// Best attempt:
	return &artifactForPlatform
}

func isArtifactForPlatform(artifact types.ArtifactMetadata, platform types.Platform) bool {
	// Invalid platform:
	if platform.Arch() == "" || platform.OS() == "" {
		return false
	}

	// As a heuristic we tokenize the name of the artifact, and return the artifact that has
	// tokens for both the OS and the Architecture.
	tokens := strings.FieldsFunc(strings.ToLower(artifact.Name), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	hasOS := false
	hasArch := false

	for _, token := range tokens {
		if matchesOS(platform, token) {
			hasOS = true
			continue
		}
		if matchesArch(platform, token) {
			hasArch = true
			continue
		}
		if hasOS && hasArch {
			return true
		}
	}
	return hasOS && hasArch
}

var alternateOSNames = map[string][]string{
	"darwin": {"macos", "mac"},
}

func matchesOS(platform types.Platform, token string) bool {
	if token == platform.OS() {
		return true
	}
	alts := alternateOSNames[platform.OS()]
	for _, alt := range alts {
		if token == alt {
			return true
		}
	}
	return false
}

var alternateArchNames = map[string][]string{
	"386":   {"i386"},
	"arm64": {"universal"},
	"amd64": {"x86_64", "universal"},
}

func matchesArch(platform types.Platform, token string) bool {
	if token == platform.Arch() {
		return true
	}
	alts := alternateArchNames[platform.Arch()]
	for _, alt := range alts {
		if token == alt {
			return true
		}
	}
	return false
}

var knownExts = []string{
	".bz2",
	".gz",
	".lz",
	".lzma",
	".lzo",
	".tar",
	".taz",
	".taZ",
	".tbz",
	".tbz2",
	".tgz",
	".tlz",
	".tz2",
	".tzst",
	".xz",
	".Z",
	".zip",
	".zst",
}

func isKnownArchive(name string) bool {
	ext := filepath.Ext(name)
	for _, knownExt := range knownExts {
		if ext == knownExt {
			return true
		}
	}
	return false
}
