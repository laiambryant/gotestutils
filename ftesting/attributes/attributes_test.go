package attributes

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetAttributeGivenTypeNil(t *testing.T) {
	attributes := NewFTAttributes()
	_, err := attributes.GetAttributeGivenType(nil)
	if err == nil {
		t.Error("expected NilTypeError")
	}
	if _, ok := err.(NilTypeError); !ok {
		t.Error("expected error to be of type NilTypeError")
	}
}

func TestGetAttributeGivenType_NilAttribute(t *testing.T) {
	attributes := FTAttributes{
		IntegerAttr: nil,
	}
	result, err := attributes.GetAttributeGivenType(reflect.TypeOf(int(0)))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result == nil {
		t.Error("expected default implementation, got nil")
	}
}

func TestGetAttributeGivenType_GetAttributesReturnsNil(t *testing.T) {
	attributes := FTAttributes{
		IntegerAttr: nilAttributeType{},
	}
	result, err := attributes.GetAttributeGivenType(reflect.TypeOf(int(0)))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if _, ok := result.(IntegerAttributesImpl[int]); !ok {
		t.Errorf("expected IntegerAttributesImpl, got %T", result)
	}
}

func TestGetAttributeGivenType_TypeOfAttributesIsNil(t *testing.T) {
	attributes := FTAttributes{
		IntegerAttr: nilTypeAttribute{},
	}
	result, err := attributes.GetAttributeGivenType(reflect.TypeOf(int(0)))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result == nil {
		t.Error("expected result, got nil")
	}
}

func TestGetDefaultForKind_UnsupportedTypes(t *testing.T) {
	attributes := NewFTAttributes()
	unsupportedKinds := []reflect.Kind{
		reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Invalid,
		reflect.Uintptr,
		reflect.UnsafePointer,
	}

	for _, kind := range unsupportedKinds {
		t.Run(kind.String(), func(t *testing.T) {
			result, err := attributes.getDefaultForKind(kind)
			if err == nil {
				t.Errorf("expected error for unsupported kind %s", kind)
			}
			if result != nil {
				t.Errorf("expected nil result for unsupported kind %s, got %v", kind, result)
			}
			if _, ok := err.(UnsupportedAttributeTypeError); !ok {
				t.Errorf("expected UnsupportedAttributeTypeError for %s, got %T: %v", kind, err, err)
			}
		})
	}
}

func TestGetDefaultForKind_AllKindsCovered(t *testing.T) {
	attributes := NewFTAttributes()
	supportedKinds := map[reflect.Kind]bool{
		reflect.Int: true, reflect.Int8: true, reflect.Int16: true, reflect.Int32: true, reflect.Int64: true,
		reflect.Uint: true, reflect.Uint8: true, reflect.Uint16: true, reflect.Uint32: true, reflect.Uint64: true,
		reflect.Float32: true, reflect.Float64: true,
		reflect.Complex64: true, reflect.Complex128: true,
		reflect.String: true, reflect.Slice: true, reflect.Bool: true,
		reflect.Map: true, reflect.Pointer: true, reflect.Struct: true, reflect.Array: true,
	}
	for kind := reflect.Invalid; kind <= reflect.UnsafePointer; kind++ {
		result, err := attributes.getDefaultForKind(kind)
		if supportedKinds[kind] {
			if err != nil {
				t.Errorf("expected no error for supported kind %s, got: %v", kind, err)
			}
			if result == nil {
				t.Errorf("expected non-nil result for supported kind %s", kind)
			}
		} else {
			if err == nil {
				t.Errorf("expected error for unsupported kind %s", kind)
			}
			if result != nil {
				t.Errorf("expected nil result for unsupported kind %s", kind)
			}
		}
	}
}

func TestNotAnAttributeTypeError_Error(t *testing.T) {
	typ := reflect.TypeOf(0)
	err := NotAnAttributeTypeError{Type: typ}
	expected := fmt.Sprintf("The passed type is not an attribute type: %v", typ)
	if err.Error() != expected {
		t.Errorf("unexpected error message: got %q, want %q", err.Error(), expected)
	}
}

func TestUnsupportedAttributeTypeError_Error(t *testing.T) {
	k := reflect.Chan
	err := UnsupportedAttributeTypeError{k}
	expected := fmt.Sprintf("The following type is not currently supported: %v", k)
	if err.Error() != expected {
		t.Errorf("unexpected error message: got %q, want %q", err.Error(), expected)
	}
}

func TestNilTypeError_Error(t *testing.T) {
	err := NilTypeError{}
	expected := "provided type is null"
	if err.Error() != expected {
		t.Errorf("unexpected error message: got %q, want %q", err.Error(), expected)
	}
}
