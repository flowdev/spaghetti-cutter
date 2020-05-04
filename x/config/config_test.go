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
			name:              "set-one-simple",
			givenPatterns:     []string{"a"},
			expectedMatches:   []string{"a"},
			expectedNoMatches: []string{"b", "aa"},
		}, {},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
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
