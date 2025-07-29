package examples

import (
	"fmt"

	"github.com/laiambryant/gotestutils/ctesting"
)

// BasicCharacterizationExample demonstrates basic characterization testing usage
func BasicCharacterizationExample() {
	// Test expecting successful execution with specific output
	_ = ctesting.NewCharacterizationTest(42, nil, func() (int, error) {
		return CalculateResult(6, 7), nil
	})

	// Test expecting an error condition
	_ = ctesting.NewCharacterizationTest(0, fmt.Errorf("division by zero"), func() (int, error) {
		return Divide(10, 0)
	})

}
