package attributes

import (
	"fmt"
	"reflect"
)

type NotAnAttributeTypeError struct {
	Type reflect.Type
}

func (naate NotAnAttributeTypeError) Error() string {
	return fmt.Sprintf("The passed type is not an attribute type: %v", naate.Type)
}

type UnsupportedAttributeTypeError struct {
	k reflect.Kind
}

func (uate UnsupportedAttributeTypeError) Error() string {
	return fmt.Sprintf("The following type is not currently supported: %v", uate.k)
}

type NilTypeError struct{}

func (nte NilTypeError) Error() string {
	return "provided type is null"
}
