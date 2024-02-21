package main

import (
	"flag"

	"github.com/mitioshi/wrapdupes/internal"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	strictness := flag.String(
		"strictness",
		"package",
		"Forbid duplicate error messages in fmt.Errorf calls based on the strictness level. Valid values are 'package', 'function'.",
	)
	flag.Parse()

	config := internal.AnalyzerConfig{Strictness: *strictness}
	analyzer := internal.NewWrapDupesAnalyzer(config)
	singlechecker.Main(&analyzer)
}
