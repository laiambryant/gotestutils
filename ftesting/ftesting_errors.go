package ftesting

import (
	"fmt"
	"reflect"
)

type NoFunctionProvidedError struct{}

func (nfpe NoFunctionProvidedError) Error() string {
	return "no function was provided to ftesting suite"
}

type NotAFunctionError struct {
	k reflect.Kind
}

func (nafe NotAFunctionError) Error() string {
	return fmt.Sprintf("f is not a function: %v", nafe.k)
}

type InputsGenerationError struct {
	err error
}

func (ige InputsGenerationError) Error() string {
	return fmt.Sprintf("error in input generation: %v", ige.err.Error())
}
