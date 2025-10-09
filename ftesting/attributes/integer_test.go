package attributes

import (
	"reflect"
	"testing"

	ctesting "github.com/laiambryant/gotestutils/ctesting"
)

func TestIntegerAttributes(t *testing.T) {
	var suite []ctesting.CharacterizationTest[bool]

	// GetAttributes tests
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := IntegerAttributesImpl[int64]{AllowNegative: true, AllowZero: true, Max: 10, Min: -5}
		got := attr.GetAttributes()
		expected := IntegerAttributesImpl[int64]{AllowNegative: true, AllowZero: true, Max: 10, Min: -5}
		return reflect.DeepEqual(got, expected), nil
	}))

	// GetReflectType tests
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := IntegerAttributesImpl[int64]{}
		got := attr.GetReflectType()
		expected := reflect.TypeOf(int64(0))
		return got == expected, nil
	}))

	// GetDefaultImplementation tests
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := IntegerAttributesImpl[int64]{}
		got := attr.GetDefaultImplementation()
		return got != nil && reflect.TypeOf(got) == reflect.TypeOf(attr), nil
	}))

	// GetRandomValue tests - basic functionality
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := IntegerAttributesImpl[int64]{Min: -10, Max: 10}
		got := attr.GetRandomValue()
		return got != nil, nil
	}))

	// Edge case: Invalid range (Max < Min)
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := IntegerAttributesImpl[int]{Max: 0, Min: -10}
		result := attr.GetRandomValue()
		return result == 0, nil
	}))

	// Edge case: Max less than Min
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := IntegerAttributesImpl[int]{Max: 5, Min: 10}
		result := attr.GetRandomValue()
		return result == 0, nil
	}))

	// Edge case: Min greater than Max exact test
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := IntegerAttributesImpl[int]{Max: 5, Min: 10, AllowNegative: true, AllowZero: true}
		result := attr.GetRandomValue()
		return result == 0, nil
	}))

	results, _ := ctesting.VerifyCharacterizationTestsAndResults(t, suite, true)
	for i, passed := range results {
		if !passed {
			t.Fatalf("IntegerAttributes test %d failed", i+1)
		}
	}
}

func TestIntegerAttributes_InvalidRange(t *testing.T) {
	attr := IntegerAttributesImpl[int]{Max: 0, Min: -10}
	result := attr.GetRandomValue()
	if result != 0 {
		t.Errorf("Expected zero value for invalid range, got %v", result)
	}
}

func TestIntegerAttributes_MaxLessThanMin(t *testing.T) {
	attr := IntegerAttributesImpl[int]{Max: 5, Min: 10}
	result := attr.GetRandomValue()
	if result != 0 {
		t.Errorf("Expected zero value when max <= min, got %v", result)
	}
}

func TestIntegerAttributes_MinGreaterThanMaxExact(t *testing.T) {
	attr := IntegerAttributesImpl[int]{Max: 5, Min: 10, AllowNegative: true, AllowZero: true}
	result := attr.GetRandomValue()
	if result != 0 {
		t.Errorf("Expected zero value when Min > Max, got %v", result)
	}
}

func TestGetAttributeGivenType_ZeroValueAttribute(t *testing.T) {
	attributes := FTAttributes{
		IntegerAttr: IntegerAttributesImpl[int]{},
	}
	result, err := attributes.GetAttributeGivenType(reflect.TypeOf(int(0)))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result == nil {
		t.Error("expected default implementation, got nil")
	}
	defaultImpl := IntegerAttributesImpl[int]{}.GetDefaultImplementation()
	if !reflect.DeepEqual(result, defaultImpl) {
		t.Errorf("expected default implementation, got different value")
	}
}

func TestGetAttributeGivenType_NonZeroValueAttribute(t *testing.T) {
	customAttr := IntegerAttributesImpl[int]{
		AllowNegative: true,
		AllowZero:     false,
		Max:           50,
		Min:           10,
	}
	attributes := FTAttributes{
		IntegerAttr: customAttr,
	}
	result, err := attributes.GetAttributeGivenType(reflect.TypeOf(int(0)))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result == nil {
		t.Error("expected custom attribute, got nil")
	}
	if !reflect.DeepEqual(result, customAttr) {
		t.Errorf("expected custom attribute to be returned as-is")
	}
}

func TestGetDefaultForKind_IntegerTypes(t *testing.T) {
	attributes := NewFTAttributes()
	intKinds := []reflect.Kind{
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
	}
	for _, kind := range intKinds {
		t.Run(kind.String(), func(t *testing.T) {
			result, err := attributes.getDefaultForKind(kind)
			if err != nil {
				t.Errorf("unexpected error for %s: %v", kind, err)
			}
			if result == nil {
				t.Errorf("expected result for %s, got nil", kind)
			}
			expected := IntegerAttributesImpl[int64]{}.GetDefaultImplementation()
			if !reflect.DeepEqual(result, expected) {
				t.Errorf("expected IntegerAttributesImpl default for %s", kind)
			}
		})
	}
}
