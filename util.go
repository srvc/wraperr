package wraperr

import (
	"go/ast"
	"go/types"
	"strings"
)

var (
	genHdr = "// Code generated "
	genFtr = " DO NOT EDIT."

	errorType *types.Interface
)

func init() {
	errorType = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)
}

func isErrorType(t types.Type) bool {
	return t != nil && types.Implements(t, errorType)
}

// https://golang.org/s/generatedcode
func isGenerated(f *ast.File) bool {
	for _, cg := range f.Comments {
		for _, c := range cg.List {
			src := c.Text
			// from https://github.com/golang/lint/blob/06c8688daad7faa9da5a0c2f163a3d14aac986ca/lint.go#L129
			if strings.HasPrefix(src, genHdr) && strings.HasSuffix(src, genFtr) && len(src) >= len(genHdr)+len(genFtr) {
				return true
			}
		}
	}
	return false
}
