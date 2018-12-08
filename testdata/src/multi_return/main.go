package multi_return

import (
	stderrors "errors"

	"github.com/pkg/errors"
)

func returnValueAndError1() (string, error) {
	return "quxquux", stderrors.New("foobarbaz")
}

func returnValueAndError2() (string, error) {
	return returnValueAndError1() // want "L14: `return returnValueAndError1()"
}

func returnValueAndError3() (string, error) {
	v, err := returnValueAndError1()
	if err != nil {
		return "", err // want "L18: `v, err := returnValueAndError1()"
	}
	return v, nil
}

func returnValueAndError4() (string, error) {
	v, err := returnValueAndError1()
	if err != nil {
		return "", errors.Wrap(err, "error occurred")
	}
	return v, nil
}
