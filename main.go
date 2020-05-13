package main

import (
	"log"
	"os"

	"github.com/flowdev/spaghetti-cutter/deps"
	"github.com/flowdev/spaghetti-cutter/parse"
	"github.com/flowdev/spaghetti-cutter/size"
	"github.com/flowdev/spaghetti-cutter/x/config"
	"github.com/flowdev/spaghetti-cutter/x/dirs"
)

func main() {
	rc := cut(os.Args[1:])
	if rc != 0 {
		os.Exit(rc)
	}
}

func cut(args []string) int {
	var err error

	cfg := config.Parse(args, dirs.FindConfig(config.File))

	cfg.Root, err = dirs.FindRoot(cfg.Root, config.File, cfg.IgnoreVendor)
	if err != nil {
		log.Printf("FATAL - %v", err)
		return 2
	}
	log.Printf("INFO - configuration God: %s", cfg.God)
	log.Printf("INFO - configuration Tool: %s", cfg.Tool)
	log.Printf("INFO - configuration DB: %s", cfg.DB)
	log.Printf("INFO - configuration Allow: %s", cfg.Allow)
	log.Printf("INFO - configuration Size: %d", cfg.Size)
	log.Printf("INFO - configuration Root: %s", cfg.Root)
	log.Printf("INFO - configuration NoGod: %t", cfg.NoGod)
	log.Printf("INFO - configuration IgnoreVendor: %t", cfg.IgnoreVendor)

	pkgs, err := parse.DirTree(cfg.Root)
	if err != nil {
		log.Printf("FATAL - %v", err)
		return 3
	}

	var errs []error
	rootPkg := parse.RootPkg(pkgs)
	log.Printf("INFO - root package: %s", rootPkg)
	for _, pkg := range pkgs {
		errs = addErrors(errs, deps.Check(pkg, rootPkg, cfg))
		errs = addErrors(errs, size.Check(pkg, rootPkg, cfg.Size))
	}

	if len(errs) > 0 {
		for _, err = range errs {
			log.Printf("ERROR - %v", err)
		}
		return 1
	}

	log.Print("INFO - No errors found.")
	return 0
}

func addErrors(errs []error, newErrs []error) []error {
	return append(errs, newErrs...)
}
