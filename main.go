package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("Root dir is: " + findRootDir())
}

func findRootDir() string {
	dir := findGoModDir()
	if dir != "" {
		return dir
	}

	dir = findVendorDir()
	if dir == "" {
		absDir, _ := filepath.Abs(".") // we checked this just inside of findVendorDir()
		log.Fatalf("FATAL: Unable to find root directory for '%s'.", absDir)
	}

	return dir
}
func findGoModDir() string {
	gomod := getOutputOfCmd("go", "env", "GOMOD")
	if gomod == os.DevNull {
		return ""
	}
	if gomod == "" {
		return ""
	}
	return path.Dir(gomod)
}
func findVendorDir() string {
	return crawlUpAndFindDirOf("vendor", ".")
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
