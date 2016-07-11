package main

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/wothing/log"
)

func readGoFile(file string) (m *Mapper, err error) {
	// var docs map[string]string
	// docs, err = parseDocs(file)
	// if err != nil {
	// 	return
	// }

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	// log.Debug(f.Doc.Text())
	//
	// for _, c := range f.Comments {
	// 	log.Info()
	// 	log.Debug(c.Text())
	// }

	for _, d := range f.Decls {
		log.Info()
		genDecl, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}

		switch genDecl.Tok {
		case token.IMPORT:
			log.Debugf("%#v", genDecl)

		case token.TYPE:
			log.Debugf("%#v", genDecl)
		}
	}

	return
}

func parseDocs(file string) (docs map[string]string, err error) {

	return
}
