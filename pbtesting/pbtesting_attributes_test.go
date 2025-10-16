package pbtesting

import (
	"testing"

	"github.com/laiambryant/gotestutils/ftesting/attributes"
)

// mockPredicateForAttrTest is a simple predicate for testing attribute functionality
type mockPredicateForAttrTest struct {
	minValue int
	maxValue int
}

func (m mockPredicateForAttrTest) Verify(val any) bool {
	if v, ok := val.(int); ok {
		return v >= m.minValue && v <= m.maxValue
	}
	return false
}

// TestRunWithAttributes_PassesAttributesToFTesting verifies that custom attributes
// are properly passed to the ftesting framework when using RunWithAttributes
func TestRunWithAttributes_PassesAttributesToFTesting(t *testing.T) {
	// Function that returns its input
	identityFunc := func(x int) int {
		return x
	}

	// Create custom attributes with constrained range
	attrs := attributes.NewFTAttributes()
	attrs.IntegerAttr = attributes.IntegerAttributesImpl[int]{
		Min:           10,
		Max:           20,
		AllowNegative: false,
		AllowZero:     false,
	}

	// Predicate to verify values are in expected range
	inRangePred := mockPredicateForAttrTest{minValue: 10, maxValue: 20}

	test := NewPBTest(identityFunc).
		WithIterations(50).
		WithPredicates(inRangePred)

	// Run with custom attributes
	results, err := test.RunWithAttributes(attrs)
	if err != nil {
		t.Fatalf("RunWithAttributes failed: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("Expected results but got empty slice")
	}

	// All results should pass the predicate if attributes were properly applied
	failures := FilterPBTTestOut(results)
	if len(failures) > 0 {
		t.Errorf("Expected all values to be in range [10, 20], but got %d failures", len(failures))
		for i, failure := range failures {
			if i < 5 { // Show first 5 failures
				t.Logf("  Failure %d: value = %v", i+1, failure.Output)
			}
		}
	}

	// Verify all outputs are in the expected range
	for i, result := range results {
		value := result.Output.(int)
		if value < 10 || value > 20 {
			t.Errorf("Iteration %d: value %d is outside expected range [10, 20]", i, value)
		}
	}
}

// TestRunWithAttributes_NilAttributes verifies that passing nil attributes
// uses default attributes (not causing errors)
func TestRunWithAttributes_NilAttributes(t *testing.T) {
	squareFunc := func(x int) int {
		return x * x
	}

	// Predicate that always passes (for this test we just want to verify execution)
	alwaysPass := mockPredicateForAttrTest{minValue: -1000000, maxValue: 1000000}

	test := NewPBTest(squareFunc).
		WithIterations(10).
		WithPredicates(alwaysPass)

	// Run with nil attributes (should use defaults)
	results, err := test.RunWithAttributes(nil)
	if err != nil {
		t.Fatalf("RunWithAttributes with nil failed: %v", err)
	}

	if len(results) != 10 {
		t.Errorf("Expected 10 results, got %d", len(results))
	}
}

// TestRun_UsesDefaultAttributes verifies that Run() (without explicit attributes)
// uses default attributes via RunWithAttributes(nil)
func TestRun_UsesDefaultAttributes(t *testing.T) {
	addFunc := func(a, b int) int {
		return a + b
	}

	test := NewPBTest(addFunc).
		WithIterations(20).
		WithT(t)

	// Run without attributes (should use defaults)
	_, err := test.Run()
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// Since we have no predicates, results will be empty
	// But the function should execute without errors
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// TestRunWithAttributes_FloatAttributes verifies attribute passing for float types
func TestRunWithAttributes_FloatAttributes(t *testing.T) {
	sqrtApprox := func(x float64) float64 {
		if x < 0 {
			return 0
		}
		return x
	}

	// Create attributes for positive floats only
	attrs := attributes.NewFTAttributes()
	attrs.FloatAttr = attributes.FloatAttributesImpl[float64]{
		Min:        0.0,
		Max:        100.0,
		NonZero:    false,
		FiniteOnly: true,
	}

	// Simple predicate: result should be non-negative
	nonNegPred := mockPredicateForFloatTest{minValue: 0.0}

	test := NewPBTest(sqrtApprox).
		WithIterations(30).
		WithPredicates(nonNegPred)

	results, err := test.RunWithAttributes(attrs)
	if err != nil {
		t.Fatalf("RunWithAttributes for float failed: %v", err)
	}

	failures := FilterPBTTestOut(results)
	if len(failures) > 0 {
		t.Errorf("Expected all values to be non-negative, got %d failures", len(failures))
	}

	// Verify inputs were in expected range (by checking outputs)
	for _, result := range results {
		value := result.Output.(float64)
		if value < 0.0 || value > 100.0 {
			t.Errorf("Value %f is outside expected range [0.0, 100.0]", value)
		}
	}
}

// mockPredicateForFloatTest is a predicate for float testing
type mockPredicateForFloatTest struct {
	minValue float64
}

func (m mockPredicateForFloatTest) Verify(val any) bool {
	if v, ok := val.(float64); ok {
		return v >= m.minValue
	}
	return false
}

// TestRunWithAttributes_StringAttributes verifies attribute passing for string types
func TestRunWithAttributes_StringAttributes(t *testing.T) {
	upperFunc := func(s string) string {
		return s
	}

	// Create attributes for strings with specific length
	attrs := attributes.NewFTAttributes()
	attrs.StringAttr = attributes.StringAttributes{
		MinLen: 5,
		MaxLen: 10,
	}

	// Predicate to check string length
	lengthPred := mockPredicateForStringTest{minLen: 5, maxLen: 10}

	test := NewPBTest(upperFunc).
		WithIterations(25).
		WithPredicates(lengthPred)

	results, err := test.RunWithAttributes(attrs)
	if err != nil {
		t.Fatalf("RunWithAttributes for string failed: %v", err)
	}

	failures := FilterPBTTestOut(results)
	if len(failures) > 0 {
		t.Errorf("Expected all strings to be length [5, 10], got %d failures", len(failures))
		for i, failure := range failures {
			if i < 3 {
				t.Logf("  Failure %d: string = %q (len=%d)", i+1, failure.Output, len(failure.Output.(string)))
			}
		}
	}
}

// mockPredicateForStringTest is a predicate for string testing
type mockPredicateForStringTest struct {
	minLen int
	maxLen int
}

func (m mockPredicateForStringTest) Verify(val any) bool {
	if v, ok := val.(string); ok {
		length := len(v)
		return length >= m.minLen && length <= m.maxLen
	}
	return false
}

// TestRunWithAttributes_MultipleParameterFunction verifies attributes work with
// functions that have multiple parameters
func TestRunWithAttributes_MultipleParameterFunction(t *testing.T) {
	maxFunc := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	// Constrain inputs to positive range
	attrs := attributes.NewFTAttributes()
	attrs.IntegerAttr = attributes.IntegerAttributesImpl[int]{
		Min:           1,
		Max:           50,
		AllowNegative: false,
		AllowZero:     false,
	}

	// Result should always be positive
	positivePred := mockPredicateForAttrTest{minValue: 1, maxValue: 50}

	test := NewPBTest(maxFunc).
		WithIterations(40).
		WithPredicates(positivePred)

	results, err := test.RunWithAttributes(attrs)
	if err != nil {
		t.Fatalf("RunWithAttributes for multi-param failed: %v", err)
	}

	failures := FilterPBTTestOut(results)
	if len(failures) > 0 {
		t.Errorf("Expected all max values to be in range [1, 50], got %d failures", len(failures))
	}
}

// TestRunWithAttributes_WithErrorReturningFunction verifies attributes work
// with functions that return errors
func TestRunWithAttributes_WithErrorReturningFunction(t *testing.T) {
	divideFunc := func(a, b int) (int, error) {
		if b == 0 {
			return 0, nil // For simplicity, just return 0
		}
		return a / b, nil
	}

	attrs := attributes.NewFTAttributes()
	attrs.IntegerAttr = attributes.IntegerAttributesImpl[int]{
		Min:           -10,
		Max:           10,
		AllowNegative: true,
		AllowZero:     true,
	}

	// Add a predicate so results are populated
	inRangePred := mockPredicateForAttrTest{minValue: -100, maxValue: 100}

	test := NewPBTest(divideFunc).
		WithIterations(20).
		WithPredicates(inRangePred)

	results, err := test.RunWithAttributes(attrs)
	if err != nil {
		t.Fatalf("RunWithAttributes with error return failed: %v", err)
	}

	if len(results) == 0 {
		t.Error("Expected results but got empty slice")
	}
}

// TestWithArgAttributes_BackwardsCompatibility verifies that the old WithArgAttributes
// method still works (even though it doesn't actually pass attributes yet)
func TestWithArgAttributes_BackwardsCompatibility(t *testing.T) {
	simpleFunc := func(x int) int {
		return x * 2
	}

	attrs := attributes.NewFTAttributes()
	attrs.IntegerAttr = attributes.IntegerAttributesImpl[int]{
		Min: 0,
		Max: 50,
	}

	test := NewPBTest(simpleFunc).
		WithIterations(10).
		WithArgAttributes(attrs) // This sets argAttrs but doesn't use them yet

	// Should still run without errors
	_, err := test.Run()
	if err != nil {
		t.Fatalf("Run with WithArgAttributes failed: %v", err)
	}

	// Results will be empty without predicates, but should execute
	if err != nil {
		t.Error("Expected successful execution")
	}
}

// TestRunWithAttributes_NilFunction verifies proper handling when function is nil
func TestRunWithAttributes_NilFunction(t *testing.T) {
	attrs := attributes.NewFTAttributes()

	test := NewPBTest(nil).WithIterations(5)

	results, err := test.RunWithAttributes(attrs)
	if err != nil {
		t.Errorf("Expected no error with nil function, got %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected empty results with nil function, got %d results", len(results))
	}
}

// TestRunWithAttributes_ZeroIterations verifies behavior with zero iterations
func TestRunWithAttributes_ZeroIterations(t *testing.T) {
	simpleFunc := func(x int) int {
		return x
	}

	attrs := attributes.NewFTAttributes()

	test := NewPBTest(simpleFunc).WithIterations(0)

	results, err := test.RunWithAttributes(attrs)
	if err != nil {
		t.Errorf("Expected no error with 0 iterations, got %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results with 0 iterations, got %d", len(results))
	}
}
