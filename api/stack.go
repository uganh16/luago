package api

func (L *LuaState) stackPush(val luaValue) {
	if len(L.stack) == cap(L.stack) {
		panic("stack overflow")
	}
	L.stack = append(L.stack, val)
}

func (L *LuaState) stackPop() luaValue {
	top := len(L.stack)
	if top == 0 {
		panic("not enough elements in the stack")
	}
	top--
	val := L.stack[top]
	L.stack[top] = nil
	L.stack = L.stack[:top]
	return val
}

func (L *LuaState) stackGet(idx int) (luaValue, bool) {
	idx = L.AbsIndex(idx)
	if 0 < idx && idx <= cap(L.stack) {
		if idx <= len(L.stack) {
			return L.stack[idx-1], true
		}
		return nil, false
	}
	panic("unacceptable index")
}

func (L *LuaState) stackSet(idx int, val luaValue) {
	idx = L.AbsIndex(idx)
	if 0 < idx && idx <= len(L.stack) {
		L.stack[idx-1] = val
		return
	}
	panic("invalid index")
}

func (L *LuaState) stackReverse(from, to int) {
	for from < to {
		L.stack[from], L.stack[to] = L.stack[to], L.stack[from]
		from++
		to--
	}
}
