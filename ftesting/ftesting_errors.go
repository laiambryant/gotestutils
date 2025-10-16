package ftesting

import (
	"fmt"
	"reflect"
)

// NoFunctionProvidedError is returned when attempting to generate inputs or execute
// a fuzz test without first setting a function using WithFunction.
//
// Example scenario:
//
//	ft := &FTesting{}
//	_, err := ft.GenerateInputs() // Returns NoFunctionProvidedError
type NoFunctionProvidedError struct{}

func (nfpe NoFunctionProvidedError) Error() string {
	return "no function was provided to ftesting suite"
}

// NotAFunctionError is returned when the value provided to WithFunction is not
// a callable function. The error includes the actual kind of the provided value.
//
// Fields:
//   - k: The reflect.Kind of the value that was incorrectly provided
//
// Example scenario:
//
//	ft := &FTesting{}
//	ft.WithFunction(42) // An integer, not a function
//	_, err := ft.GenerateInputs() // Returns NotAFunctionError{k: reflect.Int}
type NotAFunctionError struct {
	k reflect.Kind
}

func (nafe NotAFunctionError) Error() string {
	return fmt.Sprintf("f is not a function: %v", nafe.k)
}

// InputsGenerationError wraps errors that occur during random input generation
// for function parameters. This typically occurs when the attribute system cannot
// generate a value for a particular type.
//
// Fields:
//   - err: The underlying error from the attribute generation system
//
// Example scenario:
//
//	ft.WithFunction(func(x UnsupportedType) {})
//	_, err := ft.GenerateInputs() // Returns InputsGenerationError wrapping the attribute error
type InputsGenerationError struct {
	err error
}

func (ige InputsGenerationError) Error() string {
	return fmt.Sprintf("error in input generation: %v", ige.err.Error())
}
