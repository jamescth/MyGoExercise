package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
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

	visitorFunc := func(n ast.Node) bool {
		binaryExpr, ok := n.(*ast.BinaryExpr)
		if !ok {
			return true
		}

		fmt.Printf("Found binary expression at: %d:%d\n",
			fset.Position(binaryExpr.Pos()).Line,
			fset.Position(binaryExpr.Pos()).Column,
		)
		return true
	}

	ast.Inspect(node, visitorFunc)
}
