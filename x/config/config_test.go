package config_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/flowdev/spaghetti-cutter/x/config"
)

func TestParse(t *testing.T) {
	specs := []struct {
		name                 string
		givenConfigBytes     []byte
		expectedConfigString string
	}{
		{
			name: "all-empty",
			givenConfigBytes: []byte(`{
				}`),
			expectedConfigString: "{" +
				"..... ..... ... ... `main` " +
				"2048 false" +
				"}",
		}, {
			name: "scalars-only",
			givenConfigBytes: []byte(`{
				  "size": 3072,
				  "noGod": true
				}`),
			expectedConfigString: "{" +
				"..... ..... ... ... ... " +
				"3072 true" +
				"}",
		}, {
			name: "list-one",
			givenConfigBytes: []byte(`{
				  "god": ["a"]
				}`),
			expectedConfigString: "{" +
				"..... ..... " +
				"... " +
				"... " +
				"`a` " +
				"2048 false" +
				"}",
		}, {
			name: "list-many",
			givenConfigBytes: []byte(`{
				  "god": ["a", "be", "do", "ra"]
				}`),
			expectedConfigString: "{" +
				"..... ..... " +
				"... " +
				"... " +
				"`a`, `be`, `do`, `ra` " +
				"2048 false" +
				"}",
		}, {
			name: "map-simple-pair",
			givenConfigBytes: []byte(`{
				  "allowOnlyIn": {
				    "a": ["b"]
				  }
				}`),
			expectedConfigString: "{" +
				"`a`: `b` " +
				"..... " +
				"... ... `main` 2048 false" +
				"}",
		}, {
			name: "map-multiple-pairs",
			givenConfigBytes: []byte(`{
				  "allowOnlyIn": {
				    "a": ["b", "c", "do", "foo"],
				    "e": ["bar", "car"]
				  }
				}`),
			expectedConfigString: "{" +
				"`a`: `b`, `c`, `do`, `foo` ; `e`: `bar`, `car` " +
				"..... " +
				"... ... `main` 2048 false" +
				"}",
		}, {
			name: "map-one-pair-many-stars",
			givenConfigBytes: []byte(`{
				  "allowOnlyIn": {
				    "a/*/b/**": ["c/*/d/**"]
				  }
				}`),
			expectedConfigString: "{" +
				"`a/*/b/**`: `c/*/d/**` " +
				"..... " +
				"... ... `main` 2048 false" +
				"}",
		}, {
			name: "map-all-complexity",
			givenConfigBytes: []byte(`{
				  "allowAdditionally": {
				    "*/*a/**": ["*/*b/**", "b*/c*d/**"]
				  }
				}`),
			expectedConfigString: "{" +
				"..... " +
				"`*/*a/**`: `*/*b/**`, `b*/c*d/**` " +
				"... ... `main` 2048 false" +
				"}",
		}, {
			name: "maps-and-lists-only",
			givenConfigBytes: []byte(`{
					"allowOnlyIn": {
					  "github.com/lib/pq": ["a"]
					},
					"allowAdditionally": {
					  "a": ["b"]
					},
					"tool": ["x/**"],
					"db": ["pkg/db/*"],
					"god": ["main"]
				}`),
			expectedConfigString: "{" +
				"`github.com/lib/pq`: `a` " +
				"`a`: `b` " +
				"`x/**` " +
				"`pkg/db/*` " +
				"`main` " +
				"2048 " +
				"false" +
				"}",
		}, {
			name: "a-bit-of-everything",
			givenConfigBytes: []byte(`{
					"allowOnlyIn": {
					  "github.com/lib/pq": ["a", "b"]
					},
					"allowAdditionally": {
					  "a": ["b"],
					  "c": ["d"]
					},
					"tool": ["pkg/mysupertool", "pkg/x/**"],
					"db": ["pkg/db", "pkg/entities"],
					"god": ["main", "pkg/service"],
					"size": 3072,
					"noGod": true
				}`),
			expectedConfigString: "{" +
				"`github.com/lib/pq`: `a`, `b` " +
				"`a`: `b` ; `c`: `d` " +
				"`pkg/mysupertool`, `pkg/x/**` " +
				"`pkg/db`, `pkg/entities` " +
				"`main`, `pkg/service` " +
				"3072 " +
				"true" +
				"}",
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			actualConfig, err := config.Parse(spec.givenConfigBytes, spec.name)
			if err != nil {
				t.Fatalf("got unexpected error: %v", err)
			}
			actualConfigString := fmt.Sprint(actualConfig)
			if actualConfigString != spec.expectedConfigString {
				t.Errorf("expected configuration %v, actual %v",
					spec.expectedConfigString, actualConfigString)
			}
		})
	}
}

