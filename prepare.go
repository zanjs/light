package main

import (
	"os"
	"regexp"
	"strings"
)

func prepare(itf *Interface) (impl *Implement) {
	impl = &Implement{}
	impl.Source = os.Getenv("GOFILE")
	impl.Package = os.Getenv("GOPACKAGE")

	if path != "" {
		itf.Imports = append(itf.Imports, path)
	}
	for _, imp := range itf.Imports {
		if imp == "fmt" || imp == "encoding/json" || imp == "strconv" {
			continue
		}
		impl.Imports = append(impl.Imports, imp)
	}

	impl.Name = itf.Name
	if strings.HasSuffix(impl.Name, "Persister") {
		impl.Name = impl.Name[:len(impl.Name)-2]
	} else {
		impl.Name += impl.Name + "Impl"
	}

	for _, f := range itf.Methods {
		var m Method
		impl.Methods = append(impl.Methods, &m)

		prepareMethod(&m, f)
	}

	return impl
}

func prepareMethod(m *Method, f *Func) {
	m.Name = f.Name
	m.Params = f.Params
	m.Tx = calcTx(m.Params)

	m.Returns = f.Returns

	m.Type = getMethodType(f.Doc, f)

	calcInResult(m, f)

	calcFragment(m, f.Doc)

	calcMarshals(m)

	calcScans(m, f.Doc)

	calcUnmarshals(m)
}

func calcTx(ps []*VarAndType) string {
	for _, vt := range ps {
		if vt.Type == "Tx" {
			return vt.Var
		}
	}

	return db
}

func calcInResult(m *Method, f *Func) {
	switch m.Type {
	case MethodTypeAdd, MethodTypeModify, MethodTypeRemove:
		m.Result = f.Params[0]
	case MethodTypeGet, MethodTypeList:
		m.Result = f.Returns[0]
	default:
	}
}

func getMethodType(sql string, f *Func) MethodType {
	switch sql[:strings.Index(sql, " ")] {
	case "insert":
		return MethodTypeAdd

	case "update":
		return MethodTypeModify

	case "delete":
		return MethodTypeRemove

	case "select":
		if len(f.Returns) != 2 {
			panic("select method must return two result")
		}
		if f.Returns[0].Slice != "" {
			return MethodTypeList
		} else {
			countSQL := sql[strings.Index(sql, " "):]
			countSQL = strings.TrimSpace(countSQL)
			if strings.HasPrefix(countSQL, "count") {
				return MethodTypeCount
			} else if strings.HasPrefix(countSQL, "sum") {
				return MethodTypeCount
			}

			return MethodTypeGet
		}

	default:
	}
	panic("sql error: " + sql)
}

var fregmentRegexp = regexp.MustCompile(`\[\?\{(.+?)\}(.+?) \]`)

func calcFragment(m *Method, sql string) {
	matched := fregmentRegexp.FindAllStringSubmatchIndex(sql, -1)

	dollar := new(int)

	if len(matched) == 0 {
		fm := getFragment(m, sql, "", dollar)
		if fm != nil {
			m.Fragments = append(m.Fragments, fm)
		}
		return
	}

	before := 0
	for _, group := range matched {
		if before != group[0] {
			fm := getFragment(m, sql[before:group[0]], "", dollar)
			if fm != nil {
				m.Fragments = append(m.Fragments, fm)
			}
		}
		fm := getFragment(m, sql[group[4]:group[5]], sql[group[2]:group[3]], dollar)
		if fm != nil {
			m.Fragments = append(m.Fragments, fm)
		}

		before = group[1]
	}

	if before != len(sql) {
		m.Fragments = append(m.Fragments, getFragment(m, sql[before:], "", dollar))
	}
}

var placeholderRegexp = regexp.MustCompile(`\$\{(.+?)\}`)

