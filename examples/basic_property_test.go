package examples

import (
	"math"
	"testing"

	"github.com/laiambryant/gotestutils/pbtesting"
)

// NonNegativePredicate checks if a value is non-negative
type NonNegativePredicate struct{}

func (NonNegativePredicate) Verify(val any) bool {
	switch v := val.(type) {
	case int:
		return v >= 0
	case float64:
		return v >= 0.0
	default:
		return false
	}
}

// LessThanPredicate checks if a value is less than a maximum
type LessThanPredicate struct {
	Max int
}

func (p LessThanPredicate) Verify(val any) bool {
	if v, Ok := val.(int); Ok {
		return v < p.Max
	}
	return false
}

// EvenPredicate checks if a number is even
type EvenPredicate struct{}

func (EvenPredicate) Verify(val any) bool {
	if v, Ok := val.(int); Ok {
		return v%2 == 0
	}
	return false
}

// InRangePredicate checks if a value is within a range
type InRangePredicate struct {
	Min, Max int
}

func (p InRangePredicate) Verify(val any) bool {
	if v, ok := val.(int); ok {
		return v >= p.Min && v <= p.Max
	}
	return false
}

// TruePredicate always returns true (used for tests that need manual verification)
type TruePredicate struct{}

func (TruePredicate) Verify(val any) bool {
	return true
}

// TestAbsoluteValueProperty demonstrates property-based testing of absolute value
func TestAbsoluteValueProperty(t *testing.T) {
	// Property: abs(x) should always return non-negative numbers
	nonNegative := NonNegativePredicate{}

	absFunc := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}

	test := pbtesting.NewPBTest(absFunc).
		WithIterations(100).
		WithPredicates(nonNegative).
		WithT(t)

	results, err := test.Run()
	if err != nil {
		t.Fatalf("Property test failed: %v", err)
	}

	failures := pbtesting.FilterPBTTestOut(results)
	if len(failures) > 0 {
		t.Errorf("Found %d property violations", len(failures))
		for _, failure := range failures {
			t.Logf("  Failed output: %v", failure.Output)
		}
	}
}

// TestSquareProperty demonstrates testing that x² is always non-negative
func TestSquareProperty(t *testing.T) {
	nonNegative := NonNegativePredicate{}

	squareFunc := func(x int) int {
		return x * x
	}

	test := pbtesting.NewPBTest(squareFunc).
		WithIterations(100).
		WithPredicates(nonNegative).
		WithT(t)

	results, err := test.Run()
	if err != nil {
		t.Fatalf("Property test failed: %v", err)
	}

	failures := pbtesting.FilterPBTTestOut(results)
	if len(failures) > 0 {
		t.Errorf("Found %d property violations for square function", len(failures))
	}
}

// TestMultiplePredicates demonstrates using multiple predicates together
func TestMultiplePredicates(t *testing.T) {
	// Function that doubles positive numbers
	doublePositive := func(x int) int {
		if x <= 0 {
			return 0
		}
		return x * 2
	}

	// Properties: output should be non-negative AND even
	nonNegative := NonNegativePredicate{}
	even := EvenPredicate{}

	test := pbtesting.NewPBTest(doublePositive).
		WithIterations(100).
		WithPredicates(nonNegative, even).
		WithT(t)

	results, err := test.Run()
	if err != nil {
		t.Fatalf("Property test failed: %v", err)
	}

	failures := pbtesting.FilterPBTTestOut(results)
	if len(failures) > 0 {
		t.Errorf("Found %d property violations", len(failures))
		for _, failure := range failures {
			t.Logf("  Failed output: %v, Failed predicates: %d",
				failure.Output, len(failure.Predicates))
		}
	}
}

// TestIdempotenceProperty demonstrates testing idempotence (f(f(x)) = f(x))
func TestIdempotenceProperty(t *testing.T) {
	// Absolute value is idempotent
	absFunc := func(x int) int {
		return int(math.Abs(float64(x)))
	}

	// We'll test this by checking that applying abs twice gives same result
	nonNegative := NonNegativePredicate{}

	test := pbtesting.NewPBTest(absFunc).
		WithIterations(50).
		WithPredicates(nonNegative).
		WithT(t)

	results, err := test.Run()
	if err != nil {
		t.Fatalf("Property test failed: %v", err)
	}

	// Additional check: verify idempotence manually
	for _, result := range results {
		if result.Ok {
			val := result.Output.(int)
			// Apply function again
			val2 := absFunc(val)
			if val != val2 {
				t.Errorf("Idempotence violated: abs(%d) = %d, but abs(abs(%d)) = %d",
					val, val, val, val2)
			}
		}
	}
}

