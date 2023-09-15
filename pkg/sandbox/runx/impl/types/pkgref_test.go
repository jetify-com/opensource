package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.jetpack.io/pkg/sandbox/runx/impl/types"
)

func TestFromString_Valid(t *testing.T) {
	testdata := []struct {
		str string
		ref types.PkgRef
	}{
		{
			str: "foo/bar",
			ref: types.PkgRef{
				Owner:   "foo",
				Repo:    "bar",
				Version: "latest",
			},
		},
		{
			str: "foo/bar@v1.2.3",
			ref: types.PkgRef{
				Owner:   "foo",
				Repo:    "bar",
				Version: "v1.2.3",
			},
		},
	}

	for _, td := range testdata {
		t.Run(td.str, func(t *testing.T) {
			ref, err := types.NewPkgRef(td.str)
			assert.NoError(t, err)
			assert.Equal(t, td.ref, ref)
		})
	}
}

func TestFromString_Invalid(t *testing.T) {
	testdata := []struct {
		str string
	}{
		{
			str: "foobar",
		},
		{
			str: "foobar@v1.2.3",
		},
	}

	for _, td := range testdata {
		t.Run(td.str, func(t *testing.T) {
			_, err := types.NewPkgRef(td.str)
			assert.Error(t, err)
		})
	}
}
