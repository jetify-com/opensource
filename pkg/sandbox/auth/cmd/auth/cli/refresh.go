package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func RefreshCmd() *cobra.Command {
	command := &cobra.Command{
		Use:  "refresh",
		Args: cobra.ExactArgs(0),
		RunE: refreshCmd,
	}

	return command
}

func refreshCmd(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("not implemented")
}
