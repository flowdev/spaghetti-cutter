package config_test

import (
	"fmt"
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

func TestPatternMapHasKeyValue(t *testing.T) {
	specs := []struct {
		name             string
		givenJSON        string
		givenKey         string
		givenStrictKey   string
		givenValue       string
		givenStrictValue string
		expectedHasKey   bool
		expectedHasValue bool
	}{
		{
			name:             "simple-pair-full-match",
			givenJSON:        `"a": ["b"]`,
			givenKey:         "a",
			givenStrictKey:   "",
			givenValue:       "b",
			givenStrictValue: "",
			expectedHasKey:   true,
			expectedHasValue: true,
		}, {
			name:             "simple-pair-half-match",
			givenJSON:        `"a": ["b"]`,
			givenKey:         "a",
			givenStrictKey:   "",
			givenValue:       "c",
			givenStrictValue: "",
			expectedHasKey:   true,
			expectedHasValue: false,
		}, {
			name:             "simple-pair-no-match",
			givenJSON:        `"a": ["b"]`,
			givenKey:         "c",
			givenStrictKey:   "",
			givenValue:       "d",
			givenStrictValue: "",
			expectedHasKey:   false,
			expectedHasValue: false,
		}, {
			name:             "multiple-pairs-full-match",
			givenJSON:        `"a": ["b", "c", "do", "foo"]`,
			givenKey:         "",
			givenStrictKey:   "a",
			givenValue:       "do",
			givenStrictValue: "",
			expectedHasKey:   true,
			expectedHasValue: true,
		}, {
			name:             "one-pair-many-stars-full-match",
			givenJSON:        `"a/*/b/**": ["c/*/d/**"]`,
			givenKey:         "a/foo/b/bar/doo",
			givenStrictKey:   "",
			givenValue:       "c/fox/d/baz/dox",
			givenStrictValue: "",
			expectedHasKey:   true,
			expectedHasValue: true,
		}, {
			name:             "one-pair-many-stars-half-match",
			givenJSON:        `"a/*/b/**": ["c/*/d/**"]`,
			givenKey:         "a/foo/b/bar/doo",
			givenStrictKey:   "",
			givenValue:       "c/fox/d",
			givenStrictValue: "",
			expectedHasKey:   true,
			expectedHasValue: false,
		}, {
			name:             "one-pair-many-stars-no-match",
			givenJSON:        `"a/*/b/**": ["c"]`,
			givenKey:         "a/ahoi/b",
			givenStrictKey:   "",
			givenValue:       "c",
			givenStrictValue: "",
			expectedHasKey:   false,
			expectedHasValue: false,
		}, {
			name:             "shuffled-dollars-full-match",
			givenJSON:        `"a/$*/b/$**/c": ["d/$2/e/$1/f"]`,
			givenKey:         "a/foo/b/bar/car/c",
			givenStrictKey:   "ahoi",
			givenValue:       "d/bar/car/e/foo/f",
			givenStrictValue: "car",
			expectedHasKey:   true,
			expectedHasValue: true,
		}, {
			name:             "use-not-all-dollars-full-match",
			givenJSON:        `"a/$*/b/$**/c": ["d/$2/e"]`,
			givenKey:         "a/foo/b/bar/car/c",
			givenStrictKey:   "",
			givenValue:       "d/bar/car/e",
			givenStrictValue: "",
			expectedHasKey:   true,
			expectedHasValue: true,
		}, {
			name:             "use-no-dollars-full-match",
			givenJSON:        `"a/$*/b/$**/c": ["d/e/f"]`,
			givenKey:         "ahoi",
			givenStrictKey:   "a/foo/b/bar/car/c",
			givenValue:       "car",
			givenStrictValue: "d/e/f",
			expectedHasKey:   true,
			expectedHasValue: true,
		}, {
			name:             "double-use-dollar-full-match",
			givenJSON:        `"a/$**/b": ["c/$1/d/$1/e"]`,
			givenKey:         "a/foo/bar/b",
			givenStrictKey:   "",
			givenValue:       "c/foo/bar/d/foo/bar/e",
			givenStrictValue: "",
			expectedHasKey:   true,
			expectedHasValue: true,
		}, {
			name:             "competing-keys-full-match",
			givenJSON:        `"foo/bar/**": ["b"], "$*/$*ar/b$*": ["$1/$3b/$2", "fo*/c*d/**"]`,
			givenKey:         "foo/bar/by",
			givenStrictKey:   "",
			givenValue:       "foo/yb/b",
			givenStrictValue: "",
			expectedHasKey:   true,
			expectedHasValue: true,
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			cfgBytes := []byte(`{ "allowOnlyIn": { ` + spec.givenJSON + ` } }`)
			cfg, err := config.Parse(cfgBytes, spec.name)
			if err != nil {
				t.Fatalf("got unexpected error: %v", err)
			}

			actualHasKey, actualHasValue := cfg.AllowOnlyIn.HasKeyValue(
				spec.givenKey,
				spec.givenStrictKey,
				spec.givenValue,
				spec.givenStrictValue,
			)

			if spec.expectedHasKey != actualHasKey {
				t.Errorf("expected hasKey for keys %q/%q in map %v to be %t, got %t",
					spec.givenKey, spec.givenStrictKey, cfg.AllowOnlyIn, spec.expectedHasKey, actualHasKey)
			}

			if spec.expectedHasValue != actualHasValue {
				t.Errorf("expected hasValue for keys %q/%q and values %q/%q in map %v to be %t, got %t",
					spec.givenKey, spec.givenStrictKey, spec.givenValue, spec.givenStrictValue,
					cfg.AllowOnlyIn, spec.expectedHasValue, actualHasValue)
			}
		})
	}
}
