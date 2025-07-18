package examples

import (
	"testing"

	"github.com/laiambryant/gotestutils/stesting"
)

// TestStressErrorHandling demonstrates error handling in stress tests
func TestStressErrorHandling(t *testing.T) {
	// Function that will sometimes error
	errorProneFunc := func() (bool, error) {
		// This will fail and demonstrate error handling
		_, err := Divide(10, 0)
		return err == nil, err
	}

	stressTest := stesting.NewStressTest[bool, any](10, errorProneFunc, nil)

	success, err := stesting.RunStressTest(&stressTest)
	if !success {
		if ste, ok := err.(stesting.StressTestingError); ok {
			t.Logf("Test failed at iteration %d: %v", ste.Index, ste.Err)
		}
		// This is expected to fail, so we don't use t.Errorf here
	}
}
