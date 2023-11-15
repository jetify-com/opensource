package jetcloud

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"go.jetpack.io/pkg/auth/session"
	"go.jetpack.io/pkg/id"
	"go.jetpack.io/typeid"
)

var ErrProjectAlreadyInitialized = errors.New("project already initialized")

const dirName = ".jetpack.io"
const configName = "project.json"

type projectConfig struct {
	ProjectID id.ProjectID `json:"project_id"`
	OrgID     id.OrgID     `json:"org_id"`
}

func InitProject(ctx context.Context, tok *session.Token, dir string) (id.ProjectID, error) {
	if tok == nil {
		return id.ProjectID{}, errors.Errorf("Please login first")
	}
	existing, err := ProjectID(dir)
	if err == nil {
		return existing, ErrProjectAlreadyInitialized
	} else if !os.IsNotExist(err) {
		return id.ProjectID{}, err
	}

	dirPath := filepath.Join(dir, dirName)
	if err = os.MkdirAll(dirPath, 0700); err != nil {
		return id.ProjectID{}, err
	}

	if err = createGitIgnore(dir); err != nil {
		return id.ProjectID{}, err
	}

	repoURL, err := gitRepoURL(dir)
	if err != nil {
		return id.ProjectID{}, err
	}
	subdir, _ := gitSubdirectory(dir)

	projectID, err := newClient().newProjectID(ctx, tok, repoURL, subdir)
	if err != nil {
		return id.ProjectID{}, err
	}

	claims := tok.IDClaims()
	if claims == nil {
		return id.ProjectID{}, errors.Errorf("token did not contain an org")
	}

	orgID, err := typeid.Parse[id.OrgID](tok.IDClaims().OrgID)
	if err != nil {
		return id.ProjectID{}, err
	}

	cfg := projectConfig{ProjectID: projectID, OrgID: orgID}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return id.ProjectID{}, err
	}
	return projectID, os.WriteFile(filepath.Join(dirPath, configName), data, 0600)
}

func ProjectConfig(wd string) (*projectConfig, error) {
	data, err := os.ReadFile(filepath.Join(wd, dirName, configName))
	if err != nil {
		return nil, err
	}
	var cfg projectConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func ProjectID(wd string) (id.ProjectID, error) {
	cfg, err := ProjectConfig(wd)
	if err != nil {
		return id.ProjectID{}, err
	}
	return cfg.ProjectID, nil
}
