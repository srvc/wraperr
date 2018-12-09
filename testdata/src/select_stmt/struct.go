package select_stmt

import (
	"github.com/pkg/errors"
)

type StructInstance struct {
}

func (s *StructInstance) ReturnError() error {
	return nil
}

func (s *StructInstance) ReturnError1() error {
	return s.ReturnError() // want "L15: `return s.ReturnError()"
}

func (s *StructInstance) ReturnError2() error {
	return errors.WithStack(s.ReturnError())
}

func (s *StructInstance) ReturnError3() error {
	err := s.ReturnError()
	return err // want "L23: `err := s.ReturnError()"
}

func (s *StructInstance) ReturnError4() error {
	err := s.ReturnError()
	return errors.WithStack(err)
}

func (s *StructInstance) ReturnError5() error {
	if err := s.ReturnError(); err != nil {
		return err // want "L33: `if err := s.ReturnError()"
	}
	return nil
}

func (s *StructInstance) ReturnError6() error {
	if err := s.ReturnError(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *StructInstance) ReturnValueAndError1() (string, error) {
	return "", s.ReturnError() // want "L47: `return"
}

func (s *StructInstance) ReturnValueAndError2() (string, error) {
	return s.ReturnValueAndError1() // want "L51: `return s.ReturnValueAndError1()"
}

func (s *StructInstance) ReturnValueAndError3() (string, error) {
	v, err := s.ReturnValueAndError1()
	return v, errors.WithStack(err)
}

func (s *StructInstance) ReturnValueAndError4() (string, error) {
	if v, err := s.ReturnValueAndError1(); err == nil {
		return v, nil
	} else {
		return "", err // want "L60: `if v, err := s.ReturnValueAndError1()"
	}
}

func (s *StructInstance) ReturnValueAndError5() (string, error) {
	if v, err := s.ReturnValueAndError1(); err == nil {
		return v, nil
	} else {
		return "", errors.WithStack(err)
	}
}
