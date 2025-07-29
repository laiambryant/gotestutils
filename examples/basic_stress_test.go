package examples

import (
	"testing"

	"github.com/laiambryant/gotestutils/stesting"
)

// TestBasicStressExample demonstrates basic stress testing usage
func TestBasicStressExample(t *testing.T) {
	// Test a function that returns an integer
	stressTest := stesting.NewStressTest[int, any](
		10000, // 10,000 iterations
		func() (int, error) { // Function to test
			return Fibonacci(20), nil
		},
		nil, // No test variables needed
	)

	// Sequential execution
	success, err := stesting.RunStressTest(&stressTest)
	if !success {
		t.Errorf("Stress test failed: %v", err)
	}
}

// TestParallelStressExample demonstrates parallel stress testing
func TestParallelStressExample(t *testing.T) {

	iterations := uint32(1000)

	stressTest := stesting.NewStressTest[int, any](
		iterations,
		testFunction,
		nil,
	)

	maxWorkers := uint32(8)
	success, err := stesting.RunParallelStressTest(&stressTest, maxWorkers)
	if !success {
		t.Errorf("Parallel stress test failed: %v", err)
	}
}

func testFunction() (int, error) {
	return Fibonacci(20), nil
}
