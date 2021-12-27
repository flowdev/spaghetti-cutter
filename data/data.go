package data

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

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

// EnumDollar is an enumeration type for how to hanlde dollars ('$') in patterns.
type EnumDollar int

// internal enum for '$' handling in patterns
const (
	EnumDollarNone  EnumDollar = iota // '$' isn't allowed at all
	EnumDollarStar                    // '$' has to be followed by one or two '*'
	EnumDollarDigit                   // '$' has to be followed by a single digit (1-9)
)

// Pattern combines the original pattern string with a compiled regular
// expression ready for efficient evaluation.
type Pattern struct {
	Pattern    string
	Regexp     *regexp.Regexp
	DollarIdxs []int
}

// DependencyMap is mapping importing package to imported packages.
// importingPackageName -> (importedPackageNames -> PkgType)
// An imported package name could be added multiple times to the same importing
// package name due to test packages.
type DependencyMap map[string]PkgImports

// SortedPkgNames returns the sorted keys (package names) of the dependency map.
func (dm DependencyMap) SortedPkgNames() []string {
	names := make([]string, 0, len(dm))
	for pkg := range dm {
		names = append(names, pkg)
	}
	sort.Strings(names)
	return names
}

// FilterDepMap filters allMap to contain only packages matching idx and its transitive
// dependencies.  Entries matching other indices in links are filtered, too.
func FilterDepMap(allMap DependencyMap, idx int, links PatternList) DependencyMap {
	if idx < 0 || len(links) == 0 {
		return allMap
	}

	fltrMap := make(DependencyMap, len(allMap))
	for pkg := range allMap {
		if i := DocMatchStringIndex(pkg, links); i >= 0 && i == idx {
			copyDepsRecursive(allMap, pkg, fltrMap, links, idx)
		}
	}
	return fltrMap
}
func copyDepsRecursive(
	allMap DependencyMap,
	startPkg string,
	fltrMap DependencyMap,
	links PatternList,
	idx int,
) {
	if i := DocMatchStringIndex(startPkg, links); i >= 0 && i != idx {
		return
	}
	imps, ok := allMap[startPkg]
	if !ok {
		return
	}
	fltrMap[startPkg] = imps
	for pkg := range imps.Imports {
		copyDepsRecursive(allMap, pkg, fltrMap, links, idx)
	}
}

// DocMatchStringIndex matches pkg in links and returns its index.
// Only full matches are returned. If pkg doesn't match, pkg+"/" is tried.
// -1 is returned for no match.
func DocMatchStringIndex(pkg string, links PatternList) (idx int) {
	i, full := links.MatchStringIndex(pkg, nil)
	if full {
		return i
	}
	i, full = links.MatchStringIndex(pkg+"/", nil)
	if full {
		return i
	}
	return -1
}

// PkgForPattern returns the (parent) package of the given package pattern.
// If pkg doesn't contain any wildcard '*' the whole string is returned.
// Otherwise everything up to the last '/' before the wildcard or
// the empty string if there is no '/' before it.
func PkgForPattern(pkg string) string {
	i := strings.IndexRune(pkg, '*')
	if i < 0 {
		return pkg
	}
	i = strings.LastIndex(pkg[:i], "/")
	if i > 0 {
		return pkg[:i]
	}
	return ""
}

// RegexpForPattern converts the given pattern including wildcards and variables
// into a proper regular expression that can be used for matching.
func RegexpForPattern(pattern string, allowDollar EnumDollar, maxDollar int,
) (*regexp.Regexp, int, []int, error) {
	const noDollarErrorText = "a '$' has to be escaped for this configuration key"
	const dollarStarErrorText = "a '$' has to be escaped or followed by one or two unescaped '*'s"
	const dollarDigitErrorText = "a '$' has to be escaped or followed by a single digit (1-9)"
	const singleStarPattern = `(?:[^/]*)`
	const doubleStarPattern = `(?:.*)`
	re := regexp.MustCompile(`(?:\\*\$)?(?:\\*\*\*?|[1-9])?`) // constant and tested by ANY unit test
	errText := ""
	dollarCount := 0
	dollarIdxs := make([]int, 0, 9)

	pattern = re.ReplaceAllStringFunc(pattern, func(s string) string {
		if s == "" {
			return s
		}

		if len(s) == 1 {
			switch s {
			case "$":
				switch allowDollar {
				case EnumDollarNone:
					errText = noDollarErrorText
					break
				case EnumDollarStar:
					errText = dollarStarErrorText
					break
				default:
					errText = dollarDigitErrorText
				}
				return "<error>"
			case "*":
				return singleStarPattern
			default:
				return s
			}
		}

		prefix := ``
		if n := countBackslashes(s); n > 0 && len(s) > n && s[n] == '$' {
			if n%2 == 0 { // even number of `\`: '$' is NOT escaped!
				prefix = s[:n]
				s = s[n:]
			} else { // odd number of `\`: '$' is escaped
				prefix = s[:n+1]
				s = s[n+1:]
			}
		}
		if s == "" {
			return prefix
		}
		if len(s) == 1 {
			if s == "*" {
				return prefix + singleStarPattern
			}
			return prefix + s // s is a digit
		}

		if n := countBackslashes(s); n > 0 {
			if n%2 == 0 { // even number of `\`: '*' is NOT escaped!
				prefix += s[:n]
				s = s[n:]
			} else { // odd number of `\`: '*' is escaped
				prefix += s[:n+1]
				s = s[n+1:]
				if s == "" {
					return prefix
				}
			}
		}
		if s[0] == '$' {
			if allowDollar == EnumDollarNone {
				errText = noDollarErrorText
				return "<error>"
			}
			if s[1] == '\\' {
				switch allowDollar {
				case EnumDollarStar:
					errText = dollarStarErrorText
					break
				default:
					errText = dollarDigitErrorText
				}
				return `<error>`
			}
			dollarCount++
			if len(s) > 2 {
				return prefix + `(.*)`
			}
			if s[1] == '*' {
				return prefix + `([^/]*)`
			}

			// DIGIT: remember index and capture the value
			dollarIdx := int(s[1] - '1')
			if dollarIdx >= maxDollar {
				errText = fmt.Sprintf("the maximum possible dollar index is %d, found index %d", maxDollar, dollarIdx+1)
				return `<error>`
			}
			dollarIdxs = append(dollarIdxs, dollarIdx)
			return prefix + `(.*)`
		}
		if len(s) > 1 {
			return prefix + doubleStarPattern
		}
		return prefix + singleStarPattern
	})

	if errText != "" {
		return nil, 0, nil, fmt.Errorf("%s; resulting regular expression: %s", errText, pattern)
	}

	var err error
	if allowDollar == EnumDollarStar {
		re, err = regexp.Compile("^" + pattern + "$")
	} else {
		re, err = regexp.Compile("^" + pattern)
	}
	if dollarCount > 0 && allowDollar == EnumDollarDigit {
		return re, dollarCount, dollarIdxs, err
	}
	return re, dollarCount, nil, err
}
func countBackslashes(s string) int {
	count := 0
	for _, r := range s {
		if r != '\\' {
			return count
		}
		count++
	}
	return count
}
