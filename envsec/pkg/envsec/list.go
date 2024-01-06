package envsec

import (
	"context"

	"github.com/pkg/errors"
)

func (e *Envsec) List(
	ctx context.Context,
	store Store,
	envIDs ...EnvID,
) (map[EnvID][]EnvVar, error) {
	project, err := e.ProjectConfig()
	if project == nil {
		return nil, err
	}

	authClient, err := e.authClient()
	if err != nil {
		return nil, err
	}

	tok, err := authClient.LoginFlowIfNeededForOrg(ctx, project.OrgID.String())
	if err != nil {
		return nil, err
	}

	store.Identify(ctx, e, tok)

	results := map[EnvID][]EnvVar{}
	for _, envID := range envIDs {
		// TODO: parallelize
		results[envID], err = store.List(ctx, envID)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return results, nil
}
