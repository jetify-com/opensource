package ids

import (
	"go.jetpack.io/typeid"
)

type APIKeyPrefix struct{}

func (APIKeyPrefix) Prefix() string { return "apikey" }

type APIKey struct {
	typeid.TypeID[APIKeyPrefix]
}
