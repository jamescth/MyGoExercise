package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func main() {
	src := []byte(`package main

import "fmt"

func main() {
    var sum int = 3 + 2
    fmt.Println(sum)
}`)

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "demo", src, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	ast.Fprint(os.Stdout, fset, node, nil)

	// start searching for the main() function declaration
	for _, decl := range node.Decls {
		// search for main function
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// inside the main funcion body
		for _, stmt := range fn.Body.List {
			// search through statements for a declaration...
			declStmt, ok := stmt.(*ast.DeclStmt)
			if !ok {
				continue
			}

			// continue with declStmt.Decl
			genDecl, ok := declStmt.Decl.(*ast.GenDecl)
			if !ok {
				continue
			}

			// declaration can have multiple specs,
			// search for a valuespec
			for _, spec := range genDecl.Specs {
				valueSpec, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}

				// continue with valueSpec.Values
				for _, expr := range valueSpec.Values {
					// search for a binary expr
					binaryExpr, ok := expr.(*ast.BinaryExpr)
					if !ok {
						continue
					}

					// found it!
					fmt.Printf("Found binary expression at: %d:%d\n",
						fset.Position(binaryExpr.Pos()).Line,
						fset.Position(binaryExpr.Pos()).Column,
					)
				}
			}
		}
	}
}
