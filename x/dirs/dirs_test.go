package dirs_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/flowdev/spaghetti-cutter/x/dirs"
)

const testFile = ".test-file"

func TestFindRoot(t *testing.T) {
	testDataDir := mustAbs(filepath.Join("testdata", "find-root"))
	specs := []struct {
		name              string
		givenCWD          string
		givenStartDir     string
		givenIgnoreVendor bool
		expectedRoot      string
	}{
		{
			name:              "given-root",
			givenCWD:          "",
			givenStartDir:     filepath.Join("in", "some", "subdir"),
			givenIgnoreVendor: false,
			expectedRoot:      filepath.Join(testDataDir, "given-root", "in"),
		}, {
			name:              "config-file",
			givenCWD:          filepath.Join("in", "some", "subdir"),
			givenStartDir:     "",
			givenIgnoreVendor: false,
			expectedRoot:      filepath.Join(testDataDir, "config-file"),
		},
	}

	initDir := mustAbs(".")
	t.Cleanup(func() {
		mustChdir(initDir)
	})
	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			mustChdir(filepath.Join(testDataDir, spec.name, spec.givenCWD))

			actualRoot, err := dirs.FindRoot(spec.givenStartDir, testFile)
			if err != nil {
				t.Fatalf("expected no error but got: %v", err)
			}
			if actualRoot != spec.expectedRoot {
				t.Errorf("expected project root %q, actual %q",
					spec.expectedRoot, actualRoot)
			}
		})
	}
}

func mustChdir(path string) {
	err := os.Chdir(path)
	if err != nil {
		panic(err.Error())
	}
}

func mustAbs(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err.Error())
	}
	return absPath
}
