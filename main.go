package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/flowdev/spaghetti-cutter/data"
	"github.com/flowdev/spaghetti-cutter/deps"
	"github.com/flowdev/spaghetti-cutter/doc"
	"github.com/flowdev/spaghetti-cutter/parse"
	"github.com/flowdev/spaghetti-cutter/size"
	"github.com/flowdev/spaghetti-cutter/stat"
	"github.com/flowdev/spaghetti-cutter/x/config"
	"github.com/flowdev/spaghetti-cutter/x/dirs"
	"github.com/flowdev/spaghetti-cutter/x/pkgs"
)

func main() {
	rc := cut(os.Args[1:])
	if rc != 0 {
		os.Exit(rc)
	}
}

func cut(args []string) int {
	const (
		usageShort     = " (shorthand)"
		defaultRoot    = "."
		usageRoot      = "root directory of the project"
		defaultDoc     = "*"
		usageDoc       = "write '" + doc.FileName + "' for packages (separated by ','; '' for none)"
		defaultNoLinks = false
		usageNoLinks   = "don't use links in '" + doc.FileName + "' files"
		defaultStats   = "*"
		usageStats     = "write '" + stat.FileName + "' for packages (separated by ','; '' for none)"
	)
	var startDir string
	var docPkgs string
	var noLinks bool
	var statPkgs string
	fs := flag.NewFlagSet("spaghetti-cutter", flag.ExitOnError)
	fs.StringVar(&startDir, "root", defaultRoot, usageRoot)
	fs.StringVar(&startDir, "r", defaultRoot, usageRoot+usageShort)
	fs.StringVar(&docPkgs, "doc", defaultDoc, usageDoc)
	fs.StringVar(&docPkgs, "d", defaultDoc, usageDoc+usageShort)
	fs.BoolVar(&noLinks, "nolinks", defaultNoLinks, usageNoLinks)
	fs.BoolVar(&noLinks, "l", defaultNoLinks, usageNoLinks+usageShort)
	fs.StringVar(&statPkgs, "stats", defaultStats, usageStats)
	fs.StringVar(&statPkgs, "s", defaultStats, usageStats+usageShort)
	err := fs.Parse(args)
	if err != nil {
		log.Printf("FATAL - %v", err)
		return 2
	}

	root, err := dirs.FindRoot(startDir, config.File)
	if err != nil {
		log.Printf("FATAL - %v", err)
		return 3
	}
	cfgFile := filepath.Join(root, config.File)
	cfgBytes, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		log.Printf("FATAL - unable to read configuration file %q: %v", cfgFile, err)
		return 4
	}
	cfg, err := config.Parse(cfgBytes, cfgFile)
	if err != nil {
		log.Printf("FATAL - %v", err)
		return 5
	}

	log.Printf("INFO - configuration 'allowOnlyIn': %s", cfg.AllowOnlyIn)
	log.Printf("INFO - configuration 'allowAdditionally': %s", cfg.AllowAdditionally)
	log.Printf("INFO - configuration 'god': %s", cfg.God)
	log.Printf("INFO - configuration 'tool': %s", cfg.Tool)
	log.Printf("INFO - configuration 'db': %s", cfg.DB)
	log.Printf("INFO - configuration 'size': %d", cfg.Size)
	log.Printf("INFO - configuration 'noGod': %t", cfg.NoGod)
	log.Printf("INFO - documenting package(s): %s", docPkgs)
	log.Printf("INFO - no links in '"+doc.FileName+"' files: %t", noLinks)

	packs, err := parse.DirTree(root)
	if err != nil {
		log.Printf("FATAL - %v", err)
		return 6
	}

	var errs []error
	depMap := make(data.DependencyMap, 256)

	rootPkg := parse.RootPkg(packs)
	log.Printf("INFO - root package: %s", rootPkg)
	pkgInfos := pkgs.UniquePackages(packs)
	for _, pkgInfo := range pkgInfos {
		errs = addErrors(errs, deps.Check(pkgInfo.Pkg, rootPkg, cfg, &depMap))
		errs = addErrors(errs, size.Check(pkgInfo.Pkg, rootPkg, cfg.Size))
	}

	retCode := 0
	if len(errs) > 0 {
		for _, err = range errs {
			log.Printf("ERROR - %v", err)
		}
		retCode = 1
	}

	log.Print("INFO - No errors found.")

	if statPkgs != "" {
		writeStatistics(statPkgs, root, rootPkg, depMap)
	}

	if docPkgs == "" {
		log.Print("INFO - No documentation wanted.")
		return retCode
	}

	writeDocumentation(docPkgs, root, rootPkg, noLinks, depMap)

	return retCode
}

func writeStatistics(stPkgs, root, rootPkg string, depMap data.DependencyMap) {
	log.Print("INFO - Writing statistics.")
	statPkgs := findPackagesWithFileAsSlice(stat.FileName, stPkgs, root, "statistics")

	for _, statPkg := range statPkgs {
		statMD := stat.Generate(statPkg, depMap)
		if statMD == "" {
			continue
		}
		statFile := filepath.Join(statPkg, stat.FileName)
		log.Printf("INFO - Write package statistics to file: %s", statFile)
		statFile = filepath.Join(root, statFile)
		err := ioutil.WriteFile(statFile, []byte(statMD), 0644)
		if err != nil {
			log.Printf("ERROR - Unable to write package statistics to file %s: %v", statFile, err)
		}
	}
}

func writeDocumentation(docPkgs, root, rootPkg string, noLinks bool, depMap data.DependencyMap) {
	log.Print("INFO - Writing documentation.")
	dtPkgs := findPackagesWithFileAsSlice(doc.FileName, docPkgs, root, "documentation")

	linkDocPkgs := map[string]struct{}{}
	if !noLinks {
		linkDocPkgs = dirs.FindPkgsWithFile(doc.FileName, dtPkgs, root, true)
	}
	doc.WriteDocs(dtPkgs, depMap, linkDocPkgs, rootPkg, root)
}

func findPackagesWithFileAsSlice(signalFile, pkgNames, root, pkgType string) []string {
	var pkgs []string
	if pkgNames == "*" { // find all existing files
		pkgMap := dirs.FindPkgsWithFile(signalFile, nil, root, false)
		pkgs = make([]string, 0, len(pkgMap))
		for p := range pkgMap {
			pkgs = append(pkgs, p)
		}
	} else { // write explicitly given docs
		pkgs = splitPackageNames(pkgNames, pkgType)
	}
	return pkgs
}

func splitPackageNames(docPkgs, pkgType string) []string {
	splitPkgs := strings.Split(docPkgs, ",")
	retPkgs := make([]string, 0, len(splitPkgs))
	for i, splitPkg := range splitPkgs {
		pkg := strings.TrimSpace(splitPkg)
		if pkg == "" {
			log.Printf("INFO - Not writing %s for %d-th package because the name is empty.", pkgType, i+1)
			continue
		}
		retPkgs = append(retPkgs, pkg)
	}
	return retPkgs
}

func addErrors(errs []error, newErrs []error) []error {
	return append(errs, newErrs...)
}
