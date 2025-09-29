// run go run ./cmd/staticlint ./internal/... ./cmd/...
package main

import (
	"github.com/DimKa163/go-metrics/internal/analyzers/noexit"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
	"honnef.co/go/tools/unused"
)

func main() {
	var checks []*analysis.Analyzer
	checks = append(checks, printf.Analyzer)
	checks = append(checks, shadow.Analyzer)
	checks = append(checks, shift.Analyzer)
	checks = append(checks, structtag.Analyzer)
	checks = append(checks, noexit.Analyzer)
	checks = append(checks, unused.Analyzer.Analyzer)
	for _, a := range stylecheck.Analyzers {
		checks = append(checks, a.Analyzer)
	}
	for _, v := range staticcheck.Analyzers {
		checks = append(checks, v.Analyzer)
	}
	for _, a := range simple.Analyzers {
		checks = append(checks, a.Analyzer)
	}
	multichecker.Main(
		checks...,
	)
}
