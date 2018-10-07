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
		test       string
		returnedOn int
		occurredOn int
		funcname   string
		line       string
	}{
		{
			test:       "return an error directly",
			occurredOn: 28,
			returnedOn: 28,
			funcname:   "returnError2",
			line:       "return returnError()",
		},
		{
			test:       "return an error variable",
			occurredOn: 38,
			returnedOn: 39,
			funcname:   "returnError4",
			line:       "err := returnError()",
		},
		{
			test:       "return an error variable in if-statement",
			occurredOn: 50,
			returnedOn: 52,
			funcname:   "returnError6",
			line:       "err := returnError()",
		},
		{
			test:       "return an error as named return value",
			occurredOn: 68,
			returnedOn: 69,
			funcname:   "returnError8",
			line:       "err = returnError()",
		},
		{
			test:       "return a value and an error directly",
			occurredOn: 92,
			returnedOn: 92,
			funcname:   "returnValueAndError2",
			line:       "return returnValueAndError()",
		},
		{
			test:       "return a value and an error variables",
			occurredOn: 97,
			returnedOn: 99,
			funcname:   "returnValueAndError3",
			line:       "v, err := returnValueAndError()",
		},
		{
			test:       "return a value and an error as named return values",
			occurredOn: 115,
			returnedOn: 116,
			funcname:   "returnValueAndError5",
			line:       "v, err = returnValueAndError()",
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

				if got, want := unwrappedErr.Line, tc.line; got != want {
					t.Errorf("reported error is occurred on %q, want %q", got, want)
				}

				if got, want := unwrappedErr.OccurredAt.Line, tc.occurredOn; got != want {
					t.Errorf("reported error is occurred on #%d, want #%d", got, want)
				}

				if got, want := unwrappedErr.ReturnedAt.Line, tc.returnedOn; got != want {
					t.Errorf("reported error is returned on #%d, want #%d", got, want)
				}
			})
		}
	}
}
