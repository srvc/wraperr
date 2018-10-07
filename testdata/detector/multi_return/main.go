package multi_return

import (
	stderrors "errors"

	"github.com/pkg/errors"
)

// ok
func returnValueAndError1() (string, error) {
	return "quxquux", stderrors.New("foobarbaz")
}

// ng
func returnValueAndError2() (string, error) {
	return returnValueAndError1()
}

// ng
func returnValueAndError3() (string, error) {
	v, err := returnValueAndError1()
	if err != nil {
		return "", err
	}
	return v, nil
}

// ok
func returnValueAndError4() (string, error) {
	v, err := returnValueAndError1()
	if err != nil {
		return "", errors.Wrap(err, "error occurred")
	}
	return v, nil
}
