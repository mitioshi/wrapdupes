package main

import (
	"github.com/mitioshi/wrapdupes/internal"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	analyzer := internal.NewWrapDupesAnalyzer()
	singlechecker.Main(&analyzer)
}
