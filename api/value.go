package api

import (
	"fmt"
)

/**
 * basic types
 */
const (
	LUA_TNONE = iota - 1 // -1
	LUA_TNIL
	LUA_TBOOLEAN
	LUA_TLIGHTUSERDATA
	LUA_TNUMBER
	LUA_TSTRING
	LUA_TTABLE
	LUA_TFUNCTION
	LUA_TUSERDATA
	LUA_TTHREAD
)

type LuaType int

/* type of numbers in Lua */
type LuaNumber = float64

/* type for integer functions */
type LuaInteger = int64

/**
 * variant tags for strings
 */
const (
	LUA_TSHRSTR = LUA_TSTRING | (0 << 4)
	LUA_TLNGSTR = LUA_TSTRING | (1 << 4)
)

/**
 * variant tags for numbers
 */
const (
	LUA_TNUMFLT = LUA_TNUMBER | (0 << 4)
	LUA_TNUMINT = LUA_TNUMBER | (1 << 4)
)

type luaValue interface{}

func typeOf(val luaValue) LuaType {
	switch val.(type) {
	case nil:
		return LUA_TNIL
	case bool:
		return LUA_TBOOLEAN
	case float64, int64:
		return LUA_TNUMBER
	case string:
		return LUA_TSTRING
	default:
		panic(fmt.Sprintf("invalid value: %v (%T)", val, val))
	}
}
