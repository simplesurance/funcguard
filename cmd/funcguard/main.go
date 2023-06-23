package main

import (
	_ "net/http/pprof"

	"github.com/simplesurance/funcguard/funcguard"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(funcguard.NewAnalyzer().Analyzer)
}
