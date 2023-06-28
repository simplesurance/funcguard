package main

import (
	_ "net/http/pprof"

	"github.com/simplesurance/funcguard/funcguard"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(mustNewSingleCheckerAnalyzer())
}

func mustNewSingleCheckerAnalyzer() *analysis.Analyzer {
	a, err := funcguard.NewAnalyzer(funcguard.WithCmdlineFlags())
	if err != nil {
		panic("creating analyzer failed: " + err.Error())
	}

	return a.Analyzer
}
