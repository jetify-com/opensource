//go:build windows

package osdir

import (
	"golang.org/x/sys/windows"
)

var Cache = DirType{
	System: "$ProgramData",
	User:   "$XDG_CACHE_HOME", UserDefault: "$LocalAppData",
}

var Config = DirType{
	System: "$ProgramData",
	User:   "$XDG_CONFIG_HOME", UserDefault: "$AppData",
}

var Data = DirType{
	System: "$ProgramData",
	User:   "$XDG_DATA_HOME", UserDefault: "$AppData",
}

var State = DirType{
	System: "$ProgramData",
	User:   "$XDG_STATE_HOME", UserDefault: "$LocalAppData",
}

var forceSystemUser = false

func isSystemUser() bool {
	if forceSystemUser {
		return true
	}

	var token windows.Token
	err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &token)
	if err != nil {
		return false
	}
	defer token.Close()

	return token.IsElevated()
}
