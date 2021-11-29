package tree

import (
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/flowdev/spaghetti-cutter/x/pkgs"
)

const (
	File         = "dirtree.txt"
	newLine      = "\n"
	separator    = " -\t"
	emptyItem    = "    "
	middleItem   = "├── "
	continueItem = "│   "
	lastItem     = "└── "
)

func Generate(root, name string, exclude []string, packs []*pkgs.Package) (string, error) {
	sb := &strings.Builder{}

	err := generateTree(root, name, sb, "", exclude, "", packs)
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

func generateTree(root, name string, sb *strings.Builder, prefix string, exclude []string, pkg string, packs []*pkgs.Package) error {
	sb.WriteString(name)
	sb.WriteString(separator)
	sb.WriteString(docForPkg(pkg, packs))
	sb.WriteString(newLine)

	files, err := os.ReadDir(root)
	if err != nil {
		log.Printf("ERROR - Unable to read the directory %q: %v", root, err)
		return err
	}

	// reduce all files to only the items we want to include
	var items []fs.DirEntry
	for _, file := range files {
		if file.IsDir() && includeFile(file.Name(), exclude) {
			items = append(items, file)
		}
	}

	lastI := len(items) - 1
	for i, item := range items {
		if i == lastI {
			sb.WriteString(prefix + lastItem)
			err = generateTree(filepath.Join(root, item.Name()), item.Name(), sb, prefix+emptyItem, exclude, "", nil)
			if err != nil {
				return err
			}
		} else {
			sb.WriteString(prefix + middleItem)
			err = generateTree(filepath.Join(root, item.Name()), item.Name(), sb, prefix+continueItem, exclude, "", nil)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func docForPkg(pkg string, packs []*pkgs.Package) string {
	return ""
}

func includeFile(name string, exclude []string) bool {
	for _, ex := range exclude {
		if m, _ := path.Match(ex, name); m {
			return false
		}
	}
	return true
}
