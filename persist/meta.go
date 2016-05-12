package persist

type Interface struct {
	Name    string
	Imports []string
	Methods []*Func
}

type Func struct {
	Name    string
	Doc     string
	Params  []*VarAndType
	Returns []*VarAndType
}

// var ds []*domain.Demo
type VarAndType struct {
	Var    string // ds
	Scope  string // d for Status means d.Status
	Concat string // _ for status meas d_Stauts

	SQLIn bool // status in(${statuses})

	Type    string // Demo
	Alias   string // type Status int8
	Slice   string // []
	Star    string // *
	Package string // domain
	Path    string // github.com/arstd/persist/examples/domain

	Props []*VarAndType // Demo Props
}
