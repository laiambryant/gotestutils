package examples

import (
	"testing"

	"github.com/laiambryant/gotestutils/stesting"
)

// TestMemoryAllocationStress demonstrates memory allocation stress testing
func TestMemoryAllocationStress(t *testing.T) {
	memoryIntensiveFunc := func() (int, error) {
		// Allocate and process large slices
		data := make([]int, 100000)
		for i := range data {
			data[i] = i * 2
		}
		return len(data), nil
	}

	stressTest := stesting.NewStressTest[int, any](100, memoryIntensiveFunc, nil)

	// Save results to analyze memory patterns
	success, err := stesting.RunStressTestWithFilePathOut(&stressTest, "/tmp/memory_test.log")
	if !success {
		t.Errorf("Memory stress test failed: %v", err)
	}

	// Also run in parallel to test concurrent memory allocation
	success, err = stesting.RunParallelStressTest(&stressTest, 4)
	if !success {
		t.Errorf("Parallel memory stress test failed: %v", err)
	}
}
