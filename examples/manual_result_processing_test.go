package examples

import (
	"log"
	"testing"

	"github.com/laiambryant/gotestutils/ctesting"
)

// TestManualResultProcessing demonstrates manual processing of test results
func TestManualResultProcessing(t *testing.T) {
	testSuite := []ctesting.CharacterizationTest[int]{
		ctesting.NewCharacterizationTest(3, nil, func() (int, error) {
			return Sum(1, 2), nil
		}),
		ctesting.NewCharacterizationTest(5, nil, func() (int, error) {
			return Sum(2, 2), nil // This will fail: 2+2=4, expected 5
		}),
	}

	results, testSuiteRes := ctesting.VerifyCharacterizationTests(testSuite, true)

	for i, passed := range results {
		test := testSuiteRes[i]
		if !passed {
			// Custom failure handling - Note: in real code you'd need access to private fields
			// This is just for demonstration purposes
			log.Printf("Test %d failed: expected %v, got actual value from test execution",
				i+1, test.ExpectedOutput)

			// Custom error analysis
			if test.ExpectedErr != nil {
				log.Printf("Error analysis for test %d: expected error %v", i+1, test.ExpectedErr)
			}
		}
	}
}
