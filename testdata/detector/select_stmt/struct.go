package select_stmt

import (
	"github.com/pkg/errors"
)

type StructInstance struct {
}

// ok
func (s *StructInstance) ReturnError() error {
	return nil
}

// ng
func (s *StructInstance) ReturnError1() error {
	return s.ReturnError()
}

// ok
func (s *StructInstance) ReturnError2() error {
	return errors.WithStack(s.ReturnError())
}

// ng
func (s *StructInstance) ReturnError3() error {
	err := s.ReturnError()
	return err
}

// ok
func (s *StructInstance) ReturnError4() error {
	err := s.ReturnError()
	return errors.WithStack(err)
}

// ng
func (s *StructInstance) ReturnError5() error {
	if err := s.ReturnError(); err != nil {
		return err
	}
	return nil
}

// ok
func (s *StructInstance) ReturnError6() error {
	if err := s.ReturnError(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// ng
func (s *StructInstance) ReturnValueAndError1() (string, error) {
	return "", s.ReturnError()
}

// ng
func (s *StructInstance) ReturnValueAndError2() (string, error) {
	return s.ReturnValueAndError1()
}

// ok
func (s *StructInstance) ReturnValueAndError3() (string, error) {
	v, err := s.ReturnValueAndError1()
	return v, errors.WithStack(err)
}

// ng
func (s *StructInstance) ReturnValueAndError4() (string, error) {
	if v, err := s.ReturnValueAndError1(); err == nil {
		return v, nil
	} else {
		return "", err
	}
}

// ng
func (s *StructInstance) ReturnValueAndError5() (string, error) {
	if v, err := s.ReturnValueAndError1(); err == nil {
		return v, nil
	} else {
		return "", errors.WithStack(err)
	}
}
