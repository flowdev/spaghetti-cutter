package deps

import (
	"fmt"
	"sort"
	"strings"

	"github.com/flowdev/spaghetti-cutter/x/config"
	"github.com/flowdev/spaghetti-cutter/x/pkgs"
)

// pkgType can be one of: Standard, Tool, DB or God
type pkgType int

// Enum of package types: Standard, Tool, DB and God
const (
	typeStandard pkgType = iota
	typeTool
	typeDB
	typeGod
)

var typeLetters = []rune("STDG")
var typeFormats = []string{"", "_", "`", "**"}

// pkgImports contains the package type and the imported internal packages with their types.
type pkgImports struct {
	PkgType pkgType
	Imports map[string]pkgType
}

// DependencyMap is mapping importing package to imported packages.
// importingPackageName -> (importedPackageName -> struct{})
// An imported package name could be added multiple times to the same importing
// package name due to test packages.
type DependencyMap map[string]pkgImports

// GenerateTable writes the dependency matrix to a file.
func GenerateTable(depMap DependencyMap, cfg config.Config, rootPkg string) string {
	allRows := make([]string, 0, len(depMap))
	allCols := make([]string, 0, len(depMap))
	allColsMap := make(map[string]pkgType, len(depMap))

	for pkg, pkgImps := range depMap {
		allRows = append(allRows, pkg)
		for impName, impType := range pkgImps.Imports {
			if _, ok := allColsMap[impName]; !ok {
				allColsMap[impName] = impType
				allCols = append(allCols, impName)
			}
		}
	}

	sort.Strings(allRows)
	sort.Strings(allCols)

	sb := &strings.Builder{}
	intro := `# Dependency Table For: ` + rootPkg + `

| `
	sb.WriteString(intro)

	// (column) header line: | | C o l 1 - G | C o l 2 | ... | C o l N - T |
	for _, col := range allCols {
		sb.WriteString("| ")
		for _, r := range col {
			sb.WriteRune(r)
			sb.WriteRune(' ')
		}
		letter := typeLetters[allColsMap[col]]
		sb.WriteString("- ")
		sb.WriteRune(letter)
		sb.WriteRune(' ')
	}
	sb.WriteString("|\n")

	// separator line: | :- | :-: | :-: | ... | :-: |
	sb.WriteString("| :- ")
	for range allCols {
		sb.WriteString("| :-: ")
	}
	sb.WriteString("|\n")

	// normal rows: | **Row1** | **G** | | ... | **T** |
	for _, row := range allRows {
		pkgImps := depMap[row]

		sb.WriteString("| ")
		format := typeFormats[pkgImps.PkgType]
		sb.WriteString(format)
		sb.WriteString(row)
		sb.WriteString(format)
		sb.WriteRune(' ')

		for _, col := range allCols {
			sb.WriteString("| ")
			if impType, ok := pkgImps.Imports[col]; ok {
				sb.WriteString(format)
				sb.WriteRune(typeLetters[impType])
				sb.WriteString(format)
				sb.WriteRune(' ')
			}
		}
		sb.WriteString("|\n")
	}

	legend := `
### Legend

* Rows - Importing packages
* columns - Imported packages


#### Meaning Of Row And Row Header Formating

* **Bold** - God package
` + "* `Code` - Database package" + `
* _Italic_ - Tool package


#### Meaning Of Letters In Table Columns

* G - God package
* D - Database package
* T - Tool package
* S - Standard package
`
	sb.WriteString(legend)
	return sb.String()
}

// Check checks the dependencies of the given package and reports offending
// imports.
func Check(pkg *pkgs.Package, rootPkg string, cfg config.Config, depMap *DependencyMap) []error {
	relPkg, strictRelPkg := pkgs.RelativePackageName(pkg, rootPkg)
	checkSpecial := checkStandard
	pkgImps := pkgImports{}

	var fullmatch, matchTool bool
	if matchTool, fullmatch = isPackageInList(cfg.Tool, nil, relPkg, strictRelPkg); matchTool {
		if fullmatch {
			checkSpecial = checkTool
			pkgImps.PkgType = typeTool
		} else {
			checkSpecial = checkHalfTool
		}
	}
	if matchDB, fullmatch := isPackageInList(cfg.DB, nil, relPkg, strictRelPkg); matchDB {
		if fullmatch {
			checkSpecial = checkDB
			pkgImps.PkgType = typeDB
		} else if !matchTool {
			checkSpecial = checkHalfDB
		}
	}
	if _, fullmatch = isPackageInList(cfg.God, nil, relPkg, strictRelPkg); fullmatch {
		checkSpecial = checkGod
		pkgImps.PkgType = typeGod
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
	imps *pkgImports,
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

func saveDep(im map[string]pkgType, relImp, strictRelImp string, cfg config.Config) map[string]pkgType {
	if len(im) == 0 {
		im = make(map[string]pkgType, 32)
	}
	unqImp := pkgs.UniquePackageName(relImp, strictRelImp)

	if _, full := isPackageInList(cfg.God, nil, relImp, strictRelImp); full {
		im[unqImp] = typeGod
	} else if _, full := isPackageInList(cfg.DB, nil, relImp, strictRelImp); full {
		im[unqImp] = typeDB
	} else if _, full := isPackageInList(cfg.Tool, nil, relImp, strictRelImp); full {
		im[unqImp] = typeTool
	} else {
		im[unqImp] = typeStandard
	}
	return im
}
