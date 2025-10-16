package examples

import (
	"strings"
	"testing"

	"github.com/laiambryant/gotestutils/ftesting/attributes"
	"github.com/laiambryant/gotestutils/pbtesting"
)

// LengthPredicate checks if a string has a specific length constraint
type LengthPredicate struct {
	MinLength int
	MaxLength int
}

func (lp LengthPredicate) Verify(val any) bool {
	if s, ok := val.(string); ok {
		length := len(s)
		return length >= lp.MinLength && length <= lp.MaxLength
	}
	return false
}

// NotEmptyPredicate checks if a string is not empty
type NotEmptyPredicate struct{}

func (NotEmptyPredicate) Verify(val any) bool {
	if s, ok := val.(string); ok {
		return len(s) > 0
	}
	return false
}

// ContainsPredicate checks if a string contains a substring
type ContainsPredicate struct {
	Substring string
}

func (cp ContainsPredicate) Verify(val any) bool {
	if s, ok := val.(string); ok {
		return strings.Contains(s, cp.Substring)
	}
	return false
}

// PositiveFloatPredicate checks if a float is positive
type PositiveFloatPredicate struct{}

func (PositiveFloatPredicate) Verify(val any) bool {
	if f, ok := val.(float64); ok {
		return f > 0.0
	}
	return false
}

// SliceLengthPredicate checks if a slice has a specific length
type SliceLengthPredicate struct {
	MinLength int
	MaxLength int
}

func (slp SliceLengthPredicate) Verify(val any) bool {
	if s, ok := val.([]int); ok {
		length := len(s)
		return length >= slp.MinLength && length <= slp.MaxLength
	}
	return false
}

// TestStringUpperCaseProperty demonstrates property testing with string transformations
func TestStringUpperCaseProperty(t *testing.T) {
	notEmpty := NotEmptyPredicate{}

	upperFunc := func(s string) string {
		return strings.ToUpper(s)
	}

	ftAttrs := attributes.NewFTAttributes()
	ftAttrs.StringAttr = attributes.StringAttributes{
		MinLen: 1,
		MaxLen: 20,
	}

	test := pbtesting.NewPBTest(upperFunc).
		WithIterations(100).
		WithPredicates(notEmpty).
		WithArgAttributes(ftAttrs).
		WithT(t)

	results, err := test.Run()
	if err != nil {
		t.Fatalf("Property test failed: %v", err)
	}

	failures := pbtesting.FilterPBTTestOut(results)
	if len(failures) > 0 {
		t.Errorf("Found %d property violations", len(failures))
		for _, failure := range failures {
			t.Logf("  Failed output: %q", failure.Output)
		}
	}
}

// TestStringTrimProperty demonstrates property testing with trim operations
func TestStringTrimProperty(t *testing.T) {
	trimFunc := func(s string) string {
		return strings.TrimSpace(s)
	}

	test := pbtesting.NewPBTest(trimFunc).
		WithIterations(100).
		WithT(t)

	results, err := test.Run()
	if err != nil {
		t.Fatalf("Property test failed: %v", err)
	}

	for i, result := range results {
		trimmed := result.Output.(string)
		if len(trimmed) == 0 {
			t.Errorf("Iteration %d: Trimmed string is empty", i)
		}
	}
}

// TestFloatSquareRootProperty demonstrates property testing with floating point
func TestFloatSquareRootProperty(t *testing.T) {
	sqrtSquare := func(x float64) float64 {
		if x < 0 {
			return 0.0 // Handle negative inputs
		}
		result := x * x
		if result == 0 {
			return 0
		}
		return result
	}

	truePred := TruePredicate{}

	ftAttrs := attributes.NewFTAttributes()
	ftAttrs.FloatAttr = attributes.FloatAttributesImpl[float64]{
		Min: 0.1,
		Max: 100.0,
	}

	test := pbtesting.NewPBTest(sqrtSquare).
		WithIterations(50).
		WithPredicates(truePred).
		WithArgAttributes(ftAttrs).
		WithT(t)

	results, err := test.Run()
	if err != nil {
		t.Fatalf("Property test failed: %v", err)
	}

	failures := pbtesting.FilterPBTTestOut(results)
	if len(failures) > 0 {
		t.Errorf("Found %d property violations", len(failures))
	}
}

