package main

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/flowdev/spaghetti-cutter/config"
)

func TestCut(t *testing.T) {
	specs := []struct {
		name               string
		givenRoot          string
		givenConfig        string
		expectedReturnCode int
	}{
		{
			name:               "no-config-good-proj",
			givenRoot:          "good-proj",
			givenConfig:        `{}`,
			expectedReturnCode: 1,
		}, {
			name:      "strict-config-good-proj",
			givenRoot: "good-proj",
			givenConfig: `{
				"tool": ["pkg/x/*"], "db": ["pkg/db/*"],
				"allowAdditionally": {"pkg/domain4": ["pkg/domain3"]},
				"size": 16
			}`,
			expectedReturnCode: 1,
		}, {
			name:      "lenient-config-good-proj",
			givenRoot: "good-proj",
			givenConfig: `{
				"tool": ["pkg/x/*"], "db": ["pkg/db/*"],
				"allowAdditionally": {"pkg/domain4": ["pkg/domain3"]},
				"size": 1024
			}`,
			expectedReturnCode: 0,
		}, {
			name:               "no-config-bad-proj",
			givenRoot:          "bad-proj",
			givenConfig:        `{}`,
			expectedReturnCode: 6,
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			root := mustAbs(filepath.Join("testdata", spec.givenRoot))
			mustWriteFile(filepath.Join(root, config.File), []byte(spec.givenConfig))
			args := []string{"--root", root}
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

func mustWriteFile(filename string, data []byte) {
	err := ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		panic(err.Error())
	}
}
