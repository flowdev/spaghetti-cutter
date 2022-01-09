package config_test

import (
	"fmt"
	"testing"

	"github.com/flowdev/spaghetti-cutter/config"
)

func TestParseAndStringers(t *testing.T) {
	specs := []struct {
		name                 string
		givenConfigBytes     []byte
		expectedConfigString string
	}{
		{
			name: "all-empty",
			givenConfigBytes: []byte(`{
				}`),
			expectedConfigString: "{" +
				"..... ..... ... ... `main` " +
				"2048 false" +
				"}",
		}, {
			name: "scalars-only",
			givenConfigBytes: []byte(`{
				  "size": 3072,
				  "noGod": true
				}`),
			expectedConfigString: "{" +
				"..... ..... ... ... ... " +
				"3072 true" +
				"}",
		}, {
			name: "list-one",
			givenConfigBytes: []byte(`{
				  "god": ["a"]
				}`),
			expectedConfigString: "{" +
				"..... ..... " +
				"... " +
				"... " +
				"`a` " +
				"2048 false" +
				"}",
		}, {
			name: "list-many",
			givenConfigBytes: []byte(`{
				  "god": ["a", "be", "do", "ra"]
				}`),
			expectedConfigString: "{" +
				"..... ..... " +
				"... " +
				"... " +
				"`a`, `be`, `do`, `ra` " +
				"2048 false" +
				"}",
		}, {
			name: "map-simple-pair",
			givenConfigBytes: []byte(`{
				  "allowOnlyIn": {
				    "a": ["b"]
				  }
				}`),
			expectedConfigString: "{" +
				"`a`: `b` " +
				"..... " +
				"... ... `main` 2048 false" +
				"}",
		}, {
			name: "map-multiple-pairs",
			givenConfigBytes: []byte(`{
				  "allowOnlyIn": {
				    "a": ["b", "c", "do", "foo"],
				    "e": ["bar", "car"]
				  }
				}`),
			expectedConfigString: "{" +
				"`a`: `b`, `c`, `do`, `foo` ; `e`: `bar`, `car` " +
				"..... " +
				"... ... `main` 2048 false" +
				"}",
		}, {
			name: "map-one-pair-many-stars",
			givenConfigBytes: []byte(`{
				  "allowOnlyIn": {
				    "a/*/b/**": ["c/*/d/**"]
				  }
				}`),
			expectedConfigString: "{" +
				"`a/*/b/**`: `c/*/d/**` " +
				"..... " +
				"... ... `main` 2048 false" +
				"}",
		}, {
			name: "map-all-complexity",
			givenConfigBytes: []byte(`{
				  "allowAdditionally": {
				    "*/*a/**": ["*/*b/**", "b*/c*d/**"]
				  }
				}`),
			expectedConfigString: "{" +
				"..... " +
				"`*/*a/**`: `*/*b/**`, `b*/c*d/**` " +
				"... ... `main` 2048 false" +
				"}",
		}, {
			name: "maps-and-lists-only",
			givenConfigBytes: []byte(`{
					"allowOnlyIn": {
					  "github.com/lib/pq": ["a"]
					},
					"allowAdditionally": {
					  "a": ["b"]
					},
					"tool": ["x/**"],
					"db": ["pkg/db/*"],
					"god": ["main"]
				}`),
			expectedConfigString: "{" +
				"`github.com/lib/pq`: `a` " +
				"`a`: `b` " +
				"`x/**` " +
				"`pkg/db/*` " +
				"`main` " +
				"2048 " +
				"false" +
				"}",
		}, {
			name: "a-bit-of-everything",
			givenConfigBytes: []byte(`{
					"allowOnlyIn": {
					  "github.com/lib/pq": ["a", "b"]
					},
					"allowAdditionally": {
					  "a": ["b"],
					  "c": ["d"]
					},
					"tool": ["pkg/mysupertool", "pkg/x/**"],
					"db": ["pkg/db", "pkg/entities"],
					"god": ["main", "pkg/service"],
					"size": 3072,
					"noGod": true
				}`),
			expectedConfigString: "{" +
				"`github.com/lib/pq`: `a`, `b` " +
				"`a`: `b` ; `c`: `d` " +
				"`pkg/mysupertool`, `pkg/x/**` " +
				"`pkg/db`, `pkg/entities` " +
				"`main`, `pkg/service` " +
				"3072 " +
				"true" +
				"}",
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			actualConfig, err := config.Parse(spec.givenConfigBytes, spec.name)
			if err != nil {
				t.Fatalf("got unexpected error: %v", err)
			}
			actualConfigString := fmt.Sprint(actualConfig)
			if actualConfigString != spec.expectedConfigString {
				t.Errorf("expected configuration %v, actual %v",
					spec.expectedConfigString, actualConfigString)
			}
		})
	}
}
