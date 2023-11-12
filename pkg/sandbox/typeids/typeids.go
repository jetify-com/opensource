package typeids

import "go.jetpack.io/typeid"

type ProjectID struct {
	typeid.TypeID
}

func (ProjectID) AllowedPrefix() string {
	return "proj"
}

type OrgID struct {
	typeid.TypeID
}

func (OrgID) AllowedPrefix() string {
	return "org"
}
