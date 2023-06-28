package funcguard

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	a, err := NewAnalyzer()
	if err != nil {
		t.Fatal(err)
	}

	analysistest.Run(t, analysistest.TestData(), a.Analyzer)
}
