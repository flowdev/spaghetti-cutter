package goast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"

	"github.com/flowdev/spaghetti-cutter/config"
)

const (
	goTestFileName = `_test.go`
)

//
// packageDict is a simple dictionary of all known packages/paths and
// their source parts.
//

type goPackage struct {
	path       string
	deps       map[string]struct{} // or: map[string]*goPackage
	complexity int
}

type packageDict struct {
	root  string
	packs map[string]*goPackage
}

func newPackageDict(root string) *packageDict {
	return &packageDict{
		root:  root,
		packs: make(map[string]*goPackage),
	}
}

//
// Parse and process the main directory/package.
//

// WalkDirTree is walking the directory tree starting at the given root path
// looking for Go packages and analyzing them.
func WalkDirTree(cfg config.Config) error {
	fmt.Println("Parsed config:")
	fmt.Println("    Ignore:", cfg.Ignore)
	fmt.Println("    Tool:", cfg.Tool)
	fmt.Println("    DB:", cfg.DB)
	fmt.Println("    God:", cfg.God)
	fmt.Println("    Allow:", cfg.Allow)
	fmt.Println("    Root:", cfg.Root)

	root := cfg.Root
	fset := token.NewFileSet() // needed for any kind of parsing
	packDict := newPackageDict(root)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("ERROR: While walking the path %q: %v", path, err)
			return err
		}
		if info.IsDir() {
			if _, ok := cfg.Ignore[info.Name()]; ok {
				log.Printf("INFO: ignoring dir: %s", path)
				return filepath.SkipDir
			}
			return processDir(path, fset, packDict)
		}
		return nil
	})
	if err != nil {
		log.Printf("ERROR: While walking the root %q: %v", root, err)
		return err
	}
	return nil
}

func processDir(dir string, fset *token.FileSet, packDict *packageDict) error {
	fmt.Println("Parsing the whole directory:", dir)
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.AllErrors)
	if err != nil {
		return fmt.Errorf("unable to parse the directory '%s': %v", dir, err)
	}
	for _, pkg := range pkgs { // iterate over subpackages (e.g.: xxx and xxx_test)
		if err := processPackage(pkg, fset, packDict); err != nil {
			return err
		}
	}
	return nil
}

// processPackage is processing all the files of one Go package.
func processPackage(pkg *ast.Package, fset *token.FileSet, packDict *packageDict) error {
	fmt.Println("processing package:", pkg.Name)

	for name, astf := range pkg.Files {
		fmt.Println("processing file:", name)
		//baseName := goNameToBase(name)
		for _, imp := range astf.Imports {
			fmt.Println("    found import:", imp.Path.Value, "line:", lineFor(imp.Path.ValuePos, fset))
		}
	}
	return nil
}

func goNameToBase(goname string) string {
	ext := filepath.Ext(goname)
	return goname[:len(goname)-len(ext)]
}
func lineFor(p token.Pos, fset *token.FileSet) int {
	if p.IsValid() {
		pos := fset.PositionFor(p, false)
		return pos.Line
	}

	return 0
}
