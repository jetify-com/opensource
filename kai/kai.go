package kai

import "go.jetpack.io/kai/impl"

func Exec(query string) ([]string, error) {
	return impl.Exec(query)
}
