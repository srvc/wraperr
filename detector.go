package wraperr

import (
	"go/ast"
	"go/build"
	"go/token"
	"go/types"
	"sort"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/loader"
)

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
		"github.com/pkg/fail.Errorf",
		"github.com/pkg/fail.New",
		"github.com/pkg/fail.Wrap",
	}
)

type Detector interface {
	CheckPackages(paths []string) error
}

func NewDetector() Detector {
	wrapperFuncSet := make(map[string]struct{}, len(WrapperFuncList))
	for _, f := range WrapperFuncList {
		wrapperFuncSet[f] = struct{}{}
	}

	return &detectorImpl{
		ctx:            build.Default,
		wrapperFuncSet: wrapperFuncSet,
	}
}

type detectorImpl struct {
	ctx            build.Context
	wrapperFuncSet map[string]struct{}
}

func (d *detectorImpl) CheckPackages(paths []string) error {
	prog, err := d.load(paths)
	if err != nil {
		return errors.WithStack(err)
	}

	unwrappedErrs := &unwrappedErrors{}

	for _, pkgInfo := range prog.InitialPackages() {
		for _, f := range pkgInfo.Files {
			ast.Walk(newVisitor(prog, pkgInfo, unwrappedErrs, d.wrapperFuncSet), f)
		}
	}

	if len(unwrappedErrs.Errors()) == 0 {
		return nil
	}

	sort.Sort(unwrappedErrs)

	return unwrappedErrs
}

func (d *detectorImpl) load(paths []string) (*loader.Program, error) {
	lc := loader.Config{
		Build: &d.ctx,
	}

	rest, err := lc.FromArgs(paths, true) // TODO: configurable
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load paths: %v", paths)
	}
	if len(rest) != 0 {
		return nil, errors.Wrapf(err, "unhandled paths: %v", rest)
	}

	prog, err := lc.Load()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load packages")
	}

	return prog, nil
}

type visitor struct {
	prog           *loader.Program
	pkg            *loader.PackageInfo
	unwrappedErrs  UnwrappedErrors
	wrapperFuncSet map[string]struct{}
}

func newVisitor(prog *loader.Program, pkg *loader.PackageInfo, unwrappedErrs UnwrappedErrors, wrapperFuncSet map[string]struct{}) ast.Visitor {
	return &visitor{
		prog:           prog,
		pkg:            pkg,
		unwrappedErrs:  unwrappedErrs,
		wrapperFuncSet: wrapperFuncSet,
	}
}

func (v *visitor) Visit(node ast.Node) (w ast.Visitor) {
	switch stmt := node.(type) {
	case *ast.FuncDecl:
		if stmt.Type != nil && stmt.Type.Results != nil {
			if fv, ok := newFuncVisitor(v, stmt); ok {
				return fv
			}
		}
	}
	return v
}

func (v *visitor) debug(x interface{}) {
	ast.Print(v.prog.Fset, x)
}

func (v *visitor) isWrapped(call *ast.CallExpr) bool {
	switch fexpr := call.Fun.(type) {
	case *ast.SelectorExpr:
		if f, ok := v.pkg.ObjectOf(fexpr.Sel).(*types.Func); ok {
			_, ok = v.wrapperFuncSet[f.FullName()]
			return ok
		}
	}
	return false
}

type funcVisitor struct {
	*visitor

	decl        *ast.FuncDecl
	errIdents   map[string]*errIdent
	errInReturn []bool
	errNames    map[string]struct{}
}

func newFuncVisitor(parent *visitor, decl *ast.FuncDecl) (v ast.Visitor, ok bool) {
	fields := decl.Type.Results.List
	errInReturn := make([]bool, 0, 2*len(fields))
	errNames := make(map[string]struct{}, 2*len(fields))
	for _, f := range fields {
		isErr := isErrorType(parent.pkg.TypeOf(f.Type))
		ok = ok || isErr
		if len(f.Names) == 0 {
			errInReturn = append(errInReturn, isErr)
		}
		for _, n := range f.Names {
			errInReturn = append(errInReturn, isErr)
			errNames[n.Name] = struct{}{}
		}
	}
	if ok {
		v = &funcVisitor{
			visitor:     parent,
			decl:        decl,
			errIdents:   map[string]*errIdent{},
			errInReturn: errInReturn,
			errNames:    errNames,
		}
	}
	return
}

type errIdent struct {
	*ast.Ident
	wrapped bool
}

func (v *funcVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch stmt := node.(type) {
	case *ast.AssignStmt:
		var errIds []*ast.Ident
		for _, expr := range stmt.Lhs {
			if !isErrorType(v.pkg.TypeOf(expr)) {
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
					if v.isWrapped(cexpr) {
						wrapped = true
					}
				}
			}
			for _, id := range errIds {
				v.errIdents[id.Name] = &errIdent{Ident: id, wrapped: wrapped}
			}
		}
	case *ast.ReturnStmt:
		switch len(stmt.Results) {
		case 0:
			// Named return values
			for n := range v.errNames {
				if errIdent, ok := v.errIdents[n]; ok && !errIdent.wrapped {
					v.recordUnwrappedError(stmt.Return)
				}
			}
		case len(v.errInReturn):
			// Simple return
			for i, expr := range stmt.Results {
				if !v.errInReturn[i] {
					continue
				}
				switch expr := expr.(type) {
				case *ast.Ident:
					if errIdent, ok := v.errIdents[expr.Name]; ok && !errIdent.wrapped {
						v.recordUnwrappedError(expr.NamePos)
					}
				case *ast.CallExpr:
					if !v.isWrapped(expr) {
						v.recordUnwrappedError(expr.Lparen)
					}
				default:
					// TODO: should report unexpected exper
				}
			}
		case 1:
			// Return another function directly
			switch expr := stmt.Results[0].(type) {
			case *ast.CallExpr:
				if !v.isWrapped(expr) {
					v.recordUnwrappedError(expr.Pos())
				}
			default:
				// TODO: should report unexpected exper
			}
		default:
			// TODO: should report unexpected exper
		}
	}
	return v
}

func (v *funcVisitor) recordUnwrappedError(pos token.Pos) {
	v.unwrappedErrs.Add(&UnwrappedError{
		Position: v.prog.Fset.Position(pos),
		Funcname: v.decl.Name.Name,
	})
}
