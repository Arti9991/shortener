// Пакет для проверки main на наличие функции os.Exit.
package exitchecker

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Структура для анализатора.
var ExitAnalyzer = &analysis.Analyzer{
	Name: "errcheck",
	Doc:  "check for unchecked errors",
	Run:  run,
}

// run функция запуска анализатора.
func run(pass *analysis.Pass) (interface{}, error) {
	// Проверка что пакет main.
	MainPkg := func(x *ast.File) bool {
		return x.Name.Name == "main"
	}
	// Проверка что функция main.
	MainFunc := func(x *ast.FuncDecl) bool {
		return x.Name.Name == "main"
	}
	// Проверка что функция Exit из пакета os.
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
	// Проверка на сгенерированный код (в комментариях есть слово generated).
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
	// Обход файлов.
	for _, file := range pass.Files {
		// функцией ast.Inspect проходим по всем узлам AST.
		ast.Inspect(file, func(node ast.Node) bool {
			if Generated(file) {
				return true
			}
			// Проверяем типы токенов.
			switch x := node.(type) {
			// Для файла проверяем что пакет main.
			case *ast.File:
				if !MainPkg(x) {
					return true
				}
			// Для декларации функции проверяем что название main.
			case *ast.FuncDecl: // функция
				if !MainFunc(x) {
					//fmt.Printf("main function declaration: %s\n", x.Name.Name)
					return true
				}
			// Для запуска функции проверяем на os.Exit
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
