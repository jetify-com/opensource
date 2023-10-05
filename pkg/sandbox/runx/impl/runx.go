package impl

import "go.jetpack.io/pkg/sandbox/runx/impl/runopt"

type RunX struct {
	GithubAPIToken string
}

func New(opts runopt.Opts) *RunX {
	return &RunX{
		GithubAPIToken: opts.GithubAPIToken,
	}
}
