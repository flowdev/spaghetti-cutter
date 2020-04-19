package deps

import (
	"fmt"
	"strings"

	"github.com/flowdev/spaghetti-cutter/config"
	"golang.org/x/tools/go/packages"
)

const pkgMain = "main"

// Check checks the dependencies of the given package and reports offending
// imports.
func Check(pkg *packages.Package, rootPkg string, cfg config.Config) []error {
	relPkg := relativePkg(pkg, rootPkg)
	if relPkg == pkgMain {
		fmt.Println("Dependency configuration:")
		fmt.Println("    Tool:", cfg.Tool)
		fmt.Println("    DB:", cfg.DB)
		fmt.Println("    God:", cfg.God)
		fmt.Println("    Allow:", cfg.Allow)
	}

	fmt.Println(pkg.ID, pkg.Name, pkg.PkgPath)
	fmt.Println("relPkg:", relPkg)

	if _, ok := cfg.Tool[relPkg]; ok {
		fmt.Println("Check tool package")
		return checkPkg(pkg, relPkg, rootPkg, cfg, checkTool)
	}
	if _, ok := cfg.DB[relPkg]; ok {
		fmt.Println("Check DB package")
		return checkPkg(pkg, relPkg, rootPkg, cfg, checkDB)
	}
	if _, ok := cfg.God[relPkg]; ok {
		fmt.Println("Check god package")
		return nil // God packages can't have a problem by definition
	}
	fmt.Println("Check standard package")
	return checkPkg(pkg, relPkg, rootPkg, cfg, checkStandard)
}

func checkPkg(
	pkg *packages.Package,
	relPkg, rootPkg string,
	cfg config.Config,
	checkSpecial func(string, string, config.Config) error,
) (errs []error) {
	fmt.Println("Imports of:", relPkg)
	for imp, p := range pkg.Imports {
		if strings.HasPrefix(p.PkgPath, rootPkg) {
			relImp := relativePkg(p, rootPkg)
			fmt.Println(relImp, imp, p.ID, p.Name, p.PkgPath)

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
	return nil
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

// relativePkg return the relative package towards the given root package.
// It can be `/`, `main` or a path line `pkg/x/mytool`.
func relativePkg(pkg *packages.Package, rootPkg string) string {
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
