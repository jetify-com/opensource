package envcli

import (
	"os"

	"github.com/spf13/cobra"
)

func projectsCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "projects",
		Short: "envsec projects commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			// show help
			return cmd.Help()
		},
	}

	command.AddCommand(currentProjectCmd())

	return command
}

func currentProjectCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "current",
		Short: "show current project",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			workingDir, err := os.Getwd()
			if err != nil {
				return err
			}
			return defaultEnvsec(cmd, workingDir).
				DescribeCurrentProject(cmd.Context(), cmd.OutOrStdout())
		},
	}

	return command
}
