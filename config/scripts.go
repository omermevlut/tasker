package config

import (
	"path/filepath"
	"runtime"
)

// LuaScripts config
var LuaScripts = struct {
	MoveExpired string
	PopActive   string
}{
	MoveExpired: getRoot() + "/../scripts/move_expired.lua",
	PopActive:   getRoot() + "/../scripts/pop.lua",
}

func getRoot() string {
	_, b, _, _ := runtime.Caller(1)

	return filepath.Dir(b)
}
