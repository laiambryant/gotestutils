package examples

import (
	"fmt"
	"testing"

	"github.com/laiambryant/gotestutils/ctesting"
)

// TestErrorHandlingCharacterization demonstrates testing error conditions
func TestErrorHandlingCharacterization(t *testing.T) {
	testSuite := []ctesting.CharacterizationTest[int]{
		// Test division by zero
		ctesting.NewCharacterizationTest(0,
			fmt.Errorf("division by zero"), func() (int, error) {
				return Divide(10, 0)
			}),

		// Test negative input validation
		ctesting.NewCharacterizationTest(0,
			fmt.Errorf("negative numbers not allowed"), func() (int, error) {
				return ProcessPositiveNumber(-5)
			}),

		// Test successful operation
		ctesting.NewCharacterizationTest(5, nil, func() (int, error) {
			return Divide(25, 5)
		}),
	}

	ctesting.VerifyCharacterizationTestsAndResults(t, testSuite, true)
}
