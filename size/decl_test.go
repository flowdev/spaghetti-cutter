package size

import (
	"go/ast"
	"path/filepath"
	"testing"

	"github.com/flowdev/spaghetti-cutter/parse"
)

func TestSizeOfDecl(t *testing.T) {
	pkgs, err := parse.DirTree(mustAbs(filepath.Join("testdata", "decl")))
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	for _, pkg := range pkgs { // packages contain
		for _, astf := range pkg.Syntax { // files that contain
			for _, decl := range astf.Decls { // declarations that consist of
				var id string
				var expectedSize uint
				switch d := decl.(type) {
				case *ast.FuncDecl: // functions and
					id, expectedSize = nameToIDandSize(d.Name.Name)
				case *ast.GenDecl: // general declaration that have
				SpecLoop:
					for _, spec := range d.Specs { // specifications
						switch s := spec.(type) {
						case *ast.ImportSpec:
							// ignore imports
						case *ast.TypeSpec:
							id, expectedSize = nameToIDandSize(s.Name.Name)
						case *ast.ValueSpec:
							if len(s.Names) < 1 {
								t.Errorf("ValueSpec without a name: %#v", spec)
							} else {
								id, expectedSize = nameToIDandSize(s.Names[0].Name)
							}
							break SpecLoop // only look at the FIRST ValueSpec for the size
						default:
							t.Errorf("unknown decl spec: %T", spec)
						}
					}
				default:
					t.Errorf("unknown decl : %T", decl)
				}
				t.Run(id, func(t *testing.T) {
					actualSize := sizeOfDecl(decl)
					if actualSize != expectedSize {
						t.Errorf("expected size %d but got: %d", expectedSize, actualSize)
					}
				})
			}
		}
	}
}
