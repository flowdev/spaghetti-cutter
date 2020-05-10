package apitst_test

import (
	"testing"

	"github.com/flowdev/spaghetti-cutter/parse/testdata/happy-path/apitst"
)

func TestApitst(t *testing.T) {
	t.Log("Executing TestApitst")
	apitst.Apitst()
}
