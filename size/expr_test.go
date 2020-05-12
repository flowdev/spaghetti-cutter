package size

import (
	"go/ast"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/flowdev/spaghetti-cutter/parse"
)

func TestSizeOfExpr(t *testing.T) {
	pkgs, err := parse.DirTree(mustAbs(filepath.Join("testdata", "expr")))
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	for _, pkg := range pkgs { // packages contain
		for _, astf := range pkg.Syntax { // files that contain
			for _, decl := range astf.Decls { // declarations that consist of
				if d, ok := decl.(*ast.GenDecl); ok {
					for _, spec := range d.Specs { // specifications that have
						switch s := spec.(type) {
						case *ast.ValueSpec:
							if len(s.Names) != len(s.Values) {
								t.Errorf("unexpected length of decl spec: %d != %d", len(s.Names), len(s.Values))
								continue
							}
							for i, name := range s.Names { // names and values
								id, expectedSize := nameToIDandSize(name.Name)
								actualSize := sizeOfExpr(s.Values[i])
								if actualSize != expectedSize {
									t.Errorf("[%s] expected size %d but got: %d", id, expectedSize, actualSize)
								}
							}
						case *ast.ImportSpec:
							// ignore imports
						default:
							t.Errorf("unknown decl spec: %T", spec)
						}
					}
				} else {
					t.Errorf("unknown decl: %T", d)
				}
			}
		}
	}

}

func nameToIDandSize(name string) (id string, size uint) {
	parts := strings.SplitN(name, "_", 2)
	if len(parts) != 2 {
		panic("`" + name + "` doesn't contain an underscore")
	}

	i, err := strconv.Atoi(parts[1])
	if err != nil {
		panic("`" + parts[1] + "` isn't an integer number: " + err.Error())
	}
	if i < 0 {
		panic("negative size isn't allowed: " + parts[1])
	}

	return parts[0], uint(i)
}

func mustAbs(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err.Error())
	}
	return absPath
}
