package api

import "fmt"

type runtimeError string

func typeError(L *LuaState, val luaValue, op string) runtimeError {
	t := L.TypeName(typeOf(val)) // @todo
	return runtimeError(fmt.Sprintf("attempt to %s a %s value", op, t))
}

func orderError(L *LuaState, a, b luaValue) runtimeError {
	t1 := L.TypeName(typeOf(a)) // @todo
	t2 := L.TypeName(typeOf(b)) // @todo
	if t1 == t2 {
		return runtimeError(fmt.Sprintf("attempt to compare two %s values", t1))
	} else {
		return runtimeError(fmt.Sprintf("attempt to compare %s with %s", t1, t2))
	}
}
