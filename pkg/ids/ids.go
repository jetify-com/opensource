package ids

// We have a similar package in our axiom repo, ensure the prefix
// strings match (or figure out how to share a single package across both
// repos)

import (
	"fmt"

	"go.jetify.com/typeid/v2"
)

const UserPrefix = "user"

type UserID struct {
	typeid.TypeID
}

func NewUserID() (UserID, error) {
	return new[UserID](UserPrefix)
}

func ParseUserID(s string) (UserID, error) {
	return parse[UserID](s, UserPrefix)
}

const ProjectPrefix = "proj"

type ProjectID struct {
	typeid.TypeID
}

func NewProjectID() (ProjectID, error) {
	return new[ProjectID](ProjectPrefix)
}

func ParseProjectID(s string) (ProjectID, error) {
	return parse[ProjectID](s, ProjectPrefix)
}

const RepoPrefix = "repo"

type RepoID struct {
	typeid.TypeID
}

func NewRepoID() (RepoID, error) {
	return new[RepoID](RepoPrefix)
}

func ParseRepoID(s string) (RepoID, error) {
	return parse[RepoID](s, RepoPrefix)
}

const OrgPrefix = "org"

type OrgID struct {
	typeid.TypeID
}

func NewOrgID() (OrgID, error) {
	return new[OrgID](OrgPrefix)
}

func ParseOrgID(s string) (OrgID, error) {
	return parse[OrgID](s, OrgPrefix)
}

const MemberPrefix = "member"

type MemberID struct {
	typeid.TypeID
}

func NewMemberID() (MemberID, error) {
	return new[MemberID](MemberPrefix)
}

func ParseMemberID(s string) (MemberID, error) {
	return parse[MemberID](s, MemberPrefix)
}

const SecretPrefix = "secret"

type SecretID struct {
	typeid.TypeID
}

func NewSecretID() (SecretID, error) {
	return new[SecretID](SecretPrefix)
}

func ParseSecretID(s string) (SecretID, error) {
	return parse[SecretID](s, SecretPrefix)
}

const DeploymentPrefix = "deploy"

type DeploymentID struct {
	typeid.TypeID
}

func NewDeploymentID() (DeploymentID, error) {
	return new[DeploymentID](DeploymentPrefix)
}

func ParseDeploymentID(s string) (DeploymentID, error) {
	return parse[DeploymentID](s, DeploymentPrefix)
}

const CustomDomainPrefix = "customdomain"

type CustomDomainID struct {
	typeid.TypeID
}

func NewCustomDomainID() (CustomDomainID, error) {
	return new[CustomDomainID](CustomDomainPrefix)
}

func ParseCustomDomainID(s string) (CustomDomainID, error) {
	return parse[CustomDomainID](s, CustomDomainPrefix)
}

// ShortSuffix returns the last 6 characters of a TypeID's string representation.
func ShortSuffix(tid typeid.TypeID) string {
	return tid.String()[max(0, len(tid.String())-6):]
}

// IDType is a constraint for ID types that wrap typeid.TypeID
type IDType interface {
	~struct{ typeid.TypeID }
}

// These helpers exist because we decided to create subtypes for each ID type.
// With the new type id implementation, we could choose to get rid of these different
// types and just pass typeid.TypeID around. If we do that, it would be slightly less
// type safe, but we could completely get rid of these helpers.

// parse is a generic helper for parsing ID strings
func parse[T IDType](s, prefix string) (T, error) {
	var zero T
	tid, err := typeid.Parse(s)
	if err != nil {
		return zero, err
	}
	if tid.Prefix() != prefix {
		return zero, fmt.Errorf("invalid %s ID: %s", prefix, s)
	}
	return T{TypeID: tid}, nil
}

// new is a generic helper for generating new IDs
func new[T IDType](prefix string) (T, error) {
	var zero T
	tid, err := typeid.Generate(prefix)
	if err != nil {
		return zero, err
	}
	return T{TypeID: tid}, nil
}
