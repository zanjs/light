package main

import (
	"regexp"
	"strings"

	"github.com/wothing/log"
)

func prepareData() {
	for _, m := range mapper.Methods {
		m.OpType = getOpType(m)

		checkResults(m)

		calcDest(m)

		calcFragment(m)
	}
}

var brackets = regexp.MustCompile(`\[.+?\]`)

func calcFragment(m *Operation) {
	matched := brackets.FindAllStringIndex(m.Doc, -1)

	if len(matched) == 0 {
		m.Fragments = append(m.Fragments, calcArgs(m, m.Doc))
		return
	}

	var last int
	for _, a := range matched {
		stmt := m.Doc[last:a[0]]
		if stmt != "" {
			m.Fragments = append(m.Fragments, calcArgs(m, stmt))
		}
		m.Fragments = append(m.Fragments, calcArgs(m, m.Doc[a[0]:a[1]]))

		last = a[1]
	}
	if last != len(m.Doc) {
		m.Fragments = append(m.Fragments, calcArgs(m, m.Doc[last:]))
	}
}

var braces = regexp.MustCompile(`\?\{.+?\}`)
var variable = regexp.MustCompile(`\$\{.+?\}`)
var in = regexp.MustCompile(`in\s+\(\s*\$\{(.+?)\}\s*\)`)

func calcArgs(m *Operation, stmt string) (fragment *Fragment) {
	fragment = &Fragment{}

	stmt = strings.TrimSpace(stmt)
	mb := braces.FindStringIndex(stmt)

	if len(mb) == 0 {
		fragment.Stmt = stmt
	} else {
		fragment.Cond = stmt[mb[0]+2 : mb[1]-1]
		fragment.Stmt = strings.TrimSpace(stmt[mb[1]+1:len(stmt)-1]) + " "
	}

	var ins []string
	matchedIns := in.FindAllStringSubmatchIndex(fragment.Stmt, -1)
	for _, a := range matchedIns {
		ins = append(ins, fragment.Stmt[a[2]:a[3]])
	}

	matched := variable.FindAllStringIndex(fragment.Stmt, -1)
	for _, a := range matched {

		f := fragment.Stmt[a[0]:a[1]]
		log.Debug(f)
		f = strings.TrimSpace(f[2 : len(f)-1])
		log.Debug(f)

		ss := strings.Split(f, ".")
		t, ok := m.Params[ss[0]]
		if !ok {
			panic("variable ⟦" + ss[0] + "⟧ not in parameters for " + m.Name)
		}
		if len(ss) == 1 {
			var isIn bool
			for _, s := range ins {
				if s == f {
					isIn = true
					break
				}
			}
			fragment.Args = append(fragment.Args, &VarType{f, t, isIn})
		} else {
			p, ok := t.Fields[ss[1]]
			if !ok {
				panic("variable ⟦" + f + "⟧ not in parameters for " + m.Name)
			}
			var isIn bool
			for _, s := range ins {
				if s == f {
					isIn = true
					break
				}
			}
			fragment.Args = append(fragment.Args, &VarType{f, p, isIn})
		}
	}

	return
}

var parentheses = regexp.MustCompile(`\(.*\)`)

func calcDest(m *Operation) {
	var stmt string
	switch m.OpType {
	case "update", "delete":
		return

	case "insert":
		i := strings.LastIndex(m.Doc, " returning ")
		stmt = m.Doc[i+len(" returning "):]

	case "get", "list":
		var i int
		for {
			i = strings.Index(m.Doc, " from ")
			if i == -1 {
				panic(m.Name + " method sql keyword `from` not found")
			}
			if strings.Contains(m.Doc[:i], "(") == strings.Contains(m.Doc[:i], ")") {
				break
			}
		}
		stmt = m.Doc[len("select")+1 : i]
	}

	stmt = parentheses.ReplaceAllString(stmt, "")

	fs := strings.Split(stmt, ",")
	for _, s := range fs {
		f := strings.TrimSpace(s)
		f = strings.Replace(f, "_", " ", -1)
		f = strings.Title(f)
		f = strings.Replace(f, " ", "", -1)

		if m.OpType == "insert" {
			for _, t := range m.Params {
				if t.Name == "Tx" {
					continue
				}

				if p, ok := t.Fields[f]; ok {
					m.Dest = append(m.Dest, &VarType{f, p, false})
				} else {
					panic("returning field ⟦" + s + "⟧ not matched struct property")
				}
			}
		} else {
			for v, t := range m.Results {
				if v == "err" {
					continue
				}

				if f == "Count" || f == "Sum" {
					m.Dest = append(m.Dest, &VarType{strings.ToLower(f), t, false})
					break
				}

				if p, ok := t.Fields[f]; ok {
					m.Dest = append(m.Dest, &VarType{f, p, false})
				} else {
					panic("select field ⟦" + s + "⟧ not matched struct property")
				}
			}
		}
	}

}

func checkResults(m *Operation) {
	if m.OpType == "insert" {
		if len(m.Results) != 1 {
			panic(m.Name + " method must has one return only")
		}
	} else {
		if len(m.Results) != 2 {
			panic(m.Name + " method must has two returns only")
		}
	}
}

func getOpType(m *Operation) string {
	op := m.Doc[:strings.IndexAny(m.Doc, " \t")]
	switch op {
	case "insert", "update", "delete":
		return op

	case "select":
		for v, t := range m.Results {
			if v == "err" {
				continue
			}
			if t.Slice {
				return "list"
			} else {
				return "get"
			}
		}
	}
	panic("sql error, keyword not found: " + m.Doc)
}
