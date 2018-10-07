package wraperr

import (
	"fmt"
	"go/token"
	"sync"
)

type UnwrappedError struct {
	Position   token.Position
	Pkgname    string
	Funcname   string
	Line       string
	ReturnedAt token.Position
	OccurredAt token.Position
}

func (ei *UnwrappedError) less(ej *UnwrappedError) bool {
	pi, pj := ei.Position, ej.Position

	if pi.Filename != pj.Filename {
		return pi.Filename < pj.Filename
	}
	if pi.Line != pj.Line {
		return pi.Line < pj.Line
	}

	return pi.Column < pj.Column
}

type UnwrappedErrors interface {
	error

	Errors() []*UnwrappedError
	Add(*UnwrappedError)
}

func UnwrapUnwrappedErrorsError(err error) (uerr UnwrappedErrors, ok bool) {
	uerr, ok = err.(UnwrappedErrors)
	return
}

type unwrappedErrors struct {
	mutex sync.Mutex
	errs  []*UnwrappedError
}

func (e *unwrappedErrors) Error() string {
	return fmt.Sprintf("%d unchecked errors", len(e.errs))
}

func (e *unwrappedErrors) Errors() []*UnwrappedError {
	return e.errs
}

func (e *unwrappedErrors) Add(err *UnwrappedError) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.errs = append(e.errs, err)
}

func (e *unwrappedErrors) Len() int {
	return len(e.errs)
}

func (e *unwrappedErrors) Less(i, j int) bool {
	return e.errs[i].less(e.errs[j])
}

func (e *unwrappedErrors) Swap(i, j int) {
	e.errs[i], e.errs[j] = e.errs[j], e.errs[i]
}
