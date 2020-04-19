package deps

import (
	"fmt"

	"github.com/flowdev/spaghetti-cutter/config"
	"golang.org/x/tools/go/packages"
)

const pkgMain = "main"

// Check checks the dependencies of the given package and reports offending
// imports.
func Check(pkg *packages.Package, cfg config.Config, rootPkg string) []error {
	relPkg := pkg.PkgPath[len(rootPkg):]
	if pkg.Name == pkgMain {
		relPkg = pkgMain

		fmt.Println("Dependency configuration:")
		fmt.Println("    Tool:", cfg.Tool)
		fmt.Println("    DB:", cfg.DB)
		fmt.Println("    God:", cfg.God)
		fmt.Println("    Allow:", cfg.Allow)
	} else if len(relPkg) > 0 && relPkg[0] == '/' {
		relPkg = relPkg[1:]
	}

	fmt.Println(pkg.ID, pkg.Name, pkg.PkgPath)
	fmt.Println("relPkg:", relPkg)

	return nil
}
