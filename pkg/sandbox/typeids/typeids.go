package typeids

import "go.jetpack.io/typeid"

type ProjectID struct {
	typeid.TypeID
}

var _ typeid.Subtype = (*ProjectID)(nil)

func (ProjectID) AllowedPrefix() string {
	return "proj"
}

type OrgID struct {
	typeid.TypeID
}

var _ typeid.Subtype = (*OrgID)(nil)

func (OrgID) AllowedPrefix() string {
	return "org"
}
