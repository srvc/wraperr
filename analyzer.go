package wraperr

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "wraperr",
	Doc:      "Check that error return value are wrapped",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      runAnalyze,
}

var (
	// WrapperFuncList contains "path/to/pkg.Funcname" strings that adds a context to errors.
	WrapperFuncList = []string{
		// stdlibs
		"errors.New",
		"fmt.Errorf",

		// github.com/pkg/errors
		"github.com/pkg/errors.Errorf",
		"github.com/pkg/errors.New",
		"github.com/pkg/errors.WithMessage",
		"github.com/pkg/errors.WithStack",
		"github.com/pkg/errors.Wrap",
		"github.com/pkg/errors.Wrapf",

		// github.com/srvc/fail
		"github.com/srvc/fail.Errorf",
		"github.com/srvc/fail.New",
		"github.com/srvc/fail.Wrap",
	}
)

type errIdent struct {
	*ast.Ident
	wrapped bool
}

var wrapperFuncSet map[string]struct{}

func init() {
	wrapperFuncSet = make(map[string]struct{}, len(WrapperFuncList))
	for _, f := range WrapperFuncList {
		wrapperFuncSet[f] = struct{}{}
	}
}

func runAnalyze(pass *analysis.Pass) (interface{}, error) {
	r := newFileReader(pass.Fset)

	reportFunc := func(assignedAt, returnedAt token.Pos) {
		occPos := pass.Fset.Position(assignedAt)
		line := sprintInlineCode(r.GetLine(assignedAt))
		pass.Reportf(returnedAt, "the error is assigned on L%d: %s", occPos.Line, line)
	}

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		NewChecker(pass.Fset, pass.TypesInfo, n.(*ast.FuncDecl)).Check(reportFunc)
	})

	return nil, nil
}
