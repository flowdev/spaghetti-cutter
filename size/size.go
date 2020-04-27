package size

import (
	"fmt"
	"go/ast"
	"log"

	"github.com/flowdev/spaghetti-cutter/x/pkgs"
	"golang.org/x/tools/go/packages"
)

// Check checks the complexity of the given package and reports if it is too
// big.
func Check(pkg *packages.Package, rootPkg string, maxSize uint) []error {
	if pkgs.IsTestPackage(pkg) {
		return nil
	}
	relPkg := pkgs.RelativePackageName(pkg, rootPkg)
	fmt.Println("Complexity configuration - Size:", maxSize)
	fmt.Println("Package:", relPkg, pkg.Name, pkg.PkgPath)

	var realSize uint
	for _, astf := range pkg.Syntax {
		realSize += sizeOfFile(astf)
	}
	log.Printf("INFO - Size of package '%s': %d\n", relPkg, realSize)

	if realSize > maxSize {
		return []error{
			fmt.Errorf("the maximum size for package '%s' is %d but it's real size is: %d",
				relPkg, maxSize, realSize),
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
