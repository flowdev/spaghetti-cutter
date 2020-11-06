package config

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"

	"github.com/hjson/hjson-go"
)

// File is the name of the configuration file
const File = ".spaghetti-cutter.hjson"

// internal enum for '$' handling in patterns
const (
	enumNoDollar    = iota // '$' isn't allowed at all
	enumDollarStar         // '$' has to be followed by one or two '*'
	enumDollarDigit        // '$' has to be followed by a single digit (1-9)
)

// Pattern combines the original pattern string with a compiled regular
// expression ready for efficient evaluation.
type Pattern struct {
	pattern    string
	regexp     *regexp.Regexp
	dollarIdxs []int
}

// PatternList is a slice of the Pattern type.
type PatternList []Pattern

// String implements Stringer and returns the list of
// patterns, or "..." if no patterns have been added.
func (pl *PatternList) String() string {
	if pl == nil || len(*pl) <= 0 {
		return "..."
	}
	var b strings.Builder
	for i, p := range *pl {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString("`")
		b.WriteString(p.pattern)
		b.WriteString("`")
	}
	return b.String()
}

// MatchString returns true if any of the patterns in the pattern list matches
// the given string including its dollars and false otherwise.
func (pl *PatternList) MatchString(s string, dollars []string) bool {
	if pl == nil {
		return false
	}
	for _, p := range *pl {
		if m := p.regexp.FindStringSubmatch(s); len(m) > 0 {
			if matchDollars(dollars, m[1:], p.dollarIdxs) {
				return true
			}
		}
	}
	return false
}
func matchDollars(given, found []string, idxs []int) bool {
	for i, f := range found {
		if f != given[idxs[i]] {
			return false
		}
	}
	return true
}

type patternGroup struct {
	left  Pattern
	right *PatternList
}

// PatternMap is a map from a single pattern to a list of patterns.
type PatternMap map[string]patternGroup

// String implements Stringer and returns the map of patterns,
// or "....." if it is empty.
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
		b.WriteString((*pm)[left].right.String())
		b.WriteString(" ; ")
	}
	s := b.String()
	return s[:len(s)-3]
}

// MatchingList returns the PatternList and submatches in the key if any key of
// this pattern map matches the given string and nil otherwise.
func (pm *PatternMap) MatchingList(s string) (*PatternList, []string) {
	if pm == nil {
		return nil, nil
	}
	for _, group := range *pm {
		if m := group.left.regexp.FindStringSubmatch(s); len(m) > 0 {
			return group.right, m[1:]
		}
	}
	return nil, nil
}

// Config contains the parsed configuration.
type Config struct {
	AllowOnlyIn       *PatternMap
	AllowAdditionally *PatternMap
	Tool              *PatternList
	DB                *PatternList
	God               *PatternList
	Size              uint
	NoGod             bool
}

const (
	keyAllowOnlyIn       = "allowOnlyIn"
	keyAllowAdditionally = "allowAdditionally"
	keyTool              = "tool"
	keyDB                = "db"
	keyGod               = "god"
	keySize              = "size"
	keyNoGod             = "noGod"
)

type jsonConfig struct {
	AllowOnlyIn       map[string][]string `json:"allowOnlyIn,omitempty"`
	AllowAdditionally map[string][]string `json:"allowAdditionally,omitempty"`
	Tool              []string            `json:"tool,omitempty"`
	DB                []string            `json:"db,omitempty"`
	God               []string            `json:"god,omitempty"`
	Size              uint                `json:"size,omitempty"`
	NoGod             bool                `json:"noGod,omitempty"`
}

func convertFromJSON(jcfg map[string]interface{}) (Config, error) {
	var err error
	var size uint
	var noGod bool
	var pl *PatternList
	var pm *PatternMap

	cfg := Config{}

	if size, err = convertUIntFromJSON(jcfg[keySize]); err != nil {
		return Config{}, fmt.Errorf("unable to convert maximum package size from JSON: %w", err)
	}
	cfg.Size = size

	if noGod, err = convertBoolFromJSON(jcfg[keyNoGod]); err != nil {
		return Config{}, fmt.Errorf("unable to convert no-god flag from JSON: %w", err)
	}
	cfg.NoGod = noGod

	if pm, err = convertPatternMapFromJSON(jcfg[keyAllowOnlyIn], keyAllowOnlyIn); err != nil {
		return cfg, err
	}
	cfg.AllowOnlyIn = pm

	if pm, err = convertPatternMapFromJSON(jcfg[keyAllowAdditionally], keyAllowAdditionally); err != nil {
		return cfg, err
	}
	cfg.AllowAdditionally = pm

	if pl, err = convertPatternListFromJSON(jcfg[keyTool], keyTool, enumNoDollar, 0); err != nil {
		return cfg, err
	}
	cfg.Tool = pl

	if pl, err = convertPatternListFromJSON(jcfg[keyDB], keyDB, enumNoDollar, 0); err != nil {
		return cfg, err
	}
	cfg.DB = pl

	if pl, err = convertPatternListFromJSON(jcfg[keyGod], keyGod, enumNoDollar, 0); err != nil {
		return cfg, err
	}
	cfg.God = pl

	return cfg, nil
}

