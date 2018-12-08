package wraperr

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"
)

type Checker interface {
	Check(ReportFunc)
}

type ReportFunc func(assignedAt, returnedAt token.Pos)

func NewChecker(
	fset *token.FileSet,
	info *types.Info,
	f *ast.FuncDecl,
) Checker {
	return &checkerImpl{
		fset:     fset,
		info:     info,
		funcType: f.Type,
		funcBody: f.Body,
	}
}

func newClousureChecker(
	fset *token.FileSet,
	info *types.Info,
	f *ast.FuncLit,
) Checker {
	return &checkerImpl{
		fset:     fset,
		info:     info,
		funcType: f.Type,
		funcBody: f.Body,
	}
}

type checkerImpl struct {
	fset     *token.FileSet
	info     *types.Info
	funcType *ast.FuncType
	funcBody *ast.BlockStmt

	errInReturn []bool
	errNames    map[string]struct{}
	errIdents   map[string]*errIdent
}

func (c *checkerImpl) Check(report ReportFunc) {
	if !c.init() {
		return
	}

	ast.Inspect(c.funcBody, func(n ast.Node) bool {
		switch stmt := n.(type) {
		case *ast.AssignStmt:
			c.checkAssignment(stmt)
		case *ast.ReturnStmt:
			c.checkReturn(stmt, report)
		case *ast.FuncLit:
			newClousureChecker(c.fset, c.info, stmt).Check(report)
			return false
		}
		return true
	})
}

func (c *checkerImpl) init() (ok bool) {
	if c.funcType == nil || c.funcType.Results == nil {
		return
	}

	fields := c.funcType.Results.List

	c.errInReturn = make([]bool, 0, 2*len(fields))
	c.errNames = make(map[string]struct{}, 2*len(fields))
	c.errIdents = map[string]*errIdent{}

	for _, f := range fields {
		isErr := isErrorType(c.info.TypeOf(f.Type))
		ok = ok || isErr
		if len(f.Names) == 0 {
			c.errInReturn = append(c.errInReturn, isErr)
		}
		for _, n := range f.Names {
			c.errInReturn = append(c.errInReturn, isErr)
			c.errNames[n.Name] = struct{}{}
		}
	}

	return
}

func (c *checkerImpl) checkAssignment(stmt *ast.AssignStmt) {
	var errIds []*ast.Ident
	for _, expr := range stmt.Lhs {
		if !isErrorType(c.info.TypeOf(expr)) {
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
				if c.isWrapped(cexpr) {
					wrapped = true
				}
			}
		}
		for _, id := range errIds {
			c.errIdents[id.Name] = &errIdent{Ident: id, wrapped: wrapped}
		}
	}
}

func (c *checkerImpl) checkReturn(stmt *ast.ReturnStmt, report ReportFunc) {
	switch len(stmt.Results) {
	case 0:
		// Named return values
		for n := range c.errNames {
			if errIdent, ok := c.errIdents[n]; ok && !errIdent.wrapped {
				report(errIdent.Pos(), stmt.Return)
			}
		}
	case len(c.errInReturn):
		// Simple return
		for i, expr := range stmt.Results {
			if !c.errInReturn[i] {
				continue
			}
			switch expr := expr.(type) {
			case *ast.Ident:
				if errIdent, ok := c.errIdents[expr.Name]; ok && !errIdent.wrapped {
					report(errIdent.Pos(), expr.NamePos)
				}
			case *ast.CallExpr:
				if !c.isWrapped(expr) {
					report(expr.Pos(), expr.Lparen)
				}
			default:
				// TODO: should report unexpected exper
			}
		}
	case 1:
		// Return another function directly
		switch expr := stmt.Results[0].(type) {
		case *ast.CallExpr:
			if !c.isWrapped(expr) {
				report(expr.Pos(), expr.Pos())
			}
		default:
			// TODO: should report unexpected exper
		}
	default:
		// TODO: should report unexpected exper
	}
}

func (c *checkerImpl) isWrapped(call *ast.CallExpr) bool {
	switch fexpr := call.Fun.(type) {
	case *ast.SelectorExpr:
		if f, ok := c.info.ObjectOf(fexpr.Sel).(*types.Func); ok {
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
