package tool2

import (
	"log"

	"github.com/flowdev/spaghetti-cutter/deps/testdata/only-tools/x/tool"
)

// Tool2 is logging its execution.
func Tool2() {
	log.Printf("INFO - tool.Tool")
	tool.Tool() // evil dependency!
}
