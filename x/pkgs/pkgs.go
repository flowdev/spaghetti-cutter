package pkgs

import (
	"strings"

	"golang.org/x/tools/go/packages"
)

const pkgMain = "main"

const testSuffix = ".test"

// RelativePackageName return the package name relative towards the given root package.
// relPkg can be `/`, `main` or a path like `pkg/x/mytool`.
// strictRelPkg can be `/`, a path like `pkg/x/mytool` or empty.
func RelativePackageName(pkg *packages.Package, rootPkg string) (relPkg, strictRelPkg string) {
	if !strings.HasPrefix(pkg.PkgPath, rootPkg) {
		return pkg.PkgPath, ""
	}

	relPkg = pkg.PkgPath[len(rootPkg):]
	if pkg.Name == pkgMain {
		return pkgMain, strictRelativePkg(relPkg)
	}
	return strictRelativePkg(relPkg), ""
}
func strictRelativePkg(rawRelPkg string) string {
	if rawRelPkg == "" {
		return "/"
	} else if rawRelPkg[0] == '/' && len(rawRelPkg) > 1 {
		return removeDotTest(rawRelPkg[1:])
	}
	return removeDotTest(rawRelPkg)
}
func removeDotTest(pkg string) string {
	if strings.HasSuffix(pkg, testSuffix) {
		return pkg[:len(pkg)-len(testSuffix)]
	}
	return pkg
}

// UniquePackageName returns strictRelPkg if it isn't empty and relPkg otherwise.
func UniquePackageName(relPkg, strictRelPkg string) string {
	if strictRelPkg != "" {
		return strictRelPkg
	}
	return relPkg
}

// UniquePackages makes the given list of packages unique.
func UniquePackages(pkgs []*packages.Package) []*packages.Package {
	result := make([]*packages.Package, 0, len(pkgs))
	pkgNames := make(map[string]struct{}, len(pkgs))

	for _, pkg := range pkgs {
		relPkg, strictRelPkg := RelativePackageName(pkg, "")
		uniqPkg := UniquePackageName(relPkg, strictRelPkg)

		if _, ok := pkgNames[uniqPkg]; !ok {
			result = append(result, pkg)
			pkgNames[uniqPkg] = struct{}{}
		}
	}
	return result
}

// IsTestPackage returns true if the given package is a test package and false
// otherwise.
func IsTestPackage(pkg *packages.Package) bool {
	result := strings.HasSuffix(pkg.PkgPath, "_test") ||
		strings.HasSuffix(pkg.PkgPath, ".test") ||
		strings.HasSuffix(pkg.ID, ".test]") ||
		strings.HasSuffix(pkg.ID, ".test")
	return result
}
