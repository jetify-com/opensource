package api

import (
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
)

// FalseSchema returns a new Schema that fails to validate any value.
// This represents "additionalProperties": false in JSON Schema.
// In JSON Schema 2020-12, false is represented as {"not": {}}.
func FalseSchema() *jsonschema.Schema {
	return &jsonschema.Schema{Not: &jsonschema.Schema{}}
}
