package main

import "strings"

type Mapper struct {
	Source string

	Path    string
	Package string

	Imports []string

	Name string

	Methods []*Operation
}

type Operation struct {
	Name string
	Tx   string

	Params  []*VarType
	Returns []*VarType

	Type string

	Fragments []*Fragment

	Dest       []*VarType
	ReturnType []*VarType
}

func (o *Operation) Signature() string {
	return o.Name + "(" + o.ParamsCode() + ")(" + o.ReturnsCode() + ")"
}

func (o *Operation) ParamsCode() string {
	var codes []string
	for _, vt := range o.Params {
		codes = append(codes, vt.String())
	}
	return strings.Join(codes, ", ")
}

func (o *Operation) ReturnsCode() string {
	var codes []string
	for _, vt := range o.Returns {
		codes = append(codes, vt.String())
	}
	return strings.Join(codes, ", ")
}

type VarType struct {
	Var  string
	Type string

	IsIn bool
}

func (t *VarType) UnderlineVar() string {
	return strings.Replace(t.Var, ".", "_", -1)
}

func (t *VarType) String() string {
	return t.Var + " " + t.Type
}

func (t *VarType) IsPointer() bool {
	return t.Type[0] == '*'
}

func (t *VarType) IsSlice() bool {
	return t.Type[0] == '['
}

func (t *VarType) Elem() string {
	if t.Type[0] == '[' {
		if len(t.Type) < 2 {
			panic("type " + t.Type + " not supported")
		}
		if t.Type[1] != ']' {
			panic("type " + t.Type + " not supported")
		}
		if t.Type[2] == '*' {
			return t.Type[3:]
		}
		return t.Type[2:]
	}

	if t.Type[2] == '*' {
		return t.Type[3:]
	}
	return t.Type[2:]
}

func (t *VarType) IsPrimitive() bool {
	switch t.Type {
	case "bool", "byte", "error", "float32", "float64", "int", "int16",
		"int32", "int64", "int8", "rune", "string", "uint", "uint16",
		"uint32", "uint64", "uint8":
		return true

	default:
		return false
	}
}

type Fragment struct {
	Cond string
	Stmt string
	Args []*VarType
}
