package data_test

import (
	"testing"

	"github.com/flowdev/spaghetti-cutter/data"
)

func TestPatternListString(t *testing.T) {
	specs := []struct {
		name           string
		givenPatterns  []string
		expectedString string
	}{
		{
			name:           "nil",
			givenPatterns:  nil,
			expectedString: "...",
		}, {
			name:           "empty",
			givenPatterns:  []string{},
			expectedString: "...",
		}, {
			name:           "3-elements",
			givenPatterns:  []string{"a", "b**c", "d/**/e/f*"},
			expectedString: "`a`, `b**c`, `d/**/e/f*`",
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			pl, err := data.NewSimplePatternList(spec.givenPatterns, "test")
			if err != nil {
				t.Fatalf("got unexpected error: %v", err)
			}
			actualString := pl.String()
			if actualString != spec.expectedString {
				t.Errorf("expected string %q, actual %q", spec.expectedString, actualString)
			}
		})
	}
}

func TestPatternListMatchString(t *testing.T) {
	specs := []struct {
		name                string
		givenPatterns       []string
		givenDollars        []string
		expectedHalfMatches []string
		expectedFullMatches []string
		expectedNoMatches   []string
	}{
		{
			name:                "one-simple",
			givenPatterns:       []string{"a"},
			expectedFullMatches: []string{"a"},
			expectedHalfMatches: []string{"a/a", "a/"},
			expectedNoMatches:   []string{"b", "", "aa"},
		}, {
			name:                "many-simple",
			givenPatterns:       []string{"a", "be", "bed", "be/be", "do", "ra"},
			expectedFullMatches: []string{"a", "be", "bed", "be/be", "do", "ra"},
			expectedHalfMatches: []string{"a/a", "be/bed", "do/r", "ra/*"},
			expectedNoMatches:   []string{"ab", "ben", "od"},
		}, {
			name:                "one-star-1",
			givenPatterns:       []string{"a/*/b"},
			expectedFullMatches: []string{"a/cla/b", "a//b", "a/*/b"},
			expectedHalfMatches: []string{"a/cla/b/lue/b", "a//b/lue", "a/*/b/"},
			expectedNoMatches:   []string{"a/cla//b", "a/cla/blue", "a//cla/b"},
		}, {
			name:                "one-star-2",
			givenPatterns:       []string{"*/b"},
			expectedFullMatches: []string{"cla/b", "/b", "*/b"},
			expectedHalfMatches: []string{"cla/b/cd", "/b/la", "*/b/c"},
			expectedNoMatches:   []string{"bla//b", "/bla", "/cla/b"},
		}, {
			name:                "one-star-3",
			givenPatterns:       []string{"a/*"},
			expectedFullMatches: []string{"a/bla", "a/", "a/*"},
			expectedHalfMatches: []string{"a/bla/blue", "a/bla/", "a//bla", "a//"},
		}, {
			name:                "one-star-4",
			givenPatterns:       []string{"a/b*"},
			expectedFullMatches: []string{"a/bla", "a/b", "a/b*"},
			expectedHalfMatches: []string{"a/bla/blue", "a/bla/"},
		}, {
			name:                "multiple-single-stars-1",
			givenPatterns:       []string{"a/*/b/*/c"},
			expectedFullMatches: []string{"a/foo/b/bar/c", "a//b//c", "a/*/b/*/c"},
			expectedNoMatches:   []string{"a/foo//b/bar/c", "a/foo/b//bar/c", "a/bla/b///c"},
		}, {
			name:                "multiple-single-stars-2",
			givenPatterns:       []string{"a/*b/c*/d"},
			expectedFullMatches: []string{"a/foob/candy/d", "a/b/c/d"},
			expectedHalfMatches: []string{"a/foob/c/d/e"},
			expectedNoMatches:   []string{"a/foob/c/de"},
		}, {
			name:                "escaped-star",
			givenPatterns:       []string{"a/\\*/b"},
			expectedFullMatches: []string{"a/*/b"},
			expectedNoMatches:   []string{"a/bla/b", "a//b"},
		}, {
			name:                "double-stars",
			givenPatterns:       []string{"a/**"},
			expectedFullMatches: []string{"a/foob/candy/d", "a/b/c/d/..."},
			expectedNoMatches:   []string{"b/foo/b/c/d"},
		}, {
			name:                "all-stars",
			givenPatterns:       []string{"a/*/b/*/c/**"},
			expectedFullMatches: []string{"a/foo/b/bar/c/d/e/f", "a/foo/b/bar/c/d/**/f", "a//b//c/"},
			expectedNoMatches:   []string{"a/foo/b/bar/d/e/f"},
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			pl, err := data.NewSimplePatternList(spec.givenPatterns, "test")
			if err != nil {
				t.Fatalf("got unexpected error: %v", err)
			}
			for _, s := range spec.expectedFullMatches {
				if _, fullmatch := pl.MatchString(s, nil); !fullmatch {
					t.Errorf("%q should fully match one of the patterns %q", s, spec.givenPatterns)
				}
			}
			for _, s := range spec.expectedHalfMatches {
				match, fullmatch := pl.MatchString(s, nil)
				if !match {
					t.Errorf("%q should (half) match one of the patterns %q", s, spec.givenPatterns)
				}
				if fullmatch {
					t.Errorf("%q should only match HALF one of the patterns %q", s, spec.givenPatterns)
				}
			}
			for _, s := range spec.expectedNoMatches {
				if match, _ := pl.MatchString(s, nil); match {
					t.Errorf("%q should NOT match any of the patterns %q", s, spec.givenPatterns)
				}
			}
		})
	}
}
