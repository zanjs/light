package main

import (
	"fmt"
	"strings"
)

var uses = map[string]*Type{}

var mapper = &Mapper{}

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

	ParamsOrder  []*VarType
	ResultsOrder []*VarType

	Params  map[string]*Type
	Results map[string]*Type

	OpType string

	Fragments []*Fragment

	Dest   []*VarType
	Return *VarType
}

func (m *Operation) ParamsExpr() string {
	return paramsExpr(m.ParamsOrder, m.Params)
}

func (m *Operation) ResultsExpr() string {
	return paramsExpr(m.ResultsOrder, m.Results)
}

func paramsExpr(pos []*VarType, ps map[string]*Type) string {
	var slice []string
	for _, vt := range pos {
		slice = append(slice, vt.Var+" "+ps[vt.Var].String())
	}
	return strings.Join(slice, ", ")
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

// String 声明表达式，如 s *a.Bc
func (vt *VarType) String() string {
	return vt.Var + " " + vt.Type.String()
}

// AddrExpr 取该变量的地址，如 &s
func (vt *VarType) AddrExpr() string {
	if !vt.Type.Slice && vt.Type.Pointer {
		return vt.Var
	}
	return "&" + vt.Var // primitive/struct/map/slice
}

// AnotherVar 声明另一个变量，如 xa_Bc, err := json.Marshal(a.Bc)
func (vt *VarType) AnotherVar() string {
	return "x" + strings.Replace(vt.Var, ".", "_", -1)
}

// Type 类型详细描述，如果是 struct 结构，列出其属性
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

// String 声明表达式，如 var s []*a.Bc
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
		t.Name = t.Primitive
	}
	return fmt.Sprintf("%s%s%s%s", slice, star, pkg, t.Name)
}

// MakeExpr 创建表达式，如 s = []*a.Bc{}
func (t *Type) MakeExpr() string {
	// map 特殊处理
	if len(t.Type) >= 3 && t.Type[:3] == "map" {
		return t.Type + "{}"
	}

	var slice, star, pkg string
	if t.Slice {
		slice = "[]"
		if t.Pointer {
			star = "*"
		}
	} else if t.Pointer {
		star = "&"
	}
	if t.Package != "" {
		pkg = t.Package + "."
	}
	if t.Name == "" {
		t.Name = t.Primitive
	}
	return fmt.Sprintf("%s%s%s%s{}", slice, star, pkg, t.Name)
}

// MakeExpr 是否是复杂类型（SQL 不支持的类型，需要转换）
func (t *Type) IsComplex() bool {
	// map 特殊处理
	if len(t.Type) >= 3 && t.Type[:3] == "map" {
		return true
	}

	if t.Slice {
		return true
	}
	if t.Pointer {
		return true
	}
	if len(t.Type) >= 3 && t.Type[:3] == "map" {
		return true
	}
	if t.Type == "time.Time" {
		return false
	}
	return t.Primitive == ""
}

// AliasFor 其他类型的别名，type Status int8
func (t *Type) AliasFor() string {
	if t.IsComplex() {
		return ""
	}

	return t.Primitive
}

// Elem 数组元素的类型
func (t *Type) Elem() *Type {
	if t.Slice {
		return uses[t.Type[2:]]
	}
	return t
}
