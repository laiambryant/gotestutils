package attributes

import (
	"reflect"
	"testing"
)

func TestPointerAttributes_NilWhenAllowNil(t *testing.T) {
	attr := PointerAttributes{AllowNil: true, Depth: 1, Inner: IntegerAttributesImpl[int]{Max: 10, Min: 1}}
	foundNil := false
	for range 100 {
		result := attr.GetRandomValue()
		if result == nil {
			t.Error("GetRandomValue should never return nil directly, it returns a typed nil pointer")
		}
		rv := reflect.ValueOf(result)
		if rv.Kind() == reflect.Pointer && rv.IsNil() {
			foundNil = true
			break
		}
	}
	if !foundNil {
		t.Log("Warning: Did not find nil pointer in 100 attempts (this is statistically unlikely but possible)")
	}
}

func TestPointerAttributes_InvalidInnerType(t *testing.T) {
	attr := PointerAttributes{AllowNil: false, Depth: 1, Inner: "not an attribute"}
	result := attr.GetRandomValue()
	if result != nil {
		t.Errorf("Expected nil for invalid inner type, got %v", result)
	}
}

func TestPointerAttributes_MultipleDepth(t *testing.T) {
	attr := PointerAttributes{AllowNil: false, Depth: 3, Inner: IntegerAttributesImpl[int]{Max: 10, Min: 1}}
	result := attr.GetRandomValue()
	if result == nil {
		t.Error("Expected pointer result, got nil")
	}
	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Pointer {
		t.Errorf("Expected pointer type, got %v", rv.Kind())
	}
}

func TestPointerAttributes_GetReflectType_WithReflectType(t *testing.T) {
	attrs := PointerAttributes{
		Inner: reflect.TypeOf(int(0)),
		Depth: 2,
	}

	reflectType := attrs.GetReflectType()
	if reflectType == nil {
		t.Fatal("Expected non-nil reflect type for pointer")
	}

	if reflectType.Kind() != reflect.Pointer {
		t.Errorf("Expected pointer kind, got %v", reflectType.Kind())
	}

	if reflectType.Elem().Kind() != reflect.Pointer {
		t.Errorf("Expected pointer to pointer, got %v", reflectType.Elem().Kind())
	}

	if reflectType.Elem().Elem() != reflect.TypeOf(int(0)) {
		t.Errorf("Expected final type to be int, got %v", reflectType.Elem().Elem())
	}
}

func TestPointerAttributes_GetReflectType_WithNilInner(t *testing.T) {
	attrs := PointerAttributes{
		Inner: nil,
		Depth: 1,
	}

	reflectType := attrs.GetReflectType()
	if reflectType != nil {
		t.Errorf("Expected nil reflect type for pointer with nil inner, got %v", reflectType)
	}
}

func TestPointerAttributes_GetReflectType_WithZeroDepth(t *testing.T) {
	attrs := PointerAttributes{
		Inner: IntegerAttributesImpl[int]{},
		Depth: 0,
	}

	reflectType := attrs.GetReflectType()
	if reflectType == nil {
		t.Fatal("Expected non-nil reflect type for pointer")
	}

	if reflectType.Kind() != reflect.Pointer {
		t.Errorf("Expected pointer kind, got %v", reflectType.Kind())
	}
}

func TestPointerAttributes_GetReflectType_WithNegativeDepth(t *testing.T) {
	attrs := PointerAttributes{
		Inner: IntegerAttributesImpl[int]{},
		Depth: -5,
	}

	reflectType := attrs.GetReflectType()
	if reflectType == nil {
		t.Fatal("Expected non-nil reflect type for pointer")
	}

	if reflectType.Kind() != reflect.Pointer {
		t.Errorf("Expected pointer kind, got %v", reflectType.Kind())
	}
}

func TestPointerAttributes_NilInnerValueWithType(t *testing.T) {
	attrs := PointerAttributes{
		Inner:    nilReturningAttribute{},
		AllowNil: false,
		Depth:    1,
	}

	result := attrs.GetRandomValue()
	if result == nil {
		t.Fatal("Expected non-nil pointer")
	}

	ptrValue := reflect.ValueOf(result)
	if ptrValue.Kind() != reflect.Pointer {
		t.Fatalf("Expected pointer, got %v", ptrValue.Kind())
	}
}
