package dirs

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	TreeFile     = "dirtree.txt"
	newLine      = "\n"
	emptyItem    = "    "
	middleItem   = "├── "
	continueItem = "│   "
	lastItem     = "└── "
)

// FindRoot finds the root of a project.
// It looks for the configuration file: .spaghetti-cutter.json
func FindRoot(startDir, cfgFile string) (string, error) {
	if startDir == "" {
		startDir = "."
	}
	dir, err := crawlUpAndFindDirOf(startDir, cfgFile)
	if err != nil {
		return "", err
	}
	if dir == "" {
		absDir, _ := filepath.Abs(".") // we checked this just inside of crawlUpAndFindDirOf()
		return "", fmt.Errorf("unable to find root directory for: %s", absDir)
	}

	return dir, nil
}

func crawlUpAndFindDirOf(startDir string, files ...string) (string, error) {
	absDir, err := filepath.Abs(startDir)
	if err != nil {
		return "", fmt.Errorf("unable to find absolute directory (for %q): %w", startDir, err)
	}
	volName := filepath.VolumeName(absDir)
	oldDir := "" // set to impossible value first!

	for ; absDir != volName && absDir != oldDir; absDir = filepath.Dir(absDir) {
		for _, file := range files {
			path := filepath.Join(absDir, file)
			if _, err = os.Stat(path); err == nil {
				return absDir, nil
			}
		}
		oldDir = absDir
	}
	return "", nil
}

// FindPkgsWithFile is finding packages containing file on disk starting at
// 'root' and adding them to those given in 'startPkgs'.
func FindPkgsWithFile(file string, startPkgs []string, root string, excludeRoot bool) map[string]struct{} {
	val := struct{}{}
	// prefill doc packages from dtPkgs
	retPkgs := make(map[string]struct{}, 128)
	for _, p := range startPkgs {
		retPkgs[p] = val
	}

	// walk the file system to find more 'file's
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() { // we are only interested in directories
			return nil
		}
		if err != nil {
			log.Printf("WARN - Unable to list directory %q: %v", path, err)
			return filepath.SkipDir
		}
		if excludeRoot && path == root {
			return nil // don't add the root 'file'
		}

		// no valid package starts with '.' and we don't want to search in '.git' and similar
		if strings.HasPrefix(info.Name(), ".") || info.Name() == "testdata" {
			return filepath.SkipDir
		}

		if _, err := os.Lstat(filepath.Join(path, file)); err == nil {
			pkg, err := filepath.Rel(root, path)
			if err != nil {
				log.Printf("WARN - Unable to compute package for %q: %v", path, err)
				return nil // sub-directories might work
			}
			pkg = strings.ReplaceAll(pkg, "\\", "/") // packages like URLs have always '/'s
			if pkg == "." {
				retPkgs["/"] = val
			} else {
				retPkgs[pkg] = val
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("ERROR - Unable to walk the path %q: %v", root, err)
	}
	return retPkgs
}

func Tree(root, name string, exclude []string) (string, error) {
	sb := &strings.Builder{}

	err := generateTree(root, name, sb, "", exclude)
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

func generateTree(root, name string, sb *strings.Builder, prefix string, exclude []string) error {
	sb.WriteString(name)
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
			generateTree(filepath.Join(root, item.Name()), item.Name(), sb, prefix+emptyItem, exclude)
		} else {
			sb.WriteString(prefix + middleItem)
			generateTree(filepath.Join(root, item.Name()), item.Name(), sb, prefix+continueItem, exclude)
		}
	}
	return nil
}

func includeFile(name string, exclude []string) bool {
	for _, ex := range exclude {
		if m, _ := path.Match(ex, name); m {
			return false
		}
	}
	return true
}
