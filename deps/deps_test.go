package deps_test

import (
	"path/filepath"
	"testing"

	"github.com/flowdev/spaghetti-cutter/deps"
	"github.com/flowdev/spaghetti-cutter/parse"
	"github.com/flowdev/spaghetti-cutter/x/config"
)

func TestCheck(t *testing.T) {
	specs := []struct {
		name           string
		givenRoot      string
		givenArgs      []string
		expectedErrors int
	}{
		{
			name:           "no-args-one-pkg",
			givenRoot:      "one-pkg",
			givenArgs:      nil,
			expectedErrors: 0,
		}, {
			name:           "no-args-only-tools",
			givenRoot:      "only-tools",
			givenArgs:      nil,
			expectedErrors: 1,
		}, {
			name:           "allow-tool-only-tools",
			givenRoot:      "only-tools",
			givenArgs:      []string{"-allow", "x/tool2 x/tool"},
			expectedErrors: 0,
		}, {
			name:           "no-args-standard-proj",
			givenRoot:      "standard-proj",
			givenArgs:      nil,
			expectedErrors: 7,
		}, {
			name:           "standard-args-standard-proj",
			givenRoot:      "standard-proj",
			givenArgs:      []string{"--tool", "x/*", "--db", "db/*"},
			expectedErrors: 0,
		}, {
			name:           "standard-args-complex-proj",
			givenRoot:      "complex-proj",
			givenArgs:      []string{"--tool", "pkg/x/*", "--db", "pkg/db/*"},
			expectedErrors: 1,
		}, {
			name:      "explicit-args-complex-proj",
			givenRoot: "complex-proj",
			givenArgs: []string{
				"--tool", "pkg/x/*", "--db", "pkg/db/*",
				"--allow", "pkg/domain4 pkg/domain3",
				"--allow", "cmd/exe1 pkg/domain1", "--allow", "cmd/exe1 pkg/domain2",
				"--allow", "cmd/exe2 pkg/domain3", "--allow", "cmd/exe2 pkg/domain4",
				"--no-god",
			},
			expectedErrors: 0,
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			cfg := config.Parse(spec.givenArgs, "")

			pkgs, err := parse.DirTree(mustAbs(filepath.Join("testdata", spec.givenRoot)))
			if err != nil {
				t.Fatalf("Fatal parse error: %v", err)
			}

			var errs []string
			rootPkg := parse.RootPkg(pkgs)
			t.Logf("root package: %s", rootPkg)
			for _, pkg := range pkgs {
				errs = addErrors(errs, deps.Check(pkg, rootPkg, cfg))
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
