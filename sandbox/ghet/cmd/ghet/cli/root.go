package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.jetpack.io/ghet"
)

func RootCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "ghet",
		Short: "Install GitHub releases",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			return ghet.Download(args[0])
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	return command
}

func Execute(ctx context.Context, args []string) int {
	cmd := RootCmd()
	cmd.SetArgs(args)
	err := cmd.ExecuteContext(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		return 1
	}
	return 0
}

func Main() {
	code := Execute(context.Background(), os.Args[1:])
	os.Exit(code)
}
