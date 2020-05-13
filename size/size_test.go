package size_test

import (
	"path/filepath"
	"testing"

	"github.com/flowdev/spaghetti-cutter/parse"
	"github.com/flowdev/spaghetti-cutter/size"
)

func TestCheck(t *testing.T) {
	specs := []struct {
		name           string
		givenMaxSize   uint
		expectedErrors int
	}{
		{
			name:           "normal-size-no-errors",
			givenMaxSize:   1024,
			expectedErrors: 0,
		}, {
			name:           "medium-size-one-error",
			givenMaxSize:   64,
			expectedErrors: 1,
		}, {
			name:           "small-size-two-errors",
			givenMaxSize:   32,
			expectedErrors: 2,
		}, {
			name:           "tiny-size-many-errors",
			givenMaxSize:   8,
			expectedErrors: 7,
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			pkgs, err := parse.DirTree(mustAbs(filepath.Join("testdata", "size")))
			if err != nil {
				t.Fatalf("Fatal parse error: %v", err)
			}

			var errs []string
			rootPkg := parse.RootPkg(pkgs)
			t.Logf("root package: %s", rootPkg)
			for _, pkg := range pkgs {
				errs = addErrors(errs, size.Check(pkg, rootPkg, spec.givenMaxSize))
			}
			if len(errs) != spec.expectedErrors {
				t.Errorf("Expected %d errors but got %d: %q", spec.expectedErrors, len(errs), errs)
			}
		})
	}
}

func addErrors(allErrs []string, newErrs []error) []string {
	for _, err := range newErrs {
		allErrs = append(allErrs, err.Error())
	}
	return allErrs
}

func mustAbs(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err.Error())
	}
	return absPath
}
