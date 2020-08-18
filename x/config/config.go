package config

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"

	"github.com/hjson/hjson-go"
)

// File is the name of the configuration file
const File = ".spaghetti-cutter.hjson"

// Pattern combines the original pattern string with a compiled regular
// expression ready for efficient evaluation.
type Pattern struct {
	pattern string
	regexp  *regexp.Regexp
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
// the given string and false otherwise.
func (pl *PatternList) MatchString(s string) bool {
	if pl == nil {
		return false
	}
	for _, p := range *pl {
		if p.regexp.MatchString(s) {
			return true
		}
	}
	return false
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

// MatchingList returns the PatternList if any key of this pattern map matches
// the given string and nil otherwise.
func (pm *PatternMap) MatchingList(s string) *PatternList {
	if pm == nil {
		return nil
	}
	for _, group := range *pm {
		if group.left.regexp.MatchString(s) {
			return group.right
		}
	}
	return nil
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

	if pl, err = convertPatternListFromJSON(jcfg[keyTool], keyTool); err != nil {
		return cfg, err
	}
	cfg.Tool = pl

	if pl, err = convertPatternListFromJSON(jcfg[keyDB], keyDB); err != nil {
		return cfg, err
	}
	cfg.DB = pl

	if pl, err = convertPatternListFromJSON(jcfg[keyGod], keyGod); err != nil {
		return cfg, err
	}
	cfg.God = pl

	return cfg, nil
}

func convertPatternMapFromJSON(i interface{}, key string) (*PatternMap, error) {
	var err error
	var pl *PatternList
	var re *regexp.Regexp

	if i == nil {
		return nil, nil
	}

	m, ok := i.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected string map for key '%s', got type: %T", key, i)
	}

	pm := PatternMap(make(map[string]patternGroup, 16))

	for k, v := range m {
		if re, err = regexpForPattern(k); err != nil {
			return nil, fmt.Errorf("illegal left/key pattern %q for global key '%s': %w", k, key, err)
		}
		if pl, err = convertPatternListFromJSON(v, key); err != nil {
			return nil, err
		}
		pm[k] = patternGroup{
			left:  Pattern{pattern: k, regexp: re},
			right: pl,
		}
	}

	return &pm, nil
}

func convertPatternListFromJSON(i interface{}, key string) (*PatternList, error) {
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
		re, err := regexpForPattern(s)
		if err != nil {
			return nil, fmt.Errorf("unable to use pattern `%s` of key '%s': %w", s, key, err)
		}
		l[i] = Pattern{pattern: s, regexp: re}
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

func regexpForPattern(pattern string) (*regexp.Regexp, error) {
	i := strings.Index(pattern, "**")
	n2 := len(pattern) - 2
	if i >= 0 && i < n2 {
		return nil, errors.New("illegal pattern `" + pattern + "` contains `**` before the end")
	}
	if i >= 0 {
		pattern = pattern[:i]
	}
	b := strings.Builder{}
	parts := strings.Split(pattern, "*")
	n := len(parts) - 1
	for j, s := range parts {
		if j < n {
			if len(s) > 0 && s[len(s)-1] == '\\' {
				b.WriteString(regexp.QuoteMeta(s[:len(s)-1]))
				b.WriteString("\\*")
			} else {
				b.WriteString(regexp.QuoteMeta(s))
				b.WriteString("[^/]*")
			}
		} else {
			b.WriteString(regexp.QuoteMeta(s))
		}
	}
	if i >= 0 {
		b.WriteString(".*")
	}
	re := "^" + b.String() + "$"
	return regexp.Compile(re)
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
	god, _ := convertPatternListFromJSON(jsonCfg[keyGod], keyGod)
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
