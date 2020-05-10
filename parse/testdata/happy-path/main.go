package main

import (
	"log"
	"os"

	"github.com/flowdev/spaghetti-cutter/parse/testdata/happy-path/alltst"
	"github.com/flowdev/spaghetti-cutter/parse/testdata/happy-path/apitst"
)

func main() {
	doIt(os.Args[1:])
}

func doIt(args []string) {
	log.Printf("INFO - this is the main package, args: %q", args)
	apitst.Apitst()
	alltst.Alltst()
}
