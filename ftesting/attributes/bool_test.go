package attributes

import (
	"reflect"
	"testing"

	ctesting "github.com/laiambryant/gotestutils/ctesting"
)

// Test BoolAttributes (moved helper and simple tests)
func TestBoolAttributes_SimpleSuite(t *testing.T) {
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

	results, _ := ctesting.VerifyCharacterizationTestsAndResults(t, suite, true)
	for i, passed := range results {
		if !passed {
			t.Fatalf("BoolAttributes test %d failed", i+1)
		}
	}
}

func TestBoolAttributes_ForceTrue(t *testing.T) {
	attr := BoolAttributes{ForceTrue: true}
	result := attr.GetRandomValue()
	if result != true {
		t.Errorf("Expected true when ForceTrue is set, got %v", result)
	}
}

func TestBoolAttributes_ForceFalse(t *testing.T) {
	attr := BoolAttributes{ForceFalse: true}
	result := attr.GetRandomValue()
	if result != false {
		t.Errorf("Expected false when ForceFalse is set, got %v", result)
	}
}

func TestBoolAttributes_BothForced(t *testing.T) {
	attr := BoolAttributes{ForceTrue: true, ForceFalse: true}
	result := attr.GetRandomValue()
	if result != true {
		t.Errorf("Expected true when both forced (ForceTrue precedence), got %v", result)
	}
}
