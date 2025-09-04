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

type AttributeConflictError struct {
	conflict string
}

func (ace AttributeConflictError) Error() string {
	return fmt.Sprintf("there is a conflict between attributes: %s", ace.conflict)
}

type GenerationError struct {
	Op  string
	Msg string
}

func (e GenerationError) Error() string { return fmt.Sprintf("generation %s: %s", e.Op, e.Msg) }
