package data

import (
	"fmt"
	"strings"
)

// PatternList is a slice of the Pattern type.
type PatternList []Pattern

// NewSimplePatternList returns a PatternList initialized from the given patterns.
// An error is returned if a pattern isn't valid.
func NewSimplePatternList(patterns []string, key string) (PatternList, error) {
	l := make([]Pattern, len(patterns))
	for i, p := range patterns {
		re, _, dollarIdxs, err := RegexpForPattern(p, EnumDollarNone, 0)
		if err != nil {
			return nil, fmt.Errorf("unable to use pattern `%s` of key '%s': %w", p, key, err)
		}
		l[i] = Pattern{Pattern: p, Regexp: re, DollarIdxs: dollarIdxs}
	}
	pl := PatternList(l)
	return pl, nil
}

// String implements Stringer and returns the list of
// patterns, or "..." if no patterns have been added.
func (pl PatternList) String() string {
	if pl == nil || len(pl) <= 0 {
		return "..."
	}
	var b strings.Builder
	for i, p := range pl {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString("`")
		b.WriteString(p.Pattern)
		b.WriteString("`")
	}
	return b.String()
}

// MatchString returns true if any of the patterns in the pattern list matches
// the given string including its dollars and false otherwise.
// If the full string is matched full will be true and false otherwise.
func (pl PatternList) MatchString(s string, dollars []string) (atAll, full bool) {
	idx, full := pl.MatchStringIndex(s, dollars)
	return idx >= 0, full
}

// MatchStringIndex returns the index of the pattern in the list that matches
// the given string and an indicator if it was a full match.
func (pl PatternList) MatchStringIndex(s string, dollars []string) (idx int, full bool) {
	if pl == nil {
		return -1, false
	}
	idx = -1
	for i, p := range pl {
		if m := p.Regexp.FindStringSubmatch(s); len(m) > 0 {
			if matchDollars(dollars, m[1:], p.DollarIdxs) {
				lenm := len(m[0])
				if lenm >= len(s) {
					return i, true
				}
				if s[lenm] == '/' {
					idx = i
				}
			}
		}
	}
	return idx, false
}
func matchDollars(given, found []string, idxs []int) bool {
	for i, f := range found {
		if f != given[idxs[i]] {
			return false
		}
	}
	return true
}