func TestPatternListMatchString(t *testing.T) {
	specs := []struct {
		name              string
		givenPatterns     []string
		givenDollars      []string
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
			name:              "escaped-star",
			givenPatterns:     []string{"a/\\\\*/b"},
			expectedMatches:   []string{"a/*/b"},
			expectedNoMatches: []string{"a/bla/b", "a//b"},
		}, {
			name:              "double-stars",
			givenPatterns:     []string{"a/**"},
			expectedMatches:   []string{"a/foob/candy/d", "a/b/c/d/..."},
			expectedNoMatches: []string{"a/foo/candy\nd", "b/foo/b/c/d"},
		}, {
			name:              "all-stars",
			givenPatterns:     []string{"a/*/b/*/c/**"},
			expectedMatches:   []string{"a/foo/b/bar/c/d/e/f", "a/foo/b/bar/c/d/**/f", "a//b//c/"},
			expectedNoMatches: []string{"a/foo/b/bar/d/e/f"},
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			cfgBytes := []byte(`{ "db": ["` + strings.Join(spec.givenPatterns, `", "`) + `"] }`)
			cfg, err := config.Parse(cfgBytes, spec.name)
			if err != nil {
				t.Fatalf("got unexpected error: %v", err)
			}
			pl := cfg.DB
			for _, s := range spec.expectedMatches {
				if !pl.MatchString(s, nil) {
					t.Errorf("%q should match one of the patterns %q", s, spec.givenPatterns)
				}
			}
			for _, s := range spec.expectedNoMatches {
				if pl.MatchString(s, nil) {
					t.Errorf("%q should NOT match any of the patterns %q", s, spec.givenPatterns)
				}
			}
		})
	}
}

func TestPatternMap(t *testing.T) {
	specs := []struct {
		name                string
		givenJSON           string
		givenLeftPattern    string
		expectedLeftMatch   bool
		expectedRightString string
	}{
		{
			name:                "simple-pair",
			givenJSON:           `"a": ["b"]`,
			givenLeftPattern:    "a",
			expectedLeftMatch:   true,
			expectedRightString: "`b`",
		}, {
			name:                "multiple-pairs",
			givenJSON:           `"a": ["b", "c", "do", "foo"]`,
			givenLeftPattern:    "a",
			expectedLeftMatch:   true,
			expectedRightString: "`b`, `c`, `do`, `foo`",
		}, {
			name:                "one-pair-many-stars",
			givenJSON:           `"a/*/b/**": ["c/*/d/**"]`,
			givenLeftPattern:    "a/foo/b/bar/doo",
			expectedLeftMatch:   true,
			expectedRightString: "`c/*/d/**`",
		}, {
			name:                "no-match",
			givenJSON:           `"a/*/b/**": ["c"]`,
			givenLeftPattern:    "a/ahoi/b",
			expectedLeftMatch:   false,
			expectedRightString: "",
		}, {
			name:                "dollars",
			givenJSON:           `"a/$*/b/$**/c": ["d/$2/e/$1/f"]`,
			givenLeftPattern:    "a/foo/b/bar/car/c",
			expectedLeftMatch:   true,
			expectedRightString: "`d/$2/e/$1/f`",
		}, {
			name:                "all-complexity",
			givenJSON:           `"foo/bar/**": ["b"], "*/*a/**": ["*/*b/**", "b*/c*d/**"]`,
			givenLeftPattern:    "foo/bara/doo/ey",
			expectedLeftMatch:   true,
			expectedRightString: "`*/*b/**`, `b*/c*d/**`",
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			cfgBytes := []byte(`{ "allowOnlyIn": { ` + spec.givenJSON + ` } }`)
			cfg, err := config.Parse(cfgBytes, spec.name)
			if err != nil {
				t.Fatalf("got unexpected error: %v", err)
			}
			pm := cfg.AllowOnlyIn

			pl, _ := pm.MatchingList(spec.givenLeftPattern)

			if spec.expectedLeftMatch && pl == nil {
				t.Fatalf("expected left match for pattern %q in map %v", spec.givenLeftPattern, pm)
			} else if !spec.expectedLeftMatch && pl != nil {
				t.Fatalf("expected NO left match for pattern %q in map %v but got: %v", spec.givenLeftPattern, pm, pl)
			}
			if !spec.expectedLeftMatch {
				return
			}

			if spec.expectedRightString != pl.String() {
				t.Errorf("expected right string representation %q but got: %q", spec.expectedRightString, pl.String())
			}
		})
	}
}

