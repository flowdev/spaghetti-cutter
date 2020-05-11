package main

import (
	"log"
	"os"

	"github.com/flowdev/spaghetti-cutter/testdata/good-proj/pkg/db/store"
	"github.com/flowdev/spaghetti-cutter/testdata/good-proj/pkg/domain3"
	"github.com/flowdev/spaghetti-cutter/testdata/good-proj/pkg/domain4"
)

func main() {
	doIt(os.Args[1:])
}

func doIt(args []string) {
	log.Printf("INFO - this is the main package, args: %q", args)
	s := &store.Store{}
	domain3.HandleDomain3Route1(s)
	domain3.HandleDomain3Route2(s)

	domain4.HandleDomain4Route1(s)
	domain4.HandleDomain4Route2(s)
}
