package tool2_test

import (
	"testing"

	"github.com/flowdev/spaghetti-cutter/testdata/good-proj/pkg/x/tool2"
)

func TestTool(t *testing.T) {
	t.Log("Executing TestTool")
	tool2.Tool2()
}
