package deps

import (
	"fmt"
	"strings"

	"github.com/flowdev/spaghetti-cutter/x/config"
	"github.com/flowdev/spaghetti-cutter/x/pkgs"
	"golang.org/x/tools/go/packages"
)

// Check checks the dependencies of the given package and reports offending
// imports.
func Check(pkg *packages.Package, rootPkg string, cfg config.Config) []error {
	relPkg := pkgs.RelativePackageName(pkg, rootPkg)

	if isPackageInList(cfg.Tool, relPkg) {
		return checkPkg(pkg, relPkg, rootPkg, cfg, checkTool)
	}
	if isPackageInList(cfg.DB, relPkg) {
		return checkPkg(pkg, relPkg, rootPkg, cfg, checkDB)
	}
	if isPackageInList(cfg.God, relPkg) {
		return nil // God packages can't have a problem by definition
	}
	return checkPkg(pkg, relPkg, rootPkg, cfg, checkStandard)
}

func checkPkg(
	pkg *packages.Package,
	relPkg, rootPkg string,
	cfg config.Config,
	checkSpecial func(string, string, config.Config) error,
) (errs []error) {
	for _, p := range pkg.Imports {
		if strings.HasPrefix(p.PkgPath, rootPkg) {
			relImp := pkgs.RelativePackageName(p, rootPkg)
			fmt.Println(relImp, p.Name, p.PkgPath)

			// check in allow first:
			for _, group := range cfg.Allow {
				if group.Left.Regexp.MatchString(relPkg) {
					if isPackageInList(*group.Right, relPkg) {
						continue // this import is fine
					}
				}
			}

			if err := checkSpecial(relPkg, relImp, cfg); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}

func checkTool(relPkg, relImp string, cfg config.Config) error {
	if !isSubPackage(relImp, relPkg) {
		return fmt.Errorf("tool package '%s' isn't allowed to import package '%s'",
			relPkg, relImp)
	}
	return nil
}

func checkDB(relPkg, relImp string, cfg config.Config) error {
	if isPackageInList(cfg.Tool, relImp) {
		return nil
	}
	if isPackageInList(cfg.DB, relImp) {
		return nil
	}
	if !isSubPackage(relImp, relPkg) {
		return fmt.Errorf("DB package '%s' isn't allowed to import package '%s'",
			relPkg, relImp)
	}
	return nil
}

func checkStandard(relPkg, relImp string, cfg config.Config) error {
	if isPackageInList(cfg.Tool, relImp) {
		return nil
	}
	if isPackageInList(cfg.DB, relImp) {
		return nil
	}
	if !isSubPackage(relImp, relPkg) {
		return fmt.Errorf("package '%s' isn't allowed to import package '%s'",
			relPkg, relImp)
	}
	return nil
}

func isSubPackage(relImp, relPkg string) bool {
	pkg := relPkg
	if strings.HasSuffix(pkg, "_test") {
		pkg = pkg[:len(pkg)-5]
	}
	return strings.HasPrefix(relImp, pkg)
}

func isPackageInList(pl config.PatternList, pkg string) bool {
	for _, p := range pl {
		if p.Regexp.MatchString(pkg) {
			return true
		}
	}
	return false
}
