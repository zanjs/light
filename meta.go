package main

import "strings"

var uses = map[string]*Type{}

var m = &Mapper{}

type Mapper struct {
	Source string

	Path    string
	Package string

	Imports map[string]string

	Name string

	Methods []*Operation
}

type Operation struct {
	Doc string // select ... from ...

	Name string
	Tx   string

	Params  map[string]*Type
	Results map[string]*Type

	Type string

	Fragments []*Fragment

	Dest       []*Type
	ReturnType []*Type
}

func (o *Operation) Signature() string {
	return o.Name + "(" + o.ParamsCode() + ")(" + o.ResultsCode() + ")"
}

func (o *Operation) ParamsCode() string {
	var codes []string
	for v, t := range o.Params {
		codes = append(codes, v+t.String())
	}
	return strings.Join(codes, ", ")
}

func (o *Operation) ResultsCode() string {
	var codes []string
	for v, t := range o.Results {
		codes = append(codes, v+t.String())
	}
	return strings.Join(codes, ", ")
}

type Type struct {
	Type string

	Name      string
	Slice     bool
	Pointer   bool
	Package   string
	Path      string
	Primitive string

	Fields map[string]*Type

	IsIn bool
}

func (t *Type) String() string {
	return t.Type
}

func (t *Type) Alias() string {
	if t.Primitive == t.Type {
		return ""
	}
	return t.Primitive
}

func (t *Type) IsPointer() bool {
	return t.Type[0] == '*'
}

func (t *Type) IsSlice() bool {
	return t.Type[0] == '['
}

func (t *Type) Elem() string {
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

type Fragment struct {
	Cond string
	Stmt string
	Args []*Type
}
