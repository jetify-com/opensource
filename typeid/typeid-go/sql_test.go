package typeid_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.jetify.com/typeid"
)

func TestScan(t *testing.T) {
	t.Parallel()

	testdata := []struct {
		name     string
		input    any
		expected typeid.AnyID
	}{
		{"valid_text", "prefix_01jtvs4hppfp8azhhy9x703dc1", typeid.Must(typeid.FromString("prefix_01jtvs4hppfp8azhhy9x703dc1"))},
		{"valid_tuple", []byte("(prefix,0196b792-46d6-7d90-afc6-3e4f4e01b581)"), typeid.Must(typeid.FromString("prefix_01jtvs4hppfp8azhhy9x703dc1"))},
		{"nil", nil, typeid.AnyID{}},
		{"empty string", "", typeid.AnyID{}},
	}

	for _, td := range testdata {
		t.Run(td.name, func(t *testing.T) {
			t.Parallel()

			var scanned typeid.AnyID
			err := scanned.Scan(td.input)
			assert.NoError(t, err)

			assert.Equal(t, td.expected, scanned)
			assert.Equal(t, td.expected.String(), scanned.String())
		})
	}
}

func TestValuer(t *testing.T) {
	t.Parallel()
	expected := "prefix_01jtvs4hppfp8azhhy9x703dc1"
	tid := typeid.Must(typeid.FromString("prefix_01jtvs4hppfp8azhhy9x703dc1"))
	actual, err := tid.Value()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
