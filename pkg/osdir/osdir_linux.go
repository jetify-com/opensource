//go:build linux

package osdir

var Cache = DirType{
	System: "/var/cache",
	User:   "$XDG_CACHE_HOME", UserDefault: "$HOME/.cache",
}

var Config = DirType{
	System: "/etc",
	User:   "$XDG_CONFIG_HOME", UserDefault: "$HOME/.config",
	Search: "$XDG_CONFIG_DIRS", SearchDefault: "/etc/xdg",
}

var Data = DirType{
	System: "/usr/share",
	User:   "$XDG_DATA_HOME", UserDefault: "$HOME/.local/share",
	Search: "$XDG_DATA_DIRS", SearchDefault: "/usr/local/share:/usr/share",
}

var State = DirType{
	System: "/var/lib",
	User:   "$XDG_STATE_HOME", UserDefault: "$HOME/.local/state",
}
