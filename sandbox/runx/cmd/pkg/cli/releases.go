package cli

// TODO: should it be `pkg show` (like poetry) or `pkg info` (like yarn)?
// Might be moot, this CLI is mostly for testing; we'll integrate w/ devbox instead.

import (
	"github.com/k0kubun/pp"
	"github.com/spf13/cobra"
	"go.jetpack.io/runx/impl/github"
	"go.jetpack.io/runx/impl/types"
)

func ReleasesCmd() *cobra.Command {
	command := &cobra.Command{
		Use:  "releases <owner>/<repo>",
		Args: cobra.ExactArgs(1),
		RunE: releasesCmd,
	}

	return command
}

func releasesCmd(cmd *cobra.Command, args []string) error {
	ref, error := types.NewPkgRef(args[0])
	if error != nil {
		return error
	}

	gh := github.NewClient()
	releases, err := gh.ListReleases(cmd.Context(), ref.Owner, ref.Repo)
	if err != nil {
		return err
	}
	pp.Println(releases)
	return nil
}
