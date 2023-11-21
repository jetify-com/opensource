package cli

import (
	"github.com/spf13/cobra"
)

func LogoutCmd() *cobra.Command {
	command := &cobra.Command{
		Use:  "logout",
		Args: cobra.ExactArgs(0),
		RunE: logoutCmd,
	}

	return command
}

func logoutCmd(cmd *cobra.Command, args []string) error {
	return nil
}
