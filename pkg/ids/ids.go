package ids

// We have a similar package in our opensource repo, ensure the prefix
// strings match (or figure out how to share a single package across bith
// repos)

import (
	"go.jetify.com/typeid"
)

type userPrefix struct{}

func (userPrefix) Prefix() string { return "user" }

type UserID struct {
	typeid.TypeID[userPrefix]
}

type projectPrefix struct{}

func (projectPrefix) Prefix() string { return "proj" }

type ProjectID struct {
	typeid.TypeID[projectPrefix]
}

type repoPrefix struct{}

func (repoPrefix) Prefix() string { return "repo" }

type RepoID struct {
	typeid.TypeID[repoPrefix]
}

type orgPrefix struct{}

func (orgPrefix) Prefix() string { return "org" }

type OrgID struct {
	typeid.TypeID[orgPrefix]
}

type memberPrefix struct{}

func (memberPrefix) Prefix() string { return "member" }

type MemberID struct {
	typeid.TypeID[memberPrefix]
}

type secretPrefix struct{}

func (secretPrefix) Prefix() string { return "secret" }

type SecretID struct{ typeid.TypeID[secretPrefix] }

type deploymentPrefix struct{}

func (deploymentPrefix) Prefix() string { return "deploy" }

type DeploymentID struct {
	typeid.TypeID[deploymentPrefix]
}

type customDomainPrefix struct{}

func (customDomainPrefix) Prefix() string { return "customdomain" }

type CustomDomainID struct {
	typeid.TypeID[customDomainPrefix]
}

func Short[P typeid.PrefixType](tid typeid.TypeID[P]) string {
	return ShortStr(tid.String()) // NOTE: len(id.String) >= 26 always
}

// ShortStr returns the last 6 characters of a TypeID's string representation.
func ShortStr(tid string) string {
	return tid[max(0, len(tid)-6):]
}
