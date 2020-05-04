package dirs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// FindConfig finds the .spaghetti-cutter.json configuration file of a project.
func FindConfig(cfgFile string) string {
	dir := crawlUpAndFindDirOf(".", cfgFile)
	if dir != "" {
		return filepath.Join(dir, cfgFile)
	}
	return ""
}

// FindRoot finds the root of a project.
// It looks at the following things (highest priority first):
// - The given directory (used unless empty)
// - It looks for go.mod via `go env GOMOD`
// - It looks for the configuration file: .spaghetti-cutter.json
// - It looks for a `vendor` directory.
func FindRoot(dir string, cfgFile string, ignoreVendor bool) (string, error) {
	if dir != "" {
		return dir, nil
	}

	if ignoreVendor {
		dir = crawlUpAndFindDirOf(".", "go.mod", cfgFile)
	} else {
		dir = crawlUpAndFindDirOf(".", "go.mod", cfgFile, "vendor")
	}
	if dir == "" {
		absDir, _ := filepath.Abs(".") // we checked this just inside of crawlUpAndFindDirOf()
		return "", fmt.Errorf("unable to find root directory for: %s", absDir)
	}

	return dir, nil
}

func crawlUpAndFindDirOf(startDir string, files ...string) string {
	absDir, err := filepath.Abs(startDir)
	if err != nil {
		log.Fatalf("FATAL - Unable to find absolute directory (for %s): %v", startDir, err)
	}
	volName := filepath.VolumeName(absDir)
	oldDir := "" // set to impossible value first!

	for ; absDir != volName && absDir != oldDir; absDir = filepath.Dir(absDir) {
		for _, file := range files {
			path := filepath.Join(absDir, file)
			if _, err = os.Stat(path); err == nil {
				return absDir
			}
		}
		oldDir = absDir
	}
	return ""
}
