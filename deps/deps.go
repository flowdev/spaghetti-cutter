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

	if _, ok := cfg.Tool[relPkg]; ok {
		return checkPkg(pkg, relPkg, rootPkg, cfg, checkTool)
	}
	if _, ok := cfg.DB[relPkg]; ok {
		return checkPkg(pkg, relPkg, rootPkg, cfg, checkDB)
	}
	if _, ok := cfg.God[relPkg]; ok {
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
			if allowed, ok := cfg.Allow[relPkg]; ok {
				if _, ok = allowed[relImp]; ok {
					continue // this import is fine
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
	if !strings.HasPrefix(relImp, relPkg) {
		return fmt.Errorf("tool package '%s' isn't allowed to import package '%s'",
			relPkg, relImp)
	}
	return nil
}

func checkDB(relPkg, relImp string, cfg config.Config) error {
	if _, ok := cfg.Tool[relImp]; ok {
		return nil
	}
	if _, ok := cfg.DB[relImp]; ok {
		return nil
	}
	if !strings.HasPrefix(relImp, relPkg) {
		return fmt.Errorf("DB package '%s' isn't allowed to import package '%s'",
			relPkg, relImp)
	}
	return nil
}

func checkStandard(relPkg, relImp string, cfg config.Config) error {
	if _, ok := cfg.Tool[relImp]; ok {
		return nil
	}
	if _, ok := cfg.DB[relImp]; ok {
		return nil
	}
	if !strings.HasPrefix(relImp, relPkg) {
		return fmt.Errorf("package '%s' isn't allowed to import package '%s'",
			relPkg, relImp)
	}
	return nil
}
