// Package pbtesting provides a property-based testing framework for Go functions.
//
// Property-based testing validates that functions satisfy certain properties (predicates)
// across a wide range of randomly generated inputs. Unlike traditional unit tests that
// check specific input-output pairs, property-based tests verify that invariants hold
// for all inputs.
//
// The pbtesting package builds on the ftesting framework for input generation and adds
// predicate validation to ensure function outputs satisfy specified properties.
//
// Key Features:
//   - Automatic random input generation using ftesting
//   - Predicate-based validation of function outputs
//   - Configurable number of test iterations
//   - Support for functions with multiple return values
//   - Detailed failure reporting with failing predicates
//   - Integration with Go's testing framework
//
// Basic Usage:
//
//	func TestAbsoluteValue(t *testing.T) {
//	    // Property: abs(x) should always be non-negative
//	    nonNegative := predicates.Predicate{
//	        Verify: func(v any) bool { return v.(int) >= 0 },
//	    }
//
//	    test := NewPBTest(math.Abs).
//	        WithIterations(1000).
//	        WithPredicates(nonNegative).
//	        WithT(t)
//
//	    results, err := test.Run()
//	    if err != nil {
//	        t.Fatal(err)
//	    }
//
//	    // Check for failures
//	    failures := FilterPBTTestOut(results)
//	    if len(failures) > 0 {
//	        t.Errorf("Found %d failing cases", len(failures))
//	    }
//	}
//
// Advanced Usage with Custom Attributes:
//
//	func TestWithCustomInputs(t *testing.T) {
//	    attrs := attributes.NewFTAttributes()
//	    attrs.IntegerAttr = IntegerAttributesImpl[int]{Min: 1, Max: 100}
//
//	    test := NewPBTest(myFunction).
//	        WithIterations(500).
//	        WithArgAttributes(attrs).
//	        WithPredicates(pred1, pred2).
//	        WithT(t)
//
//	    results, _ := test.Run()
//	    // Process results...
//	}
//
// See the properties/predicates subpackage for predicate implementations.
package pbtesting

import (
	"reflect"
	"testing"

	"github.com/laiambryant/gotestutils/ftesting"
	p "github.com/laiambryant/gotestutils/pbtesting/properties/predicates"
	"github.com/laiambryant/gotestutils/utils"
)

// PBTest represents a property-based test suite that validates function outputs
// against a set of predicates across multiple iterations with randomly generated inputs.
//
// The test execution flow:
// 1. Generate random inputs using ftesting
// 2. Execute the function with those inputs
// 3. Validate outputs against configured predicates
// 4. Collect results including any predicate failures
//
// Fields:
//   - t: The testing.T instance for reporting results
//   - f: The function to test (can be any function signature)
//   - predicates: List of predicates that outputs must satisfy
//   - iterations: Number of test iterations to run
//   - argAttrs: Custom attributes for controlling input generation
//
// Example usage:
//
//	test := NewPBTest(myFunc).
//	    WithIterations(1000).
//	    WithPredicates(nonNegative, lessThan100).
//	    WithT(t)
type PBTest struct {
	t          *testing.T
	f          any
	predicates []p.Predicate
	iterations uint
	argAttrs   []any
}

// PBTestOut represents the result of a single property-based test iteration.
// It contains the function output, any predicates that failed, and a success flag.
//
// Fields:
//   - Output: The value returned by the function under test
//   - Predicates: List of predicates that failed for this output (nil if all passed)
//   - ok: true if all predicates passed, false if any failed
//
// Use FilterPBTTestOut to extract only the failing test cases from a slice of results.
//
// Example usage:
//
//	results, _ := test.Run()
//	for _, result := range results {
//	    if !result.ok {
//	        t.Errorf("Output %v failed predicates: %v", result.Output, result.Predicates)
//	    }
//	}
type PBTestOut struct {
	Output     any
	Predicates []p.Predicate
	ok         bool
}

// returnTypes is an internal type constraint for function return values.
// It allows functions to return either a single value or multiple values (as a slice).
type returnTypes interface {
	any | []any
}

// NewPBTest creates a new property-based test instance with the specified function.
// The function can have any signature; reflection is used to handle parameter types
// and return values.
//
// Parameters:
//   - f: The function to test (can be nil, but must be set before calling Run)
//
// Returns a PBTest instance configured with 1 iteration by default.
// Use the WithX methods to configure additional settings before calling Run.
//
// Example usage:
//
//	test := NewPBTest(func(x int) int { return x * x })
//	test.WithIterations(100).WithPredicates(nonNegative)
func NewPBTest(f any) *PBTest { return &PBTest{f: f, iterations: 1} }

