package internal

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestInspector_PackageLevel(t *testing.T) {
	testdataDir := analysistest.TestData()
	examples := []struct {
		pkg string
	}{
		{pkg: "simple"},
		{pkg: "middlewrap"},
		{pkg: "multifile"},
		{pkg: "funclevel"},
		{pkg: "dynamicerr"},
	}
	analyzer := NewWrapDupesAnalyzer(AnalyzerConfig{Strictness: PackageLevelStrictness{}})
	for _, example := range examples {
		t.Run(example.pkg, func(t *testing.T) {
			analysistest.Run(t, testdataDir, &analyzer, example.pkg)
		})
	}
}

func TestInspector_FunctionLevel(t *testing.T) {
	testdataDir := analysistest.TestData()
	examples := []struct {
		name string
		pkgs []string
	}{
		{name: "simple", pkgs: []string{"funclevel"}},
		{name: "two_functions_within_same_package", pkgs: []string{"funclvl_two_funcs"}},
		{name: "same_function_in_different_packages", pkgs: []string{"funclvl_samefunc_diff_pkgs", "funclvl_samefunc_diff_pkgs/bar"}},
	}
	analyzer := NewWrapDupesAnalyzer(AnalyzerConfig{Strictness: FunctionLevelStrictness{}})
	for _, example := range examples {
		t.Run(example.name, func(t *testing.T) {
			analysistest.Run(t, testdataDir, &analyzer, example.pkgs...)
		})
	}
}
