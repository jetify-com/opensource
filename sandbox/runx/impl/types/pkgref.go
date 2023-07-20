package types

import (
	"fmt"
	"strings"
)

type PkgRef struct {
	Owner   string
	Repo    string
	Version string
}

func NewPkgRef(pkg string) (PkgRef, error) {
	version := "latest"
	ownerrepo := pkg

	before, after, found := strings.Cut(pkg, "@")
	if found {
		ownerrepo = before
		version = after
	}

	owner, repo, found := strings.Cut(ownerrepo, "/")
	if !found {
		return PkgRef{}, fmt.Errorf("invalid package reference: %s", pkg)
	}

	return PkgRef{
		Owner:   owner,
		Repo:    repo,
		Version: version,
	}, nil
}

func (p PkgRef) String() string {
	return fmt.Sprintf("%s/%s@%s", p.Owner, p.Repo, p.Version)
}
