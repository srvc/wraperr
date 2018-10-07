package wraperr_test

import (
	"go/token"
	"path/filepath"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/srvc/wraperr"
)

func TestDetector_CheckPackages(t *testing.T) {
	errsByPkg := map[string][]*wraperr.UnwrappedError{
		"simple": []*wraperr.UnwrappedError{
			{
				Funcname:   "returnError2",
				Line:       "return returnError1()",
				OccurredAt: token.Position{Line: 16},
				ReturnedAt: token.Position{Line: 16},
			},
			{
				OccurredAt: token.Position{Line: 26},
				ReturnedAt: token.Position{Line: 27},
				Funcname:   "returnError4",
				Line:       "err := returnError1()",
			},
			{
				OccurredAt: token.Position{Line: 38},
				ReturnedAt: token.Position{Line: 40},
				Funcname:   "returnError6",
				Line:       "err := returnError1()",
			},
		},
		"named_return": []*wraperr.UnwrappedError{
			{
				OccurredAt: token.Position{Line: 17},
				ReturnedAt: token.Position{Line: 18},
				Funcname:   "returnError2",
				Line:       "err = returnError1()",
			},
		},
		"multi_return": []*wraperr.UnwrappedError{
			{
				OccurredAt: token.Position{Line: 16},
				ReturnedAt: token.Position{Line: 16},
				Funcname:   "returnValueAndError2",
				Line:       "return returnValueAndError1()",
			},
			{
				OccurredAt: token.Position{Line: 21},
				ReturnedAt: token.Position{Line: 23},
				Funcname:   "returnValueAndError3",
				Line:       "v, err := returnValueAndError1()",
			},
		},
		"multi_named_return": []*wraperr.UnwrappedError{
			{
				OccurredAt: token.Position{Line: 16},
				ReturnedAt: token.Position{Line: 17},
				Funcname:   "returnValueAndError2",
				Line:       "v, err = returnValueAndError1()",
			},
		},
	}

	pkgs := []string{"simple", "named_return", "multi_return"}
	sort.StringSlice(pkgs).Sort()

	var wantErrs []*wraperr.UnwrappedError

	for i, s := range pkgs {
		pkgs[i] = filepath.Join("github.com/srvc/wraperr/testdata/detector", s)
		for _, err := range errsByPkg[s] {
			err.Pkgname = pkgs[i]
			wantErrs = append(wantErrs, err)
		}
	}

	detector := wraperr.NewDetector()
	err := detector.CheckPackages(pkgs)
	if err == nil {
		t.Fatalf("should return an error")
	}

	gotErrs, ok := wraperr.UnwrapUnwrappedErrorsError(err)
	if !ok {
		t.Fatalf("should return an UnwrappedErrorsa, but returned %v", err)
	}

	opt := cmp.Comparer(func(x, y token.Position) bool { return x.Line == y.Line })

	if diff := cmp.Diff(gotErrs.Errors(), wantErrs, opt); diff != "" {
		t.Errorf("detected errors differs: (-want +got)\n%s", diff)
	}
}
