// Package predicates provides the Predicate interface and implementations for
// property-based testing validation.
//
// Predicates are boolean-valued functions that check whether a value satisfies
// a particular condition or property. They are the primary mechanism for validating
// function outputs in property-based tests.
//
// In the context of property-based testing, predicates allow you to express
// invariants that should hold true for all function outputs, such as:
//   - Non-negativity: output >= 0
//   - Boundedness: min <= output <= max
//   - Idempotency: f(f(x)) == f(x)
//   - Commutativity: f(a, b) == f(b, a)
//   - Type properties: output is always non-nil, non-empty, etc.
//
// Basic Usage:
//
//	// Define a predicate
//	type NonNegativePredicate struct{}
//
//	func (p NonNegativePredicate) Verify(val any) bool {
//	    if n, ok := val.(int); ok {
//	        return n >= 0
//	    }
//	    return false
//	}
//
//	// Use in a property-based test
//	test := NewPBTest(abs).
//	    WithPredicates(NonNegativePredicate{}).
//	    WithIterations(1000)
//
// Combining Predicates:
//
//	// Multiple predicates can be used together - ALL must pass
//	test := NewPBTest(myFunc).
//	    WithPredicates(
//	        NonNegativePredicate{},
//	        LessThanPredicate{Max: 100},
//	        EvenPredicate{},
//	    )
//
// Custom Predicates:
//
//	// Predicates can be parameterized for reusability
//	type InRangePredicate struct {
//	    Min, Max int
//	}
//
//	func (p InRangePredicate) Verify(val any) bool {
//	    if n, ok := val.(int); ok {
//	        return n >= p.Min && n <= p.Max
//	    }
//	    return false
//	}
//
//	// Use with different ranges
//	test1 := NewPBTest(func1).WithPredicates(InRangePredicate{0, 100})
//	test2 := NewPBTest(func2).WithPredicates(InRangePredicate{-50, 50})
//
// The Predicate interface is intentionally minimal to allow maximum flexibility
// in defining custom validation logic.
package predicates

// Predicate represents a boolean condition that can be checked against a value.
// It is the fundamental building block for property-based testing validation.
//
// A predicate encapsulates a property or invariant that should hold true for
// values produced by a function under test. During property-based testing,
// each function output is validated against all configured predicates.
//
// Methods:
//   - Verify(any) bool: Returns true if the value satisfies the predicate, false otherwise
//
// The Verify method should:
//   - Handle type assertions safely (check types before casting)
//   - Return false for unexpected types rather than panicking
//   - Be deterministic (same input always produces same output)
//   - Be side-effect free (no mutations, I/O, etc.)
//
// Example implementations:
//
//	// Predicate for non-negative numbers
//	type NonNegative struct{}
//	func (p NonNegative) Verify(val any) bool {
//	    switch v := val.(type) {
//	    case int:
//	        return v >= 0
//	    case float64:
//	        return v >= 0.0
//	    default:
//	        return false
//	    }
//	}
//
//	// Predicate for string length
//	type MinLength struct {
//	    Min int
//	}
//	func (p MinLength) Verify(val any) bool {
//	    if s, ok := val.(string); ok {
//	        return len(s) >= p.Min
//	    }
//	    return false
//	}
//
//	// Predicate for slice properties
//	type Sorted struct{}
//	func (p Sorted) Verify(val any) bool {
//	    if slice, ok := val.([]int); ok {
//	        for i := 1; i < len(slice); i++ {
//	            if slice[i] < slice[i-1] {
//	                return false
//	            }
//	        }
//	        return true
//	    }
//	    return false
//	}
type Predicate interface{ Verify(any) bool }
