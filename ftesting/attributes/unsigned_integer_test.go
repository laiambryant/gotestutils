package attributes

import (
	"reflect"
	"testing"

	ctesting "github.com/laiambryant/gotestutils/ctesting"
)

func TestUnsignedIntegerAttributes(t *testing.T) {
	var suite []ctesting.CharacterizationTest[bool]

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := UnsignedIntegerAttributesImpl[uint64]{AllowNegative: false, AllowZero: true, Max: 100, Min: 0}
		got := attr.GetAttributes()
		expected := UnsignedIntegerAttributesImpl[uint64]{AllowNegative: false, AllowZero: true, Max: 100, Min: 0}
		return reflect.DeepEqual(got, expected), nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := UnsignedIntegerAttributesImpl[uint64]{}
		got := attr.GetReflectType()
		expected := reflect.TypeOf(uint64(0))
		return got == expected, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := UnsignedIntegerAttributesImpl[uint]{Signed: true}
		reflectType := attr.GetReflectType()
		expected := reflect.TypeOf(int64(0))
		return reflectType == expected, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := UnsignedIntegerAttributesImpl[uint]{AllowNegative: true}
		reflectType := attr.GetReflectType()
		expected := reflect.TypeOf(int64(0))
		return reflectType == expected, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := UnsignedIntegerAttributesImpl[uint]{Signed: false, AllowNegative: false}
		reflectType := attr.GetReflectType()
		expected := reflect.TypeOf(uint64(0))
		return reflectType == expected, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := UnsignedIntegerAttributesImpl[uint64]{}
		got := attr.GetDefaultImplementation()
		return got != nil && reflect.TypeOf(got) == reflect.TypeOf(attr), nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := UnsignedIntegerAttributesImpl[uint64]{Min: 0, Max: 100}
		got := attr.GetRandomValue()
		return got != nil, nil
	}))

	// Edge cases
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := UnsignedIntegerAttributesImpl[uint]{Max: 0, Min: 0}
		result := attr.GetRandomValue()
		return result == uint(0), nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := UnsignedIntegerAttributesImpl[uint]{Max: 5, Min: 10}
		result := attr.GetRandomValue()
		return result == uint(0), nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := UnsignedIntegerAttributesImpl[uint]{Max: 10, Min: 10}
		result := attr.GetRandomValue()
		return result == uint(0), nil
	}))

	results, _ := ctesting.VerifyCharacterizationTestsAndResults(t, suite, true)
	for i, passed := range results {
		if !passed {
			t.Fatalf("UnsignedIntegerAttributes test %d failed", i+1)
		}
	}
}

func TestUnsignedIntegerAttributes_InvalidRange(t *testing.T) {
	attr := UnsignedIntegerAttributesImpl[uint]{Max: 0, Min: 0}
	result := attr.GetRandomValue()
	if result != uint(0) {
		t.Errorf("Expected zero value for invalid range, got %v", result)
	}
}

func TestUnsignedIntegerAttributes_MaxLessThanMin(t *testing.T) {
	attr := UnsignedIntegerAttributesImpl[uint]{Max: 5, Min: 10}
	result := attr.GetRandomValue()
	if result != uint(0) {
		t.Errorf("Expected zero value when max <= min, got %v", result)
	}
}

func TestUnsignedIntegerAttributes_DiffZero(t *testing.T) {
	attr := UnsignedIntegerAttributesImpl[uint]{Max: 10, Min: 10}
	result := attr.GetRandomValue()
	if result != uint(0) {
		t.Errorf("Expected zero value when max == min, got %v", result)
	}
}

func TestUnsignedIntegerAttributes_GetReflectType_Signed(t *testing.T) {
	attr := UnsignedIntegerAttributesImpl[uint]{Signed: true}
	reflectType := attr.GetReflectType()
	expected := reflect.TypeOf(int64(0))
	if reflectType != expected {
		t.Errorf("Expected type %v for signed, got %v", expected, reflectType)
	}
}

func TestGetDefaultForKind_UnsignedIntegerTypes(t *testing.T) {
	attributes := NewFTAttributes()
	uintKinds := []reflect.Kind{
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
	}
	for _, kind := range uintKinds {
		t.Run(kind.String(), func(t *testing.T) {
			result, err := attributes.getDefaultForKind(kind)
			if err != nil {
				t.Errorf("unexpected error for %s: %v", kind, err)
			}
			if result == nil {
				t.Errorf("expected result for %s, got nil", kind)
			}
			expected := UnsignedIntegerAttributesImpl[uint64]{}.GetDefaultImplementation()
			if !reflect.DeepEqual(result, expected) {
				t.Errorf("expected UnsignedIntegerAttributesImpl default for %s", kind)
			}
		})
	}
}

func TestUnsignedIntegerAttributes_GetReflectType_AllowNegative(t *testing.T) {
	attr := UnsignedIntegerAttributesImpl[uint]{AllowNegative: true}
	reflectType := attr.GetReflectType()
	expected := reflect.TypeOf(int64(0))
	if reflectType != expected {
		t.Errorf("Expected type %v for AllowNegative, got %v", expected, reflectType)
	}
}

func TestUnsignedIntegerAttributes_GetReflectType_Unsigned(t *testing.T) {
	attr := UnsignedIntegerAttributesImpl[uint]{Signed: false, AllowNegative: false}
	reflectType := attr.GetReflectType()
	expected := reflect.TypeOf(uint64(0))
	if reflectType != expected {
		t.Errorf("Expected type %v for unsigned, got %v", expected, reflectType)
	}
}

func TestUnsignedIntegerAttributes_ExactMaxEqualMin(t *testing.T) {
	attr := UnsignedIntegerAttributesImpl[uint]{Max: 10, Min: 19}
	result := attr.GetRandomValue()
	if result != uint(0) {
		t.Errorf("Expected zero value when Max == Min, got %v", result)
	}
}
