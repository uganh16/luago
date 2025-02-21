package binary

import (
	"fmt"
	"os"
)

const (
	LUA_SIGNATURE    = "\x1bLua"
	LUAC_VERSION     = 0x53
	LUAC_FORMAT      = 0
	LUAC_DATA        = "\x19\x93\r\n\x1a\n"
	CINT_SIZE        = 4
	CSIZET_SIZE      = 8
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE  = 8
	LUAC_INT         = 0x5678
	LUAC_NUM         = 370.5
)

const (
	LUA_TNIL     = 0x00
	LUA_TBOOLEAN = 0x01
	LUA_TNUMFLT  = 0x03
	LUA_TNUMINT  = 0x13
	LUA_TSHRSTR  = 0x04
	LUA_TLNGSTR  = 0x14
)

type binaryChunk struct {
	header
	sizeUpvalues byte
	mainFunc     *Prototype
}

type header struct {
	signature       [4]byte
	version         byte
	format          byte
	luacData        [6]byte
	cintSize        byte
	sizetSize       byte
	instructionSize byte
	luaIntegerSize  byte
	luaNumberSize   byte
	luacInt         int64
	luacNum         float64
}

type Prototype struct {
	Source          string // debug
	LineDefined     uint32
	LastLineDefined uint32
	NumParams       byte
	IsVararg        bool
	MaxStackSize    byte
	Code            []uint32
	Constants       []interface{}
	Upvalues        []Upvalue
	Protos          []*Prototype
	LineInfo        []uint32 // debug
	LocVars         []LocVar // debug
	UpvalueNames    []string // debug
}

type Upvalue struct {
	InStack byte
	Idx     byte
}

type LocVar struct {
	VarName string
	StartPC uint32
	EndPC   uint32
}

type bailout string

func Undump(file *os.File) (proto *Prototype, err error) {
	defer func() {
		switch x := recover().(type) {
		case nil:
			// no panic
		case bailout:
			err = fmt.Errorf("%s precompiled chunk", x)
		default:
			panic(x)
		}
	}()

	r := &reader{file}
	order := r.checkHeader()
	r.readByte() // size_upvalues
	proto = r.readProto(order, "")
	return
}
