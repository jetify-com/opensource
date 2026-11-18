//go:build unix

package osdir

import "os"

var forceSystemUser = false

func isSystemUser() bool {
	return forceSystemUser || os.Geteuid() == 0
}
