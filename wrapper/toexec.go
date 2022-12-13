package wrapper

import (
	"fmt"

	"github.com/ake-persson/mapslice-json"
)

func ToExecutable(config *Config) *Executable {
	return &Executable{
		Path: config.Path,
		Args: toArgs(config.Flags),
		Env:  toEnv(config.Environment),
	}
}

func toArgs(flags mapslice.MapSlice) []string {
	args := []string{}
	for _, item := range flags {
		args = append(args, fmt.Sprintf("--%s", item.Key))
		switch value := item.Value.(type) {
		case string:
			args = append(args, value)
		case []string:
			args = append(args, value...)
		}
	}
	return args
}

func toEnv(environment mapslice.MapSlice) []string {
	// TODO: give preference to values in environment
	env := []string{} //os.Environ()
	for _, item := range environment {
		env = append(env, fmt.Sprintf("%s=%s", item.Key, item.Value))
	}
	return env
}
