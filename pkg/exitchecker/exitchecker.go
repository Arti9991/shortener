package exitchecker

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "errcheck",
	Doc:  "check for unchecked errors",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	// создаём token.FileSet
	fset := token.NewFileSet()

	for _, file := range pass.Files {
		// функцией ast.Inspect проходим по всем узлам AST
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.FuncDecl: // функция
				if x.Name.Name == "main" {
					fmt.Printf("main function declaration: %s, %v", x.Name.Name, fset.Position(x.Pos()))
					return false
				}
			}
			return true
		})
	}
	return nil, nil
}
