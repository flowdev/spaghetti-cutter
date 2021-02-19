package main

import (
	"flag"
	"fmt"
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
		defaultStats   = ""
		usageStats     = "print statistics about given package"
	)
	var startDir string
	var docPkgs string
	var noLinks bool
	var stats string
	fs := flag.NewFlagSet("spaghetti-cutter", flag.ExitOnError)
	fs.StringVar(&startDir, "root", defaultRoot, usageRoot)
	fs.StringVar(&startDir, "r", defaultRoot, usageRoot+usageShort)
	fs.StringVar(&docPkgs, "doc", "/", usageDoc)
	fs.StringVar(&docPkgs, "d", "/", usageDoc+usageShort)
	fs.BoolVar(&noLinks, "nolinks", defaultNoLinks, usageNoLinks)
	fs.BoolVar(&noLinks, "l", defaultNoLinks, usageNoLinks+usageShort)
	fs.StringVar(&stats, "stats", defaultStats, usageStats)
	fs.StringVar(&stats, "s", defaultStats, usageStats+usageShort)
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

	if stats != "" {
		for _, s := range stat.Create(stats, depMap) {
			fmt.Println(s)
		}
	}

	if docPkgs == "" {
		log.Print("INFO - No documentation wanted.")
		return retCode
	}

	log.Print("INFO - Writing documentation.")
	var dtPkgs []string
	if docPkgs == "*" { // update all existing docs
		dtPkgMap := doc.FindDocPkgs(nil, root, false)
		dtPkgs = make([]string, 0, len(dtPkgMap))
		for p := range dtPkgMap {
			dtPkgs = append(dtPkgs, p)
		}
	} else { // write explicitly given docs
		dtPkgs = splitDocPackages(docPkgs)
	}
	linkDocPkgs := map[string]struct{}{}
	if !noLinks {
		linkDocPkgs = doc.FindDocPkgs(dtPkgs, root, true)
	}
	doc.WriteDocs(dtPkgs, depMap, linkDocPkgs, rootPkg, root)

	return retCode
}

func splitDocPackages(docPkgs string) []string {
	dtPkgs := strings.Split(docPkgs, ",")
	retPkgs := make([]string, 0, len(dtPkgs))
	for i, dtPkg := range dtPkgs {
		dtp := strings.TrimSpace(dtPkg)
		if dtp == "" {
			log.Printf("INFO - Not writing documentation for %d-th package because the name is empty.", i+1)
			continue
		}
		retPkgs = append(retPkgs, dtp)
	}
	return retPkgs
}

func addErrors(errs []error, newErrs []error) []error {
	return append(errs, newErrs...)
}
