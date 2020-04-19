package parse

import (
	"errors"
	"fmt"

	"golang.org/x/tools/go/packages"
)

// DirTree is parsing the whole directory tree starting at cfg.Root
// looking for Go packages and analyzing them.
func DirTree(root string) ([]*packages.Package, error) {
	fmt.Println("Parse root:", root)

	parseCfg := &packages.Config{
		Logf:  nil, // log.Printf,
		Dir:   root,
		Tests: true,
		Mode:  packages.NeedName | packages.NeedImports | packages.NeedSyntax,
	}

	pkgs, err := packages.Load(parseCfg, root+"/...")
	if err != nil {
		return nil, err
	}
	if packages.PrintErrors(pkgs) > 0 {
		return nil, errors.New("unable to parse packages")
	}
	return pkgs, nil
}
