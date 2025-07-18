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
	stressTest := stesting.NewStressTest[int, any](
		1000, // 1,000 iterations
		func() (int, error) { // Function to test
			return Sum(10, 20), nil
		},
		nil,
	)

	maxWorkers := uint32(8) // Use 8 concurrent workers
	success, err := stesting.RunParallelStressTest(&stressTest, maxWorkers)
	if !success {
		t.Errorf("Parallel stress test failed: %v", err)
	}
}
