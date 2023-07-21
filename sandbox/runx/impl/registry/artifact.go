package registry

import (
	"fmt"
	"strings"
	"unicode"

	"go.jetpack.io/runx/impl/types"
)

func findArtifactForPlatform(artifacts []types.ArtifactMetadata, platform types.Platform) *types.ArtifactMetadata {
	for _, artifact := range artifacts {
		if isArtifactForPlatform(artifact, platform) {
			return &artifact
		}
	}
	return nil
}

func isArtifactForPlatform(artifact types.ArtifactMetadata, platform types.Platform) bool {
	fmt.Println("Checking artifact: ", artifact.Name)
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
	"darwin": {"macos"},
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
	"amd64": {"x86_64"},
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
