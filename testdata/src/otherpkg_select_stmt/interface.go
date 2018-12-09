package otherpkg_select_stmt

import (
	"otherpkg_select_stmt/otherpkg"

	"github.com/pkg/errors"
)

func interfaceReturnError1() error {
	return otherpkg.InterfaceInstance.ReturnError() // want "L10: `return otherpkg.InterfaceInstance.ReturnError()"
}

func interfaceReturnError2() error {
	return errors.WithStack(otherpkg.InterfaceInstance.ReturnError())
}

func interfaceReturnError3() error {
	err := otherpkg.InterfaceInstance.ReturnError()
	return err // want "L18: `err := otherpkg.InterfaceInstance.ReturnError()"
}

func interfaceReturnError4() error {
	err := otherpkg.InterfaceInstance.ReturnError()
	return errors.WithStack(err)
}

func interfaceReturnError5() error {
	if err := otherpkg.InterfaceInstance.ReturnError(); err != nil {
		return err // want "L28: `if err := otherpkg.InterfaceInstance.ReturnError()"
	}
	return nil
}

func interfaceReturnError6() error {
	if err := otherpkg.InterfaceInstance.ReturnError(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func interfaceReturnValueAndError1() (string, error) {
	return otherpkg.InterfaceInstance.ReturnValueAndError() // want "L42: `return otherpkg.InterfaceInstance.ReturnValueAndError()"
}

func interfaceReturnValueAndError2() (string, error) {
	v, err := otherpkg.InterfaceInstance.ReturnValueAndError()
	return v, errors.WithStack(err)
}

func interfaceReturnValueAndError3() (string, error) {
	if v, err := otherpkg.InterfaceInstance.ReturnValueAndError(); err == nil {
		return v, nil
	} else {
		return "", err // want "L51: `if v, err := otherpkg.InterfaceInstance.ReturnValueAndError()"
	}
}

func interfaceReturnValueAndError4() (string, error) {
	if v, err := otherpkg.InterfaceInstance.ReturnValueAndError(); err == nil {
		return v, nil
	} else {
		return "", errors.WithStack(err)
	}
}
