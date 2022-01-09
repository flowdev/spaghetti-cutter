package data

import (
	"sort"
	"strings"
)

type PatternGroup struct {
	Left  Pattern
	Right PatternList
}

// PatternMap is a map from a single pattern to a list of patterns.
type PatternMap map[string]PatternGroup

// String implements Stringer and returns the map of patterns,
// or "....." if it is empty.
// The stringer methods are tested in the config package: TestParseAndStringers
func (pm *PatternMap) String() string {
	if pm == nil || len(*pm) <= 0 {
		return "....."
	}
	keys := make([]string, 0, len(*pm))
	for k := range *pm {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	for _, left := range keys {
		b.WriteString("`")
		b.WriteString(left)
		b.WriteString("`")
		b.WriteString(": ")
		b.WriteString((*pm)[left].Right.String())
		b.WriteString(" ; ")
	}
	s := b.String()
	return s[:len(s)-3]
}

// HasKeyValue checks if this pattern map contains the given key value pair.
// The strict versions are checked first
// (1. strictKey+strictValue, 2. strictKey+value, 3. key+strictValue, 4. key+value).
func (pm *PatternMap) HasKeyValue(key, strictKey, value, strictValue string) (hasKey, hasValue bool) {
	if pm == nil {
		return false, false
	}

	for _, k := range []string{strictKey, key} {
		if k == "" {
			continue
		}
		for _, group := range *pm {
			if m := group.Left.Regexp.FindStringSubmatch(k); len(m) > 0 {
				dollars := m[1:]

				if _, full := group.Right.MatchString(strictValue, dollars); strictValue != "" && full {
					return true, true
				}
				if _, full := group.Right.MatchString(value, dollars); value != "" && full {
					return true, true
				}
				hasKey = true
			}
		}
	}
	return hasKey, false
}
