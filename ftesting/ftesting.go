// Package ftesting provides a sophisticated fuzz testing framework for Go functions.
//
// The ftesting package enables property-based and fuzz testing by automatically generating
// random inputs for functions with arbitrary signatures. It uses reflection to inspect
// function parameter types and a comprehensive attribute system to generate type-appropriate
// random values with configurable constraints.
//
// Key Features:
//   - Automatic input generation for any function signature
//   - Type-safe random value generation using generics
//   - Configurable attributes for controlling value generation ranges and constraints
//   - Support for all Go primitive types, collections, and composite types
//   - Integration with Go's testing framework
//
// Basic Usage:
//
//	func TestMyFunction(t *testing.T) {
//	    ft := &ftesting.FTesting{t: t}
//	    ft.WithFunction(myFunction).
//	       WithIterations(1000).
//	       Verify()
//	}
//
// Advanced Usage with Custom Attributes:
//
//	func TestWithCustomAttrs(t *testing.T) {
//	    attrs := attributes.NewFTAttributes()
//	    attrs.IntegerAttr = IntegerAttributesImpl[int]{
//	        Min: 0,
//	        Max: 100,
//	        AllowZero: false,
//	    }
//	    ft := &ftesting.FTesting{t: t}
//	    ft.WithFunction(myFunction).
//	       WithAttributes(attrs).
//	       WithIterations(500).
//	       Verify()
//	}
//
// See the attributes subpackage for detailed documentation on configuring
// random value generation for specific types.
package ftesting

import (
	"fmt"
	"reflect"
	"testing"

	a "github.com/laiambryant/gotestutils/ftesting/attributes"
)

// FTesting represents a fuzz testing suite that generates random inputs
// for testing functions with arbitrary signatures.
//
// FTesting uses reflection and a sophisticated attribute system to automatically
// generate type-appropriate random values for function parameters, enabling
// property-based and fuzz testing without manual input generation.
//
// Fields:
//   - f: The function to test (can be any function signature)
//   - iterations: Number of test iterations to run
//   - attributes: Configuration for random value generation per type
//   - t: The testing.T instance for reporting results
//
// Example usage:
//
//	ft := &FTesting{}
//	ft.WithFunction(myFunc).WithIterations(100).WithAttributes(customAttrs).Verify()
type FTesting struct {
	f          any
	iterations uint
	attributes a.AttributesStruct
	t          *testing.T
}

// WithIterations sets the number of iterations for the fuzz test.
// Each iteration generates a new set of random inputs and executes the test function.
//
// Parameters:
//   - n: The number of test iterations to execute
//
// Returns the FTesting instance for method chaining.
//
// Example usage:
//
//	ft.WithIterations(1000) // Run 1000 fuzz test iterations
func (mt *FTesting) WithIterations(n uint) *FTesting { mt.iterations = n; return mt }

// WithFunction sets the function to be tested. The function can have any signature,
// and FTesting will use reflection to determine parameter types and generate
// appropriate random inputs.
//
// Parameters:
//   - f: The function to test (can be any callable function)
//
// Returns the FTesting instance for method chaining.
//
// Example usage:
//
//	ft.WithFunction(func(x int, y string) (int, error) {
//	    return len(y) + x, nil
//	})
func (mt *FTesting) WithFunction(f any) *FTesting {
	mt.f = f
	return mt
}

// WithAttributes sets custom attribute configurations for random value generation.
// Attributes control how random values are generated for each type (ranges, constraints, etc.).
//
// Parameters:
//   - a: An AttributesStruct instance containing type-specific generation rules
//
// Returns the FTesting instance for method chaining.
//
// Example usage:
//
//	attrs := attributes.NewFTAttributes()
//	attrs.IntegerAttr = IntegerAttributesImpl[int]{Min: 0, Max: 100}
//	ft.WithAttributes(attrs)
func (mt *FTesting) WithAttributes(a a.AttributesStruct) *FTesting {
	mt.attributes = a
	return mt
}

