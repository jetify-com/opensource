package ids

import (
	"go.jetify.com/typeid"
)

type APITokenPrefix struct{}

func (APITokenPrefix) Prefix() string { return "api_token" }

type APIToken struct {
	typeid.TypeID[APITokenPrefix]
}
