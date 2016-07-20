package main

import (
	"fmt"
	"strings"
)

var uses = map[string]*Type{}

var mapper = &Mapper{Tx: "db"}

type Mapper struct {
	Source string
	Tx     string

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

	ParamsOrder  []*VarType
	ResultsOrder []*VarType

	Params  map[string]*Type
	Results map[string]*Type

	OpType string

	Fragments []*Fragment

	Dest   []*VarType
	Return *VarType
}

func (m *Operation) ParamsString() string {
	return paramsString(m.ParamsOrder, m.Params)
}

func (m *Operation) ResultsString() string {
	return paramsString(m.ResultsOrder, m.Results)
}

func paramsString(pos []*VarType, ps map[string]*Type) string {
	var slice []string
	for _, vt := range pos {
		slice = append(slice, vt.Var+" "+ps[vt.Var].String())
	}
	return strings.Join(slice, ",")
}

type Fragment struct {
	Cond string
	Stmt string
	Args []*VarType
}

type VarType struct {
	Var  string
	Type *Type

	IsIn bool
}

func (vt *VarType) UnderlineVar() string {
	return "x_" + strings.Replace(vt.Var, ".", "_", -1)
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
}

func (t *Type) String() string {
	var slice, star, pkg string
	if t.Slice {
		slice = "[]"
	}
	if t.Pointer {
		star = "*"
	}
	if t.Package != "" {
		pkg = t.Package + "."
	}
	if t.Name == "" {
		t.Name = t.Type
	}
	return fmt.Sprintf("%s%s%s%s", slice, star, pkg, t.Name)
}

func (t *Type) Alias() bool {
	return t.Primitive == t.Type
}

func (t *Type) IsPointer() bool {
	return t.Type[0] == '*'
}

func (t *Type) IsSlice() bool {
	return t.Type[0] == '['
}

func (t *Type) Elem() string {
	if t.Slice {
		return uses[t.Type[2:]].Elem()
	}
	if t.Pointer {
		return t.String()[1:]
	}
	return t.String()
}

func (t *Type) NewExpression() string {
	if t.Primitive != "" {
		return t.Type
	}

	if len(t.Type) > 4 && t.Type[:4] == "map[" {
		return t.Type + "{}"
	}

	return "&" + t.Elem() + "{}"
}
