package attributes

import (
	"reflect"
	"testing"

	ctesting "github.com/laiambryant/gotestutils/ctesting"
	p "github.com/laiambryant/gotestutils/pbtesting/properties/predicates"
)

func TestSliceAttributes_FillUsesZeroWhenRandomNil(t *testing.T) {
	attrs := SliceAttributes{
		MinLen:       3,
		MaxLen:       3,
		ElementAttrs: nilReturningAttribute{},
	}

	result := attrs.GetRandomValue()
	if result == nil {
		t.Fatal("expected non-nil slice result")
	}

	sv := reflect.ValueOf(result)
	if sv.Kind() != reflect.Slice {
		t.Fatalf("expected slice, got %v", sv.Kind())
	}

	for i := 0; i < sv.Len(); i++ {
		el := sv.Index(i)
		// nilReturningAttribute reports element type int; zero value should be 0
		if el.Kind() == reflect.Int {
			if el.Int() != 0 {
				t.Errorf("expected zero int at index %d, got %v", i, el.Interface())
			}
		} else {
			// For safety, ensure it's the zero value by comparing to reflect.Zero
			if !reflect.DeepEqual(el.Interface(), reflect.Zero(el.Type()).Interface()) {
				t.Errorf("expected zero value at index %d, got %v", i, el.Interface())
			}
		}
	}
}

func TestSliceAttributes(t *testing.T) {
	var suite []ctesting.CharacterizationTest[bool]

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := SliceAttributes{MinLen: 1, MaxLen: 3, Unique: true, Sorted: true, ElementPreds: []p.Predicate{}, ElementAttrs: IntegerAttributesImpl[int64]{}}
		got := attr.GetAttributes()
		expected := SliceAttributes{MinLen: 1, MaxLen: 3, Unique: true, Sorted: true, ElementPreds: []p.Predicate{}, ElementAttrs: IntegerAttributesImpl[int64]{}}
		return reflect.DeepEqual(got, expected), nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := SliceAttributes{}
		got := attr.GetDefaultImplementation()
		return got != nil && reflect.TypeOf(got) == reflect.TypeOf(attr), nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := SliceAttributes{MinLen: 1, MaxLen: 3, ElementAttrs: IntegerAttributesImpl[int]{Min: 0, Max: 10}}
		got := attr.GetRandomValue()
		return got != nil, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := SliceAttributes{MinLen: -5, MaxLen: 10, ElementAttrs: IntegerAttributesImpl[int]{}}
		result := attr.GetRandomValue()
		return result != nil, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := SliceAttributes{MinLen: 0, MaxLen: 0, ElementAttrs: IntegerAttributesImpl[int]{}}
		result := attr.GetRandomValue()
		return result != nil, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := SliceAttributes{MinLen: 10, MaxLen: 5, ElementAttrs: IntegerAttributesImpl[int]{}}
		result := attr.GetRandomValue()
		return result != nil, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := SliceAttributes{MinLen: 1, MaxLen: 5, ElementAttrs: "not an attribute"}
		result := attr.GetRandomValue()
		return result == nil, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := SliceAttributes{
			MinLen:       1,
			MaxLen:       5,
			ElementAttrs: reflect.TypeOf(int(0)),
		}
		reflectType := attr.GetReflectType()
		if reflectType == nil {
			return false, nil
		}
		return reflectType.Kind() == reflect.Slice && reflectType.Elem() == reflect.TypeOf(int(0)), nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := SliceAttributes{
			MinLen:       1,
			MaxLen:       5,
			ElementAttrs: nil,
		}
		reflectType := attr.GetReflectType()
		return reflectType == nil, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := SliceAttributes{
			MinLen:       0,
			MaxLen:       0,
			ElementAttrs: IntegerAttributesImpl[int]{Max: 100},
		}
		result := attr.GetRandomValue()
		if result == nil {
			return false, nil
		}
		sliceValue := reflect.ValueOf(result)
		return sliceValue.Len() <= 5, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := SliceAttributes{
			MinLen:       2,
			MaxLen:       3,
			ElementAttrs: nilReturningAttribute{},
		}
		result := attr.GetRandomValue()
		if result == nil {
			return false, nil
		}
		sliceValue := reflect.ValueOf(result)
		if sliceValue.Kind() != reflect.Slice {
			return false, nil
		}
		for i := 0; i < sliceValue.Len(); i++ {
			elem := sliceValue.Index(i)
			if elem.Int() != 0 {
				return false, nil
			}
		}
		return true, nil
	}))

	// Additional tests from individual test functions
	// TestSliceAttributes_MinLenNegative
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := SliceAttributes{MinLen: -5, MaxLen: 10, ElementAttrs: IntegerAttributesImpl[int]{}}
		result := attr.GetRandomValue()
		return result != nil, nil
	}))

	// TestSliceAttributes_MaxLenZero
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := SliceAttributes{MinLen: 0, MaxLen: 0, ElementAttrs: IntegerAttributesImpl[int]{}}
		result := attr.GetRandomValue()
		return result != nil, nil
	}))

	// TestSliceAttributes_MinGreaterThanMax
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := SliceAttributes{MinLen: 10, MaxLen: 5, ElementAttrs: IntegerAttributesImpl[int]{}}
		result := attr.GetRandomValue()
		return result != nil, nil
	}))

	// TestSliceAttributes_NilElementType (converted to bool test)
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := SliceAttributes{MinLen: 1, MaxLen: 5, ElementAttrs: "not an attribute"}
		result := attr.GetRandomValue()
		return result == nil, nil
	}))

	// TestSliceAttributes_DefaultMaxLen
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attrs := SliceAttributes{
			MinLen:       0,
			MaxLen:       0,
			ElementAttrs: IntegerAttributesImpl[int]{Max: 100},
		}
		result := attrs.GetRandomValue()
		if result == nil {
			return false, nil
		}
		sliceValue := reflect.ValueOf(result)
		return sliceValue.Len() <= 5, nil
	}))

	// TestSliceAttributes_GetReflectType_WithReflectType
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attrs := SliceAttributes{
			MinLen:       1,
			MaxLen:       5,
			ElementAttrs: reflect.TypeOf(int(0)),
		}
		reflectType := attrs.GetReflectType()
		if reflectType == nil {
			return false, nil
		}
		return reflectType.Kind() == reflect.Slice && reflectType.Elem() == reflect.TypeOf(int(0)), nil
	}))

	// TestSliceAttributes_GetReflectType_WithNilElementType (converted to bool test)
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attrs := SliceAttributes{
			MinLen:       1,
			MaxLen:       5,
			ElementAttrs: nil,
		}
		reflectType := attrs.GetReflectType()
		return reflectType == nil, nil
	}))

	// TestSliceAttributes_GetReflectType_WithAttributesElement
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attrs := SliceAttributes{
			MinLen:       1,
			MaxLen:       5,
			ElementAttrs: IntegerAttributesImpl[int]{},
		}
		reflectType := attrs.GetReflectType()
		if reflectType == nil {
			return false, nil
		}
		if reflectType.Kind() != reflect.Slice {
			return false, nil
		}
		return reflectType.Elem() == reflect.TypeOf(int(0)), nil
	}))

	results, _ := ctesting.VerifyCharacterizationTestsAndResults(t, suite, true)
	for i, passed := range results {
		if !passed {
			t.Fatalf("SliceAttributes test %d failed", i+1)
		}
	}
}
