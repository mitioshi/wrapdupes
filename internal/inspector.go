package internal

import (
	"flag"
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"strings"
)

type messageKey struct {
	pkg          string
	errorMessage string
}

func NewWrapDupesAnalyzer() analysis.Analyzer {
	analyzer := analysis.Analyzer{
		Name:             "wrapdupes",
		Doc:              "wrapdupes\n\nThis linter detects duplicate error messages in fmt.Errorf calls.",
		Flags:            flag.FlagSet{},
		Run:              run,
		RunDespiteErrors: false,
	}
	return analyzer
}

func run(pass *analysis.Pass) (interface{}, error) {
	messageOccurrences := make(map[messageKey]struct{})

	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			return scanFileNode(node, pass, messageOccurrences)
		})
	}
	return nil, nil
}

func scanFileNode(node ast.Node, pass *analysis.Pass, messageOccurrences map[messageKey]struct{}) bool {
	returnNode, ok := node.(*ast.ReturnStmt)
	// We're only interested in statements like
	// return nil, fmt.Errorf("something went wrong: %w", err)
	if !ok || len(returnNode.Results) < 1 {
		return true
	}
	for _, expr := range returnNode.Results {
		callExpr, ok := expr.(*ast.CallExpr)
		if !ok {
			continue
		}
		sel, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		pkg := pass.TypesInfo.ObjectOf(sel.Sel).Pkg().Path()
		if pkg != "fmt" || sel.Sel.String() != "Errorf" {
			continue
		}
		errorMessage := callExpr.Args[0].(*ast.BasicLit).Value
		key := messageKey{errorMessage: errorMessage, pkg: pkg}
		if strings.Contains(errorMessage, "%w") {
			_, exists := messageOccurrences[key]
			if exists {
				pass.Reportf(callExpr.Pos(), "duplicate message for a wrapped error: %s", errorMessage)
			} else {
				messageOccurrences[key] = struct{}{}
			}
		}
	}
	return true
}
