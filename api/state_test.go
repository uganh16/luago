package api

import (
	"fmt"
	"testing"
)

func TestStack(t *testing.T) {
	L := NewState()
	if len(L.stack) != 0 {
		t.Errorf("Empty stack expected: %v", L.stack)
	}
	L.PushBoolean(true)
	printStack(L)
	L.PushInteger(10)
	printStack(L)
	L.PushNil()
	printStack(L)
	L.PushString("hello")
	printStack(L)
	L.PushValue(-4)
	printStack(L)
	L.Replace(3)
	printStack(L)
	L.SetTop(6)
	printStack(L)
	L.Remove(-3)
	printStack(L)
	L.SetTop(-5)
	printStack(L)
}

func TestLuaOp(t *testing.T) {
	L := NewState()
	L.PushInteger(1)
	L.PushString("2.0")
	L.PushString("3.0")
	L.PushNumber(4.0)
	printStack(L)

	L.Arith(LUA_OPADD)
	printStack(L)
	L.Arith(LUA_OPBNOT)
	printStack(L)
	L.Len(2)
	printStack(L)
	L.Concat(3)
	printStack(L)
	L.PushBoolean(L.Compare(1, 2, LUA_OPEQ))
	printStack(L)
}

func printStack(L *LuaState) {
	for idx := 1; idx <= len(L.stack); idx++ {
		t := L.Type(idx)
		switch t {
		case LUA_TBOOLEAN:
			fmt.Printf("[%t]", L.ToBoolean(idx))
		case LUA_TNUMBER:
			fmt.Printf("[%g]", L.ToNumber(idx))
		case LUA_TSTRING:
			fmt.Printf("[%q]", L.ToString(idx))
		default:
			fmt.Printf("[%s]", L.TypeName(t))
		}
	}
	fmt.Println()
}
