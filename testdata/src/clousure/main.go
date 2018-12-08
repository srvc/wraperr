package main

import (
	stderrors "errors"

	"github.com/pkg/errors"
)

func returnError() error {
	return stderrors.New("foobarbaz")
}

func returnError1() error {
	err := func() error {
		err := returnError()
		return err // want "L15: `err := returnError()"
	}()
	return err // want "L14: `err := func()"
}

func returnError2() error {
	err := func() error {
		err := returnError()
		return errors.WithStack(err)
	}()
	return err // want "L22: `err := func()"
}

func returnError3() error {
	err := func() error {
		err := returnError()
		return err // want "L31: `err := returnError()"
	}()
	return errors.WithStack(err)
}

func returnError4() error {
	err := func() error {
		err := returnError()
		return errors.WithStack(err)
	}()
	return errors.WithStack(err)
}
