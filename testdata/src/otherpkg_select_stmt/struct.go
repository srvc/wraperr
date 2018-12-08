package otherpkg_select_stmt

import (
	"otherpkg_select_stmt/otherpkg"

	"github.com/pkg/errors"
)

func structReturnError1() error {
	return otherpkg.StructInstance.ReturnError() // want "L10 `return otherpkg.StructInstance.ReturnError()"
}

func structReturnError2() error {
	return errors.WithStack(otherpkg.StructInstance.ReturnError())
}

func structReturnError3() error {
	err := otherpkg.StructInstance.ReturnError()
	return err // want "L18: `err := otherpkg.StructInstance.ReturnError()"
}

func structReturnError4() error {
	err := otherpkg.StructInstance.ReturnError()
	return errors.WithStack(err)
}

func structReturnError5() error {
	if err := otherpkg.StructInstance.ReturnError(); err != nil {
		return err // want "L28: `if err := otherpkg.StructInstance.ReturnError()"
	}
	return nil
}

func structReturnError6() error {
	if err := otherpkg.StructInstance.ReturnError(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func structReturnValueAndError1() (string, error) {
	return otherpkg.StructInstance.ReturnValueAndError() // want "L42: `return otherpkg.StructInstance.ReturnValueAndError()"
}

func structReturnValueAndError2() (string, error) {
	v, err := otherpkg.StructInstance.ReturnValueAndError()
	return v, errors.WithStack(err)
}

func structReturnValueAndError3() (string, error) {
	if v, err := otherpkg.StructInstance.ReturnValueAndError(); err == nil {
		return v, nil
	} else {
		return "", err // want "L51: `if v, err := otherpkg.StructInstance.ReturnValueAndError()"
	}
}

func structReturnValueAndError4() (string, error) {
	if v, err := otherpkg.StructInstance.ReturnValueAndError(); err == nil {
		return v, nil
	} else {
		return "", errors.WithStack(err)
	}
}
