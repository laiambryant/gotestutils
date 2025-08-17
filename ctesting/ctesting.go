package ctesting

import (
	"errors"
	"reflect"
	"testing"

	gtu "github.com/laiambryant/gotestutils/testing"
)

// CharacterizationTest represents a single test case for characterization testing.
// Characterization tests capture the current behavior of existing code by comparing
// actual outputs and errors against expected values.
//
// Type parameter t must implement TestOutputType (be comparable).
//
// Fields:
//   - err: The actual error returned by the test function (populated during test execution)
//   - ExpectedErr: The expected error that should be returned
//   - output: The actual output returned by the test function (populated during test execution)
//   - ExpectedOutput: The expected output value
//   - F: The test function to execute
//
// Example usage, this test expects sum(1,2) to return 3 with no error:
//
//	test := NewCharacterizationTest(3, nil, func() (int, error) { return sum(1, 2), nil })
type CharacterizationTest[t comparable] struct {
	err            error
	ExpectedErr    error
	output         t
	ExpectedOutput t
	F              gtu.TestFunc[t]
}

// NewCharacterizationTest creates a new CharacterizationTest instance with the specified
// expected output, expected error, and test function to execute.
//
// Parameters:
//   - expectedOutput: The value that the test function is expected to return
//   - expectedError: The error that the test function is expected to return (use nil if no error expected)
//   - function: The test function to execute, must match TestFunc[t] signature
//
// Returns a configured CharacterizationTest ready for execution.
//
// Example usage from tests:
//
//	// Test expecting successful execution with output 3
//	test := NewCharacterizationTest(3, nil, func() (int, error) { return sum(1, 2), nil })
//
//	// Test expecting an error
//	test := NewCharacterizationTest(1, fmt.Errorf("An error"), func() (int, error) { return getError() })
//
//	// Test expecting failure (wrong expected output)
//	test := NewCharacterizationTest(4, nil, func() (int, error) { return sum(1, 2), nil })
func NewCharacterizationTest[t comparable](expectedOutput t, expectedError error, function gtu.TestFunc[t]) (test CharacterizationTest[t]) {
	return CharacterizationTest[t]{
		ExpectedErr:    expectedError,
		ExpectedOutput: expectedOutput,
		F:              function,
	}
}

// VerifyCharacterizationTests executes a suite of characterization tests and returns
// the results of each test along with the updated test suite containing actual outputs.
//
// This function:
// 1. Executes each test function in the suite
// 2. Compares actual outputs/errors with expected values
// 3. Returns a boolean slice indicating pass/fail for each test
// 4. Updates the test suite with actual outputs and errors
//
// Parameters:
//   - testSuite: A slice of CharacterizationTest instances to execute
//
// Returns:
//   - []bool: A slice where each element indicates if the corresponding test passed (true) or failed (false)
//   - []CharacterizationTest[t]: The updated test suite with actual outputs and errors populated
//
// A test passes if:
//   - Both expected and actual errors are non-nil and have the same error message, OR
//   - The expected output exactly matches the actual output (using reflect.DeepEqual)
//
// Example usage from tests, results[0] will be true if sum(1,2) returns 3 with no error:
//
//	testSuite := []CharacterizationTest[int]{
//	    NewCharacterizationTest(3, nil, func() (int, error) { return sum(1, 2), nil }),
//	}
//	results, testSuiteRes := VerifyCharacterizationTests(testSuite)
func VerifyCharacterizationTests[t comparable](
	testSuite []CharacterizationTest[t], isDeepErrorCheck bool) (res []bool, _ []CharacterizationTest[t]) {
	for i, test := range testSuite {
		output, err := test.F()
		testSuite[i].output = output
		testSuite[i].err = err
		if isDeepErrorCheck {
			res = append(res, deepErrorCheck(err, test, output))
		} else {
			res = append(res, shallowErrorCheck(err, test, output))
		}
	}
	return res, testSuite
}

func deepErrorCheck[t comparable](err error, test CharacterizationTest[t], output t) (res bool) {
	if (err != nil && test.ExpectedErr != nil &&
		test.ExpectedErr.Error() == err.Error()) ||
		reflect.DeepEqual(test.ExpectedOutput, output) {
		return true
	} else {
		return false
	}
}

