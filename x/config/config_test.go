package config_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/flowdev/spaghetti-cutter/x/config"
)

func TestPatternList(t *testing.T) {
	specs := []struct {
		name              string
		givenPatterns     []string
		expectedMatches   []string
		expectedNoMatches []string
	}{
		{
			name:              "one-simple",
			givenPatterns:     []string{"a"},
			expectedMatches:   []string{"a"},
			expectedNoMatches: []string{"b", "aa"},
		}, {
			name:              "many-simple",
			givenPatterns:     []string{"a", "be", "do", "ra"},
			expectedMatches:   []string{"a", "be", "do", "ra"},
			expectedNoMatches: []string{"aa", "b", "bd", "od", "ra*"},
		}, {
			name:              "one-star-1",
			givenPatterns:     []string{"a/*/b"},
			expectedMatches:   []string{"a/bla/b", "a//b", "a/*/b"},
			expectedNoMatches: []string{"a/bla/blue/b", "a/bla//b", "a//bla/b"},
		}, {
			name:              "one-star-2",
			givenPatterns:     []string{"*/b"},
			expectedMatches:   []string{"bla/b", "/b", "*/b"},
			expectedNoMatches: []string{"bla/blue/b", "bla//b", "/bla/b"},
		}, {
			name:              "one-star-3",
			givenPatterns:     []string{"a/*"},
			expectedMatches:   []string{"a/bla", "a/", "a/*"},
			expectedNoMatches: []string{"a/bla/blue", "a/bla/", "a//bla", "a//"},
		}, {
			name:              "one-star-4",
			givenPatterns:     []string{"a/b*"},
			expectedMatches:   []string{"a/bla", "a/b", "a/b*"},
			expectedNoMatches: []string{"a/bla/blue", "a/bla/"},
		}, {
			name:              "multiple-single-stars-1",
			givenPatterns:     []string{"a/*/b/*/c"},
			expectedMatches:   []string{"a/foo/b/bar/c", "a//b//c", "a/*/b/*/c"},
			expectedNoMatches: []string{"a/foo//b/bar/c", "a/foo/b//bar/c", "a/bla/b///c"},
		}, {
			name:              "multiple-single-stars-2",
			givenPatterns:     []string{"a/*b/c*/d"},
			expectedMatches:   []string{"a/foob/candy/d", "a/b/c/d"},
			expectedNoMatches: []string{"a/foo/candy/d", "a/foob/c/de"},
		}, {
			name:              "double-stars",
			givenPatterns:     []string{"a/**"},
			expectedMatches:   []string{"a/foob/candy/d", "a/b/c/d/..."},
			expectedNoMatches: []string{"a/foo/candy\nd", "b/foo/b/c/d"},
		}, {
			name:              "all-stars",
			givenPatterns:     []string{"a/*/b/*/c/**"},
			expectedMatches:   []string{"a/foo/b/bar/c/d/e/f", "a/foo/b/bar/c/d/**/f", "a//b//c/"},
			expectedNoMatches: []string{},
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			pl := &config.PatternList{}
			for _, s := range spec.givenPatterns {
				pl.Set(s)
			}
			for _, s := range spec.expectedMatches {
				if !pl.MatchString(s) {
					t.Errorf("%q should match one of the patterns %q", s, spec.givenPatterns)
				}
			}
			for _, s := range spec.expectedNoMatches {
				if pl.MatchString(s) {
					t.Errorf("%q should NOT match any of the patterns %q", s, spec.givenPatterns)
				}
			}
		})
	}
}

func TestPatternMap(t *testing.T) {
}

func TestParse(t *testing.T) {
	specs := []struct {
		name                 string
		givenArgs            []string
		givenConfigFile      string
		expectedConfigString string
	}{
		{
			name:            "all-missing",
			givenArgs:       nil,
			givenConfigFile: "",
			expectedConfigString: "{" +
				"..... ... ... ... " +
				" " +
				"2048 false false" +
				"}",
		}, {
			name:            "all-empty",
			givenArgs:       []string{},
			givenConfigFile: "all-empty.json",
			expectedConfigString: "{" +
				"..... ... ... ... " +
				" " +
				"2048 false false" +
				"}",
		}, {
			name: "args-only",
			givenArgs: []string{
				"-allow", "a b",
				"-tool", "x/**",
				"--db", "pkg/db/*",
				"--god", "main",
				"-root", "dir/bla",
				"-size", "3072",
				"-ignore-vendor",
			},
			givenConfigFile: "",
			expectedConfigString: "{" +
				"`a`: `b` " +
				"`x/**` " +
				"`pkg/db/*` " +
				"`main` " +
				"dir/bla " +
				"3072 " +
				"false " +
				"true" +
				"}",
		}, {
			name:            "config-only",
			givenArgs:       []string{},
			givenConfigFile: "config-only.json",
			expectedConfigString: "{" +
				"`a`: `b` " +
				"`x/**` " +
				"`pkg/db/*` " +
				"`main` " +
				"dir/bla " +
				"3072 " +
				"false " +
				"true" +
				"}",
		}, {
			name: "args-and-config",
			givenArgs: []string{
				"--tool", "pkg/mysupertool",
				"-tool", "pkg/x/**",
				"--root", "dir/blue",
				"--size", "4096",
				"--ignore-vendor",
			},
			givenConfigFile: "args-and-config.json",
			expectedConfigString: "{" +
				"`a`: `b` ; `c`: `d` " +
				"`pkg/mysupertool`, `pkg/x/**` " +
				"`pkg/db`, `pkg/entities` " +
				"`main`, `pkg/service` " +
				"dir/blue " +
				"4096 " +
				"true " +
				"true" +
				"}",
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			cfgFile := spec.givenConfigFile
			if cfgFile != "" {
				cfgFile = filepath.Join("testdata", cfgFile)
			}
			actualConfig := config.Parse(spec.givenArgs, cfgFile)
			actualConfigString := fmt.Sprint(actualConfig)
			if actualConfigString != spec.expectedConfigString {
				t.Errorf("expected configuration %v, actual %v",
					spec.expectedConfigString, actualConfigString)
			}
		})
	}
}
