package data

import "log"

// PkgType can be one of: Standard, Tool, DB or God
type PkgType int

// Enum of package types: Standard, Tool, DB and God
const (
	TypeStandard PkgType = iota
	TypeTool
	TypeDB
	TypeGod
)

var typeLetters = []rune("STDG")
var typeFormats = []string{"", "_", "`", "**"}

// PkgImports contains the package type and the imported internal packages with their types.
type PkgImports struct {
	PkgType PkgType
	Imports map[string]PkgType
}

// DependencyMap is mapping importing package to imported packages.
// importingPackageName -> (importedPackageName -> PkgType)
// An imported package name could be added multiple times to the same importing
// package name due to test packages.
type DependencyMap map[string]PkgImports

// TypeLetter returns the type letter associated with package type t ('S', 'T',
// 'D' or 'G').
func TypeLetter(t PkgType) rune {
	return typeLetters[t]
}

// TypeFormat returns the formatting string associated with package type t ("",
// "_", "`" or "**").
func TypeFormat(t PkgType) string {
	return typeFormats[t]
}

// FilterDepMap filters allMap to contain only startPkg and it's transitive
// dependencies.  Entries in linkMap are filtered, too.
func FilterDepMap(allMap DependencyMap, startPkg string, linkMap map[string]struct{}) DependencyMap {
	if (startPkg == "/" || startPkg == "") && len(linkMap) == 0 {
		return allMap
	}
	if _, ok := allMap[startPkg]; !ok {
		log.Printf("ERROR - Unable to find start package %q for dependency table.", startPkg)
		return nil
	}
	fltrMap := make(DependencyMap, len(allMap))
	CopyDepsRecursive(allMap, startPkg, fltrMap, linkMap)
	return fltrMap
}

// CopyDepsRecursive copies dependencies recursively from allMap into fltrMap
// starting at startPkg and ignoring entries in linkMap.
func CopyDepsRecursive(
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
		CopyDepsRecursive(allMap, pkg, fltrMap, linkMap)
	}
}