// WithIterations sets the number of test iterations to run.
// Each iteration generates new random inputs and validates the output.
//
// Parameters:
//   - n: The number of iterations (must be > 0 for meaningful tests)
//
// Returns the PBTest instance for method chaining.
//
// Example usage:
//
//	test.WithIterations(1000) // Run 1000 property tests
func (pbt *PBTest) WithIterations(n uint) *PBTest { pbt.iterations = n; return pbt }

// WithPredicates sets the predicates that function outputs must satisfy.
// All predicates must pass for a test iteration to be considered successful.
//
// Parameters:
//   - preds: One or more predicates to validate against
//
// Returns the PBTest instance for method chaining.
//
// Example usage:
//
//	test.WithPredicates(
//	    predicates.NewNonNegative(),
//	    predicates.NewLessThan(100),
//	)
func (pbt *PBTest) WithPredicates(preds ...p.Predicate) *PBTest { pbt.predicates = preds; return pbt }

// WithArgAttributes sets custom attributes for controlling how random input values
// are generated. This allows fine-grained control over the input space explored
// during testing.
//
// Parameters:
//   - attrs: Attribute configurations for function parameters (order must match parameters)
//
// Returns the PBTest instance for method chaining.
//
// Example usage:
//
//	intAttr := IntegerAttributesImpl[int]{Min: 0, Max: 100}
//	test.WithArgAttributes(intAttr)
func (pbt *PBTest) WithArgAttributes(attrs ...any) *PBTest { pbt.argAttrs = attrs; return pbt }

// WithT sets the testing.T instance for integration with Go's testing framework.
// While not required for test execution, it's recommended for proper test reporting.
//
// Parameters:
//   - t: The testing.T instance from the test function
//
// Returns the PBTest instance for method chaining.
//
// Example usage:
//
//	func TestMyFunction(t *testing.T) {
//	    test := NewPBTest(myFunc).WithT(t)
//	    // ...
//	}
func (pbt *PBTest) WithT(t *testing.T) *PBTest { pbt.t = t; return pbt }

// WithF sets or updates the function to be tested. This allows changing the function
// after the PBTest instance was created.
//
// Parameters:
//   - f: The function to test
//
// Returns the PBTest instance for method chaining.
//
// Example usage:
//
//	test := NewPBTest(nil)
//	test.WithF(myFunction).WithIterations(100)
func (pbt *PBTest) WithF(f any) *PBTest { pbt.f = f; return pbt }

// Run executes the property-based test by performing the configured number of iterations.
// For each iteration, it:
// 1. Generates random inputs using the ftesting framework
// 2. Executes the test function with those inputs
// 3. Validates outputs against configured predicates
// 4. Collects results
//
// Returns:
//   - []PBTestOut: A slice containing results for each iteration
//   - error: An error if input generation fails or the function is invalid
//
// The returned slice includes both passing and failing iterations. Use FilterPBTTestOut
// to extract only the failures.
//
// If no predicates are configured, all iterations are marked as successful (ok=true).
// If the function is nil, returns an empty slice with no error.
//
// Example usage:
//
//	test := NewPBTest(abs).WithIterations(100).WithPredicates(nonNegative)
//	results, err := test.Run()
//	if err != nil {
//	    t.Fatal(err)
//	}
//
//	failures := FilterPBTTestOut(results)
//	if len(failures) > 0 {
//	    t.Errorf("Found %d property violations", len(failures))
//	    for _, failure := range failures {
//	        t.Logf("Output: %v, Failed predicates: %v", failure.Output, failure.Predicates)
//	    }
//	}
func (pbt *PBTest) Run() (retOut []PBTestOut, err error) {
	if pbt.f == nil {
		return []PBTestOut{}, nil
	}
	for i := uint(0); i < pbt.iterations; i++ {
		fuzzTest := (&ftesting.FTesting{}).WithFunction(pbt.f)
		inputs, err := fuzzTest.GenerateInputs()
		if err != nil {
			return nil, err
		}
		outs, _ := pbt.applyFunction(inputs...)
		if pbt.haspredicates() {
			switch ret := outs.(type) {
			case []any:
				for _, out := range ret {
					retOut = pbt.validatePredicates(retOut, out)
				}
			case any:
				retOut = pbt.validatePredicates(retOut, ret)
			}
		}
	}
	return retOut, nil
}

