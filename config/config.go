package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/peterbourgon/ff/v3"
)

// File is the name of the configuration file
const File = ".spaghetti-cutter.json"

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
	Tool   StringSet
	DB     StringSet
	God    StringSet
	Ignore StringSet
	Root   string
	Size   uint
}

// Parse parses command line arguments and configuration file
func Parse(args []string) Config {
	const (
		usageAllow  = "allowed package dependency (e.g. 'pkg/a/uses pkg/x/util')"
		usageTool   = "tool package (leave package) (e.g. 'pkg/x/**')"
		usageDB     = "common domain/database package (can only depend on tools)"
		usageGod    = "god package that can see everything (default: 'main')"
		usageIgnore = "directory to ignore"
		usageRoot   = "root directory"
		usageSize   = "maximum size of a package in \"lines\""
		defaultSize = 4096
	)
	cfg := Config{
		Allow:  make(map[string]map[string]struct{}),
		Tool:   make(map[string]struct{}),
		DB:     make(map[string]struct{}),
		God:    make(map[string]struct{}),
		Ignore: make(map[string]struct{}),
	}
	fs := flag.NewFlagSet("spaghetti-cutter", flag.ExitOnError)
	fs.Var(&cfg.Allow, "allow", usageAllow)
	fs.Var(&cfg.Tool, "tool", usageTool)
	fs.Var(&cfg.DB, "db", usageDB)
	fs.Var(&cfg.God, "god", usageGod)
	fs.Var(&cfg.Ignore, "ignore", usageIgnore)
	fs.StringVar(&cfg.Root, "root", "", usageRoot)
	fs.UintVar(&cfg.Size, "size", defaultSize, usageSize)

	err := ff.Parse(fs, os.Args[1:],
		ff.WithEnvVarPrefix("SPAGHETTI_CUTTER"),
		ff.WithConfigFile(File),
		ff.WithConfigFileParser(ff.JSONParser),
	)
	if err != nil {
		log.Fatalf("FATAL: Unable to parse command line arguments or configuration file: %v", err)
	}

	//fmt.Println("Parsed config:", cfg)
	return cfg
}
