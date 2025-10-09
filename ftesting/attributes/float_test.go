package attributes

import (
	"reflect"
	"testing"

	ctesting "github.com/laiambryant/gotestutils/ctesting"
)

func TestFloatAttributes(t *testing.T) {
	var suite []ctesting.CharacterizationTest[bool]

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := FloatAttributesImpl[float64]{Min: 1.1, Max: 2.2, NonZero: true, FiniteOnly: true, AllowNaN: true, AllowInf: true, Precision: 3}
		got := attr.GetAttributes()
		expected := FloatAttributesImpl[float64]{Min: 1.1, Max: 2.2, NonZero: true, FiniteOnly: true, AllowNaN: true, AllowInf: true, Precision: 3}
		return reflect.DeepEqual(got, expected), nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := FloatAttributesImpl[float64]{}
		got := attr.GetReflectType()
		expected := reflect.TypeOf(float64(0))
		return got == expected, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := FloatAttributesImpl[float64]{}
		got := attr.GetDefaultImplementation()
		return got != nil && reflect.TypeOf(got) == reflect.TypeOf(attr), nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := FloatAttributesImpl[float64]{Min: -1.0, Max: 1.0}
		got := attr.GetRandomValue()
		return got != nil, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := FloatAttributesImpl[float64]{Max: 1.0, Min: 2.0}
		result := attr.GetRandomValue()
		return result == float64(0), nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := FloatAttributesImpl[float64]{Max: 10.0, Min: 1.0}
		result := attr.GetRandomValue()
		return result != float64(0), nil
	}))

	results, _ := ctesting.VerifyCharacterizationTestsAndResults(t, suite, true)
	for i, passed := range results {
		if !passed {
			t.Fatalf("FloatAttributes test %d failed", i+1)
		}
	}
}

func TestFloatAttributes_InvalidRange(t *testing.T) {
	attr := FloatAttributesImpl[float64]{Max: 1.0, Min: 2.0}
	result := attr.GetRandomValue()
	if result != float64(0) {
		t.Errorf("Expected zero value for invalid range, got %v", result)
	}
}

func TestFloatAttributes_ValidRange(t *testing.T) {
	attr := FloatAttributesImpl[float64]{Max: 10.0, Min: 1.0}
	result := attr.GetRandomValue()
	if result == float64(0) {
		t.Error("Expected non-zero value for valid range")
	}
}
