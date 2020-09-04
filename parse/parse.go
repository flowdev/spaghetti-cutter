package parse

import (
	"errors"
	"strings"

	"github.com/flowdev/spaghetti-cutter/x/pkgs"
	"golang.org/x/tools/go/packages"
)

// DirTree is parsing the whole directory tree starting at root
// looking for Go packages and analyzing them.
func DirTree(root string) ([]*pkgs.Package, error) {
	parseCfg := &packages.Config{
		Logf:  nil, // log.Printf (for debug), nil (for release)
		Dir:   root,
		Tests: true,
		Mode:  packages.NeedName | packages.NeedImports | packages.NeedSyntax,
	}

	pkgs, err := packages.Load(parseCfg, root+"/...")
	if err != nil {
		return nil, err
	}
	if packages.PrintErrors(pkgs) > 0 {
		return nil, errors.New("unable to parse packages at root: " + root)
	}
	return pkgs, nil
}

// RootPkg returns the package path of the root package of all the given
// packages.
// It does so by returning the longest common prefix of all package paths.
func RootPkg(pkgs []*pkgs.Package) string {
	root := ""
	for _, pkg := range pkgs {
		if root == "" {
			root = pkg.PkgPath
		} else {
			root = commonPrefix(root, pkg.PkgPath)
		}
	}
	return root
}

func commonPrefix(s1, s2 string) string {
	sl1 := strings.Split(s1, "")
	sl2 := strings.Split(s2, "")
	n := min(len(sl1), len(sl2))
	b := strings.Builder{}

	for i := 0; i < n; i++ {
		if sl1[i] == sl2[i] {
			b.WriteString(sl1[i])
		} else {
			break
		}
	}
	return b.String()
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
