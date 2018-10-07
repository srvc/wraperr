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
		"select_stmt": []*wraperr.UnwrappedError{
			{
				OccurredAt: token.Position{Line: 17},
				ReturnedAt: token.Position{Line: 17},
				Funcname:   "ReturnError1",
				Line:       "return s.ReturnError()",
			},
			{
				OccurredAt: token.Position{Line: 27},
				ReturnedAt: token.Position{Line: 28},
				Funcname:   "ReturnError3",
				Line:       "err := s.ReturnError()",
			},
			{
				OccurredAt: token.Position{Line: 39},
				ReturnedAt: token.Position{Line: 40},
				Funcname:   "ReturnError5",
				Line:       "if err := s.ReturnError(); err != nil {",
			},
			{
				OccurredAt: token.Position{Line: 55},
				ReturnedAt: token.Position{Line: 55},
				Funcname:   "ReturnValueAndError1",
				Line:       `return "", s.ReturnError()`,
			},
			{
				OccurredAt: token.Position{Line: 60},
				ReturnedAt: token.Position{Line: 60},
				Funcname:   "ReturnValueAndError2",
				Line:       "return s.ReturnValueAndError1()",
			},
			{
				OccurredAt: token.Position{Line: 71},
				ReturnedAt: token.Position{Line: 74},
				Funcname:   "ReturnValueAndError4",
				Line:       "if v, err := s.ReturnValueAndError1(); err == nil {",
			},
		},
		"otherpkg_select_stmt": []*wraperr.UnwrappedError{
			// func
			{
				OccurredAt: token.Position{Line: 10},
				ReturnedAt: token.Position{Line: 10},
				Funcname:   "returnError1",
				Line:       "return otherpkg.ReturnErrorFunc()",
			},
			{
				OccurredAt: token.Position{Line: 20},
				ReturnedAt: token.Position{Line: 21},
				Funcname:   "returnError3",
				Line:       "err := otherpkg.ReturnErrorFunc()",
			},
			{
				OccurredAt: token.Position{Line: 32},
				ReturnedAt: token.Position{Line: 33},
				Funcname:   "returnError5",
				Line:       "if err := otherpkg.ReturnErrorFunc(); err != nil {",
			},
			{
				OccurredAt: token.Position{Line: 48},
				ReturnedAt: token.Position{Line: 48},
				Funcname:   "returnValueAndError1",
				Line:       "return otherpkg.ReturnValueAndErrorFunc()",
			},
			{
				OccurredAt: token.Position{Line: 59},
				ReturnedAt: token.Position{Line: 62},
				Funcname:   "returnValueAndError3",
				Line:       "if v, err := otherpkg.ReturnValueAndErrorFunc(); err == nil {",
			},
			// interface
			{
				OccurredAt: token.Position{Line: 10},
				ReturnedAt: token.Position{Line: 10},
				Funcname:   "interfaceReturnError1",
				Line:       "return otherpkg.InterfaceInstance.ReturnError()",
			},
			{
				OccurredAt: token.Position{Line: 20},
				ReturnedAt: token.Position{Line: 21},
				Funcname:   "interfaceReturnError3",
				Line:       "err := otherpkg.InterfaceInstance.ReturnError()",
			},
			{
				OccurredAt: token.Position{Line: 32},
				ReturnedAt: token.Position{Line: 33},
				Funcname:   "interfaceReturnError5",
				Line:       "if err := otherpkg.InterfaceInstance.ReturnError(); err != nil {",
			},
			{
				OccurredAt: token.Position{Line: 48},
				ReturnedAt: token.Position{Line: 48},
				Funcname:   "interfaceReturnValueAndError1",
				Line:       "return otherpkg.InterfaceInstance.ReturnValueAndError()",
			},
			{
				OccurredAt: token.Position{Line: 59},
				ReturnedAt: token.Position{Line: 62},
				Funcname:   "interfaceReturnValueAndError3",
				Line:       "if v, err := otherpkg.InterfaceInstance.ReturnValueAndError(); err == nil {",
			},
			// struct
			{
				OccurredAt: token.Position{Line: 10},
				ReturnedAt: token.Position{Line: 10},
				Funcname:   "structReturnError1",
				Line:       "return otherpkg.StructInstance.ReturnError()",
			},
			{
				OccurredAt: token.Position{Line: 20},
				ReturnedAt: token.Position{Line: 21},
				Funcname:   "structReturnError3",
				Line:       "err := otherpkg.StructInstance.ReturnError()",
			},
			{
				OccurredAt: token.Position{Line: 32},
				ReturnedAt: token.Position{Line: 33},
				Funcname:   "structReturnError5",
				Line:       "if err := otherpkg.StructInstance.ReturnError(); err != nil {",
			},
			{
				OccurredAt: token.Position{Line: 48},
				ReturnedAt: token.Position{Line: 48},
				Funcname:   "structReturnValueAndError1",
				Line:       "return otherpkg.StructInstance.ReturnValueAndError()",
			},
			{
				OccurredAt: token.Position{Line: 59},
				ReturnedAt: token.Position{Line: 62},
				Funcname:   "structReturnValueAndError3",
				Line:       "if v, err := otherpkg.StructInstance.ReturnValueAndError(); err == nil {",
			},
		},
	}

	keys := make([]string, 0, len(errsByPkg))
	for key := range errsByPkg {
		keys = append(keys, key)
	}
	sort.StringSlice(keys).Sort()

	var wantErrs []*wraperr.UnwrappedError
	for _, k := range keys {
		for _, err := range errsByPkg[k] {
			err.Pkgname = filepath.Join("github.com/srvc/wraperr/testdata/detector", k)
			wantErrs = append(wantErrs, err)
		}
	}

	detector := wraperr.NewDetector()
	err := detector.CheckPackages([]string{"./testdata/detector/..."})
	if err == nil {
		t.Fatalf("should return an error")
	}

	gotErrs, ok := wraperr.UnwrapUnwrappedErrorsError(err)
	if !ok {
		t.Fatalf("should return an UnwrappedErrorsa, but returned %v", err)
	}

	if want, got := len(wantErrs), len(gotErrs.Errors()); got != want {
		t.Errorf("returned %d errors, want %d errors", got, want)
	}

	opt := cmp.Comparer(func(x, y token.Position) bool { return x.Line == y.Line })

	if diff := cmp.Diff(gotErrs.Errors(), wantErrs, opt); diff != "" {
		t.Errorf("detected errors differs: (-want +got)\n%s", diff)
	}
}
