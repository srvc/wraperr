package wraperr

import "go/types"

var errorType *types.Interface

func init() {
	errorType = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)
}

func isErrorType(t types.Type) bool {
	return types.Implements(t, errorType)
}
