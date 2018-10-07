package otherpkg_select_stmt

import (
	"github.com/pkg/errors"
	"github.com/srvc/wraperr/testdata/detector/otherpkg_select_stmt/otherpkg"
)

// ng
func structReturnError1() error {
	return otherpkg.StructInstance.ReturnError()
}

// ok
func structReturnError2() error {
	return errors.WithStack(otherpkg.StructInstance.ReturnError())
}

// ng
func structReturnError3() error {
	err := otherpkg.StructInstance.ReturnError()
	return err
}

// ok
func structReturnError4() error {
	err := otherpkg.StructInstance.ReturnError()
	return errors.WithStack(err)
}

// ng
func structReturnError5() error {
	if err := otherpkg.StructInstance.ReturnError(); err != nil {
		return err
	}
	return nil
}

// ok
func structReturnError6() error {
	if err := otherpkg.StructInstance.ReturnError(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// ng
func structReturnValueAndError1() (string, error) {
	return otherpkg.StructInstance.ReturnValueAndError()
}

// ok
func structReturnValueAndError2() (string, error) {
	v, err := otherpkg.StructInstance.ReturnValueAndError()
	return v, errors.WithStack(err)
}

// ng
func structReturnValueAndError3() (string, error) {
	if v, err := otherpkg.StructInstance.ReturnValueAndError(); err == nil {
		return v, nil
	} else {
		return "", err
	}
}

// ng
func structReturnValueAndError4() (string, error) {
	if v, err := otherpkg.StructInstance.ReturnValueAndError(); err == nil {
		return v, nil
	} else {
		return "", errors.WithStack(err)
	}
}
