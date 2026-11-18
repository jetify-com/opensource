//go:build darwin

package osdir

var Cache = DirType{
	System: "/Library/Caches",
	User:   "$XDG_CACHE_HOME", UserDefault: "$HOME/Library/Caches",
}

var Config = DirType{
	System: "/Library/Application Support",
	User:   "$XDG_CONFIG_HOME", UserDefault: "$HOME/Library/Application Support",
}

var Data = DirType{
	System: "/Library/Application Support",
	User:   "$XDG_DATA_HOME", UserDefault: "$HOME/Library/Application Support",
}

var State = DirType{
	System: "/Library/Application Support",
	User:   "$XDG_STATE_HOME", UserDefault: "$HOME/Library/Application Support",
}
