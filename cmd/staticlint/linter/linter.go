package linter

import (
	"go/ast"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"golang.org/x/tools/go/analysis"
)

var ExitChecker = &analysis.Analyzer{
	Name: "ExitChecker",
	Doc:  "Doc",
	Run:  checkExit,
}

func checkExit(pass *analysis.Pass) (interface{}, error) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("starting parsing")
	for _, file := range pass.Files {
		if file.Name.String() == "main" {
			ast.Inspect(file, func(n ast.Node) bool {
				if fd, ok := n.(*ast.FuncDecl); ok {
					if fd.Name.String() != "main" {
						return false
					}
				}

				if expr, ok := n.(*ast.ExprStmt); ok {
					if call, ok := expr.X.(*ast.CallExpr); ok {
						if selector, ok := call.Fun.(*ast.SelectorExpr); ok {
							if ident, ok := selector.X.(*ast.Ident); ok {
								if ident.Name == "os" && selector.Sel.Name == "Exit" {
									log.Info().Msg("found os.Exit call in main")
									return false
								}
							}
						}
					}
				}

				return true
			})
		}
	}

	return nil, nil
}
