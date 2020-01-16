package config

// LuaScripts config
var LuaScripts = struct {
	MoveExpired string
	PopActive   string
}{
	MoveExpired: "/../scripts/move_expired.lua",
	PopActive:   "/../scripts/pop.lua",
}
