package cli

// TODO: should it be `pkg show` (like poetry) or `pkg info` (like yarn)?
// Might be moot, this CLI is mostly for testing; we'll integrate w/ devbox instead.

import (
	"os"

	"github.com/k0kubun/pp"
	"github.com/spf13/cobra"
	"go.jetify.com/pkg/runx/impl/registry"
	"go.jetify.com/pkg/runx/impl/types"
)

func ResolveCmd() *cobra.Command {
	command := &cobra.Command{
		Use:  "resolve <owner>/<repo>@<version>",
		Args: cobra.ExactArgs(1),
		RunE: resolveCmd,
	}

	return command
}

func resolveCmd(cmd *cobra.Command, args []string) error {
	ref, error := types.NewPkgRef(args[0])
	if error != nil {
		return error
	}

	registry, err := registry.NewLocalRegistry(cmd.Context(), os.Getenv("RUNX_GITHUB_API_TOKEN"))
	if err != nil {
		return err
	}

	pp.Println(registry.ResolveVersion(ref))
	return nil
}
