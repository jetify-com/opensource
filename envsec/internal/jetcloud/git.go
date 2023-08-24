package jetcloud

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func createGitIgnore(wd string) error {
	gitIgnorePath := filepath.Join(wd, dirName, ".gitignore")
	return os.WriteFile(gitIgnorePath, []byte("*"), 0600)
}

func gitRepoURL(wd string) (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func gitSubdirectory(wd string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-prefix")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return filepath.Clean(strings.TrimSpace(string(output))), nil
}
