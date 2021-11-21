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
		defaultStats   = false
		usageStats     = "write '" + stat.FileName + "' for project"
		defaultDirTree = false
		usageDirTree   = "write a directory tree (starting at the current directory) to: dirtree.txt"
		defaultNoErr   = false
		usageNoErr     = "don't report errors or exit with an error"
	)
	var startDir string
	var docPkgs string
	var noLinks bool
	var doStats bool
	var dirTree bool
	var noErr bool
	fs := flag.NewFlagSet("spaghetti-cutter", flag.ExitOnError)
	fs.StringVar(&startDir, "root", defaultRoot, usageRoot)
	fs.StringVar(&startDir, "r", defaultRoot, usageRoot+usageShort)
	fs.StringVar(&docPkgs, "doc", defaultDoc, usageDoc)
	fs.StringVar(&docPkgs, "d", defaultDoc, usageDoc+usageShort)
	fs.BoolVar(&noLinks, "nolinks", defaultNoLinks, usageNoLinks)
	fs.BoolVar(&noLinks, "l", defaultNoLinks, usageNoLinks+usageShort)
	fs.BoolVar(&doStats, "stats", defaultStats, usageStats)
	fs.BoolVar(&doStats, "s", defaultStats, usageStats+usageShort)
	fs.BoolVar(&dirTree, "dirtree", defaultDirTree, usageDirTree)
	fs.BoolVar(&dirTree, "t", defaultDirTree, usageDirTree+usageShort)
	fs.BoolVar(&noErr, "noerror", defaultNoErr, usageNoErr)
	fs.BoolVar(&noErr, "e", defaultNoErr, usageNoErr+usageShort)
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
	log.Printf("INFO - write statistics: %t", doStats)
	log.Printf("INFO - no errors are reported: %t", noErr)

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
	if !noErr {
		if len(errs) > 0 {
			for _, err = range errs {
				log.Printf("ERROR - %v", err)
			}
			retCode = 1
		} else {
			log.Print("INFO - No errors found.")
		}
	}

	if doStats {
		writeStatistics(root, rootPkg, depMap)
	} else {
		log.Print("INFO - No statistics wanted.")
	}

	if docPkgs != "" {
		writeDocumentation(docPkgs, root, rootPkg, noLinks, depMap)
	} else {
		log.Print("INFO - No documentation wanted.")
	}

	if dirTree {
		writeDirTree(".", ".")
	}

	return retCode
}

func writeStatistics(root, rootPkg string, depMap data.DependencyMap) {
	log.Print("INFO - Writing statistics.")

	statMD := stat.Generate(depMap)
	if statMD == "" {
		return
	}
	log.Printf("INFO - Writing package statistics to file: %s", stat.FileName)
	statFile := filepath.Join(root, stat.FileName)
	err := ioutil.WriteFile(statFile, []byte(statMD), 0644)
	if err != nil {
		log.Printf("ERROR - Unable to write package statistics to file %s: %v", statFile, err)
	}
}

func writeDocumentation(docPkgs, root, rootPkg string, noLinks bool, depMap data.DependencyMap) {
	log.Print("INFO - Writing documentation.")
	dtPkgs := findPackagesWithFileAsSlice(doc.FileName, docPkgs, root, "documentation")

	linkDocPkgs := map[string]struct{}{}
	if !noLinks {
		linkDocPkgs = dirs.FindPkgsWithFile(doc.FileName, dtPkgs, root, true)
		for _, p := range dtPkgs {
			linkDocPkgs[p] = struct{}{}
		}
	}
	doc.WriteDocs(dtPkgs, depMap, linkDocPkgs, rootPkg, root)
}

var mapEntry = struct{}{}

func writeDirTree(root, name string) error {
	treeFile := filepath.Join(root, dirs.TreeFile)
	log.Printf("INFO - Writing directory tree to file: %s", treeFile)
	tree, err := dirs.Tree(root, name, []string{"vendor", "testdata", ".*"})
	if err != nil {
		log.Print("ERROR - Unable to generate directory tree")
		return err
	}
	err = ioutil.WriteFile(treeFile, []byte(tree), 0644)
	if err != nil {
		log.Printf("ERROR - Unable to write directory tree to file %s: %v", treeFile, err)
		return err
	}
	return nil
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
