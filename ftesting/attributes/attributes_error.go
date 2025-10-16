package attributes

import (
	"fmt"
	"reflect"
)

// NotAnAttributeTypeError is returned when attempting to use a type that doesn't
// implement the Attributes interface in a context where an Attributes implementation
// is required.
//
// Fields:
//   - Type: The reflect.Type that was incorrectly used
//
// Example scenario:
//
//	var notAttrs someOtherType
//	// Using notAttrs where Attributes is expected would return NotAnAttributeTypeError
type NotAnAttributeTypeError struct {
	Type reflect.Type
}

func (naate NotAnAttributeTypeError) Error() string {
	return fmt.Sprintf("The passed type is not an attribute type: %v", naate.Type)
}

// UnsupportedAttributeTypeError is returned when attempting to generate random values
// for a Go type that is not currently supported by the attributes system.
//
// Fields:
//   - k: The reflect.Kind that is not supported
//
// This error occurs when FTAttributes.GetAttributeGivenType or getDefaultForKind
// encounters a type Kind that has no corresponding attribute implementation.
//
// Example scenario:
//
//	// Attempting to generate values for an unsupported type like channels
//	chanType := reflect.TypeOf(make(chan int))
//	_, err := attrs.GetAttributeGivenType(chanType)
//	// Returns UnsupportedAttributeTypeError{k: reflect.Chan}
type UnsupportedAttributeTypeError struct {
	k reflect.Kind
}

func (uate UnsupportedAttributeTypeError) Error() string {
	return fmt.Sprintf("The following type is not currently supported: %v", uate.k)
}

// NilTypeError is returned when a nil reflect.Type is passed to methods that
// require a valid type, such as GetAttributeGivenType.
//
// Example scenario:
//
//	var nilType reflect.Type
//	_, err := attrs.GetAttributeGivenType(nilType)
//	// Returns NilTypeError{}
type NilTypeError struct{}

func (nte NilTypeError) Error() string {
	return "provided type is null"
}
