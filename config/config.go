package config

import (
	"fmt"
	"math"
	"regexp"

	"github.com/flowdev/spaghetti-cutter/data"
	"github.com/hjson/hjson-go"
)

// File is the name of the configuration file
const File = ".spaghetti-cutter.hjson"

// Config contains the parsed configuration.
type Config struct {
	AllowOnlyIn       *data.PatternMap
	AllowAdditionally *data.PatternMap
	Tool              data.PatternList
	DB                data.PatternList
	God               data.PatternList
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
	var pl data.PatternList
	var pm *data.PatternMap

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

	if pl, err = convertPatternListFromJSON(jcfg[keyTool], keyTool, data.EnumDollarNone, 0); err != nil {
		return cfg, err
	}
	cfg.Tool = pl

	if pl, err = convertPatternListFromJSON(jcfg[keyDB], keyDB, data.EnumDollarNone, 0); err != nil {
		return cfg, err
	}
	cfg.DB = pl

	if pl, err = convertPatternListFromJSON(jcfg[keyGod], keyGod, data.EnumDollarNone, 0); err != nil {
		return cfg, err
	}
	cfg.God = pl

	return cfg, nil
}

func convertPatternMapFromJSON(i interface{}, key string) (*data.PatternMap, error) {
	var err error
	var pl data.PatternList
	var re *regexp.Regexp
	var dollars int

	if i == nil {
		return nil, nil
	}

	m, ok := i.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected string map for key '%s', got type: %T", key, i)
	}

	pm := data.PatternMap(make(map[string]data.PatternGroup, len(m)))

	for k, v := range m {
		if re, dollars, _, err = data.RegexpForPattern(k, data.EnumDollarStar, 0); err != nil {
			return nil, fmt.Errorf("illegal left/key pattern %q for global key '%s': %w", k, key, err)
		}
		if pl, err = convertPatternListFromJSON(v, key+": "+k, data.EnumDollarDigit, dollars); err != nil {
			return nil, err
		}
		pm[k] = data.PatternGroup{
			Left:  data.Pattern{Pattern: k, Regexp: re},
			Right: pl,
		}
	}

	return &pm, nil
}

func convertPatternListFromJSON(i interface{}, key string, allowDollar data.EnumDollar, keyDollars int) (data.PatternList, error) {
	if i == nil {
		return nil, nil
	}

	sl, ok := i.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected string list for key '%s', got type: %T", key, i)
	}

	l := make([]data.Pattern, len(sl))
	for i, v := range sl {
		s, err := convertStringFromJSON(v)
		if err != nil {
			return nil, fmt.Errorf("unable to convert list for key '%s' from JSON: %w", key, err)
		}
		re, _, dollarIdxs, err := data.RegexpForPattern(s, allowDollar, keyDollars)
		if err != nil {
			return nil, fmt.Errorf("unable to use pattern `%s` of key '%s': %w", s, key, err)
		}
		l[i] = data.Pattern{Pattern: s, Regexp: re, DollarIdxs: dollarIdxs}
	}
	return data.PatternList(l), nil
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

// Parse parses the configuration bytes and uses cfgFile only for better error
// messages.
func Parse(cfgBytes []byte, cfgFile string) (Config, error) {
	cfg := Config{}
	var jsonCfg map[string]interface{}

	if err := hjson.Unmarshal(cfgBytes, &jsonCfg); err != nil {
		return Config{}, fmt.Errorf("unable to unmarshal JSON configuration from file %q: %w", cfgFile, err)
	}

	noGod, _ := convertBoolFromJSON(jsonCfg[keyNoGod])
	god, _ := convertPatternListFromJSON(jsonCfg[keyGod], keyGod, data.EnumDollarNone, 0)
	if !noGod && (god == nil || len(god) == 0) {
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
