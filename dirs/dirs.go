package dirs

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/flowdev/spaghetti-cutter/x/config"
)

// FindRoot finds the root of a project.
// It looks at the following things (highest priority first):
// - The given directory (used unless empty)
// - It looks for go.mod via `go env GOMOD`
// - It looks for the configuration file: .spaghetti-cutter.json
// - It looks for a `vendor` directory.
func FindRoot(dir string) (string, error) {
	if dir != "" {
		return dir, nil
	}

	dir = findGoModDir()
	if dir != "" {
		return dir, nil
	}

	dir = crawlUpAndFindDirOf(config.File, ".")
	if dir != "" {
		return dir, nil
	}

	dir = crawlUpAndFindDirOf("vendor", ".")
	if dir == "" {
		absDir, _ := filepath.Abs(".") // we checked this just inside of crawlUpAndFindDirOf()
		return "", fmt.Errorf("unable to find root directory for: %s", absDir)
	}

	return dir, nil
}

func findGoModDir() string {
	gomod := getOutputOfCmd("go", "env", "GOMOD")
	if gomod == os.DevNull || gomod == "" {
		return ""
	}
	return path.Dir(gomod)
}

func getOutputOfCmd(cmd string, args ...string) string {
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		log.Fatalf("FATAL: Unable to execute command: %v", err)
	}
	return strings.TrimRight(string(out), "\r\n")
}

func crawlUpAndFindDirOf(file, startDir string) string {
	absDir, err := filepath.Abs(startDir)
	if err != nil {
		log.Fatalf("FATAL: Unable to find absolute directory (for %s): %v", startDir, err)
	}
	volName := filepath.VolumeName(absDir)
	oldDir := "" // set to impossible value first!

	for ; absDir != volName && absDir != oldDir; absDir = filepath.Dir(absDir) {
		path := filepath.Join(absDir, file)
		if _, err = os.Stat(path); err == nil {
			return absDir
		}
		oldDir = absDir
	}
	return ""
}
