package alltst_test

import (
	"testing"

	"github.com/flowdev/spaghetti-cutter/parse/testdata/happy-path/alltst"
)

func TestAlltst(t *testing.T) {
	t.Log("Executing TestAlltst")
	alltst.Alltst()
}
