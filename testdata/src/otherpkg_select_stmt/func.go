package otherpkg_select_stmt

import (
	"otherpkg_select_stmt/otherpkg"

	"github.com/pkg/errors"
)

func returnError1() error {
	return otherpkg.ReturnErrorFunc() // want "L9: `return otherpkg.ReturnErrorFunc()"
}

func returnError2() error {
	return errors.WithStack(otherpkg.ReturnErrorFunc())
}

func returnError3() error {
	err := otherpkg.ReturnErrorFunc()
	return err // want "L18: `err := otherpkg.ReturnErrorFunc()"
}

func returnError4() error {
	err := otherpkg.ReturnErrorFunc()
	return errors.WithStack(err)
}

func returnError5() error {
	if err := otherpkg.ReturnErrorFunc(); err != nil {
		return err // want "L28: `if err := otherpkg.ReturnErrorFunc()"
	}
	return nil
}

func returnError6() error {
	if err := otherpkg.ReturnErrorFunc(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func returnValueAndError1() (string, error) {
	return otherpkg.ReturnValueAndErrorFunc() // want "L42: `return otherpkg.ReturnValueAndErrorFunc()"
}

func returnValueAndError2() (string, error) {
	v, err := otherpkg.ReturnValueAndErrorFunc()
	return v, errors.WithStack(err)
}

func returnValueAndError3() (string, error) {
	if v, err := otherpkg.ReturnValueAndErrorFunc(); err == nil {
		return v, nil
	} else {
		return "", err // want "L51: `if v, err := otherpkg.ReturnValueAndErrorFunc()"
	}
}

func returnValueAndError4() (string, error) {
	if v, err := otherpkg.ReturnValueAndErrorFunc(); err == nil {
		return v, nil
	} else {
		return "", errors.WithStack(err)
	}
}
