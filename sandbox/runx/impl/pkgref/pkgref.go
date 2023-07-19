package pkgref

import (
	"fmt"
	"strings"
)

type PkgRef struct {
	Owner   string
	Repo    string
	Version string
}

func FromString(pkg string) (PkgRef, error) {
	version := "latest"
	ownerrepo := pkg

	at_splits := strings.SplitN(pkg, "@", 2)
	if len(at_splits) == 2 {
		version = at_splits[1]
		ownerrepo = at_splits[0]
	}

	slash_splits := strings.SplitN(ownerrepo, "/", 2)
	if len(slash_splits) != 2 {
		return PkgRef{}, fmt.Errorf("invalid package reference: %s", pkg)
	}

	return PkgRef{
		Owner:   slash_splits[0],
		Repo:    slash_splits[1],
		Version: version,
	}, nil
}

func (p PkgRef) String() string {
	return fmt.Sprintf("%s/%s@%s", p.Owner, p.Repo, p.Version)
}
