package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/flowdev/spaghetti-cutter/config"
	"github.com/flowdev/spaghetti-cutter/parse"
)

func main() {
	cfg := config.Parse(os.Args[1:])
	cfg.God["main"] = config.Value // the main package can always access everything
	cfg.Root = findRootDir(cfg.Root)
	pkgs, err := parse.DirTree(cfg.Root)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(2)
	}

	for _, pkg := range pkgs {
		fmt.Println(pkg.ID, pkg.Name, pkg.PkgPath)
	}
}

func findRootDir(dir string) string {
	if dir != "" {
		return dir
	}

	dir = findGoModDir()
	if dir != "" {
		return dir
	}

	dir = crawlUpAndFindDirOf(config.File, ".")
	if dir != "" {
		return dir
	}

	dir = crawlUpAndFindDirOf("vendor", ".")
	if dir == "" {
		absDir, _ := filepath.Abs(".") // we checked this just inside of findVendorDir()
		log.Fatalf("FATAL: Unable to find root directory for '%s'.", absDir)
	}

	return dir
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