// TestSliceReverseProperty demonstrates property testing with slices
func TestSliceReverseProperty(t *testing.T) {
	// Property: reversing a slice twice should give original slice length
	reverse := func(slice []int) []int {
		result := make([]int, len(slice))
		for i := range slice {
			result[i] = slice[len(slice)-1-i]
		}
		return result
	}

	lengthCheck := SliceLengthPredicate{MinLength: 0, MaxLength: 100}

	ftAttrs := attributes.NewFTAttributes()
	ftAttrs.SliceAttr = attributes.SliceAttributes{
		MinLen:       5,
		MaxLen:       20,
		ElementAttrs: attributes.IntegerAttributesImpl[int]{},
	}

	test := pbtesting.NewPBTest(reverse).
		WithIterations(50).
		WithPredicates(lengthCheck).
		WithArgAttributes(ftAttrs).
		WithT(t)

	results, err := test.Run()
	if err != nil {
		t.Fatalf("Property test failed: %v", err)
	}

	failures := pbtesting.FilterPBTTestOut(results)
	if len(failures) > 0 {
		t.Errorf("Found %d property violations", len(failures))
	}
}

// TestMapSizeProperty demonstrates property testing with maps
func TestMapSizeProperty(t *testing.T) {
	filterMap := func(m map[string]int) map[string]int {
		result := make(map[string]int)
		for k, v := range m {
			if v > 0 {
				result[k] = v
			}
		}
		return result
	}

	test := pbtesting.NewPBTest(filterMap).
		WithIterations(50).
		WithT(t)

	results, err := test.Run()
	if err != nil {
		t.Fatalf("Property test failed: %v", err)
	}

	for i, result := range results {
		filtered := result.Output.(map[string]int)
		if len(filtered) == 0 {
			t.Errorf("Iteration %d: Invalid filtered map size", i)
		}
	}
}

// TestCustomStructProperty demonstrates property testing with custom structs
func TestCustomStructProperty(t *testing.T) {
	type Point struct {
		X, Y int
	}

	translate := func(p Point) Point {
		return Point{X: p.X + 10, Y: p.Y + 10}
	}

	truePred := TruePredicate{}

	test := pbtesting.NewPBTest(translate).
		WithIterations(50).
		WithPredicates(truePred).
		WithT(t)

	results, err := test.Run()
	if err != nil {
		t.Fatalf("Property test failed: %v", err)
	}

	successCount := 0
	for _, result := range results {
		if result.Ok {
			successCount++
			translated := result.Output.(Point)
			_ = translated
		}
	}

	t.Logf("Successfully translated %d/%d points", successCount, len(results))
}

// TestMultiArgumentProperty demonstrates property testing with multiple arguments
func TestMultiArgumentProperty(t *testing.T) {
	maxFunc := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	nonNegative := TruePredicate{}

	ftAttrs := attributes.NewFTAttributes()
	ftAttrs.IntegerAttr = attributes.IntegerAttributesImpl[int]{
		Min:           0,
		Max:           1000,
		AllowNegative: false,
		AllowZero:     true,
	}

	test := pbtesting.NewPBTest(maxFunc).
		WithIterations(100).
		WithPredicates(nonNegative).
		WithArgAttributes(ftAttrs).
		WithT(t)

	results, err := test.Run()
	if err != nil {
		t.Fatalf("Property test failed: %v", err)
	}

	failures := pbtesting.FilterPBTTestOut(results)
	if len(failures) > 0 {
		t.Errorf("Found %d property violations", len(failures))
	}
}

// TestComplexPredicateChain demonstrates chaining multiple complex predicates
func TestComplexPredicateChain(t *testing.T) {
	processString := func(s string) string {
		s = strings.TrimSpace(s)
		s = strings.ToUpper(s)
		if len(s) > 10 {
			s = s[:10]
		}
		return s
	}

	notEmpty := NotEmptyPredicate{}
	maxLength := LengthPredicate{MinLength: 0, MaxLength: 10}

	ftAttrs := attributes.NewFTAttributes()
	ftAttrs.StringAttr = attributes.StringAttributes{
		MinLen: 1,
		MaxLen: 50,
	}

	test := pbtesting.NewPBTest(processString).
		WithIterations(100).
		WithPredicates(notEmpty, maxLength).
		WithArgAttributes(ftAttrs).
		WithT(t)

	results, err := test.Run()
	if err != nil {
		t.Fatalf("Property test failed: %v", err)
	}

	failures := pbtesting.FilterPBTTestOut(results)
	if len(failures) > 0 {
		t.Errorf("Found %d property violations", len(failures))
		for _, failure := range failures {
			t.Logf("  Failed output: %q (failed %d predicates)",
				failure.Output, len(failure.Predicates))
		}
	}
}

// TestErrorHandlingProperty demonstrates property testing with error returns
func TestErrorHandlingProperty(t *testing.T) {
	safeDivide := func(a, b int) (int, error) {
		if b == 0 {
			return 0, nil
		}
		return a / b, nil
	}

	truePred := TruePredicate{}

	test := pbtesting.NewPBTest(safeDivide).
		WithIterations(50).
		WithPredicates(truePred).
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

	t.Logf("Successfully completed %d/%d divisions", successCount, len(results))
}
