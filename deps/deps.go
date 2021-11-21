package deps

import (
	"fmt"
	"strings"

	"github.com/flowdev/spaghetti-cutter/data"
	"github.com/flowdev/spaghetti-cutter/x/config"
	"github.com/flowdev/spaghetti-cutter/x/pkgs"
)

// Check checks the dependencies of the given package and reports offending
// imports.
func Check(pkg *pkgs.Package, rootPkg string, cfg config.Config, depMap *data.DependencyMap) []error {
	relPkg, strictRelPkg := pkgs.RelativePackageName(pkg, rootPkg)
	checkSpecial := checkStandard
	pkgImps := data.PkgImports{}

	var fullmatch, matchDB bool
	if _, fullmatch = isPackageInList(cfg.God, nil, relPkg, strictRelPkg); fullmatch {
		checkSpecial = checkGod
		pkgImps.PkgType = data.TypeGod
	}
	if matchDB, fullmatch = isPackageInList(cfg.DB, nil, relPkg, strictRelPkg); matchDB {
		if fullmatch {
			checkSpecial = checkDB
			pkgImps.PkgType = data.TypeDB
		} else {
			checkSpecial = checkHalfDB
		}
	}
	if matchTool, fullmatch := isPackageInList(cfg.Tool, nil, relPkg, strictRelPkg); matchTool {
		if fullmatch {
			checkSpecial = checkTool
			pkgImps.PkgType = data.TypeTool
		} else if !matchDB {
			checkSpecial = checkHalfTool
		}
	}

	unqPkg := pkgs.UniquePackageName(relPkg, strictRelPkg)
	errs := checkPkg(pkg, relPkg, strictRelPkg, rootPkg, cfg, checkSpecial, &pkgImps)
	if !pkgs.IsTestPackage(pkg) && len(pkgImps.Imports) > 0 {
		(*depMap)[unqPkg] = pkgImps
	}
	return errs
}

func checkPkg(
	pkg *pkgs.Package,
	relPkg, strictRelPkg, rootPkg string,
	cfg config.Config,
	checkSpecial func(string, string, string, string, config.Config) error,
	imps *data.PkgImports,
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
			if !pkgs.IsTestPackage(pkg) {
				imps.Imports = saveDep(imps.Imports, relImp, strictRelImp, cfg)
			}

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
	if isTestPackage(relPkg, strictRelPkg) {
		return nil
	}
	return fmt.Errorf("tool package '%s' isn't allowed to import package '%s'",
		pkgs.UniquePackageName(relPkg, strictRelPkg),
		pkgs.UniquePackageName(relImp, strictRelImp))
}

func checkHalfTool(relPkg, strictRelPkg, relImp, strictRelImp string, cfg config.Config) error {
	if isTestPackage(relPkg, strictRelPkg) {
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
	if isTestPackage(relPkg, strictRelPkg) {
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
	if isTestPackage(relPkg, strictRelPkg) {
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
	if isTestPackage(relPkg, strictRelPkg) {
		return nil
	}
	return fmt.Errorf("domain package '%s' isn't allowed to import package '%s'",
		pkgs.UniquePackageName(relPkg, strictRelPkg),
		pkgs.UniquePackageName(relImp, strictRelImp))
}

func isTestPackage(rel, strict string) bool {
	p := strict
	if p == "" {
		p = rel
	}
	return strings.HasSuffix(p, "_test")
}

func isPackageInList(pl data.PatternList, dollars []string, pkg, strictPkg string) (atAll, full bool) {
	if strictPkg != "" {
		if atAll, full := pl.MatchString(strictPkg, dollars); atAll {
			return true, full
		}
	}
	return pl.MatchString(pkg, dollars)
}

func saveDep(im map[string]data.PkgType, relImp, strictRelImp string, cfg config.Config) map[string]data.PkgType {
	if len(im) == 0 {
		im = make(map[string]data.PkgType, 32)
	}
	unqImp := pkgs.UniquePackageName(relImp, strictRelImp)

	if _, full := isPackageInList(cfg.Tool, nil, relImp, strictRelImp); full {
		im[unqImp] = data.TypeTool
	} else if _, full := isPackageInList(cfg.DB, nil, relImp, strictRelImp); full {
		im[unqImp] = data.TypeDB
	} else if _, full := isPackageInList(cfg.God, nil, relImp, strictRelImp); full {
		im[unqImp] = data.TypeGod
	} else {
		im[unqImp] = data.TypeStandard
	}
	return im
}
