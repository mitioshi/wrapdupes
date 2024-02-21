package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"
	"wrapdupes/internal"
)

func main() {
	analyzer := internal.NewWrapDupesAnalyzer()
	singlechecker.Main(&analyzer)
}
