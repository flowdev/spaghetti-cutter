package parse_test

import (
	"path/filepath"
	"testing"

	"github.com/flowdev/spaghetti-cutter/parse"
	"golang.org/x/tools/go/packages"
)

func TestDirTree(t *testing.T) {
	specs := []struct {
		name          string
		givenRoot     string
		expectedPkgs  []*packages.Package
		expectedError bool
	}{
		{
			name: "happy-path",
			expectedPkgs: []*packages.Package{
				{
					Name:    "main",
					PkgPath: "github.com/flowdev/...",
				},
			},
			expectedError: false,
		}, {
			name:          "error-path",
			expectedPkgs:  nil,
			expectedError: true,
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			actualPkgs, err := parse.DirTree(filepath.Join("testdata", spec.name))
			if spec.expectedError {
				if err != nil {
					t.Logf("received expected error: %v", err)
				} else {
					t.Error("expected to receive error but didn't get one")
				}
			}
			if len(actualPkgs) != len(spec.expectedPkgs) {
				t.Errorf("expected parsed packages %v (len=%d), actual %v (len=%d)",
					spec.expectedPkgs, len(spec.expectedPkgs), actualPkgs, len(actualPkgs))
			}
		})
	}
}
func packagesAsString(pkgs []*packages.Package) string {
	return ""
}

func TestRootPkg(t *testing.T) {
	specs := []struct {
		name          string
		givenPkgPaths []string
		expectedRoot  string
	}{
		{
			name:          "empty",
			givenPkgPaths: []string{"", ""},
			expectedRoot:  "",
		}, {
			name:          "nothing-in-common",
			givenPkgPaths: []string{"a", "ba"},
			expectedRoot:  "",
		}, {
			name:          "test-packages",
			givenPkgPaths: []string{"pkg/x/a", "pkg/x/a_test", "pkg/x/a.test"},
			expectedRoot:  "pkg/x/a",
		}, {
			name:          "x-packages",
			givenPkgPaths: []string{"pkg/x/a", "pkg/x/b", "pkg/x/c"},
			expectedRoot:  "pkg/x/",
		}, {
			name:          "all-on-github",
			givenPkgPaths: []string{"github.com/org1/proj1", "github.com/org1/proj2", "github.com/borg2/proj3"},
			expectedRoot:  "github.com/",
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			actualRoot := parse.RootPkg(pkgsForPaths(spec.givenPkgPaths))
			if actualRoot != spec.expectedRoot {
				t.Errorf("expected common root %q, actual %q",
					spec.expectedRoot, actualRoot)
			}
		})
	}
}

func pkgsForPaths(paths []string) []*packages.Package {
	pkgs := make([]*packages.Package, len(paths))
	for i, path := range paths {
		pkgs[i] = &packages.Package{PkgPath: path}
	}
	return pkgs
}
