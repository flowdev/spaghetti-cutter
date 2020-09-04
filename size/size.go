package size

import (
	"fmt"
	"go/ast"
	"log"

	"github.com/flowdev/spaghetti-cutter/x/pkgs"
)

// Check checks the complexity of the given package and reports if it is too
// big.
func Check(pkg *pkgs.Package, rootPkg string, maxSize uint) []error {
	if pkgs.IsTestPackage(pkg) {
		return nil
	}
	relPkg, strictRelPkg := pkgs.RelativePackageName(pkg, rootPkg)
	uniqPkg := pkgs.UniquePackageName(relPkg, strictRelPkg)

	var realSize uint
	for _, astf := range pkg.Syntax {
		realSize += sizeOfFile(astf)
	}
	log.Printf("INFO - Size of package '%s': %d", uniqPkg, realSize)

	if realSize > maxSize {
		return []error{
			fmt.Errorf("the maximum size for package '%s' is %d but it's real size is: %d",
				uniqPkg, maxSize, realSize),
		}
	}
	return nil
}

func sizeOfFile(astf *ast.File) uint {
	var size uint

	for _, decl := range astf.Decls {
		size += sizeOfDecl(decl)
	}
	return size
}
