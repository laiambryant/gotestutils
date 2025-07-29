package examples

import (
	"fmt"
	"testing"

	"github.com/laiambryant/gotestutils/ctesting"
)

// TestBasicExecutionExample demonstrates basic test execution
func TestBasicExecutionExample(t *testing.T) {
	testSuite := []ctesting.CharacterizationTest[int]{
		ctesting.NewCharacterizationTest(3, nil, func() (int, error) {
			return Sum(1, 2), nil
		}),
		ctesting.NewCharacterizationTest(10, nil, func() (int, error) {
			return Multiply(2, 5), nil
		}),
	}

	// Execute tests with deep error checking
	results, updatedSuite := ctesting.VerifyCharacterizationTests(testSuite, true)

	// Process results manually
	for i, passed := range results {
		if passed {
			fmt.Printf("Test %d: PASSED\n", i+1)
		} else {
			fmt.Printf("Test %d: FAILED\n", i+1)
		}
	}

	_ = updatedSuite // Use the variable to avoid unused variable error
}

// TestIntegratedExecutionExample demonstrates integrated execution and reporting
func TestIntegratedExecutionExample(t *testing.T) {
	testSuite := []ctesting.CharacterizationTest[int]{
		ctesting.NewCharacterizationTest(3, nil, func() (int, error) {
			return Sum(1, 2), nil
		}),
	}

	// Automatically execute tests and report results to testing.T
	ctesting.VerifyCharacterizationTestsAndResults(t, testSuite, true)
}

// TestErrorCheckingExample demonstrates different error checking modes
func TestErrorCheckingExample(t *testing.T) {

	testSuite := []ctesting.CharacterizationTest[int]{
		ctesting.NewCharacterizationTest(3, nil, func() (int, error) {
			return Sum(1, 2), nil
		}),
	}

	// Deep error checking (exact error message comparison)
	results, _ := ctesting.VerifyCharacterizationTests(testSuite, true)

	// Shallow error checking (uses errors.Is() for error comparison)
	_, _ = ctesting.VerifyCharacterizationTests(testSuite, false)

	_ = results

}
