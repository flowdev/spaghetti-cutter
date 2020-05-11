package main

import (
	"path/filepath"
	"testing"
)

func TestCut(t *testing.T) {
	specs := []struct {
		name               string
		givenRoot          string
		givenArgs          []string
		expectedReturnCode int
	}{
		{
			name:               "no-args-good-proj",
			givenRoot:          "good-proj",
			givenArgs:          nil,
			expectedReturnCode: 1,
		}, {
			name:      "strict-args-good-proj",
			givenRoot: "good-proj",
			givenArgs: []string{
				"--tool", "pkg/x/*", "--db", "pkg/db/*",
				"--allow", "pkg/domain4 pkg/domain3",
				"--size", "16",
			},
			expectedReturnCode: 1,
		}, {
			name:      "lenient-args-good-proj",
			givenRoot: "good-proj",
			givenArgs: []string{
				"--tool", "pkg/x/*", "--db", "pkg/db/*",
				"--allow", "pkg/domain4 pkg/domain3",
				"--size", "1024",
			},
			expectedReturnCode: 0,
		}, {
			name:               "no-args-bad-proj",
			givenRoot:          "bad-proj",
			givenArgs:          nil,
			expectedReturnCode: 3,
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			args := append(spec.givenArgs, "--root", mustAbs(filepath.Join("testdata", spec.givenRoot)))
			actualReturnCode := cut(args)

			if actualReturnCode != spec.expectedReturnCode {
				t.Errorf("Expected return code %d but got: %d", spec.expectedReturnCode, actualReturnCode)
			}
		})
	}
}

func mustAbs(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err.Error())
	}
	return absPath
}
