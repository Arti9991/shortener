package exitchecker

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var ExitAnalyzer = &analysis.Analyzer{
	Name: "errcheck",
	Doc:  "check for unchecked errors",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	MainPkg := func(x *ast.File) bool {
		return x.Name.Name == "main"
	}

	MainFunc := func(x *ast.FuncDecl) bool {
		return x.Name.Name == "main"
	}

	OsExit := func(x *ast.SelectorExpr) bool {
		id, ok := x.X.(*ast.Ident)
		if !ok {
			return false
		}
		if id.Name == "os" && x.Sel.Name == "Exit" {
			pass.Reportf(id.NamePos, "call to os.Exit in main func in package main")
			return true
		}
		return false
	}

	Generated := func(f *ast.File) bool {
		for _, gr := range f.Comments {
			for _, comm := range gr.List {
				if strings.Contains(comm.Text, "generated") {
					return true
				}
			}
		}
		return false
	}
	for _, file := range pass.Files {

		// функцией ast.Inspect проходим по всем узлам AST
		ast.Inspect(file, func(node ast.Node) bool {
			if Generated(file) {
				return true
			}
			switch x := node.(type) {
			case *ast.File:
				if !MainPkg(x) {
					return true
				}
			case *ast.FuncDecl: // функция
				if !MainFunc(x) {
					//fmt.Printf("main function declaration: %s\n", x.Name.Name)
					return true
				}
			case *ast.SelectorExpr:
				if OsExit(x) {
					//fmt.Printf("%v: os.Exit called in main func in main package", file.Position(x.Pos()))
					return false
				}
			}
			return true
		})
	}
	return nil, nil
}
