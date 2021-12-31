package data

import (
	"fmt"
	"regexp"
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
