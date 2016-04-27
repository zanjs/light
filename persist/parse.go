package persist

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

func parse(gofile string) (itfData *ItfData) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, gofile, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	// ast.Print(fset, f)

	// format.Node(os.Stdout, fset, f)

	itfData = &ItfData{}

	for _, d := range f.Decls {
		// fmt.Printf("%d %#v\n", i, d)
		if gd, ok := d.(*ast.GenDecl); ok {
			for _, s := range gd.Specs {
				if ts, ok := s.(*ast.TypeSpec); ok {
					itfData.Name = ts.Name.Name + "Data"
					if it, ok := ts.Type.(*ast.InterfaceType); ok {

						for _, m := range it.Methods.List {
							funcData := FuncData{
								Params:  []*VarAndType{},
								Returns: []*VarAndType{},
							}
							itfData.Funcs = append(itfData.Funcs, &funcData)

							funcData.SQL = strings.Replace(strings.Trim(m.Doc.Text(), " \t\n"), "\n", " ", -1)
							if ft, ok := m.Type.(*ast.FuncType); ok {
								funcData.Name = m.Names[0].Name

								for _, f := range ft.Params.List {
									for _, n := range f.Names {
										if od, ok := n.Obj.Decl.(*ast.Field); ok {
											// fmt.Printf("%#v %#v\n", od.Names, od.Type)

											switch od.Type.(type) {
											case *ast.StarExpr:
												x := od.Type.(*ast.StarExpr).X.(*ast.SelectorExpr)
												y := x.X.(*ast.Ident)
												funcData.Params = append(funcData.Params,
													&VarAndType{f.Names[0].Name, "*" + y.Name + "." + x.Sel.Name})
											case *ast.SliceExpr:
												x := od.Type.(*ast.SliceExpr).X.(*ast.SelectorExpr)
												y := x.X.(*ast.Ident)
												funcData.Params = append(funcData.Params,
													&VarAndType{f.Names[0].Name, "*" + y.Name + "." + x.Sel.Name})
											case *ast.ArrayType:
												a := od.Type.(*ast.ArrayType)
												x, ok := a.Elt.(*ast.SelectorExpr)
												if ok {
													y := x.X.(*ast.Ident)
													funcData.Params = append(funcData.Params,
														&VarAndType{f.Names[0].Name, "*" + y.Name + "." + x.Sel.Name})
												} else {
													x := a.Elt.(*ast.StarExpr).X.(*ast.SelectorExpr)
													y := x.X.(*ast.Ident)
													funcData.Params = append(funcData.Params,
														&VarAndType{f.Names[0].Name, "[]*" + y.Name + "." + x.Sel.Name})
												}
											case *ast.Ident:
												x := od.Type.(*ast.Ident)
												funcData.Params = append(funcData.Params,
													&VarAndType{f.Names[0].Name, x.Name})
											}
										}
									}
								}

								for _, f := range ft.Results.List {
									for _, n := range f.Names {
										if od, ok := n.Obj.Decl.(*ast.Field); ok {
											switch od.Type.(type) {
											case *ast.StarExpr:
												x := od.Type.(*ast.StarExpr).X.(*ast.SelectorExpr)
												y := x.X.(*ast.Ident)
												funcData.Returns = append(funcData.Returns,
													&VarAndType{f.Names[0].Name, "*" + y.Name + "." + x.Sel.Name})
											case *ast.SliceExpr:
												x := od.Type.(*ast.SliceExpr).X.(*ast.SelectorExpr)
												y := x.X.(*ast.Ident)
												funcData.Returns = append(funcData.Returns,
													&VarAndType{f.Names[0].Name, "*" + y.Name + "." + x.Sel.Name})
											case *ast.ArrayType:
												a := od.Type.(*ast.ArrayType)
												x, ok := a.Elt.(*ast.SelectorExpr)
												if ok {
													y := x.X.(*ast.Ident)
													funcData.Returns = append(funcData.Returns,
														&VarAndType{f.Names[0].Name, "*" + y.Name + "." + x.Sel.Name})
												} else {
													x := a.Elt.(*ast.StarExpr).X.(*ast.SelectorExpr)
													y := x.X.(*ast.Ident)
													funcData.Returns = append(funcData.Returns,
														&VarAndType{f.Names[0].Name, "[]*" + y.Name + "." + x.Sel.Name})
												}
											case *ast.Ident:
												x := od.Type.(*ast.Ident)
												funcData.Returns = append(funcData.Returns,
													&VarAndType{f.Names[0].Name, x.Name})
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// js, _ := json.MarshalIndent(itfData, "", "    ")
	// fmt.Printf("%s\n", js)

	return itfData
}
