package stat

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/flowdev/spaghetti-cutter/data"
)

// FileName is the name of the statistics file (package_statistics.md).
const FileName = "package_statistics.md"

// Generate creates some statistics for each package in the filtered dependency
// map starting at startPkg:
// - the type of the package ('S', 'T', 'D' or 'G')
// - number of direct dependencies
// - number of dependencies including transitive dependencies
// - number of packages using it
// - maximum and minimum score for encapsulating/hiding transitive packages
func Generate(startPkg string, depMap data.DependencyMap) string {
	depMap = data.FilterDepMap(depMap, startPkg, nil)
	if len(depMap) == 0 {
		log.Printf("INFO - Won't write stats for package %q because it has no dependencies.", startPkg)
		return ""
	}

	pkgNames := sortPkgNames(depMap)
	allDeps := allDependencies(depMap)

	sb := &strings.Builder{}
	sb.WriteString(`# Package Statistics

Start package - ` + startPkg + `

max score - the sum of the packages hidden from user packages",
min score - the packages hidden from all user packages combined",

| package | type | direct deps | all deps | users | max score | min score |
| :- | :-: | -: | -: | -: | -: | -: |
`)

	for _, pkg := range pkgNames {
		pkgImps := depMap[pkg]
		users := pkgUsers(pkg, depMap)
		allImps := allDeps[pkg]
		sb.WriteString(
			fmt.Sprintf("| %s | [%c] | %d | %d | %d | %d | %d |\n",
				pkg,
				data.TypeLetter(pkgImps.PkgType),
				len(pkgImps.Imports),
				len(allImps),
				len(users),
				maxScore(pkg, allImps, users, depMap),
				minScore(pkg, allImps, users, depMap),
			),
		)
	}
	return sb.String()
}

func sortPkgNames(depMap data.DependencyMap) []string {
	names := make([]string, 0, len(depMap))
	for pkg := range depMap {
		names = append(names, pkg)
	}
	sort.Strings(names)
	return names
}

func allDependencies(depMap data.DependencyMap) map[string]map[string]struct{} {
	allDeps := make(map[string]map[string]struct{}, len(depMap))
	for pkg := range depMap {
		allPkgDeps := make(map[string]struct{}, 128)
		addRecursiveDeps(allPkgDeps, pkg, "", depMap)
		allDeps[pkg] = allPkgDeps
	}
	return allDeps
}

func addRecursiveDeps(allPkgDeps map[string]struct{}, startPkg, excludePkg string, depMap data.DependencyMap) {
	if startPkg == excludePkg {
		return
	}
	pkgImps, ok := depMap[startPkg]
	if !ok {
		return
	}
	for p := range pkgImps.Imports {
		if p == excludePkg {
			continue
		}
		allPkgDeps[p] = struct{}{}
		addRecursiveDeps(allPkgDeps, p, excludePkg, depMap)
	}
}

func pkgUsers(pkg string, depMap data.DependencyMap) []string {
	users := make([]string, 0, len(depMap))
	for p, imps := range depMap {
		if _, ok := imps.Imports[pkg]; ok {
			users = append(users, p)
		}
	}
	return users
}

func maxScore(pkg string, imps map[string]struct{}, users []string, depMap data.DependencyMap) int {
	s := 0
	is := len(imps)
	for _, u := range users {
		s += is - overlap(imps, depsWithoutPkg(u, pkg, depMap))
	}
	return s
}

func minScore(pkg string, imps map[string]struct{}, users []string, depMap data.DependencyMap) int {
	if len(users) == 0 {
		return 0
	}

	usrsDeps := make(map[string]struct{}, 128)
	for _, u := range users {
		addMap(usrsDeps, depsWithoutPkg(u, pkg, depMap))
	}
	return len(imps) - overlap(imps, usrsDeps)
}

func depsWithoutPkg(startPkg, excludePkg string, depMap data.DependencyMap) map[string]struct{} {
	usrDeps := make(map[string]struct{}, 128)
	addRecursiveDeps(usrDeps, startPkg, excludePkg, depMap)
	return usrDeps
}

func overlap(m1, m2 map[string]struct{}) int {
	o := 0
	for k := range m1 {
		if _, ok := m2[k]; ok {
			o++
		}
	}
	return o
}

func addMap(all, m map[string]struct{}) {
	for k := range m {
		all[k] = struct{}{}
	}
}
