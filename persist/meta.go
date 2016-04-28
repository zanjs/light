package persist

type Interface struct {
	Name    string
	Imports []string
	Methods []*Func
}

type Func struct {
	Name    string
	Doc     string
	Params  []*Param
	Returns []*Param
}

type Param struct {
	Name string
	Type string
}
