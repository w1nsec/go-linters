package linters

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

var OSExitCheckAnalyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "check for os.exit usage",
	Run:  osExitRun,
}

func osExitRun(pass *analysis.Pass) (interface{}, error) {
	var (
		inMain = false
	)

	osExit := func(x *ast.CallExpr) {
		if s, ok := x.Fun.(*ast.SelectorExpr); ok {
			if p, ok := s.X.(*ast.Ident); ok {
				if p.Name == "os" && s.Sel.Name == "Exit" {
					pass.Reportf(x.Pos(), "call os.Exit in main func")
				}
			}
		}
	}
	for _, file := range pass.Files {
		// функцией ast.Inspect проходим по всем узлам AST
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.Package:
				if x.Name != "main" {
					break
				}
			case *ast.FuncDecl:
				if x.Name.Name == "main" {
					inMain = true
				} else {
					inMain = false
				}
			case *ast.CallExpr:
				if inMain {
					osExit(x)
				}
			}

			return true
		})
	}
	return nil, nil
}
