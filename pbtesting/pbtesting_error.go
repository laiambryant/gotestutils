package pbtesting

import (
	"errors"
	"fmt"

	p "github.com/laiambryant/gotestutils/pbtesting/properties/predicates"
)

var (
	// ErrZeroRangeNonZeroRequested is returned when attempting to generate non-zero values
	// from a range that only contains zero (e.g., Min=0, Max=0, NonZero=true).
	ErrZeroRangeNonZeroRequested = errors.New("zero range but non-zero requested")

	// ErrMinGreaterThanMax is returned when a minimum value is configured to be greater
	// than the maximum value, making it impossible to generate valid values in the range.
	ErrMinGreaterThanMax = errors.New("minimum is greater than maximum")
)

// InvalidPropertyError is returned when a predicate fails to validate properly
// or is malformed. This typically indicates a bug in the predicate implementation.
//
// Fields:
//   - predicate: The predicate that caused the error
//
// Example scenario:
//
//	predicate := &BrokenPredicate{} // A predicate with implementation issues
//	// Using this predicate might result in InvalidPropertyError
type InvalidPropertyError struct {
	predicate p.Predicate
}

func (i InvalidPropertyError) Error() string {
	return fmt.Sprintf("invalid property: %v", i.predicate)
}

// FunctionNotProvidedError is returned when attempting to run a property-based test
// without setting a function using NewPBTest or WithF.
//
// Example scenario:
//
//	test := &PBTest{} // No function set
//	_, err := test.Run() // Returns FunctionNotProvidedError (or empty results)
type FunctionNotProvidedError struct{}

func (fnp FunctionNotProvidedError) Error() string {
	return "a function must be provided for the property-based test suite to work"
}

// InvalidFunctionProvidedError is returned when the provided function has an invalid
// signature, cannot be called with the generated arguments, or is not a function at all.
//
// Fields:
//   - f: The invalid function value that was provided
//
// This error occurs when:
//   - The provided value is not a function (wrong type)
//   - Generated arguments cannot be converted to the function's parameter types
//   - The function signature is incompatible with property-based testing
//
// Example scenario:
//
//	test := NewPBTest(42) // Not a function
//	_, err := test.Run() // Returns InvalidFunctionProvidedError{f: 42}
//
//	test2 := NewPBTest(func(x CustomUnconvertibleType) int { return 0 })
//	_, err2 := test2.Run() // May return InvalidFunctionProvidedError if arguments can't be converted
type InvalidFunctionProvidedError struct {
	f any
}

func (ifp InvalidFunctionProvidedError) Error() string {
	return fmt.Sprintf("Invalid function provided to pbt, function: [%v]", ifp.f)
}
