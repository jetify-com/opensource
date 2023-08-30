// Copyright 2023 Jetpack Technologies Inc and contributors. All rights reserved.
// Use of this source code is governed by the license in the LICENSE file.

package envcli

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.jetpack.io/envsec"
	"go.jetpack.io/envsec/internal/awsfed"
	"go.jetpack.io/envsec/internal/jetcloud"
)

// to be composed into xyzCmdFlags structs
type configFlags struct {
	projectId string
	orgId     string
	envName   string
}

func (f *configFlags) register(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(
		&f.projectId,
		"project-id",
		"",
		"Project id to namespace secrets by",
	)

	cmd.PersistentFlags().StringVar(
		&f.orgId,
		"org-id",
		"",
		"Organization id to namespace secrets by",
	)

	cmd.PersistentFlags().StringVar(
		&f.envName,
		"environment",
		"dev",
		"Environment name, such as dev or prod",
	)
}

type cmdConfig struct {
	Store envsec.Store
	EnvId envsec.EnvId
}

func (f *configFlags) genConfig(ctx context.Context) (*cmdConfig, error) {
	ssmConfig := &envsec.SSMConfig{}
	// if these two flags aren't set we try to fetch orgId and projectId from
	// jetcloud assuming user wants to the jetcloud managed version this way
	// we can accommodate both user-managed and jetcloud-managed envsec
	if f.orgId == "" || f.projectId == "" {
		user, err := newAuthenticator().GetUser()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		awsFederated := awsfed.NewAWSFed()
		ssmConfig, err = awsFederated.GetSSMConfig(user.AccessToken())
		if err != nil {
			return nil, errors.WithStack(err)
		}

		f.orgId = user.OrgId()
		wd, err := os.Getwd()
		if err != nil {
			return nil, errors.New("Failed to get current workding directory")
		}
		projectId, err := jetcloud.ProjectID(wd)
		if err != nil {
			return nil, errors.New("Failed to retrieve projectID")
		}
		f.projectId = projectId.String()

	}
	s, err := envsec.NewStore(ctx, ssmConfig)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	envid, err := envsec.NewEnvId(f.projectId, f.orgId, f.envName)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &cmdConfig{
		Store: s,
		EnvId: envid,
	}, nil
}
