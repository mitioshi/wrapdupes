package main

import (
	"flag"
	"log"

	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/mitioshi/wrapdupes/internal"
)

func main() {
	strictness := flag.String(
		"strictness",
		"function",
		"Forbid duplicate error messages in fmt.Errorf calls based on the strictness level. Valid values are 'package', 'function'.",
	)

	flag.Parse()

	var strictnessParsed internal.StrictnessLeveler

	switch *strictness {
	case "package":
		strictnessParsed = internal.PackageLevelStrictness{}
	case "function":
		strictnessParsed = internal.FunctionLevelStrictness{}
	default:
		log.Fatalln("Invalid strictness level. Valid values are 'package', 'function'.")
	}

	config := internal.AnalyzerConfig{Strictness: strictnessParsed}
	analyzer := internal.NewWrapDupesAnalyzer(config)
	singlechecker.Main(&analyzer)
}
