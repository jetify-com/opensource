package cli

// TODO: should it be `pkg show` or `pkg info`?

import (
	"github.com/k0kubun/pp"
	"github.com/spf13/cobra"
	"go.jetpack.io/runx/impl/pkgref"
	"go.jetpack.io/runx/impl/registry"
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
	ref, error := pkgref.FromString(args[0])
	if error != nil {
		return error
	}

	gh := registry.NewGithubRegistry()
	releases, err := gh.ListReleases(cmd.Context(), ref.Owner, ref.Repo)
	if err != nil {
		return err
	}
	pp.Println(releases)
	return nil
}
