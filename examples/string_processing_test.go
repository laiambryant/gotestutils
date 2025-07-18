package examples

import (
	"strings"
	"testing"

	"github.com/laiambryant/gotestutils/ctesting"
)

// TestStringProcessingCharacterization demonstrates testing string processing functions
func TestStringProcessingCharacterization(t *testing.T) {
	testSuite := []ctesting.CharacterizationTest[string]{
		// Test normal input
		ctesting.NewCharacterizationTest("HELLO", nil, func() (string, error) {
			return strings.ToUpper("hello"), nil
		}),

		// Test empty string
		ctesting.NewCharacterizationTest("", nil, func() (string, error) {
			return strings.ToUpper(""), nil
		}),

		// Test with special characters
		ctesting.NewCharacterizationTest("HELLO, WORLD!", nil, func() (string, error) {
			return strings.ToUpper("hello, world!"), nil
		}),
	}

	ctesting.VerifyCharacterizationTestsAndResults(t, testSuite, true)
}
