package examples

import (
	"testing"

	"github.com/laiambryant/gotestutils/stesting"
)

// TestStressErrorHandling demonstrates error handling in stress tests
func TestStressErrorHandling(t *testing.T) {

	stressTest := stesting.NewStressTest[bool, any](10, errorProneFunc, nil)

	success, err := stesting.RunStressTest(&stressTest)
	if !success {
		if ste, ok := err.(stesting.StressTestingError); ok {
			t.Logf("Test failed at iteration %d: %v", ste.Index, ste.Err)
		}
	}
}

func errorProneFunc() (bool, error) {
	_, err := Divide(10, 0)
	return err == nil, err
}
