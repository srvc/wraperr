package wraperr_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/srvc/wraperr"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()

	cases := []string{
		"simple",
		"named_return",
		"multi_return",
		"multi_named_return",
		"select_stmt",
		"otherpkg_select_stmt",
		"generated",
	}

	for _, tc := range cases {
		t.Run(tc, func(t *testing.T) {
			analysistest.Run(t, testdata, wraperr.Analyzer, tc)
		})
	}
}
