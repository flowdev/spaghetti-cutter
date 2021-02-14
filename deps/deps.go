package deps

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/flowdev/spaghetti-cutter/x/config"
	"github.com/flowdev/spaghetti-cutter/x/pkgs"
)

// DocFile is the name of the documentation file (package_dependencies.md).
const DocFile = "package_dependencies.md"

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

// FindDocPkgs is finding documentation packages on disk starting at 'root' and
// adding them to those given in 'dtPkgs'.
func FindDocPkgs(dtPkgs []string, root string) map[string]struct{} {
	val := struct{}{}
	// prefill doc packages from dtPkgs
	retPkgs := make(map[string]struct{}, 128)
	for _, p := range dtPkgs {
		retPkgs[p] = val
	}

	// walk the file system to find more dependency table files
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() { // we are only interested in directories
			return nil
		}
		if err != nil {
			log.Printf("WARN - Unable to list directory %q: %v", path, err)
			return filepath.SkipDir
		}

		// no valid package starts with '.' and we don't want to search in '.git' and similar
		if strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}

		if _, err := os.Lstat(filepath.Join(path, DocFile)); err == nil {
			log.Printf("DEBUG - Adding documentation package: %q", path)
			retPkgs[path] = val
		}
		return nil
	})
	if err != nil {
		log.Printf("ERROR - Unable to walk the path %q: %v", root, err)
	}
	return retPkgs
}

// WriteDocs generates documentation for the packages 'dtPkgs' and writes it to
// files.
// If linkDocPkgs is filled it will be used to link to packages instead of
// reporting all the details in one table.
func WriteDocs(
	dtPkgs []string,
	depMap DependencyMap,
	linkDocPkgs map[string]struct{},
	cfg config.Config,
	rootPkg, root string,
) {
	for _, dtPkg := range dtPkgs {
		delete(linkDocPkgs, dtPkg)
		writeDoc(dtPkg, depMap, linkDocPkgs, cfg, rootPkg, root)
		linkDocPkgs[dtPkg] = struct{}{}
	}
}

func writeDoc(
	dtPkg string,
	depMap DependencyMap,
	linkDocPkgs map[string]struct{},
	cfg config.Config,
	rootPkg, root string,
) {
	doc := GenerateTable(depMap, linkDocPkgs, cfg, rootPkg, dtPkg)
	if doc == "" {
		return
	}
	err := ioutil.WriteFile(filepath.Join(root, dtPkg, DocFile), []byte(doc), 0644)
	if err != nil {
		log.Printf("ERROR - Unable to write dependency table to file: %v", err)
	}
}

// GenerateTable generates the dependency matrix for the package 'relPkg'.
func GenerateTable(
	depMap DependencyMap,
	linkDocPkgs map[string]struct{},
	cfg config.Config,
	rootPkg, relPkg string,
) string {
	depMap = filterDepMap(depMap, relPkg, linkDocPkgs)
	if len(depMap) == 0 {
		log.Printf("INFO - Won't write doc for package %q because it has no dependencies.", relPkg)
		return ""
	}
	allRows := make([]string, 0, len(depMap))
	allCols := make([]string, 0, len(depMap)*2)
	allColsMap := make(map[string]pkgType, len(depMap)*2)

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
	intro := `# Dependency Table For: ` + path.Join(rootPkg, relPkg) + `

| `
	sb.WriteString(intro)

	// (column) header line: | | C o l 1 - G | C o l 2 | ... | C o l N - T |
	for _, col := range allCols {
		sb.WriteString("| ")
		if _, ok := linkDocPkgs[col]; ok {
			sb.WriteRune('[')
		}
		for _, r := range col {
			sb.WriteRune(r)
			sb.WriteRune(' ')
		}
		if _, ok := linkDocPkgs[col]; ok {
			sb.WriteString("](")
			sb.WriteString(path.Join(RelPath(relPkg, col), DocFile))
			sb.WriteString(") ")
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

func filterDepMap(allMap DependencyMap, startPkg string, linkMap map[string]struct{}) DependencyMap {
	if (startPkg == "/" || startPkg == "") && len(linkMap) == 0 {
		return allMap
	}
	if _, ok := allMap[startPkg]; !ok {
		log.Printf("ERROR - Unable to find start package %q for dependency table.", startPkg)
		return nil
	}
	fltrMap := make(DependencyMap, len(allMap))
	copyDepsRecursive(allMap, startPkg, fltrMap, linkMap)
	return fltrMap
}
func copyDepsRecursive(
	allMap DependencyMap,
	startPkg string,
	fltrMap DependencyMap,
	linkMap map[string]struct{},
) {
	if _, ok := linkMap[startPkg]; ok {
		return
	}
	imps, ok := allMap[startPkg]
	if !ok {
		return
	}
	fltrMap[startPkg] = imps
	for pkg := range imps.Imports {
		copyDepsRecursive(allMap, pkg, fltrMap, linkMap)
	}
}

// RelPath calculates the relative path from 'basepath' to 'targetpath'.
// This is very similar to filepath.Rel() but not OS specific but it is working
// by purely lexical processing like the path package.
func RelPath(basepath, targetpath string) string {
	base := splitPath(path.Clean(basepath))
	targ := splitPath(path.Clean(targetpath))

	n := len(base)
	m := len(targ)
	i := 0
	for i < n && i < m && base[i] == targ[i] {
		i++
	}

	ret := ""
	for j := i; j < n; j++ { // go backward for base
		ret = path.Join(ret, "..")
	}
	for j := i; j < m; j++ { // go forward for target
		ret = path.Join(ret, targ[j])
	}

	return ret
}
func splitPath(p string) []string {
	ret := make([]string, 64)
	for p != "" {
		base, last := path.Split(p)
		ret = append(ret, last)
		p = base
	}
	return reverse(ret)
}
func reverse(ss []string) []string {
	n := len(ss)
	ts := make([]string, n)
	n--
	for i, s := range ss {
		ts[n-i] = s
	}
	return ts
}

// Check checks the dependencies of the given package and reports offending
// imports.
func Check(pkg *pkgs.Package, rootPkg string, cfg config.Config, depMap *DependencyMap) []error {
	relPkg, strictRelPkg := pkgs.RelativePackageName(pkg, rootPkg)
	checkSpecial := checkStandard
	pkgImps := pkgImports{}

	var fullmatch, matchDB bool
	if _, fullmatch = isPackageInList(cfg.God, nil, relPkg, strictRelPkg); fullmatch {
		checkSpecial = checkGod
		pkgImps.PkgType = typeGod
	}
	if matchDB, fullmatch = isPackageInList(cfg.DB, nil, relPkg, strictRelPkg); matchDB {
		if fullmatch {
			checkSpecial = checkDB
			pkgImps.PkgType = typeDB
		} else {
			checkSpecial = checkHalfDB
		}
	}
	if matchTool, fullmatch := isPackageInList(cfg.Tool, nil, relPkg, strictRelPkg); matchTool {
		if fullmatch {
			checkSpecial = checkTool
			pkgImps.PkgType = typeTool
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

	if _, full := isPackageInList(cfg.Tool, nil, relImp, strictRelImp); full {
		im[unqImp] = typeTool
	} else if _, full := isPackageInList(cfg.DB, nil, relImp, strictRelImp); full {
		im[unqImp] = typeDB
	} else if _, full := isPackageInList(cfg.God, nil, relImp, strictRelImp); full {
		im[unqImp] = typeGod
	} else {
		im[unqImp] = typeStandard
	}
	return im
}
