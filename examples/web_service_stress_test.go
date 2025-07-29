package examples

import (
	"net/http"
	"testing"

	"github.com/laiambryant/gotestutils/stesting"
)

// TestWebServiceStress demonstrates stress testing a web service function
func TestWebServiceStress(t *testing.T) {
	stressTest := stesting.NewStressTest[int, any](50, simulatedWebServiceCall, nil) // Reduced iterations for example

	// Test with 5 concurrent workers to simulate real load (reduced for example)
	success, err := stesting.RunParallelStressTest(&stressTest, 5)
	if !success {
		t.Errorf("Web service stress test failed: %v", err)
	}
}

func simulatedWebServiceCall() (int, error) {
	resp, err := http.Get("http://httpbin.org/status/200")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}
