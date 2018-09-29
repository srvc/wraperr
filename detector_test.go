package wraperr_test

import (
	"testing"

	"github.com/srvc/wraperr"
)

func TestDetector_CheckPackages(t *testing.T) {
	detector := wraperr.NewDetector()
	err := detector.CheckPackages([]string{"github.com/srvc/wraperr/testdata/detector/simple"})

	if err == nil {
		t.Fatalf("should return an error")
	}

	errs, ok := wraperr.UnwrapUnwrappedErrorsError(err)
	if !ok {
		t.Fatalf("should return an UnwrappedErrorsa, but returned %v", err)
	}

	cases := []struct {
		test     string
		line     int
		funcname string
	}{
		{
			test:     "return an error directly",
			line:     28,
			funcname: "returnError2",
		},
		{
			test:     "return an error variable",
			line:     39,
			funcname: "returnError4",
		},
		{
			test:     "return an error variable in if-statement",
			line:     52,
			funcname: "returnError6",
		},
		{
			test:     "return an error as named return value",
			line:     69,
			funcname: "returnError8",
		},
		{
			test:     "return a value and an error directly",
			line:     92,
			funcname: "returnValueAndError2",
		},
		{
			test:     "return a value and an error variables",
			line:     99,
			funcname: "returnValueAndError3",
		},
		{
			test:     "return a value and an error as named return values",
			line:     116,
			funcname: "returnValueAndError5",
		},
	}

	if got, want := len(errs.Errors()), len(cases); got != want {
		t.Errorf("returned %d errors, want %d errors", got, want)
	} else {
		for i, tc := range cases {
			unwrappedErr := errs.Errors()[i]
			t.Run(tc.test, func(t *testing.T) {
				if got, want := unwrappedErr.Funcname, tc.funcname; got != want {
					t.Errorf("reported funcname is %s, want %s", got, want)
				}

				if got, want := unwrappedErr.Position.Line, tc.line; got != want {
					t.Errorf("reported line# is %d, want %d", got, want)
				}
			})
		}
	}
}
