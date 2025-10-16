package attributes

import (
	"reflect"
	"testing"

	"github.com/laiambryant/gotestutils/ctesting"
)

func TestMapAttributes(t *testing.T) {
	var suite []ctesting.CharacterizationTest[bool]

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := MapAttributes{MinSize: 1, MaxSize: 5, KeyAttrs: StringAttributes{}, ValueAttrs: IntegerAttributesImpl[int]{}}
		got := attr.GetAttributes()
		expected := MapAttributes{MinSize: 1, MaxSize: 5, KeyAttrs: StringAttributes{}, ValueAttrs: IntegerAttributesImpl[int]{}}
		return reflect.DeepEqual(got, expected), nil
	}))
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := MapAttributes{}
		got := attr.GetDefaultImplementation()
		return got != nil && reflect.TypeOf(got) == reflect.TypeOf(attr), nil
	}))
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := MapAttributes{MinSize: 1, MaxSize: 3, KeyAttrs: StringAttributes{}, ValueAttrs: IntegerAttributesImpl[int]{}}
		got := attr.GetRandomValue()
		return got != nil, nil
	}))
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := MapAttributes{MinSize: -5, MaxSize: 10, KeyAttrs: StringAttributes{}, ValueAttrs: IntegerAttributesImpl[int]{}}
		result := attr.GetRandomValue()
		return result != nil, nil
	}))
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := MapAttributes{MinSize: 0, MaxSize: 0, KeyAttrs: StringAttributes{}, ValueAttrs: IntegerAttributesImpl[int]{}}
		result := attr.GetRandomValue()
		return result != nil, nil
	}))
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := MapAttributes{MinSize: 10, MaxSize: 5, KeyAttrs: StringAttributes{}, ValueAttrs: IntegerAttributesImpl[int]{}}
		result := attr.GetRandomValue()
		return result != nil, nil
	}))
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := MapAttributes{MinSize: 1, MaxSize: 5, KeyAttrs: "not an attribute", ValueAttrs: IntegerAttributesImpl[int]{}}
		result := attr.GetRandomValue()
		return result == nil, nil
	}))
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := MapAttributes{MinSize: 1, MaxSize: 5, KeyAttrs: StringAttributes{}, ValueAttrs: "not an attribute"}
		result := attr.GetRandomValue()
		return result == nil, nil
	}))
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attrs := MapAttributes{
			MinSize:    1,
			MaxSize:    5,
			KeyAttrs:   reflect.TypeOf(""),
			ValueAttrs: reflect.TypeOf(0),
		}
		reflectType := attrs.GetReflectType()
		if reflectType == nil || reflectType.Kind() != reflect.Map {
			return false, nil
		}
		return reflectType.Key() == reflect.TypeOf("") && reflectType.Elem() == reflect.TypeOf(0), nil
	}))
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attrs := MapAttributes{
			MinSize:    1,
			MaxSize:    5,
			KeyAttrs:   nil,
			ValueAttrs: IntegerAttributesImpl[int]{},
		}
		reflectType := attrs.GetReflectType()
		return reflectType == nil, nil
	}))
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attrs := MapAttributes{
			MinSize:    1,
			MaxSize:    5,
			KeyAttrs:   StringAttributes{},
			ValueAttrs: nil,
		}
		reflectType := attrs.GetReflectType()
		return reflectType == nil, nil
	}))
	results, _ := ctesting.VerifyCharacterizationTestsAndResults(t, suite, true)
	for i, passed := range results {
		if !passed {
			t.Fatalf("MapAttributes test %d failed", i+1)
		}
	}
}

func TestMapAttributes_MaxSizeZero(t *testing.T) {
	attr := MapAttributes{MinSize: 0, MaxSize: 0, KeyAttrs: StringAttributes{}, ValueAttrs: IntegerAttributesImpl[int]{}}
	result := attr.GetRandomValue()
	if result == nil {
		t.Error("Expected map result, got nil")
	}
}

func TestMapAttributes_MinGreaterThanMax(t *testing.T) {
	attr := MapAttributes{MinSize: 10, MaxSize: 5, KeyAttrs: StringAttributes{}, ValueAttrs: IntegerAttributesImpl[int]{}}
	result := attr.GetRandomValue()
	if result == nil {
		t.Error("Expected map result, got nil")
	}
}

func TestMapAttributes_InvalidKeyType(t *testing.T) {
	attr := MapAttributes{MinSize: 1, MaxSize: 5, KeyAttrs: "not an attribute", ValueAttrs: IntegerAttributesImpl[int]{}}
	result := attr.GetRandomValue()
	if result != nil {
		t.Errorf("Expected nil for invalid key type, got %v", result)
	}
}

func TestMapAttributes_InvalidValueType(t *testing.T) {
	attr := MapAttributes{MinSize: 1, MaxSize: 5, KeyAttrs: StringAttributes{}, ValueAttrs: "not an attribute"}
	result := attr.GetRandomValue()
	if result != nil {
		t.Errorf("Expected nil for invalid value type, got %v", result)
	}
}

func TestMapAttributes_NilKeyValue(t *testing.T) {
	attrs := MapAttributes{
		MinSize:    2,
		MaxSize:    3,
		KeyAttrs:   nilReturningAttribute{},
		ValueAttrs: nilReturningAttribute{},
	}
	result := attrs.GetRandomValue()
	if result == nil {
		t.Fatal("Expected non-nil map")
	}
	mapValue := reflect.ValueOf(result)
	if mapValue.Kind() != reflect.Map {
		t.Fatalf("Expected map, got %v", mapValue.Kind())
	}
	if mapValue.Len() < 1 {
		t.Errorf("Expected at least 1 entry, got %d", mapValue.Len())
	}
}

