package main

import (
	"bytes"
	"go/token"
	"strings"
	"testing"

	"github.com/srvc/wraperr"
)

func TestPrintUnwrappedErrors(t *testing.T) {
	w := new(bytes.Buffer)
	errs := &dummyUnwrappedErrors{
		errs: []*wraperr.UnwrappedError{
			{
				Line:       "err := foo()",
				ReturnedAt: token.Position{Filename: "foo.go", Line: 12, Column: 4},
				OccurredAt: token.Position{Filename: "foo.go", Line: 10, Column: 8},
			},
			{
				Line:       "_, err := bar()",
				ReturnedAt: token.Position{Filename: "bar.go", Line: 19, Column: 8},
				OccurredAt: token.Position{Filename: "baz.go", Line: 13, Column: 6},
			},
			{
				Line:       "err := qux(`quux`)",
				ReturnedAt: token.Position{Filename: "qux.go", Line: 25, Column: 12},
				OccurredAt: token.Position{Filename: "qux.go", Line: 21, Column: 10},
			},
		},
	}

	fprintUnwrappedErrors(w, errs)

	wantLines := []string{
		"foo.go:12:4:	foo.go:10:8:	`err := foo()`",
		"bar.go:19:8:	baz.go:13:6:	`_, err := bar()`",
		"qux.go:25:12:	qux.go:21:10:	``err := qux(`quux`)``",
	}

	if got, want := w.String(), strings.Join(wantLines, "\n")+"\n"; got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

type dummyUnwrappedErrors struct {
	wraperr.UnwrappedErrors
	errs []*wraperr.UnwrappedError
}

func (e *dummyUnwrappedErrors) Errors() []*wraperr.UnwrappedError {
	return e.errs
}
