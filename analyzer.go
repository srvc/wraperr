package wraperr

import (
	"bufio"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"strings"

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

var wrapperFuncSet map[string]struct{}

func init() {
	wrapperFuncSet = make(map[string]struct{}, len(WrapperFuncList))
	for _, f := range WrapperFuncList {
		wrapperFuncSet[f] = struct{}{}
	}
}

func runAnalyze(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	isWrapped := func(call *ast.CallExpr) bool {
		switch fexpr := call.Fun.(type) {
		case *ast.SelectorExpr:
			if f, ok := pass.TypesInfo.ObjectOf(fexpr.Sel).(*types.Func); ok {
				chunks := []string{}
				for _, chunk := range strings.Split(f.FullName(), "/") {
					if chunk == "vendor" {
						chunks = []string{}
					} else {
						chunks = append(chunks, chunk)
					}
				}
				_, ok = wrapperFuncSet[strings.Join(chunks, "/")]
				return ok
			}
		}
		return false
	}

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		f := n.(*ast.FuncDecl)
		if f.Type == nil || f.Type.Results == nil {
			return
		}

		fields := f.Type.Results.List
		errInReturn := make([]bool, 0, 2*len(fields))
		errNames := make(map[string]struct{}, 2*len(fields))
		errIdents := map[string]*errIdent{}
		var ok bool
		for _, f := range fields {
			isErr := isErrorType(pass.TypesInfo.TypeOf(f.Type))
			ok = ok || isErr
			if len(f.Names) == 0 {
				errInReturn = append(errInReturn, isErr)
			}
			for _, n := range f.Names {
				errInReturn = append(errInReturn, isErr)
				errNames[n.Name] = struct{}{}
			}
		}
		if !ok {
			return
		}

		lines := make(map[string][]string)

		// ref: https://github.com/kisielk/errcheck/blob/1787c4bee836470bf45018cfbc783650db3c6501/internal/errcheck/errcheck.go#L488-L498
		getLine := func(tp token.Pos) string {
			pos := pass.Fset.Position(tp)
			foundLines, ok := lines[pos.Filename]

			if !ok {
				f, err := os.Open(pos.Filename)
				if err == nil {
					sc := bufio.NewScanner(f)
					for sc.Scan() {
						foundLines = append(foundLines, sc.Text())
					}
					lines[pos.Filename] = foundLines
					f.Close()
				}
			}

			line := "??"
			if pos.Line-1 < len(foundLines) {
				line = strings.TrimSpace(foundLines[pos.Line-1])
			}

			return line
		}

		recordUnwrappedError := func(occurredAt, returnedAt token.Pos) {
			occPos := pass.Fset.Position(occurredAt)
			pass.Reportf(returnedAt, "the error is assigned on L%d: %s", occPos.Line, sprintInlineCode(getLine(occurredAt)))
		}

		ast.Inspect(f, func(n ast.Node) bool {
			switch stmt := n.(type) {
			case *ast.AssignStmt:
				var errIds []*ast.Ident
				for _, expr := range stmt.Lhs {
					if !isErrorType(pass.TypesInfo.TypeOf(expr)) {
						continue
					}
					if id, ok := expr.(*ast.Ident); ok {
						errIds = append(errIds, id)
					}
				}
				if len(errIds) > 0 {
					// Detect wrapped error assignment
					var wrapped bool
					for _, expr := range stmt.Rhs {
						if cexpr, ok := expr.(*ast.CallExpr); ok {
							if isWrapped(cexpr) {
								wrapped = true
							}
						}
					}
					for _, id := range errIds {
						errIdents[id.Name] = &errIdent{Ident: id, wrapped: wrapped}
					}
				}
			case *ast.ReturnStmt:
				switch len(stmt.Results) {
				case 0:
					// Named return values
					for n := range errNames {
						if errIdent, ok := errIdents[n]; ok && !errIdent.wrapped {
							recordUnwrappedError(errIdent.Pos(), stmt.Return)
						}
					}
				case len(errInReturn):
					// Simple return
					for i, expr := range stmt.Results {
						if !errInReturn[i] {
							continue
						}
						switch expr := expr.(type) {
						case *ast.Ident:
							if errIdent, ok := errIdents[expr.Name]; ok && !errIdent.wrapped {
								recordUnwrappedError(errIdent.Pos(), expr.NamePos)
							}
						case *ast.CallExpr:
							if !isWrapped(expr) {
								recordUnwrappedError(expr.Pos(), expr.Lparen)
							}
						default:
							// TODO: should report unexpected exper
						}
					}
				case 1:
					// Return another function directly
					switch expr := stmt.Results[0].(type) {
					case *ast.CallExpr:
						if !isWrapped(expr) {
							recordUnwrappedError(expr.Pos(), expr.Pos())
						}
					default:
						// TODO: should report unexpected exper
					}
				default:
					// TODO: should report unexpected exper
				}
			}
			return true
		})
	})

	return nil, nil
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
