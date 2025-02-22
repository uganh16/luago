package main

import (
	"fmt"
	"os"

	"github.com/uganh16/luago/binary"
)

func main() {
	for _, file := range os.Args[1:] {
		var p *binary.Prototype
		f, err := os.Open(file)
		if err == nil {
			p, err = binary.Undump(f)
			f.Close()
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", file, err)
			continue
		}
		list(p)
	}
}

func list(p *binary.Prototype) {
	printHeader(p)
	printCode(p)
	printDebug(p)
	for _, p := range p.Protos {
		list(p)
	}
}

func printHeader(p *binary.Prototype) {
	funcType := "main"
	if p.LineDefined > 0 {
		funcType = "function"
	}

	source := p.Source
	if source == "" {
		source = "=?"
	}
	if source[0] == '@' || source[0] == '=' {
		source = source[1:]
	} else if source[0] == binary.LUA_SIGNATURE[0] {
		source = "(bstring)"
	} else {
		source = "(string)"
	}

	varargFlag := ""
	if p.IsVararg {
		varargFlag = "+"
	}

	fmt.Printf("\n%s <%s:%d,%d> (%d instruction%s)\n", funcType, source, p.LineDefined, p.LastLineDefined, len(p.Code), ss(len(p.Code)))
	fmt.Printf("%d%s param%s, %d slot%s, %d upvalue%s, %d local%s, %d constant%s, %d function%s\n", p.NumParams, varargFlag, ss(int(p.NumParams)), p.MaxStackSize, ss(int(p.MaxStackSize)), len(p.Upvalues), ss(len(p.Upvalues)), len(p.LocVars), ss(len(p.LocVars)), len(p.Constants), ss(len(p.Constants)), len(p.Protos), ss(len(p.Protos)))
}

func printCode(p *binary.Prototype) {
	for pc, i := range p.Code {
		line := "-"
		if len(p.LineInfo) > pc {
			line = fmt.Sprintf("%d", p.LineInfo[pc])
		}
		fmt.Printf("\t%d\t[%s]\t0x%08X\n", pc+1, line, i)
	}
}

func printDebug(p *binary.Prototype) {
	fmt.Printf("constants (%d):\n", len(p.Constants))
	for i, k := range p.Constants {
		s := "?"
		switch k.(type) {
		case nil:
			s = "nil"
		case bool:
			s = fmt.Sprintf("%t", k)
		case int64:
			s = fmt.Sprintf("%d", k)
		case float64:
			s = fmt.Sprintf("%g", k)
		case string:
			s = fmt.Sprintf("%q", k)
		}
		fmt.Printf("\t%d\t%s\n", i+1, s)
	}

	fmt.Printf("locals (%d):\n", len(p.LocVars))
	for i, locVar := range p.LocVars {
		fmt.Printf("\t%d\t%s\t%d\t%d\n", i, locVar.VarName, locVar.StartPC+1, locVar.EndPC+1)
	}

	fmt.Printf("upvalues (%d):\n", len(p.Upvalues))
	for i, upvalue := range p.Upvalues {
		upvalueName := "-"
		if len(p.UpvalueNames) > 0 {
			upvalueName = p.UpvalueNames[i]
		}
		fmt.Printf("\t%d\t%s\t%d\t%d\n", i, upvalueName, upvalue.InStack, upvalue.Idx)
	}
}

func ss(n int) string {
	if n != 1 {
		return "s"
	}
	return ""
}
