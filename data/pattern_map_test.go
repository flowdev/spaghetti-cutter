package data_test

import (
	"testing"

	"github.com/flowdev/spaghetti-cutter/config"
)

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
				t.Errorf("expected hasKey for keys %q/%q in map %s to be %t, got %t",
					spec.givenKey, spec.givenStrictKey, cfg.AllowOnlyIn.String(), spec.expectedHasKey, actualHasKey)
			}

			if spec.expectedHasValue != actualHasValue {
				t.Errorf("expected hasValue for keys %q/%q and values %q/%q in map %s to be %t, got %t",
					spec.givenKey, spec.givenStrictKey, spec.givenValue, spec.givenStrictValue,
					cfg.AllowOnlyIn.String(), spec.expectedHasValue, actualHasValue)
			}
		})
	}
}
