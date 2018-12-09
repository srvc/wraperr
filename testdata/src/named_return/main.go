package named_return

import (
	stderrors "errors"

	"github.com/pkg/errors"
)

func returnError1() (err error) {
	err = stderrors.New("foobarbaz")
	return
}

func returnError2() (err error) {
	err = returnError1()
	return // want "L15: `err = returnError1()"
}

func returnError3() (err error) {
	err = errors.WithStack(returnError1())
	return
}

func returnError4() (err error) {
	err = returnError1()
	err = errors.WithStack(err)
	return
}

func returnError5() (err error) {
	err = returnError1()
	if err != nil {
		err = errors.WithStack(err)
	}
	return
}
