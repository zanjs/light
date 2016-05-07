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
	MethodTypeGet    = "get"
	MethodTypeList   = "list"
	MethodTypeRemove = "remove"
	MethodTypePage   = "page"
)

type Method struct {
	Type MethodType

	Name    string
	Params  []*Param
	Returns []*Param

	In         string
	Result     string
	ResultType string

	Prefix   string
	Suffix   string
	Args     []string
	Marshals []string

	Fragments []*Fragment

	Scans       []string
	Unmarshals  []string
	Unmarshals1 []*Param
}

type Fragment struct {
	Cond string // if xxx != yyy
	Args []*VarAndType
	Stmt string
}
