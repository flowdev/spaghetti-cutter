package deps

import (
	"fmt"
	"strings"

	"github.com/flowdev/spaghetti-cutter/x/config"
	"github.com/flowdev/spaghetti-cutter/x/pkgs"
)

// DependencyMap is mapping importing package to imported packages.
// importingPackageName -> (importedPackageName -> struct{})
// An imported package name could be added multiple times to the same importing
// package name due to test packages.
type DependencyMap map[string]map[string]struct{}

var mapValue = struct{}{}

// Check checks the dependencies of the given package and reports offending
// imports.
func Check(pkg *pkgs.Package, rootPkg string, cfg config.Config, depMap DependencyMap) []error {
	relPkg, strictRelPkg := pkgs.RelativePackageName(pkg, rootPkg)
	checkSpecial := checkStandard

	var fullmatch, matchTool bool
	if matchTool, fullmatch = isPackageInList(cfg.Tool, nil, relPkg, strictRelPkg); matchTool {
		if fullmatch {
			checkSpecial = checkTool
		} else {
			checkSpecial = checkHalfTool
		}
	}
	if matchDB, fullmatch := isPackageInList(cfg.DB, nil, relPkg, strictRelPkg); matchDB {
		if fullmatch {
			checkSpecial = checkDB
		} else if !matchTool {
			checkSpecial = checkHalfDB
		}
	}
	if _, fullmatch = isPackageInList(cfg.God, nil, relPkg, strictRelPkg); fullmatch {
		checkSpecial = checkGod
	}

	return checkPkg(pkg, relPkg, strictRelPkg, rootPkg, cfg, checkSpecial, depMap)
}

func checkPkg(
	pkg *pkgs.Package,
	relPkg, strictRelPkg, rootPkg string,
	cfg config.Config,
	checkSpecial func(string, string, string, string, config.Config) error,
	depMap DependencyMap,
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

		unqPkg := pkgs.UniquePackageName(relPkg, strictRelPkg)
		unqImp := pkgs.UniquePackageName(relImp, strictRelImp)
		pl, dollars := cfg.AllowOnlyIn.MatchingList(strictRelImp)
		if pl == nil {
			pl, dollars = cfg.AllowOnlyIn.MatchingList(relImp)
		}
		if pl != nil {
			if _, full := isPackageInList(pl, dollars, relPkg, strictRelPkg); !full {
				errs = append(errs, fmt.Errorf(
					"package '%s' isn't allowed to import package '%s' (because of allowOnlyIn)",
					unqPkg, unqImp))
			}
			continue
		}

		if internal {
			saveDep(depMap, unqPkg, unqImp) // remember dependency

			// check in allow first:
			pl = nil
			if strictRelPkg != "" {
				pl, dollars = cfg.AllowAdditionally.MatchingList(strictRelPkg)
			}
			if pl == nil {
				pl, dollars = cfg.AllowAdditionally.MatchingList(relPkg)
			}
			if _, full := isPackageInList(pl, dollars, relImp, strictRelImp); full {
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
	if _, full := isPackageInList(cfg.Tool, nil, relImp, strictRelImp); full {
		return nil
	}
	if _, full := isPackageInList(cfg.DB, nil, relImp, strictRelImp); full {
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
	if _, full := isPackageInList(cfg.Tool, nil, relImp, strictRelImp); full {
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
	if _, full := isPackageInList(cfg.Tool, nil, relImp, strictRelImp); full {
		return nil
	}
	if _, full := isPackageInList(cfg.DB, nil, relImp, strictRelImp); full {
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

func isPackageInList(pl *config.PatternList, dollars []string, pkg, strictPkg string) (atAll, full bool) {
	if strictPkg != "" {
		if atAll, full := pl.MatchString(strictPkg, dollars); atAll {
			return true, full
		}
	}
	return pl.MatchString(pkg, dollars)
}

func saveDep(dm DependencyMap, pkg, imp string) {
	im := dm[pkg]
	if len(im) == 0 {
		im = make(map[string]struct{}, 32)
		dm[pkg] = im
	}
	im[imp] = mapValue
}
