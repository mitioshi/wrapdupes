package internal

import (
	"go/ast"
	"slices"
	"strings"
	"sync"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/types/typeutil"
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

var shortFormWraps = []string{
	"\"%w: %v\"", // fmt.Errorf("%w: %v", errClosing, s.readErr)
	"\"%w: %s\"", // fmt.Errorf("%w: %s", ErrValidation, err.Error())
	"\"%w: %w\"", // fmt.Errorf("%w: %w", err, errNoLabelName)
	"\"%s: %w\"", // pre1.19 fmt.Errorf("%s: %w", err.Error(), ErrUniqueAErr)
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

		fn := typeutil.StaticCallee(runner.pass.TypesInfo, callExpr)
		if fn == nil {
			return false
		}

		if fn.Pkg().Path() != "fmt" || fn.Name() != "Errorf" {
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

	// skip messages like "%w: %v", otherwise it'd give false positives for actually different errors
	if strings.Contains(errorMessageLiteral.Value, "%w") &&
		!slices.Contains(shortFormWraps, errorMessageLiteral.Value) {
		_, exists := runner.messageOccurrences[key]

		if exists {
			return true
		}

		runner.messageOccurrences[key] = struct{}{}
	}

	return false
}
