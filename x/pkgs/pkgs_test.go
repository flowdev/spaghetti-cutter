package pkgs_test

import (
	"testing"

	"github.com/flowdev/spaghetti-cutter/x/pkgs"
	"golang.org/x/tools/go/packages"
)

func TestRelativePackageName(t *testing.T) {
	givenRootPkg := "github.com/flowdev/spaghetti-cutter"
	specs := []struct {
		name                 string
		givenPkgPath         string
		givenPkgName         string
		expectedRelPkg       string
		expectedStrictRelPkg string
	}{
		{
			name:                 "foreign-package",
			givenPkgPath:         "golang.org/x/tools/go/packages",
			givenPkgName:         "packages",
			expectedRelPkg:       "golang.org/x/tools/go/packages",
			expectedStrictRelPkg: "",
		}, {
			name:                 "main-package",
			givenPkgPath:         "github.com/flowdev/spaghetti-cutter",
			givenPkgName:         "main",
			expectedRelPkg:       "main",
			expectedStrictRelPkg: "/",
		}, {
			name:                 "root-package",
			givenPkgPath:         "github.com/flowdev/spaghetti-cutter",
			givenPkgName:         "spaghetti-cutter",
			expectedRelPkg:       "/",
			expectedStrictRelPkg: "",
		}, {
			name:                 "standard-package",
			givenPkgPath:         "github.com/flowdev/spaghetti-cutter/x/config",
			givenPkgName:         "config",
			expectedRelPkg:       "x/config",
			expectedStrictRelPkg: "",
		}, {
			name:                 "test-package",
			givenPkgPath:         "github.com/flowdev/spaghetti-cutter/x/config_test",
			givenPkgName:         "config_test",
			expectedRelPkg:       "x/config_test",
			expectedStrictRelPkg: "",
		}, {
			name:                 "test-package-2",
			givenPkgPath:         "github.com/flowdev/spaghetti-cutter/x/config.test",
			givenPkgName:         "config.test",
			expectedRelPkg:       "x/config",
			expectedStrictRelPkg: "",
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			pkg := &packages.Package{
				PkgPath: spec.givenPkgPath,
				Name:    spec.givenPkgName,
			}
			actualRelPkg, actualStrictRelPkg := pkgs.RelativePackageName(pkg, givenRootPkg)
			if actualRelPkg != spec.expectedRelPkg {
				t.Errorf("expected relative package %q, actual %q",
					spec.expectedRelPkg, actualRelPkg)
			}
			if actualStrictRelPkg != spec.expectedStrictRelPkg {
				t.Errorf("expected strict relative package %q, actual %q",
					spec.expectedStrictRelPkg, actualStrictRelPkg)
			}
		})
	}
}

func TestUniquePackageName(t *testing.T) {
	specs := []struct {
		name              string
		givenRelPkg       string
		givenStrictRelPkg string
		expectedUniqPkg   string
	}{
		{
			name:              "empty strict",
			givenRelPkg:       "x/config",
			givenStrictRelPkg: "",
			expectedUniqPkg:   "x/config",
		}, {
			name:              "with strict",
			givenRelPkg:       "main",
			givenStrictRelPkg: "/",
			expectedUniqPkg:   "/",
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			actualUniqPkg := pkgs.UniquePackageName(spec.givenRelPkg, spec.givenStrictRelPkg)
			if actualUniqPkg != spec.expectedUniqPkg {
				t.Errorf("expected %q, actual %q", spec.expectedUniqPkg, actualUniqPkg)
			}
		})
	}
}

func TestUniquePackages(t *testing.T) {
	specs := []struct {
		name          string
		givenPkgNames []string
		givenPkgPaths []string
		expectedLen   int
	}{
		{
			name:          "standard package",
			givenPkgNames: []string{"config", "config_test", "main", "config"},
			givenPkgPaths: []string{"x/config", "x/config_test", "cmd/my_service/main.test", "x/config"},
			expectedLen:   3,
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			packs := make([]*packages.Package, len(spec.givenPkgNames))
			for i, name := range spec.givenPkgNames {
				packs[i] = &packages.Package{
					ID:      name,
					Name:    name,
					PkgPath: spec.givenPkgPaths[i],
				}
			}
			uniqPkgs := pkgs.UniquePackages(packs)
			if len(uniqPkgs) != spec.expectedLen {
				t.Errorf("expected length %d, actual %d %v", spec.expectedLen, len(uniqPkgs), uniqPkgs)
			}
		})
	}
}

func TestIsTestPackage(t *testing.T) {
	specs := []struct {
		name           string
		givenPkgPath   string
		expectedIsTest bool
	}{
		{
			name:           "standard package",
			givenPkgPath:   "x/config",
			expectedIsTest: false,
		}, {
			name:           "normal test package",
			givenPkgPath:   "x/config_test",
			expectedIsTest: true,
		}, {
			name:           "main test package",
			givenPkgPath:   "cmd/my_service/main.test",
			expectedIsTest: true,
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			pkg := &packages.Package{
				PkgPath: spec.givenPkgPath,
			}
			actualIsTest := pkgs.IsTestPackage(pkg)
			if actualIsTest != spec.expectedIsTest {
				t.Errorf("expected %t, actual %t", spec.expectedIsTest, actualIsTest)
			}
		})
	}
}
