package typeid_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.jetify.com/typeid"
	"gopkg.in/yaml.v2"
)

//go:embed testdata/valid.yml
var validEncodingYML []byte

//go:embed testdata/invalid.yml
var invalidEncodingYML []byte

func TestJSONValid(t *testing.T) {
	var testdata []ValidExample
	err := yaml.Unmarshal(validEncodingYML, &testdata)
	require.NoError(t, err)

	for _, td := range testdata {
		t.Run(td.Name, func(t *testing.T) {
			// Test MarshalText via JSON encoding
			tid := typeid.Must(typeid.Parse(td.Tid))
			encoded, err := json.Marshal(tid)
			assert.NoError(t, err)
			assert.Equal(t, `"`+td.Tid+`"`, string(encoded))

			// Test UnmarshalText via JSON decoding
			var decoded typeid.TypeID
			err = json.Unmarshal(encoded, &decoded)
			assert.NoError(t, err)
			assert.Equal(t, tid, decoded)
			assert.Equal(t, td.Tid, decoded.String())
		})
	}
}

func TestJSONInvalid(t *testing.T) {
	var testdata []InvalidExample
	err := yaml.Unmarshal(invalidEncodingYML, &testdata)
	require.NoError(t, err)

	for _, td := range testdata {
		t.Run(td.Name, func(t *testing.T) {
			// Test UnmarshalText with invalid TypeID strings
			var decoded typeid.TypeID
			invalidJSON := `"` + td.Tid + `"`
			err := json.Unmarshal([]byte(invalidJSON), &decoded)
			assert.Error(t, err, "JSON unmarshal should fail for invalid typeid: %s", td.Tid)
		})
	}
}
