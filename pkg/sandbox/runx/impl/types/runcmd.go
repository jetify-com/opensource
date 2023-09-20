package types

// TODO: is this the best name for this struct?

type RunCmd struct {
	Packages []PkgRef
	App      string
	Args     []string
}
