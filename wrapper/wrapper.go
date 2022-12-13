package wrapper

import (
	"encoding/json"
	"os"
	"syscall"
)

type Wrapper struct {
	Config *Config
}

func New(config *Config) *Wrapper {
	return &Wrapper{
		Config: config,
	}
}

func FromPath(path string) (*Wrapper, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return FromBytes(bytes)
}

func FromBytes(bytes []byte) (*Wrapper, error) {
	config := &Config{}
	err := json.Unmarshal(bytes, config)
	if err != nil {
		return nil, err
	}
	return New(config), nil
}

func (w *Wrapper) Exec() error {
	exe := ToExecutable(w.Config)
	return syscall.Exec(exe.Path, exe.Args, exe.Env)
}
