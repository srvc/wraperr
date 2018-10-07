package otherpkg_select_stmt

import (
	"github.com/pkg/errors"
	"github.com/srvc/wraperr/testdata/detector/otherpkg_select_stmt/otherpkg"
)

// ng
func returnError1() error {
	return otherpkg.ReturnErrorFunc()
}

// ok
func returnError2() error {
	return errors.WithStack(otherpkg.ReturnErrorFunc())
}

// ng
func returnError3() error {
	err := otherpkg.ReturnErrorFunc()
	return err
}

// ok
func returnError4() error {
	err := otherpkg.ReturnErrorFunc()
	return errors.WithStack(err)
}

// ng
func returnError5() error {
	if err := otherpkg.ReturnErrorFunc(); err != nil {
		return err
	}
	return nil
}

// ok
func returnError6() error {
	if err := otherpkg.ReturnErrorFunc(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// ng
func returnValueAndError1() (string, error) {
	return otherpkg.ReturnValueAndErrorFunc()
}

// ok
func returnValueAndError2() (string, error) {
	v, err := otherpkg.ReturnValueAndErrorFunc()
	return v, errors.WithStack(err)
}

// ng
func returnValueAndError3() (string, error) {
	if v, err := otherpkg.ReturnValueAndErrorFunc(); err == nil {
		return v, nil
	} else {
		return "", err
	}
}

// ng
func returnValueAndError4() (string, error) {
	if v, err := otherpkg.ReturnValueAndErrorFunc(); err == nil {
		return v, nil
	} else {
		return "", errors.WithStack(err)
	}
}
