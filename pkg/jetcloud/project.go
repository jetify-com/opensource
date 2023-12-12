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

const (
	dirName       = ".jetpack.io"
	configName    = "project.json"
	devConfigName = "dev.project.json"
)

type projectConfig struct {
	ProjectID id.ProjectID `json:"project_id"`
	OrgID     id.OrgID     `json:"org_id"`
}

type InitProjectArgs struct {
	Dir   string
	Force bool
	Token *session.Token
}

func (c *Client) InitProject(
	ctx context.Context,
	args InitProjectArgs,
) (id.ProjectID, error) {
	if args.Token == nil {
		return id.ProjectID{}, errors.Errorf("Please login first")
	}
	existing, err := c.projectID(args.Dir)
	if err == nil && args.Force {
		if err := c.removeConfig(args.Dir); err != nil {
			return id.ProjectID{}, err
		}
	} else if err == nil {
		return existing, ErrProjectAlreadyInitialized
	} else if !os.IsNotExist(err) {
		return id.ProjectID{}, err
	}

	dirPath := filepath.Join(args.Dir, dirName)
	if err = os.MkdirAll(dirPath, 0700); err != nil {
		return id.ProjectID{}, err
	}

	if err = createGitIgnore(args.Dir); err != nil {
		return id.ProjectID{}, err
	}

	repoURL, err := gitRepoURL(args.Dir)
	if err != nil {
		return id.ProjectID{}, err
	}
	subdir, _ := gitSubdirectory(args.Dir)

	projectID, err := c.newProjectID(ctx, args.Token, repoURL, subdir)
	if err != nil {
		return id.ProjectID{}, err
	}

	claims := args.Token.IDClaims()
	if claims == nil {
		return id.ProjectID{}, errors.Errorf("token did not contain an org")
	}

	orgID, err := typeid.Parse[id.OrgID](args.Token.IDClaims().OrgID)
	if err != nil {
		return id.ProjectID{}, err
	}

	cfg := projectConfig{ProjectID: projectID, OrgID: orgID}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return id.ProjectID{}, err
	}
	return projectID, os.WriteFile(filepath.Join(dirPath, c.configName()), data, 0600)
}

func (c *Client) ProjectConfig(wd string) (*projectConfig, error) {
	data, err := os.ReadFile(c.configPath(wd))
	if err != nil {
		return nil, err
	}
	var cfg projectConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Client) projectID(wd string) (id.ProjectID, error) {
	cfg, err := c.ProjectConfig(wd)
	if err != nil {
		return id.ProjectID{}, err
	}
	return cfg.ProjectID, nil
}

func (c *Client) configName() string {
	if c.IsDev {
		return devConfigName
	}
	return configName
}

func (c *Client) configPath(wd string) string {
	return filepath.Join(wd, dirName, c.configName())
}

func (c *Client) removeConfig(wd string) error {
	return os.Remove(c.configPath(wd))
}
