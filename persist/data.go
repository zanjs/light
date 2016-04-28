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

	Name       string
	Params     []*Param
	Returns    []*Param
	Result     string
	ResultType string

	Prefix    string
	Optionals []*Optional
	Suffix    string

	Args     []string
	Marshals []string

	Scans      []string
	Unmarshals []string
}

type Optional struct {
	Condition string
	Field     string
	Stmt      string
}
