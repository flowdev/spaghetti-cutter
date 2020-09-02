package tool

import (
	"log"

	"github.com/flowdev/spaghetti-cutter/deps/testdata/half-pkgs-proj/x/tool/subtool"
)

// Tool is logging its execution.
func Tool() {
	log.Printf("INFO - tool.Tool")
	subtool.Subtool()
}
