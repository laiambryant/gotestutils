package generation

import (
	"fmt"
	"reflect"
)

type UnknownTypeError struct {
	rt reflect.Type
}

func (ute UnknownTypeError) Error() string {
	return fmt.Sprintf("The type is not supported: %v", ute.rt)
}

type EmptyArrayError struct {
}

func (eae EmptyArrayError) Error() string {
	return "The array is empty"
}
