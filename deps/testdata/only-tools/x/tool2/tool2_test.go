package tool2_test

import (
	"testing"

	"github.com/flowdev/spaghetti-cutter/deps/testdata/only-tools/x/tool2"
)

func TestTool(t *testing.T) {
	t.Log("Executing TestTool")
	tool2.Tool2()
}
