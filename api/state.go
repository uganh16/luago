package api

import (
	"fmt"
)

type LuaState struct {
	stack []luaValue
}

/**
 * state manipulation
 */

func NewState() *LuaState {
	return &LuaState{
		stack: make([]luaValue, 0, 20),
	}
}

/**
 * basic stack manipulation
 */

func (L *LuaState) AbsIndex(idx int) int {
	if idx > 0 {
		return idx
	}
	return idx + len(L.stack) + 1
}

func (L *LuaState) GetTop() int {
	return len(L.stack)
}

func (L *LuaState) SetTop(idx int) {
	top := len(L.stack)
	if idx >= 0 {
		if idx > cap(L.stack) {
			panic("new top too large")
		}
		for top < idx {
			L.stack = append(L.stack, nil)
			top++
		}
	} else {
		if -idx > top {
			panic("invalid new top")
		}
		idx = top + idx + 1
	}
	for top > idx {
		top--
		L.stack[top] = nil
	}
	L.stack = L.stack[:top]
}

func (L *LuaState) PushValue(idx int) {
	val, _ := L.stackGet(idx)
	L.stackPush(val)
}

func (L *LuaState) Rotate(idx, n int) {
	t := len(L.stack) - 1    // end of stack segment being rotated
	p := L.AbsIndex(idx) - 1 // start of segment
	if p < 0 || p > t {
		panic("index not in the stack")
	}
	var m int // end of prefix
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}
	if m < p || m > t {
		panic("invalid 'n'")
	}
	L.stackReverse(p, m)   // reverse the prefix with length 'n'
	L.stackReverse(m+1, t) // reverse the suffix
	L.stackReverse(p, t)   // reverse the entire segment
}

func (L *LuaState) Copy(srcIdx, dstIdx int) {
	val, _ := L.stackGet(srcIdx)
	L.stackSet(dstIdx, val)
}

func (L *LuaState) CheckStack(n int) bool {
	if cap(L.stack)-len(L.stack) < n {
		newSize := cap(L.stack) * 2
		if newSize < len(L.stack)+n {
			newSize = len(L.stack) + n
		}
		newStack := make([]luaValue, len(L.stack), newSize)
		copy(newStack, L.stack)
		L.stack = newStack
	}
	return true
}

/**
 * access functions (stack -> Go)
 */

func (L *LuaState) IsNumber(idx int) bool {
	_, ok := L.ToNumberX(idx)
	return ok
}

func (L *LuaState) IsString(idx int) bool {
	t := L.Type(idx)
	return t == LUA_TSTRING || t == LUA_TNUMBER
}

func (L *LuaState) IsInteger(idx int) bool {
	val, _ := L.stackGet(idx)
	_, ok := val.(int64)
	return ok
}

func (L *LuaState) Type(idx int) LuaType {
	if val, ok := L.stackGet(idx); ok {
		return typeOf(val)
	}
	return LUA_TNONE
}

func (L *LuaState) TypeName(t LuaType) string {
	switch t {
	case LUA_TNONE:
		return "no value"
	case LUA_TNIL:
		return "nil"
	case LUA_TBOOLEAN:
		return "boolean"
	case LUA_TLIGHTUSERDATA:
		return "userdata"
	case LUA_TNUMBER:
		return "number"
	case LUA_TSTRING:
		return "string"
	case LUA_TTABLE:
		return "table"
	case LUA_TFUNCTION:
		return "function"
	case LUA_TUSERDATA:
		return "userdata"
	case LUA_TTHREAD:
		return "thread"
	default:
		panic("invalid tag")
	}
}

func (L *LuaState) ToNumberX(idx int) (float64, bool) {
	val, _ := L.stackGet(idx)
	switch val := val.(type) {
	case float64:
		return val, true
	case int64:
		return float64(val), true
	/* @todo string convertible to number? */
	default:
		return 0.0, false
	}
}

func (L *LuaState) ToIntegerX(idx int) (int64, bool) {
	val, _ := L.stackGet(idx)
	switch val := val.(type) {
	case int64:
		return val, true
	/* @todo try to convert a value to an integer */
	default:
		return 0, false
	}
}

func (L *LuaState) ToBoolean(idx int) bool {
	val, _ := L.stackGet(idx)
	switch val := val.(type) {
	case nil:
		return false
	case bool:
		return val
	default:
		return true
	}
}

func (L *LuaState) ToStringX(idx int) (string, bool) {
	val, _ := L.stackGet(idx)
	switch val := val.(type) {
	case string:
		return val, true
	case float64, int64:
		str := fmt.Sprintf("%v", val)
		L.stackSet(idx, str)
		return str, true
	default:
		return "", false
	}
}

/**
 * push functions (Go -> stack)
 */

func (L *LuaState) PushNil() {
	L.stackPush(nil)
}

func (L *LuaState) PushNumber(n float64) {
	L.stackPush(n)
}

func (L *LuaState) PushInteger(n int64) {
	L.stackPush(n)
}

func (L *LuaState) PushString(s string) {
	L.stackPush(s)
}

func (L *LuaState) PushBoolean(b bool) {
	L.stackPush(b)
}

/**
 * some useful macros
 */

func (L *LuaState) ToNumber(idx int) float64 {
	val, _ := L.ToNumberX(idx)
	return val
}

func (L *LuaState) ToInteger(idx int) int64 {
	val, _ := L.ToIntegerX(idx)
	return val
}

func (L *LuaState) Pop(n int) {
	L.SetTop(-n - 1)
}

func (L *LuaState) IsNil(idx int) bool {
	return L.Type(idx) == LUA_TNIL
}

func (L *LuaState) IsBoolean(idx int) bool {
	return L.Type(idx) == LUA_TBOOLEAN
}

func (L *LuaState) IsNone(idx int) bool {
	return L.Type(idx) == LUA_TNONE
}

func (L *LuaState) IsNoneOrNil(idx int) bool {
	return L.Type(idx) <= LUA_TNIL
}

func (L *LuaState) ToString(idx int) string {
	val, _ := L.ToStringX(idx)
	return val
}

func (L *LuaState) Insert(idx int) {
	L.Rotate(idx, 1)
}

func (L *LuaState) Remove(idx int) {
	L.Rotate(idx, -1)
	L.Pop(1)
}

func (L *LuaState) Replace(idx int) {
	L.Copy(-1, idx)
	L.Pop(1)
}
