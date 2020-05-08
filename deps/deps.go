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
	relPkg, strictRelPkg := pkgs.RelativePackageName(pkg, rootPkg)

	if isPackageInList(cfg.Tool, relPkg, strictRelPkg) {
		return checkPkg(pkg, relPkg, strictRelPkg, rootPkg, cfg, checkTool)
	}
	if isPackageInList(cfg.DB, relPkg, strictRelPkg) {
		return checkPkg(pkg, relPkg, strictRelPkg, rootPkg, cfg, checkDB)
	}
	if isPackageInList(cfg.God, relPkg, strictRelPkg) {
		return nil // God packages can't have a problem by definition
	}
	return checkPkg(pkg, relPkg, strictRelPkg, rootPkg, cfg, checkStandard)
}

func checkPkg(
	pkg *packages.Package,
	relPkg, strictRelPkg, rootPkg string,
	cfg config.Config,
	checkSpecial func(string, string, string, string, config.Config) error,
) (errs []error) {
	for _, p := range pkg.Imports {
		if strings.HasPrefix(p.PkgPath, rootPkg) {
			relImp, strictRelImp := pkgs.RelativePackageName(p, rootPkg)
			fmt.Println("checkPkg - imp:", relImp, strictRelImp, p.Name, p.PkgPath)

			// check in allow first:
			var pl *config.PatternList
			if strictRelPkg != "" {
				pl = cfg.Allow.MatchingList(strictRelPkg)
			}
			if pl == nil {
				pl = cfg.Allow.MatchingList(relPkg)
			}
			if isPackageInList(pl, relImp, strictRelImp) {
				continue // this import is fine
			}

			if err := checkSpecial(relPkg, strictRelPkg, relImp, strictRelImp, cfg); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}

func checkTool(relPkg, strictRelPkg, relImp, strictRelImp string, cfg config.Config) error {
	if !isSubPackage(relImp, strictRelImp, relPkg, strictRelPkg) {
		return fmt.Errorf("tool package '%s' isn't allowed to import package '%s'",
			pkgs.UniquePackageName(relPkg, strictRelPkg),
			pkgs.UniquePackageName(relImp, strictRelImp))
	}
	return nil
}

func checkDB(relPkg, strictRelPkg, relImp, strictRelImp string, cfg config.Config) error {
	if isPackageInList(cfg.Tool, relImp, strictRelImp) {
		return nil
	}
	if isPackageInList(cfg.DB, relImp, strictRelImp) {
		return nil
	}
	if !isSubPackage(relImp, strictRelImp, relPkg, strictRelPkg) {
		return fmt.Errorf("DB package '%s' isn't allowed to import package '%s'",
			pkgs.UniquePackageName(relPkg, strictRelPkg),
			pkgs.UniquePackageName(relImp, strictRelImp))
	}
	return nil
}

func checkStandard(relPkg, strictRelPkg, relImp, strictRelImp string, cfg config.Config) error {
	if isPackageInList(cfg.Tool, relImp, strictRelImp) {
		return nil
	}
	if isPackageInList(cfg.DB, relImp, strictRelImp) {
		return nil
	}
	if !isSubPackage(relImp, strictRelImp, relPkg, strictRelPkg) {
		return fmt.Errorf("package '%s' isn't allowed to import package '%s'",
			pkgs.UniquePackageName(relPkg, strictRelPkg),
			pkgs.UniquePackageName(relImp, strictRelImp))
	}
	return nil
}

func isSubPackage(relImp, strictRelImp, relPkg, strictRelPkg string) bool {
	pkg := strictRelPkg
	if pkg == "" {
		pkg = relPkg
	}
	if strings.HasSuffix(pkg, "_test") {
		pkg = pkg[:len(pkg)-5]
	}

	imp := strictRelImp
	if imp == "" {
		imp = relImp
	}
	return strings.HasPrefix(imp, pkg)
}

func isPackageInList(pl *config.PatternList, pkg, strictPkg string) bool {
	if strictPkg != "" && pl.MatchString(strictPkg) {
		return true
	}
	return pl.MatchString(pkg)
}
