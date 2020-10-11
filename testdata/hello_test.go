package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"testing"
	"unicode"
)

func TestHello(t *testing.T) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "track1.go", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Imports:")
	for _, i := range node.Imports {
		fmt.Println("\t", i.Path.Value)
	}

	fmt.Println("Comments:")
	for _, c := range node.Comments {
		fmt.Print("\t", c.Text())
	}

	fmt.Println("Functions:")
	for _, f := range node.Decls {
		fn, ok := f.(*ast.FuncDecl)
		if !ok {
			continue
		}

		fmt.Printf("\t%s(%s):\t%s", fn.Name.Name, printFuncParams(fn.Type.Params.List), fn.Doc.Text())

		if unicode.IsUpper(rune(fn.Name.Name[0])) {
			fmt.Println("\t\texported")
		}
	}
}

func printFuncParams(params []*ast.Field) string {
	out := ""

	for _, param := range params {
		for _, name := range param.Names {
			out += name.Name
		}
	}

	return out
}
