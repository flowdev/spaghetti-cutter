package main

import (
	"log"
	"os"

	"github.com/flowdev/spaghetti-cutter/deps/testdata/only-tools/x/tool"
	"github.com/flowdev/spaghetti-cutter/deps/testdata/only-tools/x/tool2"
)

func main() {
	doIt(os.Args[1:])
}

func doIt(args []string) {
	log.Printf("INFO - this is the main package, args: %q", args)
	tool.Tool()
	tool2.Tool2()
}
