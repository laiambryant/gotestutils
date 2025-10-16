package attributes

import (
	"reflect"
	"testing"

	ctesting "github.com/laiambryant/gotestutils/ctesting"
)

func TestBoolAttributes(t *testing.T) {
	var suite []ctesting.CharacterizationTest[bool]

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := BoolAttributes{ForceTrue: true}
		got := attr.GetAttributes()
		expected := BoolAttributes{ForceTrue: true}
		return reflect.DeepEqual(got, expected), nil
	}))
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := BoolAttributes{}
		got := attr.GetReflectType()
		expected := reflect.TypeOf(true)
		return got == expected, nil
	}))
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := BoolAttributes{}
		got := attr.GetDefaultImplementation()
		return got != nil && reflect.TypeOf(got) == reflect.TypeOf(attr), nil
	}))
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := BoolAttributes{}
		got := attr.GetRandomValue()
		return got != nil, nil
	}))
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := BoolAttributes{ForceTrue: true}
		result := attr.GetRandomValue()
		if b, ok := result.(bool); ok {
			return b, nil
		}
		return false, nil
	}))
	suite = append(suite, ctesting.NewCharacterizationTest(false, nil, func() (bool, error) {
		attr := BoolAttributes{ForceFalse: true}
		result := attr.GetRandomValue()
		if b, ok := result.(bool); ok {
			return b, nil
		}
		return false, nil
	}))
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