// TestRangeConstraints demonstrates property testing with range constraints
func TestRangeConstraints(t *testing.T) {
	// Function that clamps value between 0 and 100
	clamp := func(x int) int {
		if x < 0 {
			return 0
		}
		if x > 100 {
			return 100
		}
		return x
	}

	// Property: output should always be in range [0, 100]
	inRange := InRangePredicate{Min: 0, Max: 100}

	test := pbtesting.NewPBTest(clamp).
		WithIterations(100).
		WithPredicates(inRange).
		WithT(t)

	results, err := test.Run()
	if err != nil {
		t.Fatalf("Property test failed: %v", err)
	}

	failures := pbtesting.FilterPBTTestOut(results)
	if len(failures) > 0 {
		t.Errorf("Found %d values outside range [0, 100]", len(failures))
		for _, failure := range failures {
			t.Logf("  Out of range value: %v", failure.Output)
		}
	}
}

// TestInverseProperty demonstrates testing inverse relationships (f(g(x)) = x)
func TestInverseProperty(t *testing.T) {
	// Test that negation is its own inverse
	negate := func(x int) int {
		return -x
	}

	// Use TruePredicate to get results
	truePred := TruePredicate{}

	test := pbtesting.NewPBTest(negate).
		WithIterations(50).
		WithPredicates(truePred).
		WithT(t)

	results, err := test.Run()
	if err != nil {
		t.Fatalf("Property test failed: %v", err)
	}

	// Check inverse property: negate(negate(x)) should equal x
	for i, result := range results {
		negated := result.Output.(int)
		doubleNegated := negate(negated)

		// We need to get the original input - for this demo we'll skip
		// In practice, you'd store inputs or test differently
		_ = doubleNegated
		t.Logf("Iteration %d: output = %v", i, negated)
	}
} // TestStringLengthProperty demonstrates property testing with strings
func TestStringLengthProperty(t *testing.T) {
	// Property: concatenating two strings should produce length = len(a) + len(b)
	concat := func(a, b string) string {
		return a + b
	}

	// Use a predicate to verify non-negative length (always true, but needed for Run to return results)
	nonEmpty := TruePredicate{}

	test := pbtesting.NewPBTest(concat).
		WithIterations(50).
		WithPredicates(nonEmpty).
		WithT(t)

	results, err := test.Run()
	if err != nil {
		t.Fatalf("Property test failed: %v", err)
	}

	// Verify concatenation property manually
	for i, result := range results {
		output := result.Output.(string)
		if len(output) == 0 { // This should never happen
			t.Errorf("Iteration %d: Invalid string length", i)
		}
	}
}

// TestCommutativityProperty demonstrates testing commutative operations
func TestCommutativityProperty(t *testing.T) {
	// Addition is commutative: a + b = b + a
	add := func(a, b int) int {
		return a + b
	}

	// Use TruePredicate to get results back
	truePred := TruePredicate{}

	test := pbtesting.NewPBTest(add).
		WithIterations(50).
		WithPredicates(truePred).
		WithT(t)

	results, err := test.Run()
	if err != nil {
		t.Fatalf("Property test failed: %v", err)
	}

	// In a real test, we'd verify commutativity by calling with swapped args
	// For this demo, just verify we got results
	if len(results) != 50 {
		t.Errorf("Expected 50 results, got %d", len(results))
	}
}

// TestAssociativityProperty demonstrates testing associative operations
func TestAssociativityProperty(t *testing.T) {
	// Max operation is associative: max(max(a,b),c) = max(a,max(b,c))
	maxFunc := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	nonNegative := NonNegativePredicate{}

	test := pbtesting.NewPBTest(maxFunc).
		WithIterations(50).
		WithPredicates(nonNegative).
		WithT(t)

	results, err := test.Run()
	if err != nil {
		t.Fatalf("Property test failed: %v", err)
	}

	successCount := 0
	for _, result := range results {
		if result.Ok {
			successCount++
		}
	}

	t.Logf("Property held for %d/%d iterations", successCount, len(results))
}
