package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path/filepath"
	"strconv"
	"strings"
)

func parseFile(gofile string) (itf *Interface, err error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, gofile, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	// ast.Print(fset, f)
	// format.Node(os.Stdout, fset, f)

	itf = &Interface{}

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		switch genDecl.Tok {
		case token.IMPORT:
			parseImports(genDecl, itf)

		case token.TYPE:
			parseInterface(genDecl, itf)
		}
	}

	deepParse(itf)

	return itf, nil
}

func parseImports(genDecl *ast.GenDecl, itf *Interface) {
Next:
	for _, spec := range genDecl.Specs {
		importSpec, ok := spec.(*ast.ImportSpec)
		if !ok {
			continue
		}

		imp := ""
		if importSpec.Name != nil {
			imp += importSpec.Name.Name + " "
		}
		imp += importSpec.Path.Value

		for _, ix := range itf.Imports {
			if ix == imp || ix == "fmt" {
				continue Next
			}
		}
		itf.Imports = append(itf.Imports, imp)
	}
}

func parseInterface(genDecl *ast.GenDecl, itf *Interface) {
	for _, spec := range genDecl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		itf.Name = typeSpec.Name.Name

		interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
		if !ok {
			continue
		}

		parseMethods(interfaceType, itf)
	}
}

func parseMethods(interfaceType *ast.InterfaceType, itf *Interface) {
	for _, m := range interfaceType.Methods.List {
		var f Func
		itf.Methods = append(itf.Methods, &f)

		f.Doc = getDoc(m.Doc)

		f.Name = m.Names[0].Name

		funcType, ok := m.Type.(*ast.FuncType)
		if !ok {
			continue
		}

		parseFuncType(funcType, &f)
	}
}

func parseFuncType(funcType *ast.FuncType, f *Func) {
	for _, field := range funcType.Params.List {
		f.Params = append(f.Params, parseField(field)...)
	}

	for _, field := range funcType.Results.List {
		f.Returns = append(f.Returns, parseField(field)...)
	}
}

func parseField(field *ast.Field) (rets []*VarAndType) {
	var vat VarAndType
	vat.Type = parseExpr(&vat, field.Type)

	for _, name := range field.Names {
		tmp := vat
		tmp.Var = name.Name
		rets = append(rets, &tmp)
	}
	if len(field.Names) == 0 {
		rets = append(rets, &vat)
	}
	return rets
}

func parseExpr(vat *VarAndType, expr ast.Expr) string {
	switch expr.(type) {
	case *ast.Ident:
		ident := expr.(*ast.Ident)
		return ident.Name

	case *ast.StarExpr:
		starExpr := expr.(*ast.StarExpr)
		vat.Star = "*"
		return parseExpr(vat, starExpr.X)

	case *ast.ArrayType:
		arrayType := expr.(*ast.ArrayType)
		vat.Slice = "[]"
		return parseExpr(vat, arrayType.Elt)

	case *ast.SelectorExpr:
		selectorExpr := expr.(*ast.SelectorExpr)
		vat.Package = parseExpr(vat, selectorExpr.X)
		return parseExpr(vat, selectorExpr.Sel)

	case *ast.StructType:
		// ast.Print(nil, expr)
		// structType := expr.(*ast.StructType)
		return ""

	case *ast.MapType:
		mapType := expr.(*ast.MapType)
		k := parseExpr(&VarAndType{}, mapType.Key)
		v := parseExpr(&VarAndType{}, mapType.Value)
		return "map[" + k + "]" + v

	default:
		panic(fmt.Sprintf("unsupported type: %#v", expr))
	}
}

func getDoc(g *ast.CommentGroup) (doc string) {
	if g == nil || len(g.List) == 0 {
		panic("no comment")
	}

	for _, comment := range g.List {
		text := strings.TrimSpace(comment.Text)
		text = strings.Trim(text, "/")
		text = strings.TrimSpace(text)
		doc += " " + text
	}
	return doc[1:]
}

func deepParse(itf *Interface) {
	for _, f := range itf.Methods {
		for _, p := range f.Params {
			parseStruct(itf, p)
		}
		for _, r := range f.Returns {
			parseStruct(itf, r)
		}
	}
}

func parseStruct(itf *Interface, vat *VarAndType) {
	if isBuiltin(vat.Type) {
		return
	}

	if vat.Package == "sql" && vat.Type == "Tx" {
		return
	}

	fillTypePath(itf, vat)

	typeSpec := getTypeSpec(itf, vat.Path, vat.Type)

	stype, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		typ := parseExpr(vat, typeSpec.Type)
		if isBuiltin(typ) {
			vat.Alias = typ
			return
		}
		return
	}

	for _, f := range stype.Fields.List {
		if len(f.Names) == 0 {
			// Embedded struct: recurse
			// TODO
			continue
		}
		var prop VarAndType
		prop.Type = parseExpr(&prop, f.Type)
		prop.Scope = vat.Var
		prop.Concat = "."

		if !isBuiltin(prop.Type) && !strings.HasPrefix(prop.Type, "map") {
			if prop.Package == "" {
				prop.Package = vat.Package
				prop.Path = vat.Path
			} else {
				fillTypePath(itf, &prop)
			}

			typeSpec := getTypeSpec(itf, prop.Path, prop.Type)
			typ := parseExpr(vat, typeSpec.Type)
			if typ != "" && isBuiltin(typ) {
				prop.Alias = typ
			}
		}
		for _, nam := range f.Names {
			tmp := prop
			tmp.Var = nam.Name
			vat.Props = append(vat.Props, &tmp)
		}
	}
}

func isBuiltin(typ string) bool {
	if typ == "" || strings.HasPrefix(typ, "map") {
		return false
	}
	if 'a' <= typ[0] && typ[0] <= 'z' {
		return true
	}
	return false
}

func fillTypePath(itf *Interface, vat *VarAndType) {
	var err error
	for _, imp := range itf.Imports {
		if strings.HasPrefix(imp, vat.Package+" ") {
			vat.Path = imp[len(vat.Package)+1:]
			vat.Path, err = strconv.Unquote(vat.Type)
			if err != nil {
				panic(fmt.Sprintf("unquote %s error: %s", imp[len(vat.Package)+1:], err))
			}
			return
		}

		if strings.HasSuffix(imp, vat.Package+`"`) {
			vat.Path, err = strconv.Unquote(imp)
			if err != nil {
				panic(fmt.Sprintf("unquote %s error: %s", imp, err))
			}
			return
		}
	}
}

func getTypeSpec(itf *Interface, path, name string) (typeSpec *ast.TypeSpec) {
	pkg, _ := build.Import(path, "", 0)
	fset := token.NewFileSet() // share one fset across the whole package
	for _, file := range pkg.GoFiles {
		f, err := parser.ParseFile(fset, filepath.Join(pkg.Dir, file), nil, 0)
		if err != nil {
			continue
		}

		// ast.Print(fset, f)

		for _, decl := range f.Decls {
			decl, ok := decl.(*ast.GenDecl)
			if ok && decl.Tok == token.IMPORT {
				parseImports(decl, itf)
				continue
			}

			if !ok || decl.Tok != token.TYPE {
				continue
			}

			for _, spec := range decl.Specs {
				spec := spec.(*ast.TypeSpec)
				if spec.Name.Name != name {
					continue
				}
				return spec
			}
		}
	}

	panic(fmt.Sprintf("%s.%s not exist", path, name))
}
