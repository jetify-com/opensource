package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func WhoAmICmd() *cobra.Command {
	command := &cobra.Command{
		Use:  "whoami",
		Args: cobra.ExactArgs(0),
		RunE: whoamiCmd,
	}

	return command
}

func whoamiCmd(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("not implemented")
}
