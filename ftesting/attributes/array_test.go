package attributes

import (
	"reflect"
	"testing"
)

// constIntAttr is a small Attributes implementation used in tests that
// always returns the integer 7 as the random value.
type constIntAttr struct{}

func (c constIntAttr) GetAttributes() any                   { return c }
func (c constIntAttr) GetReflectType() reflect.Type         { return reflect.TypeOf(int(0)) }
func (c constIntAttr) GetRandomValue() any                  { return 7 }
func (c constIntAttr) GetDefaultImplementation() Attributes { return c }

func TestArrayAttributes_GetReflectType_WithReflectType(t *testing.T) {
	attrs := ArrayAttributes{
		Length:       5,
		ElementAttrs: reflect.TypeOf(int(0)),
	}

	reflectType := attrs.GetReflectType()
	if reflectType == nil {
		t.Fatal("Expected non-nil reflect type for array")
	}

	if reflectType.Kind() != reflect.Array {
		t.Errorf("Expected array kind, got %v", reflectType.Kind())
	}

	if reflectType.Len() != 5 {
		t.Errorf("Expected array length 5, got %d", reflectType.Len())
	}

	if reflectType.Elem() != reflect.TypeOf(int(0)) {
		t.Errorf("Expected int element type, got %v", reflectType.Elem())
	}
}

func TestArrayAttributes_GetReflectType_WithNegativeLength(t *testing.T) {
	// Test with negative length
	attrs := ArrayAttributes{
		Length:       -5,
		ElementAttrs: IntegerAttributesImpl[int]{},
	}

	reflectType := attrs.GetReflectType()
	if reflectType != nil {
		t.Errorf("Expected nil reflect type for array with negative length, got %v", reflectType)
	}
}

func TestArrayAttributes_GetReflectType_WithNilElementType(t *testing.T) {
	// Test with nil element attrs
	attrs := ArrayAttributes{
		Length:       5,
		ElementAttrs: nil,
	}

	reflectType := attrs.GetReflectType()
	if reflectType != nil {
		t.Errorf("Expected nil reflect type for array with nil element attrs, got %v", reflectType)
	}
}

func TestArrayAttributes_GetReflectType_WithNilReturningAttribute(t *testing.T) {
	attrs := ArrayAttributes{
		Length:       5,
		ElementAttrs: nilTypeReturningAttribute{},
	}

	reflectType := attrs.GetReflectType()
	if reflectType != nil {
		t.Errorf("Expected nil reflect type for array with nil-returning element attrs, got %v", reflectType)
	}
}

func TestArrayAttributes_ZeroLength(t *testing.T) {
	attr := ArrayAttributes{Length: 0, ElementAttrs: IntegerAttributesImpl[int]{}}
	result := attr.GetRandomValue()
	if result != nil {
		t.Errorf("Expected nil for zero length array, got %v", result)
	}
}

func TestArrayAttributes_NegativeLength(t *testing.T) {
	attr := ArrayAttributes{Length: -5, ElementAttrs: IntegerAttributesImpl[int]{}}
	result := attr.GetRandomValue()
	if result != nil {
		t.Errorf("Expected nil for negative length array, got %v", result)
	}
}

func TestArrayAttributes_InvalidElementType(t *testing.T) {
	attr := ArrayAttributes{Length: 5, ElementAttrs: "not an attribute"}
	result := attr.GetRandomValue()
	if result != nil {
		t.Errorf("Expected nil for invalid element type, got %v", result)
	}
}

func TestArrayAttributes_FillUsesZeroWhenRandomNil(t *testing.T) {
	attrs := ArrayAttributes{
		Length:       3,
		ElementAttrs: nilReturningAttribute{},
	}

	result := attrs.GetRandomValue()
	if result == nil {
		t.Fatal("expected non-nil array result")
	}

	av := reflect.ValueOf(result)
	if av.Kind() != reflect.Array {
		t.Fatalf("expected array, got %v", av.Kind())
	}

	for i := 0; i < av.Len(); i++ {
		el := av.Index(i)
		if el.Kind() == reflect.Int {
			if el.Int() != 0 {
				t.Errorf("expected zero int at index %d, got %v", i, el.Interface())
			}
		} else {
			if !reflect.DeepEqual(el.Interface(), reflect.Zero(el.Type()).Interface()) {
				t.Errorf("expected zero value at index %d, got %v", i, el.Interface())
			}
		}
	}
}

func TestArrayAttributes_FillWithRandomInts(t *testing.T) {
	attrs := ArrayAttributes{
		Length:       4,
		ElementAttrs: constIntAttr{},
	}

	result := attrs.GetRandomValue()
	if result == nil {
		t.Fatal("expected non-nil array result")
	}

	av := reflect.ValueOf(result)
	for i := 0; i < av.Len(); i++ {
		el := av.Index(i)
		if el.Kind() != reflect.Int {
			t.Fatalf("expected int elements, got %v", el.Kind())
		}
		if el.Int() != 7 {
			t.Errorf("expected 7 at index %d, got %v", i, el.Interface())
		}
	}
}
