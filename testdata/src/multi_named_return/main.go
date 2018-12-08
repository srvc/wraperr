package multi_named_return

import (
	stderrors "errors"

	"github.com/pkg/errors"
)

func returnValueAndError1() (v string, err error) {
	return "quxquux", stderrors.New("foobarbaz")
}

func returnValueAndError2() (v string, err error) {
	v, err = returnValueAndError1()
	return // want "L14: `v, err = returnValueAndError1()"
}

func returnValueAndError3() (v string, err error) {
	v, err = returnValueAndError1()
	if err != nil {
		err = errors.Wrap(err, "error occurred")
		return
	}
	return
}
