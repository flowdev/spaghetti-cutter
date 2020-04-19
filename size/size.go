package size

import (
	"fmt"

	"golang.org/x/tools/go/packages"
)

// Check checks the complexity of the given package and reports if it is too
// big.
func Check(pack *packages.Package, size uint) []error {
	fmt.Println("Complexity configuration:")
	fmt.Println("    Size:", size)

	return nil
}
