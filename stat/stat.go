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

var mapValue struct{} = struct{}{}

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
	sb2 := &strings.Builder{}
	sb.WriteString(`# Package Statistics

Start package - ` + startPkg + `

| package | type | direct deps | all deps | users | max score | min score |
| :- | :-: | -: | -: | -: | -: | -: |
`)

	for _, pkg := range pkgNames {
		pkgImps := depMap[pkg]
		allImps := allDeps[pkg]
		users := pkgUsers(pkg, depMap)
		maxSc, _ := maxScore(pkg, allImps, users, depMap)
		minScMap := minScore(pkg, allImps, users, depMap)
		sb.WriteString(
			fmt.Sprintf("| %s | [%c] | %d | %d | %d | %d | %d |\n",
				pkg,
				data.TypeLetter(pkgImps.PkgType),
				len(pkgImps.Imports),
				len(allImps),
				len(users),
				maxSc,
				len(minScMap),
			),
		)
		sb2.WriteString(`
### ` + title(pkg) + `

#### Direct Dependencies (Imports)

`)
		for imp := range pkgImps.Imports {
			if _, ok := depMap[imp]; ok {
				sb2.WriteString(`* [` + imp + `](` + fragmentLink(imp) + `)
`)
			} else {
				sb2.WriteString("* `" + imp + "`\n")
			}
		}

		sb2.WriteString(`
#### All Dependencies (Imports) Including Transitive Dependencies
`)
		sb2.WriteString(`
#### Packages Using (Importing) This Package
`)

		sb2.WriteString(`
#### Packages Not Imported By Users
`)

		sb2.WriteString(`
#### Packages Not Imported By All Users Combined
`)
	}

	sb.WriteString(`
### Legend

* package - name of the internal package without the part common to all packages.
* type - type of the package:
  * [G] - God package (can use all packages)
  * [D] - Database package (can only use tool and other database packages)
  * [T] - Tool package (foundational, no dependencies)
  * [S] - Standard package (can only use tool and database packages)
* direct deps - number of internal packages directly imported by this one.
* all deps - number of transitive internal packages imported by this package.
* users - number of internal packages that import this one.
* max score - sum of the numbers of packages hidden from user packages.
* min score - number of packages hidden from all user packages combined.
`)
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
		allPkgDeps[p] = mapValue
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

func maxScore(pkg string, imps map[string]struct{}, users []string, depMap data.DependencyMap,
) (int, map[string]map[string]struct{}) {
	sc := 0
	sm := make(map[string]map[string]struct{}, len(users))
	for _, u := range users {
		m := minus(imps, overlap(imps, depsWithoutPkg(u, pkg, depMap)))
		sc += len(m)
		sm[u] = m
	}
	return sc, sm
}

func minScore(pkg string, imps map[string]struct{}, users []string, depMap data.DependencyMap) map[string]struct{} {
	if len(users) == 0 {
		return nil
	}

	usrsDeps := make(map[string]struct{}, 128)
	for _, u := range users {
		addMap(usrsDeps, depsWithoutPkg(u, pkg, depMap))
	}
	return minus(imps, overlap(imps, usrsDeps))
}

func depsWithoutPkg(startPkg, excludePkg string, depMap data.DependencyMap) map[string]struct{} {
	usrDeps := make(map[string]struct{}, 128)
	addRecursiveDeps(usrDeps, startPkg, excludePkg, depMap)
	return usrDeps
}

func minus(m1, m2 map[string]struct{}) map[string]struct{} {
	m := make(map[string]struct{}, len(m2))
	for k := range m1 {
		if _, ok := m2[k]; !ok {
			m[k] = mapValue
		}
	}
	return m
}

func overlap(m1, m2 map[string]struct{}) map[string]struct{} {
	o := make(map[string]struct{}, len(m2))
	for k := range m1 {
		if _, ok := m2[k]; ok {
			o[k] = mapValue
		}
	}
	return o
}

func addMap(all, m map[string]struct{}) {
	for k := range m {
		all[k] = mapValue
	}
}

func title(pkg string) string {
	if pkg == "/" {
		return "Root Package"
	}
	return "Package " + pkg
}

func fragmentLink(pkg string) string {
	return "#" + strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ToLower(
				title(pkg),
			),
			" ",
			"-",
		),
		"/",
		"",
	)
}
