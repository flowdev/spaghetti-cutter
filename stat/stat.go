package stat

import (
	"fmt"
	"sort"

	"github.com/flowdev/spaghetti-cutter/data"
)

// Create creates some statistics for each package in the filtered dependency
// map starting at startPkg:
// - the type of the package ('S', 'T', 'D' or 'G')
// - number of direct dependencies
// - number of dependencies including transitive dependencies
// - number of packages using it
// - maximum and minimum score for encapsulating/hiding transitive packages
func Create(startPkg string, depMap data.DependencyMap) []string {
	const pkgHead = "package"
	depMap = data.FilterDepMap(depMap, startPkg, nil)
	stats := make([]string, 0, 9+len(depMap))

	maxPkgLen := maxPkgNameLen(depMap)
	maxPkgLen = max(maxPkgLen, len(pkgHead))
	pkgNames := sortPkgNames(depMap)
	allDeps := allDependencies(depMap)
	stats = append(stats,
		" S T A T I S T I C S",
		" ===================",
		"",
		" Start package - "+startPkg,
		"",
		" max score - the sum of the packages hidden from user packages",
		" min score - the packages hidden from all user packages combined",
		"",
		fmt.Sprintf(" %-*s | type | direct deps | all deps | usages | max score | min score",
			maxPkgLen, pkgHead),
	)
	for _, pkg := range pkgNames {
		pkgImps := depMap[pkg]
		users := pkgUsers(pkg, depMap)
		allImps := allDeps[pkg]
		stats = append(stats,
			fmt.Sprintf(" %-*s |  [%c] | %11d | %8d | %6d | %9d | %9d",
				maxPkgLen, pkg,
				data.TypeLetter(pkgImps.PkgType),
				len(pkgImps.Imports),
				len(allDeps[pkg]),
				len(users),
				maxScore(pkg, allImps, users, depMap),
				minScore(pkg, allImps, users, depMap),
			),
		)
	}
	return stats
}

func maxPkgNameLen(depMap data.DependencyMap) int {
	maxPkgLen := 0
	for pkg := range depMap {
		if len(pkg) > maxPkgLen {
			maxPkgLen = len(pkg)
		}
	}
	return maxPkgLen
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
		addRecursiveDeps(allPkgDeps, pkg, depMap)
		allDeps[pkg] = allPkgDeps
	}
	return allDeps
}

func addRecursiveDeps(allPkgDeps map[string]struct{}, pkg string, depMap data.DependencyMap) {
	pkgImps, ok := depMap[pkg]
	if !ok {
		return
	}
	for p := range pkgImps.Imports {
		allPkgDeps[p] = struct{}{}
		addRecursiveDeps(allPkgDeps, p, depMap)
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

func depsWithoutPkg(user, pkg string, depMap data.DependencyMap) map[string]struct{} {
	usrDeps := make(map[string]struct{}, 128)
	for p := range depMap[user].Imports {
		if p != pkg {
			usrDeps[p] = struct{}{}
			addRecursiveDeps(usrDeps, p, depMap)
		}
	}
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

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
