package doc

import (
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/flowdev/spaghetti-cutter/data"
)

// FileName is the name of the documentation file (package_dependencies.md).
const FileName = "package_dependencies.md"

// WriteDocs generates documentation for the packages 'dtPkgs' and writes it to
// files.
// If linkDocPkgs is filled it will be used to link to packages instead of
// reporting all the details in one table.
func WriteDocs(
	dtPkgs []string,
	depMap data.DependencyMap,
	linkDocPkgs map[string]struct{},
	rootPkg, root string,
) {
	for _, dtPkg := range dtPkgs {
		delete(linkDocPkgs, dtPkg)
		writeDoc(dtPkg, depMap, linkDocPkgs, rootPkg, root)
		linkDocPkgs[dtPkg] = struct{}{}
	}
}

func writeDoc(
	dtPkg string,
	depMap data.DependencyMap,
	linkDocPkgs map[string]struct{},
	rootPkg, root string,
) {
	doc := GenerateTable(depMap, linkDocPkgs, rootPkg, dtPkg)
	if doc == "" {
		return
	}
	docFile := filepath.Join(dtPkg, FileName)
	log.Printf("INFO - Write dependency table to file: %s", docFile)
	docFile = filepath.Join(root, docFile)
	err := ioutil.WriteFile(docFile, []byte(doc), 0644)
	if err != nil {
		log.Printf("ERROR - Unable to write dependency table to file %s: %v", docFile, err)
	}
}

// GenerateTable generates the dependency matrix for the package 'relPkg'.
func GenerateTable(
	depMap data.DependencyMap,
	linkDocPkgs map[string]struct{},
	rootPkg, relPkg string,
) string {
	depMap = data.FilterDepMap(depMap, relPkg, linkDocPkgs)
	if len(depMap) == 0 {
		log.Printf("INFO - Won't write doc for package %q because it has no dependencies.", relPkg)
		return ""
	}
	allRows := make([]string, 0, len(depMap))
	allCols := make([]string, 0, len(depMap)*2)
	allColsMap := make(map[string]data.PkgType, len(depMap)*2)

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
		letter := data.TypeLetter(allColsMap[col])
		sb.WriteString("- ")
		sb.WriteRune(letter)
		if _, ok := linkDocPkgs[col]; ok {
			sb.WriteString("](")
			sb.WriteString(path.Join(RelPath(relPkg, col), FileName))
			sb.WriteString(") ")
		}
		sb.WriteRune(' ')
	}
	sb.WriteString("|\n")

	// separator line: | :- | :-: | :-: | ... | :-: |
	sb.WriteString("| :- ")
	for range allCols {
		sb.WriteString("| :- ")
	}
	sb.WriteString("|\n")

	// normal rows: | **Row1** | **G** | | ... | **T** |
	for _, row := range allRows {
		pkgImps := depMap[row]

		sb.WriteString("| ")
		format := data.TypeFormat(pkgImps.PkgType)
		sb.WriteString(format)
		sb.WriteString(row)
		sb.WriteString(format)
		sb.WriteRune(' ')

		for _, col := range allCols {
			sb.WriteString("| ")
			if impType, ok := pkgImps.Imports[col]; ok {
				sb.WriteString(format)
				sb.WriteRune(data.TypeLetter(impType))
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

* **Bold** - God package (can use all packages)
` + "* `Code` - Database package (can only use tool and other database packages)" + `
* _Italic_ - Tool package (foundational, no dependencies)
* no formatting - Standard package (can only use tool and database packages)


#### Meaning Of Letters In Table Columns

* G - God package (can use all packages)
* D - Database package (can only use tool and other database packages)
* T - Tool package (foundational, no dependencies)
* S - Standard package (can only use tool and database packages)
`
	sb.WriteString(legend)
	return sb.String()
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
	for j := i; j < n; j++ { // go backward from base
		ret = path.Join(ret, "..")
	}
	for j := i; j < m; j++ { // go forward to target
		ret = path.Join(ret, targ[j])
	}

	return ret
}
func splitPath(p string) []string {
	ret := make([]string, 0, 64)
	for p != "" {
		base, last := path.Split(p)
		ret = append(ret, last)
		p = removeTrailingSlash(base)
	}
	return reverse(ret)
}
func removeTrailingSlash(p string) string {
	if !strings.HasSuffix(p, "/") {
		return p
	}
	return p[:len(p)-1]
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
