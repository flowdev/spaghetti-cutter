package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"
)

// File is the name of the configuration file
const File = ".spaghetti-cutter.json"

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
	IgnoreVendor      bool
}

type jsonConfig struct {
	AllowOnlyIn       map[string][]string `json:"allowOnlyIn,omitempty"`
	AllowAdditionally map[string][]string `json:"allowAdditionally,omitempty"`
	Tool              []string            `json:"tool,omitempty"`
	DB                []string            `json:"db,omitempty"`
	God               []string            `json:"god,omitempty"`
	Size              uint                `json:"size,omitempty"`
	NoGod             bool                `json:"noGod,omitempty"`
	IgnoreVendor      bool                `json:"ignoreVendor,omitempty"`
}

func convertFromJSON(jcfg jsonConfig) (Config, error) {
	var err error
	var pl *PatternList
	var pm *PatternMap

	cfg := Config{
		Size:         jcfg.Size,
		NoGod:        jcfg.NoGod,
		IgnoreVendor: jcfg.IgnoreVendor,
	}

	if pm, err = convertPatternMapFromJSON(jcfg.AllowOnlyIn); err != nil {
		return cfg, err
	}
	cfg.AllowOnlyIn = pm

	if pm, err = convertPatternMapFromJSON(jcfg.AllowAdditionally); err != nil {
		return cfg, err
	}
	cfg.AllowAdditionally = pm

	if pl, err = convertPatternListFromJSON(jcfg.Tool); err != nil {
		return cfg, err
	}
	cfg.Tool = pl

	if pl, err = convertPatternListFromJSON(jcfg.DB); err != nil {
		return cfg, err
	}
	cfg.DB = pl

	if pl, err = convertPatternListFromJSON(jcfg.God); err != nil {
		return cfg, err
	}
	cfg.God = pl

	return cfg, nil
}

func convertPatternMapFromJSON(m map[string][]string) (*PatternMap, error) {
	var err error
	var pl *PatternList
	var re *regexp.Regexp
	pm := PatternMap(make(map[string]patternGroup, 16))

	for k, v := range m {
		if re, err = regexpForPattern(k); err != nil {
			return nil, fmt.Errorf("unable to set left/key pattern %q: %w", k, err)
		}
		if pl, err = convertPatternListFromJSON(v); err != nil {
			return nil, err
		}
		pm[k] = patternGroup{
			left:  Pattern{pattern: k, regexp: re},
			right: pl,
		}
	}

	return &pm, nil
}

func convertPatternListFromJSON(s []string) (*PatternList, error) {
	var err error
	var pl PatternList
	var re *regexp.Regexp

	l := make([]Pattern, len(s))
	for i, t := range s {
		if re, err = regexpForPattern(t); err != nil {
			return nil, fmt.Errorf("unable to use pattern `%s`: %w", t, err)
		}
		l[i] = Pattern{pattern: t, regexp: re}
	}
	pl = PatternList(l)
	return &pl, nil
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

// Parse parses the configuration file
func Parse(cfgFile string) (Config, error) {
	cfg := Config{}
	jsonCfg := jsonConfig{}
	cfgBytes, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return Config{}, fmt.Errorf("unable to read configuration file %q: %w", cfgFile, err)
	}
	if err = json.Unmarshal(cfgBytes, &jsonCfg); err != nil {
		return Config{}, fmt.Errorf("unable to unmarshal JSON configuration from file %q: %w", cfgFile, err)
	}

	if !jsonCfg.NoGod && len(jsonCfg.God) == 0 {
		jsonCfg.God = []string{"main"} // default
	}
	if jsonCfg.Size == 0 {
		jsonCfg.Size = 2048
	}

	if cfg, err = convertFromJSON(jsonCfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
