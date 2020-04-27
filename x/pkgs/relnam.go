package pkgs

import (
	"strings"

	"golang.org/x/tools/go/packages"
)

const pkgMain = "main"

// RelativePackageName return the package name relative towards the given root package.
// It can be `/`, `main` or a path like `pkg/x/mytool`.
func RelativePackageName(pkg *packages.Package, rootPkg string) string {
	if !strings.HasPrefix(pkg.PkgPath, rootPkg) {
		return pkg.PkgPath
	}

	relPkg := pkg.PkgPath[len(rootPkg):]
	if pkg.Name == pkgMain {
		return pkgMain
	} else if relPkg == "" {
		return "/"
	} else if relPkg[0] == '/' {
		return relPkg[1:]
	}
	return relPkg
}
