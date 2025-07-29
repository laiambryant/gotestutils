package examples

import (
	"fmt"
	"testing"

	"github.com/laiambryant/gotestutils/ctesting"
)

type UserProfile struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// TestUserProfileCreation demonstrates testing complex data structures
func TestUserProfileCreation(t *testing.T) {
	expectedProfile := UserProfile{ID: 1, Name: "John Doe", Age: 30}

	testSuite := []ctesting.CharacterizationTest[UserProfile]{
		ctesting.NewCharacterizationTest(expectedProfile, nil, func() (UserProfile, error) {
			return createUserProfile("John Doe", 30)
		}),

		// Test with edge case - empty name should return error
		ctesting.NewCharacterizationTest(UserProfile{},
			fmt.Errorf("name cannot be empty"), func() (UserProfile, error) {
				return createUserProfile("", 30)
			}),
	}

	ctesting.VerifyCharacterizationTestsAndResults(t, testSuite, true)
}

// Helper function for creating user profiles
func createUserProfile(name string, age int) (UserProfile, error) {
	if name == "" {
		return UserProfile{}, fmt.Errorf("name cannot be empty")
	}
	return UserProfile{
		ID:   1,
		Name: name,
		Age:  age,
	}, nil
}