func TestMapAttributes_EqualMinMaxSize(t *testing.T) {
	attrs := MapAttributes{
		MinSize:    5,
		MaxSize:    5,
		KeyAttrs:   StringAttributes{MinLen: 5, MaxLen: 10},
		ValueAttrs: IntegerAttributesImpl[int]{Max: 100},
	}
	result := attrs.GetRandomValue()
	if result == nil {
		t.Fatal("Expected non-nil map")
	}
	mapValue := reflect.ValueOf(result)
	if mapValue.Len() < 3 || mapValue.Len() > 5 {
		t.Errorf("Expected map size around 5, got %d", mapValue.Len())
	}
}

func TestMapAttributes_DefaultSizes(t *testing.T) {
	attrs := MapAttributes{
		MinSize:    -1,
		MaxSize:    0,
		KeyAttrs:   IntegerAttributesImpl[int]{Max: 100},
		ValueAttrs: StringAttributes{MaxLen: 10},
	}
	result := attrs.GetRandomValue()
	if result == nil {
		t.Fatal("Expected non-nil map")
	}
	mapValue := reflect.ValueOf(result)
	if mapValue.Len() > 5 {
		t.Errorf("Expected max size 5 (default), got %d", mapValue.Len())
	}
}

func TestGetAttributeGivenType_KindMapHit(t *testing.T) {
	attributes := NewFTAttributes()
	testCases := []struct {
		name     string
		typ      reflect.Type
		expected Attributes
	}{
		{"int", reflect.TypeOf(int(0)), attributes.IntegerAttr},
		{"int32", reflect.TypeOf(int32(0)), attributes.IntegerAttr},
		{"uint", reflect.TypeOf(uint(0)), attributes.UIntegerAttr},
		{"float64", reflect.TypeOf(float64(0)), attributes.FloatAttr},
		{"complex128", reflect.TypeOf(complex128(0)), attributes.ComplexAttr},
		{"string", reflect.TypeOf(""), attributes.StringAttr},
		{"bool", reflect.TypeOf(true), attributes.BoolAttr},
		{"slice", reflect.TypeOf([]int{}), attributes.SliceAttr},
		{"map", reflect.TypeOf(map[string]int{}), attributes.MapAttr},
		{"pointer", reflect.TypeOf(new(int)), attributes.PointerAttr},
		{"struct", reflect.TypeOf(struct{}{}), attributes.StructAttr},
		{"array", reflect.TypeOf([3]int{}), attributes.ArrayAttr},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := attributes.GetAttributeGivenType(tc.typ)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if reflect.TypeOf(result) != reflect.TypeOf(tc.expected) {
				t.Errorf("expected type %T, got %T", tc.expected, result)
			}
		})
	}
}

func TestGetAttributeGivenType_KindNotInMap(t *testing.T) {
	attributes := NewFTAttributes()
	testCases := []struct {
		name string
		typ  reflect.Type
	}{
		{"chan", reflect.TypeOf(make(chan int))},
		{"func", reflect.TypeOf(func() {})},
		{"interface", reflect.TypeOf((*any)(nil)).Elem()},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := attributes.GetAttributeGivenType(tc.typ)
			if err == nil {
				t.Error("expected error for unsupported type")
			}
			if _, ok := err.(UnsupportedAttributeTypeError); !ok {
				t.Errorf("expected UnsupportedAttributeTypeError, got %T: %v", err, err)
			}
		})
	}
}

func TestMapAttributes_GetReflectType_WithReflectTypes(t *testing.T) {
	attrs := MapAttributes{
		MinSize:    1,
		MaxSize:    5,
		KeyAttrs:   reflect.TypeOf(""),
		ValueAttrs: reflect.TypeOf(0),
	}

	reflectType := attrs.GetReflectType()
	if reflectType == nil {
		t.Fatal("Expected non-nil reflect type for map")
	}
	if reflectType.Kind() != reflect.Map {
		t.Errorf("Expected map kind, got %v", reflectType.Kind())
	}
	if reflectType.Key() != reflect.TypeOf("") {
		t.Errorf("Expected string key type, got %v", reflectType.Key())
	}
	if reflectType.Elem() != reflect.TypeOf(0) {
		t.Errorf("Expected int value type, got %v", reflectType.Elem())
	}
}

func TestMapAttributes_GetReflectType_WithNilKeyType(t *testing.T) {
	attrs := MapAttributes{
		MinSize:    1,
		MaxSize:    5,
		KeyAttrs:   nil,
		ValueAttrs: IntegerAttributesImpl[int]{},
	}

	reflectType := attrs.GetReflectType()
	if reflectType != nil {
		t.Errorf("Expected nil reflect type for map with nil key attrs, got %v", reflectType)
	}
}

func TestMapAttributes_GetReflectType_WithNilValueType(t *testing.T) {
	attrs := MapAttributes{
		MinSize:    1,
		MaxSize:    5,
		KeyAttrs:   StringAttributes{},
		ValueAttrs: nil,
	}

	reflectType := attrs.GetReflectType()
	if reflectType != nil {
		t.Errorf("Expected nil reflect type for map with nil value attrs, got %v", reflectType)
	}
}
