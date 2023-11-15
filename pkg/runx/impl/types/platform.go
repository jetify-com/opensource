package types

import (
	"fmt"
	"runtime"
	"strings"
)

type Platform struct {
	os   string
	arch string
}

func CurrentPlatform() Platform {
	return Platform{
		os:   runtime.GOOS,
		arch: runtime.GOARCH,
	}
}

func NewPlatform(os, arch string) Platform {
	if os == "" {
		os = runtime.GOOS
	}

	if arch == "" {
		arch = runtime.GOARCH
	}

	return Platform{
		os:   os,
		arch: arch,
	}
}

func ParsePlatform(s string) (Platform, error) {
	os, arch, found := strings.Cut(s, "/")
	if !found {
		return Platform{}, fmt.Errorf("invalid platform string: %s", s)
	}
	return NewPlatform(os, arch), nil
}

func (p Platform) OS() string {
	return p.os
}

func (p Platform) Arch() string {
	return p.arch
}

func (p Platform) String() string {
	return fmt.Sprintf("%s/%s", p.os, p.arch)
}
