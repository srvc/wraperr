package named_return

import (
	stderrors "errors"

	"github.com/pkg/errors"
)

// ok
func returnError1() (err error) {
	err = stderrors.New("foobarbaz")
	return
}

// ng
func returnError2() (err error) {
	err = returnError1()
	return
}

// ok
func returnError3() (err error) {
	err = errors.WithStack(returnError1())
	return
}

// ok
func returnError4() (err error) {
	err = returnError1()
	err = errors.WithStack(err)
	return
}

// ok
func returnError5() (err error) {
	err = returnError1()
	if err != nil {
		err = errors.WithStack(err)
	}
	return
}
