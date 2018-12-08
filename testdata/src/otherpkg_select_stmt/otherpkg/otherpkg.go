package otherpkg

func ReturnErrorFunc() error {
	return nil
}

func ReturnValueAndErrorFunc() (string, error) {
	return "", nil
}

type InterfaceType interface {
	ReturnError() error
	ReturnValueAndError() (string, error)
}

var (
	InterfaceInstance InterfaceType
	StructInstance    StructType
)

type StructType struct {
}

func (s *StructType) ReturnError() error {
	return nil
}

func (s *StructType) ReturnValueAndError() (string, error) {
	return "", nil
}
