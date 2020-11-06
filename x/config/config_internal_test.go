package config

import (
	"reflect"
	"testing"
)

func TestRegexpForPattern(t *testing.T) {
	specs := []struct {
		name             string
		givenPattern     string
		givenAllowDollar int
		givenMaxDollar   int
		expectedDollars  int
		expectedIdxs     []int
		expectedRegexp   string
		expectedError    bool
	}{
		// DOLLAR STAR CASES:
		{
			name:             "empty",
			givenPattern:     ``,
			givenAllowDollar: enumDollarStar,
			expectedDollars:  0,
			expectedRegexp:   `^$`,
			expectedError:    false,
		}, {
			name:             "simple",
			givenPattern:     `abcd`,
			givenAllowDollar: enumDollarStar,
			expectedDollars:  0,
			expectedRegexp:   `^abcd$`,
			expectedError:    false,
		}, {
			name:             "double-backslash",
			givenPattern:     `ab\\\\\\cd`,
			givenAllowDollar: enumDollarStar,
			expectedDollars:  0,
			expectedRegexp:   `^ab\\\\\\cd$`,
			expectedError:    false,
		}, {
			name:             "one-star",
			givenPattern:     `ab*cd`,
			givenAllowDollar: enumDollarStar,
			expectedDollars:  0,
			expectedRegexp:   `^ab(?:[^/]*)cd$`,
			expectedError:    false,
		}, {
			name:             "two-stars",
			givenPattern:     `ab**cd`,
			givenAllowDollar: enumDollarStar,
			expectedDollars:  0,
			expectedRegexp:   `^ab(?:.*)cd$`,
			expectedError:    false,
		}, {
			name:             "unescaped-one-star",
			givenPattern:     `ab\\\\*cd`,
			givenAllowDollar: enumDollarStar,
			expectedDollars:  0,
			expectedRegexp:   `^ab\\\\(?:[^/]*)cd$`,
			expectedError:    false,
		}, {
			name:             "escaped-one-star",
			givenPattern:     `ab\\\\\*cd`,
			givenAllowDollar: enumDollarStar,
			expectedDollars:  0,
			expectedRegexp:   `^ab\\\\\*cd$`,
			expectedError:    false,
		}, {
			name:             "escaped-two-stars",
			givenPattern:     `ab\\\\\*\\\*cd`,
			givenAllowDollar: enumDollarStar,
			expectedDollars:  0,
			expectedRegexp:   `^ab\\\\\*\\\*cd$`,
			expectedError:    false,
		}, {
			name:             "escaped-dollar",
			givenPattern:     `ab\\\\\$\\\$cd`,
			givenAllowDollar: enumDollarStar,
			expectedDollars:  0,
			expectedRegexp:   `^ab\\\\\$\\\$cd$`,
			expectedError:    false,
		}, {
			name:             "dollar-one-star",
			givenPattern:     `ab$*cd`,
			givenAllowDollar: enumDollarStar,
			expectedDollars:  1,
			expectedRegexp:   `^ab([^/]*)cd$`,
			expectedError:    false,
		}, {
			name:             "dollar-two-stars",
			givenPattern:     `ab$**cd`,
			givenAllowDollar: enumDollarStar,
			expectedDollars:  1,
			expectedRegexp:   `^ab(.*)cd$`,
			expectedError:    false,
		}, {
			name:             "escaped-dollar-escaped-star",
			givenPattern:     `ab\$\*cd`,
			givenAllowDollar: enumDollarStar,
			expectedDollars:  0,
			expectedRegexp:   `^ab\$\*cd$`,
			expectedError:    false,
		}, {
			name:             "escaped-dollar-escaped-star-star",
			givenPattern:     `ab\$\**cd`,
			givenAllowDollar: enumDollarStar,
			expectedDollars:  0,
			expectedRegexp:   `^ab\$\*(?:[^/]*)cd$`,
			expectedError:    false,
		}, {
			name:             "escaped-dollar-star",
			givenPattern:     `ab\$*cd`,
			givenAllowDollar: enumDollarStar,
			expectedDollars:  0,
			expectedRegexp:   `^ab\$(?:[^/]*)cd$`,
			expectedError:    false,
		}, {
			name:             "escaped-dollar-escaped-star-escaped-star",
			givenPattern:     `ab\$\*\*cd`,
			givenAllowDollar: enumDollarStar,
			expectedDollars:  0,
			expectedRegexp:   `^ab\$\*\*cd$`,
			expectedError:    false,
		}, {
			name:             "all-good-cases",
			givenPattern:     `**se*a**fo\*o\**l/do\*\*or a\nd wi$*do$**ws or n\$\*ot a\$*ll th\$\**at g\$re\$\*\*at at al\$**l`,
			givenAllowDollar: enumDollarStar,
			expectedDollars:  2,
			expectedRegexp:   `^(?:.*)se(?:[^/]*)a(?:.*)fo\*o\*(?:[^/]*)l/do\*\*or a\nd wi([^/]*)do(.*)ws or n\$\*ot a\$(?:[^/]*)ll th\$\*(?:[^/]*)at g\$re\$\*\*at at al\$(?:.*)l$`,
			expectedError:    false,

			// NO DOLLAR CASES:
		}, {
			name:             "one-star",
			givenPattern:     `ab*cd`,
			givenAllowDollar: enumNoDollar,
			expectedDollars:  0,
			expectedRegexp:   `^ab(?:[^/]*)cd$`,
			expectedError:    false,
		}, {
			name:             "two-stars",
			givenPattern:     `ab**cd`,
			givenAllowDollar: enumNoDollar,
			expectedDollars:  0,
			expectedRegexp:   `^ab(?:.*)cd$`,
			expectedError:    false,
		}, {
			name:             "escaped-dollar-one-star",
			givenPattern:     `ab\$*cd`,
			givenAllowDollar: enumNoDollar,
			expectedDollars:  0,
			expectedRegexp:   `^ab\$(?:[^/]*)cd$`,
			expectedError:    false,
		}, {
			name:             "escaped-dollar-two-stars",
			givenPattern:     `ab\$**cd`,
			givenAllowDollar: enumNoDollar,
			expectedDollars:  0,
			expectedRegexp:   `^ab\$(?:.*)cd$`,
			expectedError:    false,

			// DOLLAR DIGIT CASES:
		}, {
			name:             "dollar-digit",
			givenPattern:     `ab$1cd`,
			givenAllowDollar: enumDollarDigit,
			givenMaxDollar:   1,
			expectedDollars:  1,
			expectedIdxs:     []int{0},
			expectedRegexp:   `^ab(.*)cd$`,
			expectedError:    false,
		}, {
			name:             "unescaped-dollar-digit",
			givenPattern:     `ab\\\\$1cd`,
			givenAllowDollar: enumDollarDigit,
			givenMaxDollar:   1,
			expectedDollars:  1,
			expectedIdxs:     []int{0},
			expectedRegexp:   `^ab\\\\(.*)cd$`,
			expectedError:    false,
		}, {
			name:             "dollar-double-digit",
			givenPattern:     `ab$11cd`,
			givenAllowDollar: enumDollarDigit,
			givenMaxDollar:   1,
			expectedDollars:  1,
			expectedIdxs:     []int{0},
			expectedRegexp:   `^ab(.*)1cd$`,
			expectedError:    false,
		}, {
			name:             "escaped-dollar-digit",
			givenPattern:     `ab\\\$3cd`,
			givenAllowDollar: enumDollarDigit,
			givenMaxDollar:   1,
			expectedDollars:  0,
			expectedRegexp:   `^ab\\\$3cd$`,
			expectedError:    false,
		}, {
			name:             "many-dollars",
			givenPattern:     `a$3b$1c$2d`,
			givenAllowDollar: enumDollarDigit,
			givenMaxDollar:   3,
			expectedDollars:  3,
			expectedIdxs:     []int{2, 0, 1},
			expectedRegexp:   `^a(.*)b(.*)c(.*)d$`,
			expectedError:    false,

			// ERROR CASES:
		}, {
			name:             "illegal-regexp",
			givenPattern:     `ab[cd`,
			givenAllowDollar: enumDollarStar,
			expectedDollars:  0,
			expectedRegexp:   ``,
			expectedError:    true,
		}, {
			name:             "unescaped-allowed-dollar",
			givenPattern:     `ab$cd`,
			givenAllowDollar: enumDollarStar,
			expectedDollars:  0,
			expectedRegexp:   ``,
			expectedError:    true,
		}, {
			name:             "unescaped-unallowed-dollar",
			givenPattern:     `ab$cd`,
			givenAllowDollar: enumNoDollar,
			expectedDollars:  0,
			expectedRegexp:   ``,
			expectedError:    true,
		}, {
			name:             "dollar-escaped-star",
			givenPattern:     `ab$\*cd`,
			givenAllowDollar: enumDollarStar,
			expectedDollars:  0,
			expectedRegexp:   ``,
			expectedError:    true,
		}, {
			name:             "escaped-backslash-dollar-escaped-star",
			givenPattern:     `ab\\$\*cd`,
			givenAllowDollar: enumDollarStar,
			expectedDollars:  0,
			expectedRegexp:   ``,
			expectedError:    true,
		}, {
			name:             "dollar-escaped-star-star",
			givenPattern:     `ab$\**cd`,
			givenAllowDollar: enumDollarStar,
			expectedDollars:  0,
			expectedRegexp:   ``,
			expectedError:    true,
		}, {
			name:             "dollar-not-allowed",
			givenPattern:     `ab$*cd`,
			givenAllowDollar: enumNoDollar,
			expectedDollars:  0,
			expectedRegexp:   ``,
			expectedError:    true,
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			t.Logf("pattern=%q, allowedDollar=%d", spec.givenPattern, spec.givenAllowDollar)
			actualRegexp, actualDollars, actualIdxs, err := regexpForPattern(spec.givenPattern, spec.givenAllowDollar, spec.givenMaxDollar)
			testOn := checkError(t, err, spec.expectedError)
			if !testOn {
				return
			}

			if actualDollars != spec.expectedDollars {
				t.Errorf("expected %d dollars, actual %d",
					spec.expectedDollars, actualDollars)
			}
			if !reflect.DeepEqual(actualIdxs, spec.expectedIdxs) {
				t.Errorf("expected dollar indices to be %v, got %v", spec.expectedIdxs, actualIdxs)
			}
			if actualRegexp.String() != spec.expectedRegexp {
				t.Errorf("expected regular expression %q, actual %q",
					spec.expectedRegexp, actualRegexp)
			}
		})
	}
}

func checkError(t *testing.T, err error, expectedError bool) bool {
	t.Helper()
	if err != nil {
		if !expectedError {
			t.Fatalf("got unexpected error: %v", err)
		}
		return false
	}
	if expectedError {
		t.Fatal("expected an error but got none")
		return false
	}
	return true
}
