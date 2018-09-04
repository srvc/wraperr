package main

import (
	stderrors "errors"
	"log"

	"github.com/pkg/errors"
)

func main() {
	run()
}

func run() {
	err := returnError()
	if err != nil {
		log.Fatalln(err)
	}
}

// ok
func returnError() error {
	return stderrors.New("foobarbaz")
}

// ng
func returnError2() error {
	return returnError()
}

// ok
func returnError3() error {
	return errors.Wrap(returnError(), "error occurred")
}

// ng
func returnError4() error {
	err := returnError()
	return err
}

// ok
func returnError5() error {
	err := errors.WithStack(returnError())
	return err
}

// ng
func returnError6() error {
	err := returnError()
	if err != nil {
		return err
	}
	return nil
}

// ok
func returnError7() error {
	err := returnError()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// ng
func returnError8() (err error) {
	err = returnError()
	return
}

// ok
func returnError9() (err error) {
	err = errors.WithStack(returnError())
	return
}

// ok
func returnError10() (err error) {
	err = returnError()
	err = errors.WithStack(err)
	return
}

// ok
func returnValueAndError() (string, error) {
	return "quxquux", stderrors.New("foobarbaz")
}

// ng
func returnValueAndError2() (string, error) {
	return returnValueAndError()
}

// ng
func returnValueAndError3() (string, error) {
	v, err := returnValueAndError()
	if err != nil {
		return "", err
	}
	return v, nil
}

// ok
func returnValueAndError4() (string, error) {
	v, err := returnValueAndError()
	if err != nil {
		return "", errors.Wrap(err, "error occurred")
	}
	return v, nil
}

// ng
func returnValueAndError5() (v string, err error) {
	v, err = returnValueAndError()
	return
}

// ok
func returnValueAndError6() (v string, err error) {
	v, err = returnValueAndError()
	if err != nil {
		err = errors.Wrap(err, "error occurred")
		return
	}
	return
}
