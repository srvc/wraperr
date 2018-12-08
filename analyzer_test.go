package wraperr_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/srvc/wraperr"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, wraperr.Analyzer, "simple")
}
