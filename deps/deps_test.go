package deps_test

import (
	"path/filepath"
	"testing"

	"github.com/flowdev/spaghetti-cutter/config"
	"github.com/flowdev/spaghetti-cutter/deps"
	"github.com/flowdev/spaghetti-cutter/parse"
	"github.com/flowdev/spaghetti-cutter/x/pkgs"
)

func TestCheck(t *testing.T) {
	specs := []struct {
		name           string
		givenRoot      string
		givenConfig    string
		expectedErrors int
	}{
		{
			name:           "no-config-one-pkg",
			givenRoot:      "one-pkg",
			givenConfig:    `{}`,
			expectedErrors: 0,
		}, {
			name:           "no-config-only-tools",
			givenRoot:      "only-tools",
			givenConfig:    `{}`,
			expectedErrors: 1,
		}, {
			name:           "allow-tool-only-tools",
			givenRoot:      "only-tools",
			givenConfig:    `{"allowAdditionally": {"x/tool2": ["x/tool"]} }`,
			expectedErrors: 0,
		}, {
			name:           "no-config-standard-proj",
			givenRoot:      "standard-proj",
			givenConfig:    `{}`,
			expectedErrors: 7,
		}, {
			name:           "standard-config-standard-proj",
			givenRoot:      "standard-proj",
			givenConfig:    `{"tool": ["x/*"], "db": ["db/*"], "allowAdditionally": {"db/store": ["db/model"]} }`,
			expectedErrors: 0,
		}, {
			name:           "standard-config-complex-proj",
			givenRoot:      "complex-proj",
			givenConfig:    `{"tool": ["pkg/x/*"], "db": ["pkg/db/*"], "allowAdditionally": {"pkg/db/store": ["pkg/db/model"]} }`,
			expectedErrors: 1,
		}, {
			name:      "allowOnlyIn-config-complex-proj",
			givenRoot: "complex-proj",
			givenConfig: `{
						"allowOnlyIn": {"pkg/domain3": ["pkg/domain4", "cmd/exe2"]},
						"tool": ["pkg/x/*"], "db": ["pkg/db/*"],
						"allowAdditionally": {"pkg/db/store": ["pkg/db/model"]}
					}`,
			expectedErrors: 0,
		}, {
			name:      "bad-allowOnlyIn-config-complex-proj",
			givenRoot: "complex-proj",
			givenConfig: `{
						"allowOnlyIn": {"pkg/domain3": ["pkg/domain1", "cmd/exe2"]},
						"tool": ["pkg/x/*"], "db": ["pkg/db/*"]
						"allowAdditionally": {"pkg/db/store": ["pkg/db/model"]}
					}`,
			expectedErrors: 1,
		}, {
			name:      "explicit-config-complex-proj",
			givenRoot: "complex-proj",
			givenConfig: `{
						"tool": ["pkg/x/*"], "db": ["pkg/db/*"],
						"allowAdditionally": {
						  "pkg/domain4": ["pkg/domain3"],
						  "cmd/exe1": ["pkg/domain1", "pkg/domain2"],
						  "cmd/exe2": ["pkg/domain3", "pkg/domain4"],
						  "pkg/db/store": ["pkg/db/model"]
						},
						"noGod": true
					}`,
			expectedErrors: 0,
		}, {
			name:           "standard-config-half-pkgs-proj",
			givenRoot:      "half-pkgs-proj",
			givenConfig:    `{tool: ["x/*"], db: ["db/*"], "allowAdditionally": {"db/store": ["db/model"]} }`,
			expectedErrors: 2,
		}, {
			name:      "explicit-config-half-pkgs-proj",
			givenRoot: "half-pkgs-proj",
			givenConfig: `{
						tool: ["x/*"]
						db: ["db/*"]
						allowAdditionally: {
							// they just grew out of bounds!
							"x/tool": ["x/tool/subtool"]
							"db/store": ["db/model", "db/store/substore"]
						}
					}`,
			expectedErrors: 0,
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			cfg, err := config.Parse([]byte(spec.givenConfig), spec.name)
			if err != nil {
				t.Fatalf("got unexpected error: %v", err)
			}

			packs, err := parse.DirTree(mustAbs(filepath.Join("testdata", spec.givenRoot)))
			if err != nil {
				t.Fatalf("Fatal parse error: %v", err)
			}

			var errs []string
			rootPkg := parse.RootPkg(packs)
			t.Logf("root package: %s", rootPkg)
			pkgInfos := pkgs.UniquePackages(packs)
			for _, pkgInfo := range pkgInfos {
				errs = addErrors(errs, deps.Check(pkgInfo.Pkg, rootPkg, cfg))
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
