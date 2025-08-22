package attributes

import (
	"fmt"
	"reflect"
)

type NotAnAttributeTypeError struct {
	Type reflect.Type
}

func (naae NotAnAttributeTypeError) Error() string {
	return fmt.Sprintf("The passed type is not an attribute type: %v", naae.Type)
}
