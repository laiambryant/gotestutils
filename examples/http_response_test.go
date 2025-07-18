package examples

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/laiambryant/gotestutils/ctesting"
)

type APIResponse struct {
	StatusCode int
	Body       string
}

// TestAPIResponseCharacterization demonstrates testing HTTP response handling
func TestAPIResponseCharacterization(t *testing.T) {
	// Mock HTTP server setup
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	expectedResponse := APIResponse{
		StatusCode: 200,
		Body:       `{"status": "ok"}`,
	}

	testSuite := []ctesting.CharacterizationTest[APIResponse]{
		ctesting.NewCharacterizationTest(expectedResponse, nil, func() (APIResponse, error) {
			resp, err := http.Get(server.URL)
			if err != nil {
				return APIResponse{}, err
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return APIResponse{}, err
			}

			return APIResponse{
				StatusCode: resp.StatusCode,
				Body:       string(body),
			}, nil
		}),
	}

	ctesting.VerifyCharacterizationTestsAndResults(t, testSuite, true)
}
