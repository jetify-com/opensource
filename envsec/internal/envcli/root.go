// Copyright 2022 Jetpack Technologies Inc and contributors. All rights reserved.
// Use of this source code is governed by the license in the LICENSE file.

package envcli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

var globalFlags = &rootCmdFlags{}

type rootCmdFlags struct {
	jsonErrors bool
	idToken    string
}

func RootCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "envsec",
		Short: "Manage environment variables and secrets",
		Long: heredoc.Doc(`
			Manage environment variables and secrets

			Securely stores and retrieves environment variables on the cloud.
			Environment variables are always encrypted, which makes it possible to
			store values that contain passwords and other secrets.
		`),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if globalFlags.jsonErrors {
				// Don't print anything to stderr so we can print the error in json
				cmd.SetErr(io.Discard)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		// we're manually showing usage
		SilenceUsage: true,
		// We manually capture errors so we can print different formats
		SilenceErrors: true,
	}

	command.PersistentFlags().BoolVar(
		&globalFlags.jsonErrors,
		"json-errors", false, "Print errors in json format",
	)
	command.Flag("json-errors").Hidden = true
	command.PersistentFlags().StringVar(
		&globalFlags.idToken,
		"id-token", "", "ID token to use for authentication",
	)
	command.Flag("id-token").Hidden = true

	command.AddCommand(authCmd())
	command.AddCommand(downloadCmd())
	command.AddCommand(execCmd())
	command.AddCommand(initCmd())
	command.AddCommand(listCmd())
	command.AddCommand(removeCmd())
	command.AddCommand(setCmd())
	command.AddCommand(uploadCmd())
	command.SetUsageFunc(UsageFunc)
	return command
}

func Execute(ctx context.Context) int {
	cmd := RootCmd()
	err := cmd.ExecuteContext(ctx)
	if err == nil {
		return 0
	}
	if globalFlags.jsonErrors {
		var jsonErr struct {
			Error string `json:"error"`
		}
		jsonErr.Error = err.Error()
		b, err := json.Marshal(jsonErr)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			fmt.Println(string(b))
		}
	} else {
		fmt.Fprintln(os.Stderr, err)
	}
	return 1

}
