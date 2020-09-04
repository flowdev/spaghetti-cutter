package deps

import (
	"fmt"
	"strings"

	"github.com/flowdev/spaghetti-cutter/x/config"
	"github.com/flowdev/spaghetti-cutter/x/pkgs"
)

// Check checks the dependencies of the given package and reports offending
// imports.
func Check(pkg *pkgs.Package, rootPkg string, cfg config.Config) []error {
	relPkg, strictRelPkg := pkgs.RelativePackageName(pkg, rootPkg)

	checkSpecial := checkStandard
	var fullmatch, fullfound bool
	var halfmatchTool, halfmatchDB string

	if fullmatch, halfmatchTool = isHalfPackageInList(cfg.Tool, relPkg, strictRelPkg); fullmatch {
		checkSpecial = checkTool
		fullfound = true
	}
	if fullmatch, halfmatchDB = isHalfPackageInList(cfg.DB, relPkg, strictRelPkg); fullmatch {
		checkSpecial = checkDB
		fullfound = true
	}
	if isFullPackageInList(cfg.God, relPkg, strictRelPkg) {
		checkSpecial = checkGod
		fullfound = true
	}

	if !fullfound {
		if len(halfmatchTool) > 0 && len(halfmatchTool) >= len(halfmatchDB) {
			checkSpecial = checkHalfTool
		} else if len(halfmatchDB) > 0 {
			checkSpecial = checkHalfDB
		}
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

		pl := cfg.AllowOnlyIn.MatchingList(strictRelImp)
		if pl == nil {
			pl = cfg.AllowOnlyIn.MatchingList(relImp)
		}
		if pl != nil {
			if !isFullPackageInList(pl, relPkg, strictRelPkg) {
				errs = append(errs, fmt.Errorf(
					"package '%s' isn't allowed to import package '%s' (because of allowOnlyIn)",
					pkgs.UniquePackageName(relPkg, strictRelPkg),
					pkgs.UniquePackageName(relImp, strictRelImp)))
			}
			continue
		}

		if internal {
			// check in allow first:
			var pl *config.PatternList
			if strictRelPkg != "" {
				pl = cfg.AllowAdditionally.MatchingList(strictRelPkg)
			}
			if pl == nil {
				pl = cfg.AllowAdditionally.MatchingList(relPkg)
			}
			if isFullPackageInList(pl, relImp, strictRelImp) {
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

func checkHalfTool(relPkg, strictRelPkg, relImp, strictRelImp string, cfg config.Config) error {
	if isTestPackage(relPkg, strictRelPkg, relImp, strictRelImp) {
		return nil
	}
	return fmt.Errorf("tool sub-package '%s' isn't allowed to import package '%s'",
		pkgs.UniquePackageName(relPkg, strictRelPkg),
		pkgs.UniquePackageName(relImp, strictRelImp))
}

func checkDB(relPkg, strictRelPkg, relImp, strictRelImp string, cfg config.Config) error {
	if isFullPackageInList(cfg.Tool, relImp, strictRelImp) {
		return nil
	}
	if isFullPackageInList(cfg.DB, relImp, strictRelImp) {
		return nil
	}
	if isTestPackage(relPkg, strictRelPkg, relImp, strictRelImp) {
		return nil
	}
	return fmt.Errorf("DB package '%s' isn't allowed to import package '%s'",
		pkgs.UniquePackageName(relPkg, strictRelPkg),
		pkgs.UniquePackageName(relImp, strictRelImp))
}

func checkHalfDB(relPkg, strictRelPkg, relImp, strictRelImp string, cfg config.Config) error {
	if isFullPackageInList(cfg.Tool, relImp, strictRelImp) {
		return nil
	}
	if isTestPackage(relPkg, strictRelPkg, relImp, strictRelImp) {
		return nil
	}
	return fmt.Errorf("DB sub-package '%s' isn't allowed to import package '%s'",
		pkgs.UniquePackageName(relPkg, strictRelPkg),
		pkgs.UniquePackageName(relImp, strictRelImp))
}

func checkGod(relPkg, strictRelPkg, relImp, strictRelImp string, cfg config.Config) error {
	return nil // God never fails ;-)
}

func checkStandard(relPkg, strictRelPkg, relImp, strictRelImp string, cfg config.Config) error {
	if isFullPackageInList(cfg.Tool, relImp, strictRelImp) {
		return nil
	}
	if isFullPackageInList(cfg.DB, relImp, strictRelImp) {
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

func isFullPackageInList(pl *config.PatternList, pkg, strictPkg string) bool {
	if strictPkg != "" {
		if pl.MatchString(strictPkg) {
			return true
		}
	}
	return pl.MatchString(pkg)
}

func isHalfPackageInList(pl *config.PatternList, pkg, strictPkg string) (full bool, match string) {
	if strictPkg != "" {
		full, match = pl.FindString(strictPkg)
		if full {
			return true, ""
		}
	}
	return pl.FindString(pkg)
}
