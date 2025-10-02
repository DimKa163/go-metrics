package noexit

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"strings"
)

var Analyzer = &analysis.Analyzer{
	Name: "noexit",
	Doc:  "This analyzer reports no-exit functions.",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if isGenerated(file, pass.Fset) {
			continue
		}
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := fun.X.(*ast.Ident); ok && ident.Name == "os" && fun.Sel.Name == "Exit" {
					pass.Reportf(fun.Pos(), "do not use os.Exit")
				}
			}
			return true
		})
	}
	return nil, nil
}

// isGenerated search phrase "Code generated"
// band-aid solution
// without this method checker starts scanning autogen files
func isGenerated(file *ast.File, fset *token.FileSet) bool {
	if file.Comments == nil {
		return false
	}
	for _, cg := range file.Comments {
		for _, c := range cg.List {
			if strings.Contains(c.Text, "Code generated") && strings.Contains(c.Text, "DO NOT EDIT") {
				return true
			}
		}
	}
	return false
}