// validatePredicates checks if an output value satisfies all configured predicates
// and appends the result to the output slice.
//
// Parameters:
//   - retOut: The accumulating slice of test results
//   - out: The output value to validate
//
// Returns the updated slice with the new test result appended.
//
// This method is called internally by Run for each function output.
func (pbt PBTest) validatePredicates(retOut []PBTestOut, out any) []PBTestOut {
	if ok, failedpredicates := pbt.satisfyAll(out); !ok {
		retOut = append(retOut, PBTestOut{
			Output:     out,
			Predicates: failedpredicates,
			ok:         false,
		})
	} else {
		retOut = append(retOut, PBTestOut{
			Output:     out,
			Predicates: nil,
			ok:         true,
		})
	}
	return retOut
}

// applyFunction executes the test function with the given arguments and returns the result(s).
// This method handles various function signatures using reflection and type conversion.
//
// Parameters:
//   - args: The arguments to pass to the function
//
// Returns:
//   - returnTypes: The function's return value(s), either as a single value or []any for multiple returns
//   - error: An error if the function is invalid or arguments don't match the signature
//
// The method handles:
//   - Functions with no return values (returns nil)
//   - Functions with a single return value
//   - Functions with multiple return values (returned as []any)
//   - Type conversion when arguments are convertible to expected types
//
// This method is called internally by Run for each iteration.
func (pbt *PBTest) applyFunction(args ...any) (returnTypes, error) {
	if pbt.f == nil {
		return nil, nil
	}
	switch fn := pbt.f.(type) {
	case func(any) any:
		return fn(args), nil
	case func(...any) any:
		return fn(args...), nil
	}
	fValue := reflect.ValueOf(pbt.f)
	fType := fValue.Type()

	if fType.Kind() != reflect.Func {
		return nil, &InvalidFunctionProvidedError{pbt.f}
	}
	reflectArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		argValue := reflect.ValueOf(arg)
		expectedType := fType.In(i)
		if argValue.Type() != expectedType {
			if argValue.Type().ConvertibleTo(expectedType) {
				argValue = argValue.Convert(expectedType)
			} else {
				return nil, &InvalidFunctionProvidedError{pbt.f}
			}
		}
		reflectArgs[i] = argValue
	}
	results := fValue.Call(reflectArgs)
	if len(results) == 0 {
		return nil, nil
	} else if len(results) == 1 {
		return results[0].Interface(), nil
	} else {
		retSlice := make([]any, len(results))
		for i, result := range results {
			retSlice[i] = result.Interface()
		}
		return retSlice, nil
	}
}

// satisfyAll checks if a value satisfies all configured predicates.
//
// Parameters:
//   - val: The value to check against predicates
//
// Returns:
//   - ok: true if all predicates pass, false if any fail
//   - failedpredicates: A slice of predicates that failed (nil if all passed)
//
// If no predicates are configured, returns (true, nil).
//
// This method is called internally by validatePredicates.
func (pbt *PBTest) satisfyAll(val any) (ok bool, failedpredicates []p.Predicate) {
	if len(pbt.predicates) == 0 {
		return true, nil
	}
	for _, predicate := range pbt.predicates {
		if !predicate.Verify(val) {
			failedpredicates = append(failedpredicates, predicate)
		}
	}
	if len(failedpredicates) > 0 {
		return false, failedpredicates
	}
	return true, nil
}

// haspredicates checks if any predicates are configured for this test.
//
// Returns true if predicates have been set with WithPredicates, false otherwise.
//
// This method is used internally to determine if predicate validation should be performed.
func (pbt *PBTest) haspredicates() bool {
	return pbt.predicates != nil
}

// FilterPBTTestOut filters a slice of test results to return only the failing cases.
// This is a convenience function for extracting property violations from test results.
//
// Parameters:
//   - in: A slice of PBTestOut results from Run()
//
// Returns a new slice containing only the results where ok is false (i.e., where
// at least one predicate failed).
//
// Example usage:
//
//	results, _ := test.Run()
//	failures := FilterPBTTestOut(results)
//
//	if len(failures) > 0 {
//	    t.Errorf("Found %d property violations:", len(failures))
//	    for i, failure := range failures {
//	        t.Errorf("  Failure %d: Output=%v, Failed predicates=%v",
//	            i+1, failure.Output, failure.Predicates)
//	    }
//	}
func FilterPBTTestOut(in []PBTestOut) []PBTestOut {
	return utils.Filter(in, func(po PBTestOut) bool {
		return !po.ok
	})
}
