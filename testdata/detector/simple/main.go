package simple

import (
	stderrors "errors"

	"github.com/pkg/errors"
)

// ok
func returnError1() error {
	return stderrors.New("foobarbaz")
}

// ng
func returnError2() error {
	return returnError1()
}

// ok
func returnError3() error {
	return errors.Wrap(returnError1(), "error occurred")
}

// ng
func returnError4() error {
	err := returnError1()
	return err
}

// ok
func returnError5() error {
	err := errors.WithStack(returnError1())
	return err
}

// ng
func returnError6() error {
	err := returnError1()
	if err != nil {
		return err
	}
	return nil
}

// ok
func returnError7() error {
	err := returnError1()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
