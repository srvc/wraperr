package wraperr

import (
	"go/ast"
	"go/build"
	"go/token"
	"go/types"
	"sort"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
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
		"github.com/srvc/fail.Errorf",
		"github.com/srvc/fail.New",
		"github.com/srvc/fail.Wrap",
	}
)

type Detector interface {
	CheckPackages(paths []string) error
}

type Package struct {
	*types.Package
	*types.Info
	Files []*ast.File
}

func NewDetector() Detector {
	wrapperFuncSet := make(map[string]struct{}, len(WrapperFuncList))
	for _, f := range WrapperFuncList {
		wrapperFuncSet[f] = struct{}{}
	}

	return &detectorImpl{
		ctx:            build.Default,
		fset:           token.NewFileSet(),
		wrapperFuncSet: wrapperFuncSet,
	}
}

type detectorImpl struct {
	ctx            build.Context
	fset           *token.FileSet
	wrapperFuncSet map[string]struct{}
}

func (d *detectorImpl) CheckPackages(paths []string) error {
	pkgs, err := d.load(paths)
	if err != nil {
		return errors.WithStack(err)
	}

	var wg sync.WaitGroup
	unwrappedErrs := &unwrappedErrors{}

	for _, pkg := range pkgs {
		wg.Add(1)
		go func(pkg *Package) {
			defer wg.Done()
			for _, f := range pkg.Files {
				if isGenerated(f) {
					continue
				}
				ast.Walk(newVisitor(d.fset, pkg, unwrappedErrs, d.wrapperFuncSet), f)
			}
		}(pkg)
	}

	wg.Wait()

	if len(unwrappedErrs.Errors()) == 0 {
		return nil
	}

	sort.Sort(unwrappedErrs)

	return unwrappedErrs
}

func (d *detectorImpl) load(paths []string) ([]*Package, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode:  packages.LoadAllSyntax,
		Fset:  d.fset,
		Tests: true,
	}, paths...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load paths: %v", paths)
	}

	result := make([]*Package, len(pkgs), len(pkgs))

	for i, pkg := range pkgs {
		result[i] = &Package{
			Package: pkg.Types,
			Info:    pkg.TypesInfo,
			Files:   pkg.Syntax,
		}
	}

	return result, nil
}

type visitor struct {
	fset           *token.FileSet
	pkg            *Package
	unwrappedErrs  UnwrappedErrors
	wrapperFuncSet map[string]struct{}
}

func newVisitor(fset *token.FileSet, pkg *Package, unwrappedErrs UnwrappedErrors, wrapperFuncSet map[string]struct{}) ast.Visitor {
	return &visitor{
		fset:           fset,
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
	ast.Print(v.fset, x)
}

func (v *visitor) isWrapped(call *ast.CallExpr) bool {
	switch fexpr := call.Fun.(type) {
	case *ast.SelectorExpr:
		if f, ok := v.pkg.ObjectOf(fexpr.Sel).(*types.Func); ok {
			fullName := f.FullName()
			if idx := strings.LastIndex(fullName, "/vendor/"); idx != -1 {
				fullName = fullName[idx+len("/vendor/"):]
			}
			_, ok = v.wrapperFuncSet[fullName]
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
		Position: v.fset.Position(pos),
		Pkgname:  v.pkg.Path(),
		Funcname: v.decl.Name.Name,
	})
}
