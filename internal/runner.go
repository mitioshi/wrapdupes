package internal

import (
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

type Runner struct {
	pass               *analysis.Pass
	messageOccurrences map[messageKey]struct{}
	config             AnalyzerConfig
	parents            []ast.Node
	mx                 sync.Mutex
}

func (runner *Runner) ScanNode(node ast.Node) bool {
	if node == nil {
		runner.parents = runner.parents[:len(runner.parents)-1]
	} else {
		runner.parents = append(runner.parents, node)
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

		pkg := runner.pass.TypesInfo.ObjectOf(sel.Sel).Pkg()

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

		if !runner.isDuplicatingMessage(errorMessageLiteral) {
			continue
		}

		runner.pass.Reportf(
			callExpr.Pos(),
			"duplicate message for a wrapped error: %s",
			errorMessageLiteral.Value,
		)
	}

	return true
}

func (runner *Runner) isDuplicatingMessage(errorMessageLiteral *ast.BasicLit) bool {
	var key messageKey
	switch runner.config.Strictness.(type) {
	case PackageLevelStrictness:
		key = messageKey{errorMessage: errorMessageLiteral.Value, pkg: runner.pass.Pkg.Path()}
	case FunctionLevelStrictness:
		for parentIdx := len(runner.parents) - 1; parentIdx > 0; parentIdx-- {
			parent := runner.parents[parentIdx]
			if parentFunc, ok := parent.(*ast.FuncDecl); ok {
				key = messageKey{errorMessage: errorMessageLiteral.Value, pkg: runner.pass.Pkg.Path(), fn: parentFunc.Name.String()}
				break
			}
		}
	}

	runner.mx.Lock()
	defer runner.mx.Unlock()

	if strings.Contains(errorMessageLiteral.Value, "%w") {
		_, exists := runner.messageOccurrences[key]

		if exists {
			return true
		}

		runner.messageOccurrences[key] = struct{}{}
	}

	return false
}
