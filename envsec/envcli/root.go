// Copyright 2022 Jetpack Technologies Inc and contributors. All rights reserved.
// Use of this source code is governed by the license in the LICENSE file.

package envcli

import (
	"context"

	"github.com/MakeNowJust/heredoc"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.jetpack.io/envsec"
)

type cmdConfig struct {
	Store envsec.Store
	EnvId envsec.EnvId
}

// TODO: These should not be global
type globalFlags struct {
	projectId string
	orgId     string
	envName   string
}

func EnvCmd() *cobra.Command {
	flags := &globalFlags{}
	var provider configProvider

	command := &cobra.Command{
		Use:   "envsec",
		Short: "Manage environment variables and secrets",
		Long: heredoc.Doc(`
			Manage environment variables and secrets

			Securely stores and retrieves environment variables on the cloud.
			Environment variables are always encrypted, which makes it possible to
			store values that contain passwords and other secrets.
		`),
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			provider = newConfigProvider(flags)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		// we're manually showing usage
		SilenceUsage: true,
	}

	registerFlags(command, flags)

	command.AddCommand(DownloadCmd(provider))
	command.AddCommand(ExecCmd(provider))
	command.AddCommand(ListCmd(provider))
	command.AddCommand(RemoveCmd(provider))
	command.AddCommand(SetCmd(provider))
	command.AddCommand(UploadCmd(provider))
	command.AddCommand(authCmd())
	command.SetUsageFunc(UsageFunc)
	return command
}

func registerFlags(cmd *cobra.Command, opts *globalFlags) {
	cmd.PersistentFlags().StringVar(
		&opts.projectId,
		"project-id",
		"",
		"Project id to namespace secrets by",
	)

	cmd.PersistentFlags().StringVar(
		&opts.orgId,
		"org-id",
		"",
		"Organization id to namespace secrets by",
	)

	cmd.PersistentFlags().StringVar(
		&opts.envName,
		"environment",
		"dev",
		"Environment name, such as dev or prod",
	)
}

func Execute(ctx context.Context) {
	cmd := EnvCmd()
	_ = cmd.ExecuteContext(ctx)
}

type configProvider func(ctx context.Context) (*cmdConfig, error)

func newConfigProvider(flags *globalFlags) configProvider {
	return func(ctx context.Context) (*cmdConfig, error) {
		s, err := envsec.NewStore(ctx, &envsec.SSMConfig{})
		if err != nil {
			return nil, errors.WithStack(err)
		}

		envid, err := envsec.NewEnvId(flags.projectId, flags.orgId, flags.envName)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return &cmdConfig{
			Store: s,
			EnvId: envid,
		}, nil
	}
}