func convertPatternMapFromJSON(i interface{}, key string) (*PatternMap, error) {
	var err error
	var pl *PatternList
	var re *regexp.Regexp
	var dollars int

	if i == nil {
		return nil, nil
	}

	m, ok := i.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected string map for key '%s', got type: %T", key, i)
	}

	pm := PatternMap(make(map[string]patternGroup, len(m)))

	for k, v := range m {
		if re, dollars, _, err = regexpForPattern(k, enumDollarStar, 0); err != nil {
			return nil, fmt.Errorf("illegal left/key pattern %q for global key '%s': %w", k, key, err)
		}
		if pl, err = convertPatternListFromJSON(v, key+": "+k, enumDollarDigit, dollars); err != nil {
			return nil, err
		}
		pm[k] = patternGroup{
			left:  Pattern{pattern: k, regexp: re},
			right: pl,
		}
	}

	return &pm, nil
}

func convertPatternListFromJSON(i interface{}, key string, allowDollar int, keyDollars int) (*PatternList, error) {
	if i == nil {
		return nil, nil
	}

	sl, ok := i.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected string list for key '%s', got type: %T", key, i)
	}

	l := make([]Pattern, len(sl))
	for i, v := range sl {
		s, err := convertStringFromJSON(v)
		if err != nil {
			return nil, fmt.Errorf("unable to convert list for key '%s' from JSON: %w", key, err)
		}
		re, _, dollarIdxs, err := regexpForPattern(s, allowDollar, keyDollars)
		if err != nil {
			return nil, fmt.Errorf("unable to use pattern `%s` of key '%s': %w", s, key, err)
		}
		l[i] = Pattern{pattern: s, regexp: re, dollarIdxs: dollarIdxs}
	}
	pl := PatternList(l)
	return &pl, nil
}

func convertUIntFromJSON(i interface{}) (uint, error) {
	var f float64
	var ok bool

	if i == nil {
		return 0, nil
	}

	if f, ok = i.(float64); !ok {
		return 0, fmt.Errorf("expected positive integer value, got type: %T", i)
	}

	if f < 0.0 {
		return 0, fmt.Errorf("expected positive integer value, got negative: %f", f)
	}

	if f > math.MaxUint32 {
		return 0, fmt.Errorf("expected unsigned integer value, got too large: %f", f)
	}

	if f != math.Trunc(f) {
		return 0, fmt.Errorf("expected unsigned integer value, got float: %f", f)
	}

	return uint(f), nil
}

func convertBoolFromJSON(i interface{}) (bool, error) {
	var b, ok bool

	if i == nil {
		return false, nil
	}

	if b, ok = i.(bool); !ok {
		return false, fmt.Errorf("expected boolean value, got type: %T", i)
	}

	return b, nil
}

func convertStringFromJSON(i interface{}) (string, error) {
	var s string
	var ok bool

	if i == nil {
		return "", nil
	}

	if s, ok = i.(string); !ok {
		return "", fmt.Errorf("expected string value, got type: %T", i)
	}

	return s, nil
}

func regexpForPattern(pattern string, allowDollar int, maxDollar int) (*regexp.Regexp, int, []int, error) {
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
				case enumNoDollar:
					errText = noDollarErrorText
					break
				case enumDollarStar:
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
			if allowDollar == enumNoDollar {
				errText = noDollarErrorText
				return "<error>"
			}
			if s[1] == '\\' {
				switch allowDollar {
				case enumDollarStar:
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
	re, err := regexp.Compile("^" + pattern + "$")
	if dollarCount > 0 && allowDollar == enumDollarDigit {
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

// Parse parses the configuration bytes and uses cfgFile only for better error
// messages.
func Parse(cfgBytes []byte, cfgFile string) (Config, error) {
	cfg := Config{}
	var jsonCfg map[string]interface{}

	if err := hjson.Unmarshal(cfgBytes, &jsonCfg); err != nil {
		return Config{}, fmt.Errorf("unable to unmarshal JSON configuration from file %q: %w", cfgFile, err)
	}

	noGod, _ := convertBoolFromJSON(jsonCfg[keyNoGod])
	god, _ := convertPatternListFromJSON(jsonCfg[keyGod], keyGod, enumNoDollar, 0)
	if !noGod && (god == nil || len(*god) == 0) {
		jsonCfg[keyGod] = []interface{}{"main"} // default
	}

	if size, _ := convertUIntFromJSON(jsonCfg[keySize]); size == 0 {
		jsonCfg[keySize] = 2048.0
	}

	cfg, err := convertFromJSON(jsonCfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
