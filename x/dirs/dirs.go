package dirs

import (
	"fmt"
	"os"
	"path/filepath"
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