func shallowErrorCheck[t comparable](err error, test CharacterizationTest[t], output t) (res bool) {
	if ((err != nil && test.ExpectedErr != nil) && (errors.Is(err, test.ExpectedErr) || err.Error() == test.ExpectedErr.Error())) ||
		reflect.DeepEqual(test.ExpectedOutput, output) {
		return true
	} else {
		return false
	}
}

// VerifyResults processes the results from VerifyCharacterizationTests and reports
// test outcomes using the provided testing.T instance. For failed tests, it logs
// detailed error information including expected vs actual values and errors.
// For successful tests, it logs success information.
//
// This function should be called after VerifyCharacterizationTests to properly
// report test results in the testing framework.
//
// Parameters:
//   - t: A testing.T instance used for logging results and reporting failures
//   - results: Boolean slice from VerifyCharacterizationTests indicating pass/fail status
//   - testSuiteRes: Updated test suite from VerifyCharacterizationTests with actual outputs
//
// Behavior:
//   - For failed tests (results[i] == false): Calls t.Errorf with detailed comparison
//   - For successful tests (results[i] == true): Calls t.Logf with success information
//
// Example usage from tests:
//
//	// Standard usage with real testing.T
//	results, testSuiteRes := VerifyCharacterizationTests(testSuite)
//	VerifyResults(t, results, testSuiteRes)
//
//	// Usage with mock testing.T for testing failure scenarios
//	mockT := testing.T{}
//	VerifyResults(&mockT, results, testSuiteRes)
func VerifyResults[T comparable](t *testing.T, results []bool, testSuiteRes []CharacterizationTest[T]) {
	for i, result := range results {
		if !result {
			t.Errorf("test number %d: ERROR [ERRORS] got error {%v}, expected {%v}, [VALUES] got {%v} expected {%v}",
				i+1, testSuiteRes[i].err, testSuiteRes[i].ExpectedErr, testSuiteRes[i].output, testSuiteRes[i].ExpectedOutput)
		} else {
			t.Logf("test number %d: SUCCESS [ERRORS] got error {%v}, expected {%v}, [VALUES] got {%v} expected {%v}",
				i+1, testSuiteRes[i].err, testSuiteRes[i].ExpectedErr, testSuiteRes[i].output, testSuiteRes[i].ExpectedOutput)
		}
	}
}

// VerifyCharacterizationTestsAndResults is a convenience function that combines
// VerifyCharacterizationTests and VerifyResults into a single call. This function
// executes the test suite and immediately reports the results using the provided
// testing.T instance.
//
// This function simplifies the common pattern of running characterization tests
// and immediately reporting their results, reducing boilerplate code in test functions.
//
// Parameters:
//   - t: A testing.T instance used for logging results and reporting failures
//   - testSuite: A slice of CharacterizationTest instances to execute
//
// Returns:
//   - []bool: A slice where each element indicates if the corresponding test passed (true) or failed (false)
//   - []CharacterizationTest[T]: The updated test suite with actual outputs and errors populated
//
// This function is equivalent to calling:
//
//	results, testSuiteRes := VerifyCharacterizationTests(testSuite)
//	VerifyResults(t, results, testSuiteRes)
//	return results, testSuiteRes
//
// Example usage from tests, test results are automatically reported to t, and you still get the return values if needed:
//
//	testSuite := []CharacterizationTest[int]{
//	    NewCharacterizationTest(3, nil, func() (int, error) { return sum(1, 2), nil }),
//	}
//	results, testSuiteRes := VerifyCharacterizationTestsAndResults(t, testSuite)
func VerifyCharacterizationTestsAndResults[T comparable](t *testing.T, testSuite []CharacterizationTest[T], deepErrorCheck bool) ([]bool, []CharacterizationTest[T]) {
	results, testSuiteRes := VerifyCharacterizationTests(testSuite, deepErrorCheck)
	VerifyResults(t, results, testSuiteRes)
	return results, testSuiteRes
}
