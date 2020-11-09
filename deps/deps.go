package deps

import (
	"fmt"
	"strings"

	"github.com/flowdev/spaghetti-cutter/x/config"
	"github.com/flowdev/spaghetti-cutter/x/pkgs"
)

// Check checks the dependencies of the given package and reports offending
// imports.
func Check(pkg *pkgs.Package, rootPkg string, cfg config.Config, pkgInfos map[string]*pkgs.PackageInfo) []error {
	relPkg, strictRelPkg := pkgs.RelativePackageName(pkg, rootPkg)
	checkSpecial := checkStandard

	if isPackageInList(cfg.Tool, nil, relPkg, strictRelPkg) {
		checkSpecial = checkTool
	} else if isPackageInList(cfg.DB, nil, relPkg, strictRelPkg) {
		checkSpecial = checkDB
	} else if isPackageInList(cfg.God, nil, relPkg, strictRelPkg) {
		checkSpecial = checkGod
	}

	return checkPkg(pkg, relPkg, strictRelPkg, rootPkg, cfg, checkSpecial)
}

func checkPkg(
	pkg *pkgs.Package,
	relPkg, strictRelPkg, rootPkg string,
	cfg config.Config,
	checkSpecial func(string, string, string, string, config.Config) error,
) (errs []error) {
	for _, p := range pkg.Imports {
		relImp, strictRelImp := "", ""
		internal := false

		if strings.HasPrefix(p.PkgPath, rootPkg) {
			relImp, strictRelImp = pkgs.RelativePackageName(p, rootPkg)
			internal = true
		} else {
			strictRelImp = p.PkgPath
		}

		pl, dollars := cfg.AllowOnlyIn.MatchingList(strictRelImp)
		if pl == nil {
			pl, dollars = cfg.AllowOnlyIn.MatchingList(relImp)
		}
		if pl != nil {
			if !isPackageInList(pl, dollars, relPkg, strictRelPkg) {
				errs = append(errs, fmt.Errorf(
					"package '%s' isn't allowed to import package '%s' (because of allowOnlyIn)",
					pkgs.UniquePackageName(relPkg, strictRelPkg),
					pkgs.UniquePackageName(relImp, strictRelImp)))
			}
			continue
		}

		if internal {
			// check in allow first:
			pl = nil
			if strictRelPkg != "" {
				pl, dollars = cfg.AllowAdditionally.MatchingList(strictRelPkg)
			}
			if pl == nil {
				pl, dollars = cfg.AllowAdditionally.MatchingList(relPkg)
			}
			if isPackageInList(pl, dollars, relImp, strictRelImp) {
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
	if isTestPackage(relPkg, strictRelPkg, relImp, strictRelImp) {
		return nil
	}
	return fmt.Errorf("tool package '%s' isn't allowed to import package '%s'",
		pkgs.UniquePackageName(relPkg, strictRelPkg),
		pkgs.UniquePackageName(relImp, strictRelImp))
}

func checkDB(relPkg, strictRelPkg, relImp, strictRelImp string, cfg config.Config) error {
	if isPackageInList(cfg.Tool, nil, relImp, strictRelImp) {
		return nil
	}
	if isPackageInList(cfg.DB, nil, relImp, strictRelImp) {
		return nil
	}
	if isTestPackage(relPkg, strictRelPkg, relImp, strictRelImp) {
		return nil
	}
	return fmt.Errorf("DB package '%s' isn't allowed to import package '%s'",
		pkgs.UniquePackageName(relPkg, strictRelPkg),
		pkgs.UniquePackageName(relImp, strictRelImp))
}

func checkGod(relPkg, strictRelPkg, relImp, strictRelImp string, cfg config.Config) error {
	return nil // God never fails ;-)
}

func checkStandard(relPkg, strictRelPkg, relImp, strictRelImp string, cfg config.Config) error {
	if isPackageInList(cfg.Tool, nil, relImp, strictRelImp) {
		return nil
	}
	if isPackageInList(cfg.DB, nil, relImp, strictRelImp) {
		return nil
	}
	if isTestPackage(relPkg, strictRelPkg, relImp, strictRelImp) {
		return nil
	}
	return fmt.Errorf("domain package '%s' isn't allowed to import package '%s'",
		pkgs.UniquePackageName(relPkg, strictRelPkg),
		pkgs.UniquePackageName(relImp, strictRelImp))
}

func isTestPackage(relPkg, strictRelPkg, relImp, strictRelImp string) bool {
	pkg := prodPkg(relPkg, strictRelPkg)
	imp := prodPkg(relImp, strictRelImp)
	return pkg == imp
}
func prodPkg(rel, strict string) string {
	p := strict
	if p == "" {
		p = rel
	}
	if strings.HasSuffix(p, "_test") {
		p = p[:len(p)-5]
	}
	return p
}

func isPackageInList(pl *config.PatternList, dollars []string, pkg, strictPkg string) bool {
	if strictPkg != "" {
		if pl.MatchString(strictPkg, dollars) {
			return true
		}
	}
	return pl.MatchString(pkg, dollars)
}
