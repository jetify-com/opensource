package wrapper

import "github.com/ake-persson/mapslice-json"

type Config struct {
	Path        string            `json:"path"`
	Environment mapslice.MapSlice `json:"environment"`
	Flags       mapslice.MapSlice `json:"flags"`
}

type Executable struct {
	Path string
	Args []string
	Env  []string
}
