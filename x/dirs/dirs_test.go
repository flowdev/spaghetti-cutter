package dirs_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/flowdev/spaghetti-cutter/x/dirs"
)

const testFile = ".test-file"

func TestFindConfig(t *testing.T) {
	testDataDir := mustAbs(filepath.Join("testdata", "find-config"))
	specs := []struct {
		name                string
		givenStartDir       string
		expectedCfgFilePath string
	}{
		{
			name:                "in-current-dir",
			givenStartDir:       "",
			expectedCfgFilePath: filepath.Join(testDataDir, "in-current-dir", testFile),
		}, {
			name:                "in-far-away-parent-dir",
			givenStartDir:       filepath.Join("deep", "down", "in", "the", "rabbit", "hole"),
			expectedCfgFilePath: filepath.Join(testDataDir, "in-far-away-parent-dir", testFile),
		}, {
			name:                "does-not-exist",
			givenStartDir:       filepath.Join("in", "some", "subdir"),
			expectedCfgFilePath: "",
		},
	}

	initDir := mustAbs(".")
	t.Cleanup(func() {
		mustChdir(initDir)
	})
	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			mustChdir(filepath.Join(testDataDir, spec.name, spec.givenStartDir))

			actualCfgFilePath := dirs.FindConfig(testFile)
			if actualCfgFilePath != spec.expectedCfgFilePath {
				t.Errorf("expected configuration file path %q, actual %q",
					spec.expectedCfgFilePath, actualCfgFilePath)
			}
		})
	}
}

func TestFindRoot(t *testing.T) {
	testDataDir := mustAbs(filepath.Join("testdata", "find-root"))
	givenStartDir := filepath.Join("in", "some", "subdir")
	specs := []struct {
		name         string
		givenRoot    string
		expectedRoot string
		expectedErr  bool
	}{
		{
			name:         "go-mod",
			givenRoot:    "",
			expectedRoot: filepath.Join(testDataDir, "go-mod"),
		}, {
			name:         "given-root",
			givenRoot:    "/my/given/root/dir",
			expectedRoot: "/my/given/root/dir",
		}, {
			name:         "config-file",
			givenRoot:    "",
			expectedRoot: filepath.Join(testDataDir, "config-file"),
		}, {
			name:         "vendor-dir",
			givenRoot:    "",
			expectedRoot: filepath.Join(testDataDir, "vendor-dir"),
		},
	}

	initDir := mustAbs(".")
	t.Cleanup(func() {
		mustChdir(initDir)
	})
	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			mustChdir(filepath.Join(testDataDir, spec.name, givenStartDir))

			actualRoot, err := dirs.FindRoot(spec.givenRoot, testFile)
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
