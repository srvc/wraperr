package otherpkg_select_stmt

import (
	"github.com/pkg/errors"
	"github.com/srvc/wraperr/testdata/detector/otherpkg_select_stmt/otherpkg"
)

// ng
func interfaceReturnError1() error {
	return otherpkg.InterfaceInstance.ReturnError()
}

// ok
func interfaceReturnError2() error {
	return errors.WithStack(otherpkg.InterfaceInstance.ReturnError())
}

// ng
func interfaceReturnError3() error {
	err := otherpkg.InterfaceInstance.ReturnError()
	return err
}

// ok
func interfaceReturnError4() error {
	err := otherpkg.InterfaceInstance.ReturnError()
	return errors.WithStack(err)
}

// ng
func interfaceReturnError5() error {
	if err := otherpkg.InterfaceInstance.ReturnError(); err != nil {
		return err
	}
	return nil
}

// ok
func interfaceReturnError6() error {
	if err := otherpkg.InterfaceInstance.ReturnError(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// ng
func interfaceReturnValueAndError1() (string, error) {
	return otherpkg.InterfaceInstance.ReturnValueAndError()
}

// ok
func interfaceReturnValueAndError2() (string, error) {
	v, err := otherpkg.InterfaceInstance.ReturnValueAndError()
	return v, errors.WithStack(err)
}

// ng
func interfaceReturnValueAndError3() (string, error) {
	if v, err := otherpkg.InterfaceInstance.ReturnValueAndError(); err == nil {
		return v, nil
	} else {
		return "", err
	}
}

// ng
func interfaceReturnValueAndError4() (string, error) {
	if v, err := otherpkg.InterfaceInstance.ReturnValueAndError(); err == nil {
		return v, nil
	} else {
		return "", errors.WithStack(err)
	}
}
