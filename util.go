package wraperr

import (
	"bufio"
	"go/ast"
	"go/token"
	"go/types"
	"os"
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
func isGeneratedFile(f *ast.File) bool {
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

func isTestFile(fset *token.FileSet, f *ast.File) bool {
	return strings.HasSuffix(fset.File(f.Pos()).Name(), "_test.go")
}

func sprintInlineCode(s string) string {
	cc := 1
	c := cc
	for _, r := range s {
		if r == '`' {
			cc++
			if cc > c {
				c = cc
			}
		} else {
			cc = 1
		}
	}
	q := strings.Repeat("`", c)
	return q + s + q
}

type fileReader struct {
	fset  *token.FileSet
	lines map[string][]string
}

func newFileReader(fset *token.FileSet) *fileReader {
	return &fileReader{
		fset:  fset,
		lines: make(map[string][]string),
	}
}

// ref: https://github.com/kisielk/errcheck/blob/1787c4bee836470bf45018cfbc783650db3c6501/internal/errcheck/errcheck.go#L488-L498
func (r *fileReader) GetLine(tp token.Pos) string {
	pos := r.fset.Position(tp)
	foundLines, ok := r.lines[pos.Filename]

	if !ok {
		f, err := os.Open(pos.Filename)
		if err == nil {
			sc := bufio.NewScanner(f)
			for sc.Scan() {
				foundLines = append(foundLines, sc.Text())
			}
			r.lines[pos.Filename] = foundLines
			f.Close()
		}
	}

	line := "??"
	if pos.Line-1 < len(foundLines) {
		line = strings.TrimSpace(foundLines[pos.Line-1])
	}

	return line
}
