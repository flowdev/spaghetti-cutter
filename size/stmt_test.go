package size

import (
	"go/ast"
	"path/filepath"
	"testing"

	"github.com/flowdev/spaghetti-cutter/parse"
)

func TestSizeOfStmt(t *testing.T) {
	pkgs, err := parse.DirTree(mustAbs(filepath.Join("testdata", "stmt")))
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	for _, pkg := range pkgs { // packages contain
		for _, astf := range pkg.Syntax { // files that contain
			for _, decl := range astf.Decls { // declarations that consist of
				switch d := decl.(type) {
				case *ast.FuncDecl: // functions and
					id, expectedSize := nameToIDandSize(d.Name.Name)
					t.Run(id, func(t *testing.T) {
						actualSize := uint(0)
						for _, stmt := range d.Body.List {
							actualSize += sizeOfStmt(stmt)
						}
						if actualSize != expectedSize {
							t.Errorf("expected size %d but got: %d", expectedSize, actualSize)
						}
					})
				case *ast.GenDecl: // general declaration that have
					for _, spec := range d.Specs { // specifications
						switch spec.(type) {
						case *ast.ImportSpec:
							// ignore imports
						default:
							t.Errorf("unknown decl spec: %T", spec)
						}
					}
				default:
					t.Errorf("unknown decl : %T", decl)
				}
			}
		}
	}
}
