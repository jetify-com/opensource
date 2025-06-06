package interpreter

import (
	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"
	"go.jetify.com/tyson/internal/tsembed"
)

func Eval(entrypoint string) (goja.Value, error) {
	return tsembed.Eval(entrypoint, tsembed.Options{
		Plugins: []api.Plugin{
			tsonTransform,
		},
	})
}
