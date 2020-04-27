package size

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/flowdev/spaghetti-cutter/x/toolpkg"
	"golang.org/x/tools/go/packages"
)

// Check checks the complexity of the given package and reports if it is too
// big.
func Check(pkg *packages.Package, rootPkg string, maxSize uint) []error {
	relPkg := toolpkg.RelativePackageName(pkg, rootPkg)
	fmt.Println("Complexity configuration - Size:", maxSize)
	fmt.Println("Package:", relPkg, pkg.Name, pkg.PkgPath)

	var realSize uint
	for _, astf := range pkg.Syntax {
		realSize += sizeOfFile(astf)
	}
	fmt.Println("Size of package:", relPkg, realSize)
	if realSize > maxSize {
		return []error{
			fmt.Errorf("the maximum size for package '%s' is %d but it's real size is: %d",
				relPkg, maxSize, realSize),
		}
	}
	return nil
}

func sizeOfFile(astf *ast.File) uint {
	var size uint = 3

	for _, idecl := range astf.Decls {
		switch decl := idecl.(type) {
		case *ast.FuncDecl:
			size += sizeOfFunc(decl)
		case *ast.GenDecl:
			size += sizeOfDecl(decl)
		}
	}
	return size
}

func sizeOfDecl(decl *ast.GenDecl) uint {
	var size uint

	switch decl.Tok {
	case token.TYPE:
		for _, s := range decl.Specs {
			ts := s.(*ast.TypeSpec)
			size++
			size += sizeOfExpr(ts.Type)
			name := ts.Name.Name
			fmt.Println("Size of type:", name, size)
		}
	case token.VAR, token.CONST:
		for _, s := range decl.Specs {
			vs := s.(*ast.ValueSpec)
			size += uint(len(vs.Names))
			for _, v := range vs.Values {
				size += sizeOfExpr(v)
			}
			fmt.Println("Size of values:", vs.Names, size)
		}
	}
	return size
}

func sizeOfFunc(fun *ast.FuncDecl) uint {
	var size uint = 1

	name := fun.Name.Name
	fmt.Println("Size of func:", name, size)
	return size
}

func sizeOfExpr(expr ast.Expr) uint {
	return 1
}
