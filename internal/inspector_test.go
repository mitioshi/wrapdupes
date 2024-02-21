package internal

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func TestInspector(t *testing.T) {
	testdataDir := analysistest.TestData()
	examples := []struct {
		pkg string
	}{
		{pkg: "simple"},
		{pkg: "middlewrap"},
		{pkg: "multifile"},
		{pkg: "complex"},
		{pkg: "dynamicerr"},
	}
	analyzer := NewWrapDupesAnalyzer()
	for _, example := range examples {
		t.Run(example.pkg, func(t *testing.T) {
			analysistest.Run(t, testdataDir, &analyzer, example.pkg)
		})
	}
}
