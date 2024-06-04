package id

import "go.jetify.com/typeid"

type ProjectPrefix struct{}

func (ProjectPrefix) Prefix() string {
	return "proj"
}

type ProjectID struct {
	typeid.TypeID[ProjectPrefix]
}

type OrgPrefix struct{}

func (OrgPrefix) Prefix() string {
	return "org"
}

type OrgID struct {
	typeid.TypeID[OrgPrefix]
}
