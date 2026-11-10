package types_test

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.jetify.com/pkg/runx/impl/types"
)

func TestCurrentPlatform(t *testing.T) {
	p := types.CurrentPlatform()
	assert.Equal(t, runtime.GOOS, p.OS())
	assert.Equal(t, runtime.GOARCH, p.Arch())
	assert.Equal(t, runtime.GOOS+"/"+runtime.GOARCH, p.String())
}

func TestNewPlatform(t *testing.T) {
	tests := []struct {
		os, arch         string
		wantOS, wantArch string
	}{
		{"linux", "arm64", "linux", "arm64"},
		{"", "arm64", runtime.GOOS, "arm64"},
		{"windows", "", "windows", runtime.GOARCH},
		{"", "", runtime.GOOS, runtime.GOARCH},
	}
	for _, tt := range tests {
		p := types.NewPlatform(tt.os, tt.arch)
		assert.Equal(t, tt.wantOS, p.OS(), "NewPlatform(%q, %q).OS()", tt.os, tt.arch)
		assert.Equal(t, tt.wantArch, p.Arch(), "NewPlatform(%q, %q).Arch()", tt.os, tt.arch)
	}
}

func TestParsePlatform(t *testing.T) {
	p, err := types.ParsePlatform("darwin/arm64")
	require.NoError(t, err)
	assert.Equal(t, "darwin", p.OS())
	assert.Equal(t, "arm64", p.Arch())
	assert.Equal(t, "darwin/arm64", p.String())

	// Empty halves fall back to the current platform.
	p, err = types.ParsePlatform("/")
	require.NoError(t, err)
	assert.Equal(t, types.CurrentPlatform(), p)

	p, err = types.ParsePlatform("linux/")
	require.NoError(t, err)
	assert.Equal(t, "linux", p.OS())
	assert.Equal(t, runtime.GOARCH, p.Arch())
}

func TestParsePlatformInvalid(t *testing.T) {
	for _, s := range []string{"", "linux", "linux-amd64"} {
		_, err := types.ParsePlatform(s)
		assert.Error(t, err, "ParsePlatform(%q)", s)
	}
}

func TestPkgRefString(t *testing.T) {
	ref := types.PkgRef{Owner: "o", Repo: "r", Version: "v1"}
	assert.Equal(t, "o/r@v1", ref.String())

	parsed, err := types.NewPkgRef(ref.String())
	require.NoError(t, err)
	assert.Equal(t, ref, parsed)
}
