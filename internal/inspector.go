package internal

import (
	"flag"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

//sumtype:decl
type StrictnessLeveler interface {
	sealed()
}

type PackageLevelStrictness struct{}
type FunctionLevelStrictness struct{}

func (PackageLevelStrictness) sealed()  {}
func (FunctionLevelStrictness) sealed() {}

type AnalyzerConfig struct {
	Strictness StrictnessLeveler
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

func runWithConfig(config AnalyzerConfig) func(pass *analysis.Pass) (interface{}, error) {
	var messageOccurrences = make(map[messageKey]struct{})

	return func(pass *analysis.Pass) (interface{}, error) {
		for _, file := range pass.Files {
			// record the parent expression to always know which function we're in
			var parents []ast.Node
			runner := Runner{
				pass:               pass,
				messageOccurrences: messageOccurrences,
				config:             config,
				parents:            parents,
			}

			ast.Inspect(file, runner.ScanNode)
		}

		return nil, nil
	}
}
