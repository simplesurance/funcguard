package funcguard

import (
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestDefaultRules(t *testing.T) {
	a, err := NewAnalyzer(WithLogger(t.Logf))
	if err != nil {
		t.Fatal(err)
	}

	analysistest.Run(t, filepath.Join(analysistest.TestData(), "default"), a.Analyzer)
}

func TestCustomRule(t *testing.T) {
	a, err := NewAnalyzer(WithConfig(&Config{
		Rules: []*Rule{
			{
				FunctionPath: "fmt.Println",
				ErrorMsg:     "fmt.Println is not allowed",
			},
		},
	}), WithLogger(t.Logf))
	if err != nil {
		t.Fatal(err)
	}

	analysistest.Run(t, filepath.Join(analysistest.TestData(), "custom"), a.Analyzer)
}
