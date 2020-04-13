package config

import (
	"flag"
	"os"
	"strings"

	"github.com/peterbourgon/ff/v3"
)

// StringSlice is a flag.Value that collects each Set string
// into a slice, allowing for repeated flags.
type StringSlice []string

// Set implements flag.Value and appends the string to the slice.
func (ss *StringSlice) Set(s string) error {
	(*ss) = append(*ss, s)
	return nil
}

// String implements flag.Value and returns the list of
// strings, or "..." if no strings have been added.
func (ss *StringSlice) String() string {
	if len(*ss) <= 0 {
		return "..."
	}
	return strings.Join(*ss, ", ")
}

// Config contains the parsed configuration.
type Config struct {
	Allow  StringSlice
	Tools  StringSlice
	DB     StringSlice
	Ignore StringSlice
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
	c := Config{}
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
