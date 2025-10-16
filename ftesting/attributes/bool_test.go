package attributes

import (
	"reflect"
	"testing"

	ctesting "github.com/laiambryant/gotestutils/ctesting"
)

func TestBoolAttributes(t *testing.T) {
	var suite []ctesting.CharacterizationTest[bool]

	// GetAttributes tests
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := BoolAttributes{ForceTrue: true}
		got := attr.GetAttributes()
		expected := BoolAttributes{ForceTrue: true}
		return reflect.DeepEqual(got, expected), nil
	}))

	// GetReflectType tests
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := BoolAttributes{}
		got := attr.GetReflectType()
		expected := reflect.TypeOf(true)
		return got == expected, nil
	}))

	// GetDefaultImplementation tests
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := BoolAttributes{}
		got := attr.GetDefaultImplementation()
		return got != nil && reflect.TypeOf(got) == reflect.TypeOf(attr), nil
	}))

	// GetRandomValue tests - basic functionality
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := BoolAttributes{}
		got := attr.GetRandomValue()
		return got != nil, nil
	}))

	// ForceTrue tests
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := BoolAttributes{ForceTrue: true}
		result := attr.GetRandomValue()
		if b, ok := result.(bool); ok {
			return b, nil
		}
		return false, nil
	}))

	// ForceFalse tests
	suite = append(suite, ctesting.NewCharacterizationTest(false, nil, func() (bool, error) {
		attr := BoolAttributes{ForceFalse: true}
		result := attr.GetRandomValue()
		if b, ok := result.(bool); ok {
			return b, nil
		}
		return false, nil
	}))

	// BothForced tests
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := BoolAttributes{ForceTrue: true, ForceFalse: true}
		result := attr.GetRandomValue()
		if b, ok := result.(bool); ok {
			return b, nil
		}
		return false, nil
	}))

	results, _ := ctesting.VerifyCharacterizationTestsAndResults(t, suite, true)
	for i, passed := range results {
		if !passed {
			t.Fatalf("BoolAttributes test %d failed", i+1)
		}
	}
}
