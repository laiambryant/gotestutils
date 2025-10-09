package examples

import (
	"testing"

	"github.com/laiambryant/gotestutils/stesting"
)

// TestMemoryAllocationStress demonstrates memory allocation stress testing
func TestMemoryAllocationStress(t *testing.T) {
	memoryIntensiveFunc := func() (int, error) {
		data := make([]int, 100000)
		for i := range data {
			data[i] = i * 2
		}
		return len(data), nil
	}
	var dummyVar int
	stressTest := stesting.NewStressTest[int, int](100, memoryIntensiveFunc, &dummyVar)
	success, err := stesting.RunStressTestWithFilePathOut(&stressTest, "memory_test.log")
	if !success {
		t.Errorf("Memory stress test failed: %v", err)
	}
	success, err = stesting.RunParallelStressTest(&stressTest, 4)
	if !success {
		t.Errorf("Parallel memory stress test failed: %v", err)
	}
}
