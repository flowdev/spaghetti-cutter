package config

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/peterbourgon/ff/v3"
)

var (
	// Value is the set value for own maps that are really sets.
	Value = struct{}{}
)

// StringSet is a flag.Value that collects each Set string
// into a set, allowing for repeated flags.
type StringSet map[string]struct{}

// Set implements flag.Value and appends the string to the slice.
func (ss *StringSet) Set(s string) error {
	(*ss)[s] = Value
	return nil
}

// String implements flag.Value and returns the list of
// strings, or "..." if no strings have been added.
func (ss *StringSet) String() string {
	return setToString(*ss)
}

// MapSet is a flag.Value that collects each Set string
// into a map of sets, allowing for repeated flags.
type MapSet map[string]map[string]struct{}

// Set implements flag.Value and adds key and value (seperated by space in s)
// to the MapSet.
func (ms *MapSet) Set(s string) error {
	var key, value string
	_, err := fmt.Sscan(s, &key, &value)
	if err != nil {
		return fmt.Errorf("unable to split '%s' into key and value: %w", s, err)
	}

	set := (*ms)[key]
	if set == nil {
		set = make(map[string]struct{})
		(*ms)[key] = set
	}
	set[value] = Value
	return nil
}

// String implements flag.Value and returns the map of
// string sets, or "..." if no strings have been added.
func (ms *MapSet) String() string {
	if len(*ms) <= 0 {
		return "..."
	}
	var b strings.Builder
	for key, set := range *ms {
		b.WriteString(key)
		b.WriteString(": ")
		b.WriteString(setToString(set))
		b.WriteString(" ; ")
	}
	s := b.String()
	return s[:len(s)-3]
}

func setToString(set map[string]struct{}) string {
	if len(set) <= 0 {
		return "..."
	}
	var b strings.Builder
	for s := range set {
		b.WriteString(s)
		b.WriteString(", ")
	}
	s := b.String()
	return s[:len(s)-2]
}

// Config contains the parsed configuration.
type Config struct {
	Allow  MapSet
	Tools  StringSet
	DB     StringSet
	Ignore StringSet
}

// ParseConfig parses command line arguments and configuration file
func ParseConfig(args []string) Config {
	const (
		allowUsage     = "allowed package dependency (e.g. 'pkg/a/uses pkg/x/util')"
		toolUsage      = "tool package (leave package) (e.g. 'pkg/x/**')"
		dbUsage        = "common domain/database package (can only depend on tools)"
		ignoreUsage    = "directory to ignore"
		shorthandUsage = " (shorthand)"
	)
	c := Config{
		Allow:  make(map[string]map[string]struct{}),
		Tools:  make(map[string]struct{}),
		DB:     make(map[string]struct{}),
		Ignore: make(map[string]struct{}),
	}
	fs := flag.NewFlagSet("spaghetti-cutter", flag.ExitOnError)
	fs.Var(&c.Allow, "a", allowUsage+shorthandUsage)
	fs.Var(&c.Allow, "allow", allowUsage)
	fs.Var(&c.Tools, "t", toolUsage+shorthandUsage)
	fs.Var(&c.Tools, "tool", toolUsage)
	fs.Var(&c.DB, "d", dbUsage+shorthandUsage)
	fs.Var(&c.DB, "db", dbUsage)
	fs.Var(&c.Ignore, "i", ignoreUsage+shorthandUsage)
	fs.Var(&c.Ignore, "ignore", ignoreUsage)

	ff.Parse(fs, os.Args[1:],
		ff.WithEnvVarPrefix("SPAGHETTI_CUTTER"),
		ff.WithConfigFile(".spaghetti-cutter"),
		ff.WithConfigFileParser(ff.JSONParser),
	)
	return c
}
