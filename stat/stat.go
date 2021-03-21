package stat

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/flowdev/spaghetti-cutter/data"
)

// FileName is the name of the statistics file (package_statistics.md).
const FileName = "package_statistics.md"

var mapValue struct{} = struct{}{}

// Constants for the title of the sub-sections of the stat docs.
const (
	titleImps     = `Direct Dependencies (Imports) Of `
	titleAllImps  = `All (Including Transitive) Dependencies (Imports) Of `
	titleUsers    = `Packages Using (Importing) `
	titleMaxScore = `Packages Shielded From Users Of `
	titleMinScore = `Packages Shielded From All Users Of `
)

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
		maxScoreInt, maxScoreMap := maxScore(pkg, allImps, users, depMap)
		minScoreMap := minScore(pkg, allImps, users, depMap)

		pkgTitle := title(pkg)
		linkPkg := `[` + pkg + `](` + fragment(pkgTitle) + `)`
		linkType := `[ \[` + string(data.TypeLetter(pkgImps.PkgType)) + `\] ](#legend)`

		linkImps := `0`
		if len(pkgImps.Imports) > 0 {
			linkImps = `[` + strconv.Itoa(len(pkgImps.Imports)) + `](` + fragment(titleImps+pkgTitle) + `)`
		}

		linkAllImps := `0`
		if len(allImps) > 0 {
			linkAllImps = `[` + strconv.Itoa(len(allImps)) + `](` + fragment(titleAllImps+pkgTitle) + `)`
		}

		linkUsers := `0`
		if len(users) > 0 {
			linkUsers = `[` + strconv.Itoa(len(users)) + `](` + fragment(titleUsers+pkgTitle) + `)`
		}

		linkMaxScore := `0`
		if maxScoreInt > 0 {
			linkMaxScore = `[` + strconv.Itoa(maxScoreInt) + `](` + fragment(titleMaxScore+pkgTitle) + `)`
		}

		linkMinScore := `0`
		if len(minScoreMap) > 0 {
			linkMinScore = `[` + strconv.Itoa(len(minScoreMap)) + `](` + fragment(titleMinScore+pkgTitle) + `)`
		}

		sb.WriteString(
			fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s |\n",
				linkPkg,
				linkType,
				linkImps,
				linkAllImps,
				linkUsers,
				linkMaxScore,
				linkMinScore,
			),
		)

		allSectionsEmpty := true
		sb2.WriteString(`

### ` + pkgTitle + `

`)
		if len(allImps) > 0 {
			sb2.WriteString(`
#### ` + titleImps + pkgTitle + `
`)
			writeImportLinks(sb2, pkgImps.Imports, depMap)
			allSectionsEmpty = false
		}
		if len(allImps) > 0 {
			sb2.WriteString(`

#### ` + titleAllImps + pkgTitle + `
`)
			writePackageLinks(sb2, allImps, depMap)
			allSectionsEmpty = false
		}
		if len(users) > 0 {
			sb2.WriteString(`

#### ` + titleUsers + pkgTitle + `
`)
			writeFragmentLinks(sb2, users, depMap)
			allSectionsEmpty = false
		}
		if len(maxScoreMap) > 0 {
			sb2.WriteString(`

#### ` + titleMaxScore + pkgTitle + `
`)
			for p, noImps := range maxScoreMap {
				sb2.WriteString("* " + fragmentLink(p) + ": ")
				writePackageLinks(sb2, noImps, depMap)
				sb2.WriteString("\n")
			}
			allSectionsEmpty = false
		}
		if len(minScoreMap) > 0 {
			sb2.WriteString(`

#### ` + titleMinScore + pkgTitle + `
`)
			writePackageLinks(sb2, minScoreMap, depMap)
			allSectionsEmpty = false
		}
		if allSectionsEmpty {
			sb2.WriteString(`No additional data.
`)

		}
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
	sb.WriteString(sb2.String())
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
	sort.Strings(users)
	return users
}

func maxScore(pkg string, imps map[string]struct{}, users []string, depMap data.DependencyMap,
) (int, map[string]map[string]struct{}) {
	sc := 0
	sm := make(map[string]map[string]struct{}, len(users))
	for _, u := range users {
		m := minus(imps, overlap(imps, depsWithoutPkg(u, pkg, depMap)))
		n := len(m)
		if n > 0 {
			sc += n
			sm[u] = m
		}
	}
	return sc, sm
}

func minScore(pkg string, imps map[string]struct{}, users []string, depMap data.DependencyMap) map[string]struct{} {
	if len(users) == 0 {
		return nil
	}

	usrsDeps := make(map[string]struct{}, 128)
	for _, u := range users {
		addToFirst(usrsDeps, depsWithoutPkg(u, pkg, depMap))
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

func addToFirst(all, m map[string]struct{}) {
	for k := range m {
		all[k] = mapValue
	}
}

func writeImportLinks(sb *strings.Builder, imps map[string]data.PkgType, depMap data.DependencyMap) {
	sl := make([]string, 0, len(imps))
	for imp := range imps {
		sl = append(sl, imp)
	}
	writeFragmentLinks(sb, sl, depMap)
}

func writePackageLinks(sb *strings.Builder, pkgs map[string]struct{}, depMap data.DependencyMap) {
	sl := make([]string, 0, len(pkgs))
	for pkg := range pkgs {
		sl = append(sl, pkg)
	}
	writeFragmentLinks(sb, sl, depMap)
}

func writeFragmentLinks(sb *strings.Builder, pkgs []string, depMap data.DependencyMap) {
	sort.Strings(pkgs)
	for i, p := range pkgs {
		if i > 0 {
			sb.WriteString(", ")
		}
		if _, ok := depMap[p]; ok {
			sb.WriteString(fragmentLink(p))
		} else {
			sb.WriteString("`" + p + "`")
		}
	}
}

func title(pkg string) string {
	if pkg == "/" {
		return "Root Package"
	}
	return "Package " + pkg
}

func pkgName(pkg string) string {
	if pkg == "/" {
		return "root"
	}
	return pkg
}

func fragmentLink(pkg string) string {
	return `[` + pkgName(pkg) + `](` + fragment(title(pkg)) + `)`
}

var (
	notAlphaNums = regexp.MustCompile(`[^a-z0-9 ]+`) // is constant and would blow up at first test
	spaces       = regexp.MustCompile(`[ ]+`)        // is constant and would blow up at first test
)

func fragment(s string) string {
	return `#` + spaces.ReplaceAllString(
		notAlphaNums.ReplaceAllString(
			strings.ToLower(s),
			"",
		),
		"-",
	)
}
