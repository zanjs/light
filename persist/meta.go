package persist

import (
	"strconv"
	"strings"

	"github.com/wothing/log"
)

type ItfData struct {
	Name  string      // DemoPersisterData
	Funcs []*FuncData // {Get,Modify,Add,List,Remove}
}

type FuncData struct {
	SQL     string        // select * from demo
	Params  []*VarAndType // {id: int}
	Returns []*VarAndType // {d: *domain.Demo, err: error}
	// internal
	Name       string // Get
	Type       string // add/modify/get/list/remove
	Stmt       string
	Args       []string
	Result     string
	ResultType string
	Scans      []string
	Returning  []string
}

type VarAndType struct {
	Var  string
	Type string
}

func (meta *ItfData) Explain() {
	for _, f := range meta.Funcs {
		// parse variabled sql to prepared sql
		matched := re.FindAllStringIndex(f.SQL, -1)
		if len(matched) > 0 {
			start, s, e := 0, 0, 0
			for i := 0; i < len(matched); i++ {
				s, e = matched[i][0], matched[i][1]
				f.Stmt += f.SQL[start:s] + "$" + strconv.Itoa(i+1)
				f.Args = append(f.Args, f.SQL[s+2:e-1])

				start = e
			}
			f.Stmt += f.SQL[start:]
		}

		// fmt.Println(f.Stmt)

		// classify methods
		switch f.SQL[:strings.Index(f.SQL, " ")] {
		case "insert":
			f.Type = "add"
			f.Returning = []string{f.Params[0].Var + ".Id"}
		case "update":
			f.Type = "modify"
		case "delete":
			f.Type = "remove"
		case "select":
			if f.Returns[0].Type[:3] == "[]*" {
				f.Type = "list"
				f.ResultType = f.Returns[0].Type[3:]
			} else {
				f.Type = "get"
				f.ResultType = f.Returns[0].Type[1:]
			}

			fieldsSQl := f.SQL[7 : strings.Index(f.SQL, "from")-1]
			fields := strings.Split(fieldsSQl, ",")
			for _, fi := range fields {
				fi = strings.TrimSpace(fi)
				fi = strings.Replace(fi, "_", " ", -1)
				fi = strings.Title(fi)
				fi = strings.Replace(fi, " ", "", -1)
				f.Scans = append(f.Scans, fi)
			}
			f.Result = f.Returns[0].Var

		default:
			log.Errorf("unsupported for `%s`", f.SQL)
		}
	}
}
