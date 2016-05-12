package persist

type Implement struct {
	Source  string
	Package string
	Imports []string

	Name    string
	Methods []*Method
}

type MethodType string

const (
	MethodTypeAdd    = "add"
	MethodTypeModify = "modify"
	MethodTypeRemove = "remove"
	MethodTypeGet    = "get"
	MethodTypeCount  = "count"
	MethodTypeList   = "list"
)

type Method struct {
	Type MethodType

	Name    string
	Params  []*VarAndType
	Returns []*VarAndType

	Marshals  []*VarAndType
	Fragments []*Fragment

	Result     *VarAndType
	Scans      []*VarAndType
	Unmarshals []*VarAndType
}

type Fragment struct {
	Cond string // if xxx != yyy
	Args []*VarAndType
	Stmt string
}
