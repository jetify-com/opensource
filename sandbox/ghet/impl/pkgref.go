package impl

type PkgRef struct {
	Owner   string
	Repo    string
	Version string
}

func FromString(pkg string) (PkgRef, error) {
	owner, repo, version := parsePkg(pkg)
	return PkgRef{
		Owner:   owner,
		Repo:    repo,
		Version: version,
	}
}

func parsePkg(pkg string) (owner, repo, version string) {
}

func (p PkgRef) ToString() string {
}
