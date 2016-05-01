package persist

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

func prepare(itf *Interface) (impl *Implement) {
	impl = &Implement{}
	impl.Source = os.Getenv("GOFILE")
	impl.Package = os.Getenv("GOPACKAGE")

	impl.Imports = itf.Imports

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
	m.Returns = f.Returns

	sql := strings.Trim(f.Doc, " \t")
	m.Type = getMethodType(sql, f)

	calcInResult(m, f)

	calcFragment(m, sql)

	calcMarshals(m)

	calcScans(m, sql)

	calcUnmarshals(m)
}

func calcInResult(m *Method, f *Func) {
	if m.Type == MethodTypeGet {
		m.In = f.Params[0].Name
		m.Result, m.ResultType = f.Returns[0].Name, f.Returns[0].Type[1:]
	} else if m.Type == MethodTypeList {
		m.In = f.Params[0].Name
		m.Result, m.ResultType = f.Returns[0].Name, f.Returns[0].Type[3:]
	} else if m.Type == MethodTypePage {
		m.In = f.Params[0].Name
		m.Result, m.ResultType = f.Returns[1].Name, f.Returns[1].Type[3:]
	} else {
		m.In = f.Params[0].Name
		m.Result, m.ResultType = f.Params[0].Name, f.Params[0].Type[1:]
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
		if len(f.Returns) == 3 {
			return MethodTypePage
		} else {
			if len(f.Returns) > 0 && strings.HasPrefix(f.Returns[0].Type, "[]") {
				return MethodTypeList
			} else {
				return MethodTypeGet
			}
		}
	default:
		panic("sql error: " + sql)
	}
}

var fregmentRegexp = regexp.MustCompile(`\[\?\((.+?)\)(.+?)\]`)

func calcFragment(m *Method, sql string) {
	matched := fregmentRegexp.FindAllStringSubmatchIndex(sql, -1)

	dollar := new(int)

	if len(matched) == 0 {
		m.Fragments = append(m.Fragments, getFragment(sql, "", dollar))
		return
	}

	before := 0
	for _, group := range matched {
		if before != group[0] {
			m.Fragments = append(m.Fragments, getFragment(sql[before:group[0]], "", dollar))
		}
		m.Fragments = append(m.Fragments,
			getFragment(sql[group[4]:group[5]], sql[group[2]:group[3]], dollar))

		before = group[1]
	}

	if before != len(sql) {
		m.Fragments = append(m.Fragments, getFragment(sql[before:], "", dollar))
	}
}

var placeholderRegexp = regexp.MustCompile(`\$\{(.+?)\}`)

func getFragment(sql string, cond string, dollar *int) (fm *Fragment) {
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

		//= select ... from demos where id < $1
		*dollar++
		fm.Stmt += "$" + strconv.Itoa(*dollar)

		fm.Args = append(fm.Args, getVarAndType(sql[group[2]:group[3]]))

		before = group[1]
	}

	//= select ... from demos where id < $1 ... and demo_status=1
	fm.Stmt += sql[before:]

	return fm
}

func getVarAndType(v string) (t *VarAndType) {
	t = &VarAndType{Var: v}

	// TODO
	// Var string
	//
	// Type    string
	// Slice   string
	// Star    string
	// Package string
	return t
}

func calcScans(m *Method, sql string) {
	var start, end int
	switch m.Type {
	case MethodTypeGet, MethodTypeList, MethodTypePage:
		start, end = 6, strings.Index(sql, " from ")

	case MethodTypeAdd, MethodTypeModify, MethodTypeRemove:
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
		panic("unreachable code")
	}

	fields := strings.Split(sql[start:end], ",")
	for _, f := range fields {
		f = strings.Trim(f, " \t")
		f = strings.Replace(f, "_", " ", -1)
		f = strings.Title(f)
		f = strings.Replace(f, " ", "", -1)
		m.Scans = append(m.Scans, "x."+f)
	}
}

func calcMarshals(m *Method) {
	for _, p := range m.Params {
		for _, prop := range p.Props {
			switch prop.Type {
			case "int", "int64", "int32", "int16", "int8":
			case "uint", "uint64", "uint32", "uin16", "uint8", "byte":
			case "string":
			default:
				for _, fm := range m.Fragments {
					for i, arg := range fm.Args {
						if arg.Var == m.In+"."+prop.Name {
							fm.Args[i].Var = m.In + "_" + prop.Name
							m.Marshals = append(m.Marshals, prop.Name)
						}
					}
				}
			}
		}
	}
}

func calcUnmarshals(m *Method) {
	for _, p := range m.Returns {
		for _, prop := range p.Props {
			switch prop.Type {
			case "int", "int64", "int32", "int16", "int8":
			case "uint", "uint64", "uint32", "uin16", "uint8", "byte":
			case "string":
			default:
				for i, scan := range m.Scans {
					if scan == "x."+prop.Name {
						m.Scans[i] = "x_" + prop.Name
						m.Unmarshals = append(m.Unmarshals, prop.Name)
						switch prop.Type[0] {
						case 'm':
						case '*':
							prop.Type = "&" + prop.Type[1:]
						case '[':
							prop.Type = "[]*" + m.ResultType[:strings.Index(m.ResultType, ".")+1] + prop.Type[3:]
						default:
							prop.Type = "&" + m.ResultType[:strings.Index(m.ResultType, ".")+1] + prop.Type
						}
						m.Unmarshals1 = append(m.Unmarshals1, prop)
					}
				}
			}
		}
	}
}
