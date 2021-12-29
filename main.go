package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/flowdev/spaghetti-cutter/deps"
	"github.com/flowdev/spaghetti-cutter/parse"
	"github.com/flowdev/spaghetti-cutter/size"
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
		usageShort   = " (shorthand)"
		defaultRoot  = "."
		usageRoot    = "root directory of the project"
		defaultNoErr = false
		usageNoErr   = "don't report errors or exit with an error"
	)
	var startDir string
	var noErr bool
	fs := flag.NewFlagSet("spaghetti-cutter", flag.ExitOnError)
	fs.StringVar(&startDir, "root", defaultRoot, usageRoot)
	fs.StringVar(&startDir, "r", defaultRoot, usageRoot+usageShort)
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
	log.Printf("INFO - no errors are reported: %t", noErr)

	packs, err := parse.DirTree(root)
	if err != nil {
		log.Printf("FATAL - %v", err)
		return 6
	}

	rootPkg := parse.RootPkg(packs)
	log.Printf("INFO - root package: %s", rootPkg)
	pkgInfos := pkgs.UniquePackages(packs)

	var errs []error
	for _, pkgInfo := range pkgInfos {
		errs = addErrors(errs, deps.Check(pkgInfo.Pkg, rootPkg, cfg))
		errs = addErrors(errs, size.Check(pkgInfo.Pkg, rootPkg, cfg.Size))
	}

	retCode := 0
	if len(errs) > 0 {
		for _, err = range errs {
			log.Printf("ERROR - %v", err)
		}
		if !noErr {
			retCode = 1
		}
	} else {
		log.Print("INFO - No errors found.")
	}

	return retCode
}

func addErrors(errs []error, newErrs []error) []error {
	return append(errs, newErrs...)
}