// GenerateInputs creates a slice of random input values matching the parameter types
// of the configured test function. This method uses reflection to inspect the function
// signature and the attribute system to generate type-appropriate values.
//
// Returns:
//   - []any: A slice of generated input values, one per function parameter
//   - error: An error if the function is nil, not a function, or if input generation fails
//
// Errors returned:
//   - NoFunctionProvidedError: When no function has been set with WithFunction
//   - NotAFunctionError: When the provided value is not a callable function
//   - Attribute-related errors: When random value generation fails for a parameter type
//
// The method automatically initializes default attributes if none were provided.
//
// Example usage:
//
//	ft.WithFunction(func(x int, y string) int { return x + len(y) })
//	inputs, err := ft.GenerateInputs()
//	// inputs might be: []any{42, "hello"}
func (mt *FTesting) GenerateInputs() ([]any, error) {
	if mt.f == nil {
		return nil, &NoFunctionProvidedError{}
	}
	if reflect.TypeOf(mt.f).Kind() != reflect.Func {
		return nil, &NotAFunctionError{}
	}
	if mt.attributes == nil {
		mt.attributes = a.NewFTAttributes()
	}
	fType := reflect.TypeOf(mt.f)
	args := make([]any, fType.NumIn())
	for i := 0; i < fType.NumIn(); i++ {
		argType := fType.In(i)
		v, err := mt.attributes.GetAttributeGivenType(argType)
		if err != nil {
			return nil, err
		}
		args[i] = v.GetRandomValue()
	}
	return args, nil
}

// ApplyFunction generates random inputs and executes the configured test function
// with those inputs. This method combines input generation and function execution
// into a single operation.
//
// Returns:
//   - bool: true if the function executed successfully, false otherwise
//   - error: An error if input generation fails or if the function is not set
//
// The method uses reflection to call the function with generated arguments and
// discards the return values. The focus is on whether the function can execute
// without panicking, not on validating return values.
//
// Example usage:
//
//	ft.WithFunction(func(x int) int { return x * 2 })
//	success, err := ft.ApplyFunction()
//	if err != nil {
//	    // Handle error
//	}
func (mt *FTesting) ApplyFunction() (bool, error) {
	if mt.f == nil {
		return false, fmt.Errorf("function is nil")
	}
	inputs, err := mt.GenerateInputs()
	if err != nil {
		return false, fmt.Errorf("failed to generate inputs: %w", err)
	}
	args := make([]reflect.Value, len(inputs))
	for i, input := range inputs {
		args[i] = reflect.ValueOf(input)
	}
	fValue := reflect.ValueOf(mt.f)
	_ = fValue.Call(args)
	return true, nil
}

// Verify executes the fuzz test and reports results using the configured testing.T instance.
// This is the primary entry point for running fuzz tests. It calls ApplyFunction and
// reports any errors to the test framework.
//
// If no testing.T instance was provided, this method returns silently without testing.
// For successful execution, the method returns without reporting. For failures, it calls
// t.Errorf or t.Error as appropriate.
//
// Behavior:
//   - Returns early if testing.T is nil (no-op)
//   - Calls ApplyFunction to generate inputs and execute the function
//   - Reports errors via t.Errorf with detailed error messages
//   - Reports general failures via t.Error
//
// This method should be called after configuring the FTesting instance with
// WithFunction, WithIterations (optional), and WithAttributes (optional).
//
// Example usage:
//
//	ft := &FTesting{t: t}
//	ft.WithFunction(myFunc).WithIterations(100).Verify()
//	// Test results are automatically reported to testing.T
func (mt *FTesting) Verify() {
	if mt.t == nil {
		return
	}
	ok, err := mt.ApplyFunction()
	if err != nil {
		mt.t.Errorf("Test Failed with error: [%s]", err.Error())
	}
	if !ok {
		mt.t.Error("Test Failed")
	}
}
