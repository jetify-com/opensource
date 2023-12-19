package flow

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/samber/lo"
	"go.jetpack.io/envsec/internal/git"
	"go.jetpack.io/pkg/api/api"
	projectsv1alpha1 "go.jetpack.io/pkg/api/gen/priv/projects/v1alpha1"
	"go.jetpack.io/pkg/auth/session"
	"go.jetpack.io/pkg/id"
	"go.jetpack.io/typeid"
)

// flow:
// 0. Ask if you want to overwrite existing config [y/N]
// 1. Link to an existing project? [Y/n]
// 2a. What project would you like to link to? (sorted by repo/dir match)
// 2b. What’s the name of your new project?

type Init struct {
	Client                *api.Client
	PromptOverwriteConfig bool
	Token                 *session.Token
	WorkingDir            string
}

func (i *Init) Run(ctx context.Context) (id.ProjectID, error) {
	if i.PromptOverwriteConfig {
		overwrite, err := i.overwriteConfigPrompt()
		if err != nil {
			return id.ProjectID{}, err
		}
		if !overwrite {
			return id.ProjectID{}, errors.New("aborted")
		}
	}

	// TODO: printTeamNotice will be a team picker once that is implemented.
	if err := i.printTeamNotice(ctx); err != nil {
		return id.ProjectID{}, err
	}
	linkToExisting, err := i.linkToExistingPrompt()
	if err != nil {
		return id.ProjectID{}, err
	}
	if linkToExisting {
		return i.showExistingListPrompt(ctx)
	}
	return i.createNewPrompt(ctx)
}

func (i *Init) overwriteConfigPrompt() (bool, error) {
	return boolPrompt(
		"Project already exists. Overwrite existing project config",
		"n",
	)
}

func (i *Init) printTeamNotice(ctx context.Context) error {
	member, err := i.Client.GetMember(ctx, i.Token.IDClaims().Subject)
	if err != nil {
		return err
	}
	fmt.Fprintf(
		os.Stderr,
		"Initializing project for %s\n",
		member.Organization.Name,
	)
	return nil
}

func (i *Init) linkToExistingPrompt() (bool, error) {
	return boolPrompt("Link to an existing project", "y")
}

func (i *Init) showExistingListPrompt(
	ctx context.Context,
) (id.ProjectID, error) {
	orgID, err := typeid.Parse[id.OrgID](i.Token.IDClaims().OrgID)
	if err != nil {
		return id.ProjectID{}, err
	}

	projects, err := i.Client.ListProjects(ctx, orgID)
	if err != nil {
		return id.ProjectID{}, err
	}

	repo, err := git.GitRepoURL(i.WorkingDir)
	if err != nil {
		return id.ProjectID{}, err
	}

	directory, err := git.GitSubdirectory(i.WorkingDir)
	if err != nil {
		return id.ProjectID{}, err
	}

	sort.SliceStable(projects, func(i, j int) bool {
		return projects[i].GetRepo() == repo &&
			projects[i].GetDirectory() == directory
	})

	prompt := promptui.Select{
		Label: "What project would you like to link to",
		Items: lo.Map(projects, func(p *projectsv1alpha1.Project, _ int) string {
			item := strings.TrimSpace(p.GetName())
			if item == "" {
				item = "unnamed project"
			}
			if p.GetRepo() != "" {
				item += " repo: " + p.GetRepo()
			}
			if p.GetDirectory() != "" && p.GetDirectory() != "." {
				item += " dir: " + p.GetDirectory()
			}
			return item
		}),
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return id.ProjectID{}, err
	}

	return typeid.Parse[id.ProjectID](projects[idx].GetId())
}

func (i *Init) createNewPrompt(ctx context.Context) (id.ProjectID, error) {
	prompt := promptui.Prompt{
		Label:   "What’s the name of your new project",
		Default: filepath.Base(i.WorkingDir),
		Validate: func(name string) error {
			if name == "" {
				return errors.New("project name cannot be empty")
			}
			return nil
		},
	}

	name, err := prompt.Run()
	if err != nil {
		return id.ProjectID{}, err
	}

	orgID, err := typeid.Parse[id.OrgID](i.Token.IDClaims().OrgID)
	if err != nil {
		return id.ProjectID{}, err
	}

	repo, err := git.GitRepoURL(i.WorkingDir)
	if err != nil {
		return id.ProjectID{}, err
	}

	directory, err := git.GitSubdirectory(i.WorkingDir)
	if err != nil {
		return id.ProjectID{}, err
	}

	project, err := i.Client.CreateProject(
		ctx,
		orgID,
		repo,
		directory,
		name,
	)
	if err != nil {
		return id.ProjectID{}, err
	}
	return typeid.Parse[id.ProjectID](project.GetId())
}

func boolPrompt(label, defaultResult string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
		Default:   defaultResult,
	}

	result, err := prompt.Run()
	// promptui.ErrAbort is returned when user enters "n" which is valid.
	if err != nil && !errors.Is(err, promptui.ErrAbort) {
		return false, err
	}
	if result == "" {
		result = defaultResult
	}

	return strings.ToLower(result) == "y", nil
}
