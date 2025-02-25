package api

import (
	"fmt"
	"math"

	"github.com/uganh16/luago/number"
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
	return toNumber(val)
}

func (L *LuaState) ToIntegerX(idx int) (int64, bool) {
	val, _ := L.stackGet(idx)
	return toInteger(val)
}

func (L *LuaState) ToBoolean(idx int) bool {
	val, _ := L.stackGet(idx)
	return toBoolean(val)
}

func (L *LuaState) ToStringX(idx int) (string, bool) {
	val, _ := L.stackGet(idx)
	if str, ok := val.(string); ok {
		return str, ok
	} else if str, ok = toString(val); ok {
		L.stackSet(idx, str)
		return str, ok
	}
	return "", false
}

/**
 * comparison and arithmetic functions
 */

const (
	LUA_OPADD  = iota // +
	LUA_OPSUB         // -
	LUA_OPMUL         // *
	LUA_OPMOD         // %
	LUA_OPPOW         // ^
	LUA_OPDIV         // /
	LUA_OPIDIV        // //
	LUA_OPBAND        // &
	LUA_OPBOR         // |
	LUA_OPBXOR        // ~
	LUA_OPSHL         // <<
	LUA_OPSHR         // >>
	LUA_OPUNM         // - (unary minus)
	LUA_OPBNOT        // ~
)

type ArithOp = int

func (L *LuaState) Arith(op ArithOp) {
	var a, b, r luaValue
	b = L.stackPop()
	if op != LUA_OPUNM && op != LUA_OPBNOT {
		a = L.stackPop()
	} else {
		a = b
	}

	var iFunc func(int64, int64) int64
	var fFunc func(float64, float64) float64

	switch op {
	case LUA_OPADD:
		iFunc = func(a, b int64) int64 { return a + b }
		fFunc = func(a, b float64) float64 { return a + b }
	case LUA_OPSUB:
		iFunc = func(a, b int64) int64 { return a - b }
		fFunc = func(a, b float64) float64 { return a - b }
	case LUA_OPMUL:
		iFunc = func(a, b int64) int64 { return a * b }
		fFunc = func(a, b float64) float64 { return a * b }
	case LUA_OPMOD:
		iFunc = number.IMod
		fFunc = number.FMod
	case LUA_OPPOW:
		fFunc = math.Pow
	case LUA_OPDIV:
		fFunc = func(a, b float64) float64 { return a / b }
	case LUA_OPIDIV:
		iFunc = number.IFloorDiv
		fFunc = number.FFloorDiv
	case LUA_OPBAND:
		iFunc = func(a, b int64) int64 { return a & b }
	case LUA_OPBOR:
		iFunc = func(a, b int64) int64 { return a | b }
	case LUA_OPBXOR:
		iFunc = func(a, b int64) int64 { return a ^ b }
	case LUA_OPSHL:
		iFunc = number.ShiftLeft
	case LUA_OPSHR:
		iFunc = number.ShiftRight
	case LUA_OPUNM:
		iFunc = func(a, _ int64) int64 { return -a }
		fFunc = func(a, _ float64) float64 { return -a }
	case LUA_OPBNOT:
		iFunc = func(a, _ int64) int64 { return ^a }
	default:
		panic(fmt.Sprintf("invalid arith op: %d", op))
	}

	if fFunc == nil { // bitwise operation
		if a, ok := toInteger(a); ok {
			if b, ok := toInteger(b); ok {
				r = iFunc(a, b)
			}
		}
	} else {
		if iFunc != nil {
			if a, ok := a.(int64); ok {
				if b, ok := b.(int64); ok {
					r = iFunc(a, b)
				}
			}
		}

		if r == nil {
			if a, ok := toNumber(a); ok {
				if b, ok := toNumber(b); ok {
					r = fFunc(a, b)
				}
			}
		}
	}

	if r != nil {
		L.stackPush(r)
	} else {
		switch op {
		case LUA_OPBAND, LUA_OPBOR, LUA_OPBXOR, LUA_OPSHL, LUA_OPSHR:
			_, ok1 := toNumber(a)
			_, ok2 := toNumber(b)
			if ok1 && ok2 {
				panic(runtimeError("number has no integer representation")) // @todo
			} else {
				if !ok1 {
					b = a
				}
				panic(typeError(L, b, "perform bitwise operation on"))
			}
		default:
			if _, ok := toNumber(a); !ok {
				b = a
			}
			panic(typeError(L, b, "perform arithmetic on"))
		}
	}
}

const (
	LUA_OPEQ = iota // ==
	LUA_OPLT        // <
	LUA_OPLE        // <=
)

type CompareOp = int

func (L *LuaState) Compare(idx1, idx2 int, op CompareOp) bool {
	a, ok1 := L.stackGet(idx1)
	b, ok2 := L.stackGet(idx2)
	if !ok1 || !ok2 {
		return false
	}
	switch op {
	case LUA_OPEQ:
		return equal(a, b)
	case LUA_OPLT:
		return lessThan(L, a, b)
	case LUA_OPLE:
		return lessEqual(L, a, b)
	default:
		panic(fmt.Sprintf("invalid compare op: %d", op))
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
 * miscellaneous functions
 */

func (L *LuaState) Concat(n int) {
	if n == 0 {
		L.stackPush("")
	}
	if n >= 2 {
		b := L.stackPop()
		for n > 1 {
			a := L.stackPop()
			n--
			if s1, ok := toString(a); ok {
				if s2, ok := toString(b); ok {
					b = s1 + s2
					continue
				}
			}
			// @todo
			if _, ok := toString(a); ok {
				a = b
			}
			panic(typeError(L, a, "concatenate"))
		}
		L.stackPush(b)
	}
}

func (L *LuaState) Len(idx int) {
	val, _ := L.stackGet(idx)
	if str, ok := val.(string); ok {
		L.stackPush(int64(len(str)))
	} else {
		typeError(L, val, "get length of")
	}
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
