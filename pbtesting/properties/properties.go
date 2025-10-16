// Package properties defines the core Property interface for property-based testing.
//
// Properties represent invariants or conditions that should hold true for function
// outputs across all valid inputs. This package provides the foundational interface
// that property implementations must satisfy.
//
// The Property interface is similar to the Predicate interface in the predicates
// subpackage, and in practice, predicates are the primary implementation of properties
// used in the pbtesting framework.
//
// Example implementation:
//
//	type NonNegativeProperty struct{}
//
//	func (p NonNegativeProperty) Verify(val any) bool {
//	    if n, ok := val.(int); ok {
//	        return n >= 0
//	    }
//	    return false
//	}
//
// See the predicates subpackage for concrete implementations.
package properties

// Property represents an invariant or condition that should hold true for values.
// It defines a single method for verifying whether a value satisfies the property.
//
// Property is a general interface for any kind of value validation. In the context
// of property-based testing, properties are used to validate function outputs to
// ensure they meet expected invariants regardless of the inputs.
//
// Methods:
//   - Verify(any) bool: Returns true if the value satisfies the property, false otherwise
//
// Example implementations:
//
//	// Property that checks if a number is positive
//	type PositiveProperty struct{}
//	func (p PositiveProperty) Verify(val any) bool {
//	    if n, ok := val.(int); ok {
//	        return n > 0
//	    }
//	    return false
//	}
//
//	// Property that checks if a string is non-empty
//	type NonEmptyString struct{}
//	func (p NonEmptyString) Verify(val any) bool {
//	    if s, ok := val.(string); ok {
//	        return len(s) > 0
//	    }
//	    return false
//	}
type Property interface {
	Verify(any) bool
}