func getFragment(m *Method, sql string, cond string, dollar *int) (fm *Fragment) {
	sql = strings.TrimSpace(sql)
	if sql == "" {
		return nil
	}

	fm = &Fragment{Cond: cond}
	matched := placeholderRegexp.FindAllStringSubmatchIndex(sql, -1)
	if len(matched) == 0 {
		fm.Stmt = sql
		return fm
	}

	// i.e.
	// select id, demo_name, demo_status
	// from demos
	// where id < ${id} and demo_name=${d.demeName} and demo_status=1
	var before int
	for _, group := range matched {
		//= select ... from demos where id <
		fm.Stmt += sql[before:group[0]]

		vt := calcArgs(m, fm, sql[group[2]:group[3]])
		if vt.Slice == "[]" && strings.HasSuffix(fm.Stmt, "(") {
			vt.SQLIn = true
		}

		//= select ... from demos where id < $1
		*dollar++
		fm.Stmt += "$%d"
		fm.Args = append(fm.Args, vt)

		before = group[1]
	}

	//= select ... from demos where id < $1 ... and demo_status=1
	fm.Stmt += sql[before:]

	return fm
}

func calcArgs(m *Method, fm *Fragment, v string) *VarAndType {
	sv := strings.Split(v, ".")

	if len(sv) == 1 {
		for _, vt := range m.Params {
			if vt.Var == v {
				return vt
			}
		}
		panic("variable " + v + " not in params")
	}

	if len(sv) > 2 {
		panic("unsupport x.y.z variable")
	}

	f := strings.TrimSpace(sv[1])
	f = strings.Replace(f, "_", " ", -1)
	f = strings.Title(f)
	f = strings.Replace(f, " ", "", -1)

	for _, p := range m.Params {
		if p.Var != sv[0] {
			continue
		}
		for _, prop := range p.Props {
			tmp := *prop
			tmp.Scope = p.Var
			if tmp.Var == f {
				return &tmp
			}
		}
	}
	panic("variable " + v + " not in params")
}

func calcScans(m *Method, sql string) {
	var start, end int
	ret := m.Returns
	switch m.Type {
	case MethodTypeCount:
		m.Scans = append(m.Scans, m.Returns[0])
		return

	case MethodTypeGet, MethodTypeList:
		start, end = 6, strings.LastIndex(sql, " from ")
		if end == -1 {
			panic("sql error: from not found")
		}

	case MethodTypeAdd, MethodTypeModify, MethodTypeRemove:
		ret = m.Params
		start = strings.Index(sql, " returning ")
		if start == -1 {
			return
		}
		start += 11
		end = strings.Index(sql, " on conflict ")
		if end == -1 {
			end = len(sql)
		}

	default:
		panic("unreachable code for " + string(m.Type))
	}

	fields := strings.Split(sql[start:end], ",")
next:
	for _, f := range fields {
		// abc_def => AbcDef
		f := strings.TrimSpace(f)
		f = f[strings.LastIndex(f, " ")+1:]
		f = strings.Replace(f, "_", " ", -1)
		f = strings.Title(f)
		f = strings.Replace(f, " ", "", -1)

		for _, p := range ret {
			if p.Package == "" || p.Package == "sql" {
				continue
			}
			for _, prop := range p.Props {
				if prop.Var == f {
					tmp := *prop
					tmp.Scope = p.Var
					tmp.Concat = "."
					m.Scans = append(m.Scans, &tmp)
					continue next
				}
			}
		}
	}
}

func calcMarshals(m *Method) {
	for _, fm := range m.Fragments {
		for _, prop := range fm.Args {
			if isBuiltin(prop.Type) {
				continue
			}
			if prop.Alias != "" {
				continue
			}
			if prop.Type == "Time" && prop.Package == "time" {
				continue
			}

			prop.Marshal = true

			prop.Concat = "_"

			tmp := *prop
			tmp.Concat = "_"
			m.Marshals = append(m.Marshals, &tmp)
		}
	}
}

func calcUnmarshals(m *Method) {
	for _, prop := range m.Scans {
		// prop.Scope = p.Var
		if prop.Alias != "" || isBuiltin(prop.Type) {
			continue
		}
		prop.Concat = "_"
		m.Unmarshals = append(m.Unmarshals, prop)
	}
}
