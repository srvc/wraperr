package simple

import (
	stderrors "errors"

	"github.com/pkg/errors"
)

func returnError1() error {
	return stderrors.New("foobarbaz")
}

func returnError2() error {
	return returnError1() // want "L14: `return returnError1()"
}

func returnError3() error {
	return errors.Wrap(returnError1(), "error occurred")
}

func returnError4() error {
	err := returnError1()
	return err // want "L22: `err := returnError1()"
}

func returnError5() error {
	err := errors.WithStack(returnError1())
	return err
}

func returnError6() error {
	err := returnError1()
	if err != nil {
		return err // want "L32: `err := returnError1()"
	}
	return nil
}

func returnError7() error {
	err := returnError1()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
