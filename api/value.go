package api

import (
	"fmt"

	"github.com/uganh16/luago/number"
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

func toBoolean(val luaValue) bool {
	switch val := val.(type) {
	case nil:
		return false
	case bool:
		return val
	default:
		return true
	}
}

func toNumber(val luaValue) (float64, bool) {
	switch val := val.(type) {
	case float64:
		return val, true
	case int64:
		return float64(val), true
	case string:
		return number.ParseFloat(val)
	default:
		return 0.0, false
	}
}

func toInteger(val luaValue) (int64, bool) {
	switch val := val.(type) {
	case int64:
		return val, true
	case float64:
		return number.FloatToInteger(val)
	case string:
		if val, ok := number.ParseInteger(val); ok {
			return val, ok
		}
		if val, ok := number.ParseFloat(val); ok {
			return number.FloatToInteger(val)
		}
	}
	return 0, false
}

func toString(val luaValue) (string, bool) {
	switch val := val.(type) {
	case string:
		return val, true
	case float64, int64:
		return fmt.Sprintf("%v", val), true
	default:
		return "", false
	}
}

func equal(a, b luaValue) bool {
	switch a := a.(type) {
	case nil:
		return b == nil
	case bool:
		b, ok := b.(bool)
		return ok && a == b
	case string:
		b, ok := b.(string)
		return ok && a == b
	case int64:
		switch b := b.(type) {
		case int64:
			return a == b
		case float64:
			return float64(a) == b
		default:
			return false
		}
	case float64:
		switch b := b.(type) {
		case float64:
			return a == b
		case int64:
			return a == float64(b)
		default:
			return false
		}
	default:
		return a == b
	}
}

func lessThan(L *LuaState, a, b luaValue) bool {
	switch a := a.(type) {
	case string:
		if b, ok := b.(string); ok {
			return a < b
		}
	case int64:
		switch b := b.(type) {
		case int64:
			return a < b
		case float64:
			return float64(a) < b
		}
	case float64:
		switch b := b.(type) {
		case float64:
			return a < b
		case int64:
			return a < float64(b)
		}
	}
	panic(orderError(L, a, b))
}

func lessEqual(L *LuaState, a, b luaValue) bool {
	switch a := a.(type) {
	case string:
		if b, ok := b.(string); ok {
			return a <= b
		}
	case int64:
		switch b := b.(type) {
		case int64:
			return a <= b
		case float64:
			return float64(a) <= b
		}
	case float64:
		switch b := b.(type) {
		case float64:
			return a <= b
		case int64:
			return a <= float64(b)
		}
	}
	panic(orderError(L, a, b))
}
