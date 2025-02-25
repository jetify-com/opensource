package registry

import (
	"path/filepath"
	"strings"

	"go.jetify.com/pkg/runx/impl/types"
)

func findArtifactForPlatform(artifacts []types.ArtifactMetadata, platform types.Platform) (types.ArtifactMetadata, error) {
	var artifactForPlatform types.ArtifactMetadata
	for _, artifact := range artifacts {
		if isArtifactForPlatform(artifact.Name, platform) {
			artifactForPlatform = artifact
			if isKnownArchive(artifact.Name) {
				// We prefer known archives because sometimes releases contain multiple files
				// for the same platform. Some times those files are alternative installation methods
				// like `.dmg`, `.msi`, or `.deb`, and sometimes they are metadata files like `.sha256`
				// or a `.sig` file. We don't want to install those.
				// If no known archive is found, but the platform is supported, we return the last artifact
				// with a matching platform. This is useful when releases don't have extensions.
				return artifact, nil
			}
		}
	}
	if artifactForPlatform.Name != "" {
		return artifactForPlatform, nil
	}
	return types.ArtifactMetadata{}, types.ErrPlatformNotSupported
}

func isArtifactForPlatform(artifactName string, platform types.Platform) bool {
	// Invalid platform:
	if platform.Arch() == "" || platform.OS() == "" {
		return false
	}

	hasOS := false
	hasArch := false

	// We just check that the artifact name, forced to lowercase,
	// contains the OS and architecture of the invoking system
	if matchesOS(platform, strings.ToLower(artifactName)) {
		hasOS = true
	}
	if matchesArch(platform, strings.ToLower(artifactName)) {
		hasArch = true
	}
	return hasOS && hasArch
}

var alternateOSNames = map[string][]string{
	"darwin": {"macos", "mac"},
}

func matchesOS(platform types.Platform, artifactName string) bool {
	alts := alternateOSNames[platform.OS()]
	for _, alt := range alts {
		if strings.Contains(artifactName, alt) {
			return true
		}
	}
	return strings.Contains(artifactName, platform.OS())
}

var alternateArchNames = map[string][]string{
	"386":   {"i386"},
	"arm64": {"universal"},
	"amd64": {"x86_64", "universal"},
}

func matchesArch(platform types.Platform, artifactName string) bool {
	alts := alternateArchNames[platform.Arch()]
	for _, alt := range alts {
		if strings.Contains(artifactName, alt) {
			return true
		}
	}
	return strings.Contains(artifactName, platform.Arch())
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
