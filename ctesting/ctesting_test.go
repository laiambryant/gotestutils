package ctesting

import (
	"fmt"
	"testing"
)

const testErrorMessage = "An error"

func sum(a int, b int) (ret int) {
	return a + b
}

func getError() (int, error) {
	return 1, fmt.Errorf("%s", testErrorMessage)
}

// Tests if output is correct
func TestUtilsExample(t *testing.T) {
	testSuite := []CharacterizationTest[int]{
		NewCharacterizationTest(3, nil, func() (int, error) { return sum(1, 2), nil }),
	}
	results, testSuiteRes := VerifyCharacterizationTests(testSuite, true)
	VerifyResults(t, results, testSuiteRes)
}

// Tests if error is correct
func TestCharacterizationTestOnError(t *testing.T) {
	testSuite := []CharacterizationTest[int]{
		NewCharacterizationTest(1, fmt.Errorf("%s", testErrorMessage), func() (int, error) { return getError() }),
	}
	results, testSuiteRes := VerifyCharacterizationTests(testSuite, true)
	VerifyResults(t, results, testSuiteRes)
}

// Tests if false is put in the response slice should the two outputs differ
func TestCharacterizationTestFailure(t *testing.T) {
	testSuite := []CharacterizationTest[int]{
		NewCharacterizationTest(4, nil, func() (int, error) { return sum(1, 2), nil }),
	}
	VerifyCharacterizationTests(testSuite, true)
}

// Test only for coverage. Covers the case in which a test fails using a mock testing.T
func TestErrorWithMockT(t *testing.T) {
	mockT := testing.T{}
	testSuite := []CharacterizationTest[int]{
		NewCharacterizationTest(4, nil, func() (int, error) { return sum(1, 2), nil }),
	}
	results, testSuiteRes := VerifyCharacterizationTests(testSuite, true)
	VerifyResults(&mockT, results, testSuiteRes)
}

// TestVerifyCharacterizationTestsAndResults tests the convenience function that combines
// test execution and result reporting in a single call.
func TestVerifyCharacterizationTestsAndResults(t *testing.T) {
	testSuite := []CharacterizationTest[int]{
		NewCharacterizationTest(3, nil, func() (int, error) { return sum(1, 2), nil }),
		NewCharacterizationTest(1, fmt.Errorf("%s", testErrorMessage), func() (int, error) { return getError() }),
	}

	VerifyCharacterizationTestsAndResults(t, testSuite, true)
}

// Tests to cover both cases should you not desire to do a deepEquals in verifying the characterization test
func TestVerifyCharacterizationTestsErrors(t *testing.T) {
	testSuite := []CharacterizationTest[int]{
		NewCharacterizationTest(4, nil, func() (int, error) { return sum(1, 2), nil }),
		NewCharacterizationTest(3, nil, func() (int, error) { return sum(1, 2), nil }),
	}
	results, _ := VerifyCharacterizationTests(testSuite, false)
	if len(results) < 2 || results[0] || !results[1] {
		t.Error("The results are incorrect")
	}
}
