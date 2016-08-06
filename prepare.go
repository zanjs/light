package main

import (
	"regexp"
	"strings"

	"github.com/wothing/log"
)

func prepareData() {
	mapper.Name += "Impl"

	if *path != "" {
		mapper.Imports[*path] = `"` + *path + `"`
	}

	for _, m := range mapper.Methods {
		if m.Tx == "" {
			m.Tx = *db
		}

		log.Debug("getOpType")
		m.OpType = getOpType(m)

		log.Debug("checkResults")
		checkResults(m)

		log.Debug("calcDest")
		calcDest(m)

		log.Debug("calcFragment")
		calcFragment(m)
	}
}

var brackets = regexp.MustCompile(`\[.+? \]`)

func calcFragment(m *Operation) {
	if strings.Count(m.Doc, "[") != strings.Count(m.Doc, "]") {
		panic("brackets `[` and `]` not matched")
	}

	matched := brackets.FindAllStringIndex(m.Doc, -1)

	if len(matched) == 0 {
		m.Fragments = append(m.Fragments, calcArgs(m, m.Doc))
		return
	}

	var last int
	for _, a := range matched {
		stmt := strings.TrimSpace(m.Doc[last:a[0]])
		if stmt != "" {
			m.Fragments = append(m.Fragments, calcArgs(m, stmt))
		}
		stmt = strings.TrimSpace(m.Doc[a[0]:a[1]])
		if stmt != "" {
			m.Fragments = append(m.Fragments, calcArgs(m, stmt))
		}
		last = a[1]
	}
	if last != len(m.Doc) {
		stmt := strings.TrimSpace(m.Doc[last:])
		if stmt != "" {
			m.Fragments = append(m.Fragments, calcArgs(m, stmt))
		}
	}
}

var braces = regexp.MustCompile(`\?\{.+?\}`)
var variable = regexp.MustCompile(`\$\{.+?\}`)
var in = regexp.MustCompile(`(?:in\s+|array)[\[\(]\s*\$\{(.+?)\}\s*[\]\)]`)

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
	stmt = fragment.Stmt
	for _, a := range matched {
		s := fragment.Stmt[a[0]:a[1]]
		f := strings.TrimSpace(s[2 : len(s)-1])
		// log.Debug(f)

		ss := strings.Split(f, ".")
		t, ok := m.Params[ss[0]]
		if !ok {
			panic("variable `" + ss[0] + "` not in parameters for " + m.Name)
		}
		if len(ss) != 1 {
			p, ok := t.Fields[ss[1]]
			if !ok {
				panic("variable `" + f + "` not in parameters for " + m.Name)
			}
			t = p
		}

		var isIn bool
		for _, s := range ins {
			if s == f {
				isIn = true
				break
			}
		}
		if !isIn {
			stmt = strings.Replace(stmt, s, "%s", -1)
		}

		fragment.Args = append(fragment.Args, &VarType{f, t, isIn})
	}
	fragment.Stmt = stmt

	return
}

var parentheses = regexp.MustCompile(`\([^\(]*?\)`)

func calcDest(m *Operation) {
	log.Debug("get dest stmt")

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
			tmp := strings.Index(m.Doc[i:], " from ")
			if tmp == -1 {
				panic(m.Name + " method sql keyword `from` not found")
			}
			i += tmp
			if strings.Count(m.Doc[:i], "(") == strings.Count(m.Doc[:i], ")") {
				break
			}
			i += len(" from ")
		}
		stmt = m.Doc[len("select")+1 : i]
	}

	stmt = parentheses.ReplaceAllString(stmt, "")
	stmt = parentheses.ReplaceAllString(stmt, "")
	stmt = parentheses.ReplaceAllString(stmt, "")

	fs := strings.Split(stmt, ",")
	for _, s := range fs {
		s = strings.TrimSpace(s)
		f := s[strings.LastIndexAny(s, " \t")+1:]
		f = strings.Replace(f, "_", " ", -1)
		f = strings.Title(f)
		f = strings.Replace(f, " ", "", -1)

		if m.OpType == "insert" {
			for v, t := range m.Params {
				if t.Name == "Tx" {
					m.Tx = v
					continue
				}

				if p, ok := t.Fields[f]; ok {
					m.Dest = append(m.Dest, &VarType{v + "." + f, p, false})
				} else {
					panic("returning field `" + s + "` not matched struct property")
				}
			}
		} else {
			for v, t := range m.Results {
				if v == "err" {
					t.Primitive = "error"
					continue
				}
				if t.Name == "Tx" {
					m.Tx = v
					continue
				}

				if f == "Count" || f == "Sum" {
					m.Dest = append(m.Dest, &VarType{v, t, false})
					break
				}

				if s == "photo_status" {
					log.Debug(f)
					log.JSONIndent(t)
				}

				if m.OpType == "get" {
					if p, ok := t.Fields[f]; ok {
						m.Dest = append(m.Dest, &VarType{v + "." + f, p, false})
					} else {
						panic("select field `" + s + "` not matched struct property for " + m.Name)
					}
				} else if m.OpType == "list" {
					if p, ok := t.Fields[f]; ok {
						m.Dest = append(m.Dest, &VarType{"x." + f, p, false})
					} else {
						panic("select field `" + s + "` not matched struct property for " + m.Name)
					}
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

	if m.OpType == "get" || m.OpType == "list" {
		for _, vt := range m.ResultsOrder {
			if vt.Var == "err" {
				continue
			}

			m.Return = vt
		}
	}
}

func getOpType(m *Operation) string {
	op := m.Doc[:strings.IndexAny(m.Doc, " \t")]
	switch op {
	case "insert":
		i := strings.LastIndex(m.Doc, " returning ")
		if i == -1 {
			return "update"
		}
		return op

	case "update", "delete":
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
