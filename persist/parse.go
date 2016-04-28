package persist

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/wothing/log"
)

func parseFile(gofile string) (i *Interface, err error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, gofile, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	// ast.Print(fset, f)
	// format.Node(os.Stdout, fset, f)

	i = &Interface{}

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		switch genDecl.Tok {
		case token.IMPORT:
			parseImports(genDecl, i)

		case token.TYPE:
			parseType(genDecl, i)
		}
	}

	return i, nil
}

func parseImports(genDecl *ast.GenDecl, i *Interface) {
	for _, spec := range genDecl.Specs {
		importSpec, ok := spec.(*ast.ImportSpec)
		if !ok {
			continue
		}

		path := ""
		if importSpec.Name != nil {
			path += importSpec.Name.Name + " "
		}
		path += importSpec.Path.Value

		i.Imports = append(i.Imports, path)
	}
}

func parseType(genDecl *ast.GenDecl, i *Interface) {
	for _, spec := range genDecl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		i.Name = typeSpec.Name.Name

		interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
		if !ok {
			continue
		}

		parseMethods(interfaceType, i)
	}
}

func parseMethods(interfaceType *ast.InterfaceType, i *Interface) {
	for _, m := range interfaceType.Methods.List {
		var f Func
		i.Methods = append(i.Methods, &f)

		f.Doc = getDoc(m.Doc)

		f.Name = m.Names[0].Name

		funcType, ok := m.Type.(*ast.FuncType)
		if !ok {
			continue
		}

		parseFuncType(funcType, &f)
	}
}

func parseFuncType(funcType *ast.FuncType, i *Func) {
	for _, p := range funcType.Params.List {
		var param Param
		i.Params = append(i.Params, &param)

		parseField(p, &param)
	}

	for _, p := range funcType.Results.List {
		var param Param
		i.Returns = append(i.Returns, &param)

		parseField(p, &param)
	}
}

func parseField(field *ast.Field, param *Param) {
	param.Name = field.Names[0].Name

	param.Type = parseExpr(field.Type)
}

func parseExpr(expr ast.Expr) (x string) {
	switch expr.(type) {
	case *ast.Ident:
		ident := expr.(*ast.Ident)
		return ident.Name

	case *ast.StarExpr:
		starExpr := expr.(*ast.StarExpr)
		return "*" + parseExpr(starExpr.X)

	case *ast.SelectorExpr:
		selectorExpr := expr.(*ast.SelectorExpr)
		return parseExpr(selectorExpr.X) + "." + selectorExpr.Sel.Name

	case *ast.ArrayType:
		arrayType := expr.(*ast.ArrayType)
		return "[]" + parseExpr(arrayType.Elt)

	default:
		panic("not implemented")
	}
}

func getDoc(g *ast.CommentGroup) (doc string) {
	for _, comment := range g.List {
		doc += " " + strings.TrimLeft(comment.Text, " /")
	}

	return doc
}
