package deps

import (
	"fmt"

	"github.com/flowdev/spaghetti-cutter/config"
	"golang.org/x/tools/go/packages"
)

// Check checks the dependencies of the given package and reports offending
// imports.
func Check(pkg *packages.Package, cfg config.Config, rootPkg string) []error {
	fmt.Println("Dependency configuration:")
	fmt.Println("    Tool:", cfg.Tool)
	fmt.Println("    DB:", cfg.DB)
	fmt.Println("    God:", cfg.God)
	fmt.Println("    Allow:", cfg.Allow)

	fmt.Println(pkg.ID, pkg.Name, pkg.PkgPath)

	return nil
}
