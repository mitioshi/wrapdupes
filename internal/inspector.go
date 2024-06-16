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
	return func(pass *analysis.Pass) (interface{}, error) {
		runner := Runner{
			pass:               pass,
			messageOccurrences: make(map[messageKey]struct{}),
			config:             config,
		}

		for _, file := range pass.Files {
			// record the parent expression to always know which function we're in
			runner.parents = runner.parents[:0]

			ast.Inspect(file, runner.ScanNode)
		}

		return nil, nil
	}
}
