package internal

import (
	"flag"
	"go/ast"
	"strings"
	"sync"

	"golang.org/x/tools/go/analysis"
)

type messageKey struct {
	errorMessage string
	pkg          string
	fn           string
}

type AnalyzerConfig struct {
	Strictness string
}

func NewWrapDupesAnalyzer(config AnalyzerConfig) analysis.Analyzer {
	analyzer := analysis.Analyzer{
		Name:             "wrapdupes",
		Doc:              "wrapdupes\n\nThis linter detects duplicate error messages in fmt.Errorf calls.",
		Flags:            flag.FlagSet{},
		Run:              runWithConfig(config),
		RunDespiteErrors: false,
	}

	return analyzer
}

var mx = sync.Mutex{}

func runWithConfig(config AnalyzerConfig) func(pass *analysis.Pass) (interface{}, error) {
	var messageOccurrences = make(map[messageKey]struct{})

	return func(pass *analysis.Pass) (interface{}, error) {
		for _, file := range pass.Files {
			// record the parent expression to always know which function we're in
			var parents []ast.Node

			ast.Inspect(file, func(node ast.Node) bool {
				return scanFileNode(node, pass, messageOccurrences, config, &parents)
			})
		}

		return nil, nil
	}
}

//nolint:cyclop,funlen // This is a big function
func scanFileNode(
	node ast.Node,
	pass *analysis.Pass,
	messageOccurrences map[messageKey]struct{},
	config AnalyzerConfig, parents *[]ast.Node,
) bool {
	if node == nil {
		*parents = (*parents)[:len(*parents)-1]
	} else {
		*parents = append(*parents, node)
	}

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

		pkg := pass.TypesInfo.ObjectOf(sel.Sel).Pkg()

		// pkg is nil for method calls on local variables
		if pkg == nil || pkg.Path() != "fmt" || sel.Sel.String() != "Errorf" {
			continue
		}

		errorMessageLiteral, ok := callExpr.Args[0].(*ast.BasicLit)

		if !ok { // i.e. fmt.Errorf(functionCall(...))
			// although this can produce a duplicate wrapper message, it cannot be realistically detected
			// Thus, let's skip this case
			continue
		}

		var key messageKey

		if config.Strictness == "package" {
			key = messageKey{errorMessage: errorMessageLiteral.Value, pkg: pass.Pkg.Path()}
		} else if config.Strictness == "function" {
			for parentIdx := len(*parents) - 1; parentIdx > 0; parentIdx-- {
				parent := (*parents)[parentIdx]
				if parentFunc, ok := parent.(*ast.FuncDecl); ok {
					key = messageKey{errorMessage: errorMessageLiteral.Value, pkg: pass.Pkg.Path(), fn: parentFunc.Name.String()}
					break
				}
			}
		}

		if strings.Contains(errorMessageLiteral.Value, "%w") {
			mx.Lock()

			_, exists := messageOccurrences[key]

			if exists {
				pass.Reportf(callExpr.Pos(), "duplicate message for a wrapped error: %s", errorMessageLiteral.Value)
			} else {
				messageOccurrences[key] = struct{}{}
			}
			mx.Unlock()
		}
	}

	return true
}