func TestDollars(t *testing.T) {
	specs := []struct {
		name                string
		givenJSON           string
		givenLeftPattern    string
		givenRightPattern   string
		expectedKeyDollars  []string
		expectedRightString string
	}{
		{
			name:                "simple-star",
			givenJSON:           `"a/$*/b": ["c/$1/d"]`,
			givenLeftPattern:    "a/foo/b",
			givenRightPattern:   "c/foo/d",
			expectedKeyDollars:  []string{"foo"},
			expectedRightString: "`c/$1/d`",
		}, {
			name:                "double-star",
			givenJSON:           `"a/$**/b": ["c/$1/d"]`,
			givenLeftPattern:    "a/foo/bar/b",
			givenRightPattern:   "c/foo/bar/d",
			expectedKeyDollars:  []string{"foo/bar"},
			expectedRightString: "`c/$1/d`",
		}, {
			name:                "many-dollars",
			givenJSON:           `"a/$*/b/$**/c": ["d/$1/e/$2/f"]`,
			givenLeftPattern:    "a/foo/b/bar/car/c",
			givenRightPattern:   "d/foo/e/bar/car/f",
			expectedKeyDollars:  []string{"foo", "bar/car"},
			expectedRightString: "`d/$1/e/$2/f`",
		}, {
			name:                "shuffled-dollars",
			givenJSON:           `"a/$*/b/$**/c": ["d/$2/e/$1/f"]`,
			givenLeftPattern:    "a/foo/b/bar/car/c",
			givenRightPattern:   "d/bar/car/e/foo/f",
			expectedKeyDollars:  []string{"foo", "bar/car"},
			expectedRightString: "`d/$2/e/$1/f`",
		}, {
			name:                "use-not-all-dollars",
			givenJSON:           `"a/$*/b/$**/c": ["d/$2/e"]`,
			givenLeftPattern:    "a/foo/b/bar/car/c",
			givenRightPattern:   "d/bar/car/e",
			expectedKeyDollars:  []string{"foo", "bar/car"},
			expectedRightString: "`d/$2/e`",
		}, {
			name:                "use-no-dollars",
			givenJSON:           `"a/$*/b/$**/c": ["d/e/f"]`,
			givenLeftPattern:    "a/foo/b/bar/car/c",
			givenRightPattern:   "d/e/f",
			expectedKeyDollars:  []string{"foo", "bar/car"},
			expectedRightString: "`d/e/f`",
		}, {
			name:                "double-use-dollar",
			givenJSON:           `"a/$**/b": ["c/$1/d/$1/e"]`,
			givenLeftPattern:    "a/foo/bar/b",
			givenRightPattern:   "c/foo/bar/d/foo/bar/e",
			expectedKeyDollars:  []string{"foo/bar"},
			expectedRightString: "`c/$1/d/$1/e`",
		}, {
			name:                "all-at-once",
			givenJSON:           `"a/$*/b/$**/c": ["d/$2/e/$1/f/$1/g/$2/h"]`,
			givenLeftPattern:    "a/foo/b/bar/car/c",
			givenRightPattern:   "d/bar/car/e/foo/f/foo/g/bar/car/h",
			expectedKeyDollars:  []string{"foo", "bar/car"},
			expectedRightString: "`d/$2/e/$1/f/$1/g/$2/h`",
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			cfgBytes := []byte(`{ "allowOnlyIn": { ` + spec.givenJSON + ` } }`)
			cfg, err := config.Parse(cfgBytes, spec.name)
			if err != nil {
				t.Fatalf("got unexpected error: %v", err)
			}
			pm := cfg.AllowOnlyIn

			pl, actualKeyDollars := pm.MatchingList(spec.givenLeftPattern)

			if pl == nil {
				t.Fatalf("expected left match for pattern %q in map %v", spec.givenLeftPattern, pm)
			}

			if !reflect.DeepEqual(actualKeyDollars, spec.expectedKeyDollars) {
				t.Errorf("expected dollar matches in key to be %q, got %q", spec.expectedKeyDollars, actualKeyDollars)
			}
			if spec.expectedRightString != pl.String() {
				t.Errorf("expected right string representation %q but got: %q", spec.expectedRightString, pl.String())
			}
			if !pl.MatchString(spec.givenRightPattern, spec.expectedKeyDollars) {
				t.Errorf("right pattern %q didn't match with dollars %q",
					spec.givenRightPattern, spec.expectedKeyDollars)
			}
		})
	}
}
