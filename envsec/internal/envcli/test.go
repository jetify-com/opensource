// Copyright 2023 Jetpack Technologies Inc and contributors. All rights reserved.
// Use of this source code is governed by the license in the LICENSE file.

package envcli

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.jetpack.io/envsec/internal/auth"
)

func testCmd() *cobra.Command {

	command := &cobra.Command{
		Use:     "test",
		Aliases: []string{"test"},
		Short:   "test",
		Long:    "test.",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runTestCmd(cmd)
		},
	}
	return command
}

func runTestCmd(cmd *cobra.Command) error {
	params, err := auth.GetParameters()
	if err != nil {
		return err
	}
	fmt.Printf("params: %v\n", params)
	return nil
}
