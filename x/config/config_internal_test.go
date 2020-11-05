package config

import (
	"testing"
)

func TestRegexpForPattern(t *testing.T) {
	specs := []struct {
		name           string
		givenPattern   string
		expectedRegexp string
		expectedError  bool
	}{
		{
			name:           "empty",
			givenPattern:   "",
			expectedRegexp: "^$",
			expectedError:  false,
		}, {
			name:           "simple",
			givenPattern:   "abcd",
			expectedRegexp: "^abcd$",
			expectedError:  false,
		}, {
			name:           "one-star",
			givenPattern:   "ab*cd",
			expectedRegexp: "^ab(?:[^/]*)cd$",
			expectedError:  false,
		}, {
			name:           "two-stars",
			givenPattern:   "ab**cd",
			expectedRegexp: "^ab(?:.*)cd$",
			expectedError:  false,
		}, {
			name:           "escaped-one-star",
			givenPattern:   "ab\\*cd",
			expectedRegexp: "^ab\\*cd$",
			expectedError:  false,
		}, {
			name:           "escaped-two-stars",
			givenPattern:   "ab\\*\\*cd",
			expectedRegexp: "^ab\\*\\*cd$",
			expectedError:  false,
		}, {
			name:           "escaped-dollar",
			givenPattern:   "ab\\$\\$cd",
			expectedRegexp: "^ab\\$\\$cd$",
			expectedError:  false,
		}, {
			name:           "dollar-one-star",
			givenPattern:   "ab$*cd",
			expectedRegexp: "^ab([^/]*)cd$",
			expectedError:  false,
		}, {
			name:           "dollar-two-stars",
			givenPattern:   "ab$**cd",
			expectedRegexp: "^ab(.*)cd$",
			expectedError:  false,
		}, {
			name:           "escaped-dollar-escaped-star",
			givenPattern:   "ab\\$\\*cd",
			expectedRegexp: "^ab\\$\\*cd$",
			expectedError:  false,
		}, {
			name:           "escaped-dollar-escaped-star-star",
			givenPattern:   "ab\\$\\**cd",
			expectedRegexp: "^ab\\$\\*(?:[^/]*)cd$",
			expectedError:  false,
		}, {
			name:           "escaped-dollar-star",
			givenPattern:   "ab\\$*cd",
			expectedRegexp: "^ab\\$(?:[^/]*)cd$",
			expectedError:  false,
		}, {
			name:           "escaped-dollar-escaped-star-escaped-star",
			givenPattern:   "ab\\$\\*\\*cd",
			expectedRegexp: "^ab\\$\\*\\*cd$",
			expectedError:  false,
		}, {
			name:           "illegal-regexp",
			givenPattern:   "ab[cd",
			expectedRegexp: "",
			expectedError:  true,
		}, {
			name:           "unescaped-dollar",
			givenPattern:   "ab$cd",
			expectedRegexp: "",
			expectedError:  true,
		}, {
			name:           "dollar-escaped-star",
			givenPattern:   "ab$\\*cd",
			expectedRegexp: "",
			expectedError:  true,
		}, {
			name:           "dollar-escaped-star-star",
			givenPattern:   "ab$\\**cd",
			expectedRegexp: "",
			expectedError:  true,
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			actualRegexp, err := regexpForPattern(spec.givenPattern)
			testOn := checkError(t, err, spec.expectedError)
			if !testOn {
				return
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
