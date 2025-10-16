package examples

import (
	"testing"

	"github.com/laiambryant/gotestutils/ftesting"
	"github.com/laiambryant/gotestutils/ftesting/attributes"
)

// TestBasicFuzzExample demonstrates basic fuzz testing with default attributes
func TestBasicFuzzExample(t *testing.T) {
	// Test a simple function that should never panic
	ft := &ftesting.FTesting{}
	ft.WithFunction(func(x int, y int) int {
		return x + y
	}).WithIterations(100)

	// Generate and execute with random inputs
	for i := uint(0); i < 100; i++ {
		ok, err := ft.ApplyFunction()
		if err != nil {
			t.Errorf("Iteration %d failed: %v", i, err)
		}
		if !ok {
			t.Errorf("Iteration %d returned false", i)
		}
	}
}

// TestFuzzWithCustomAttributes demonstrates fuzz testing with custom attributes
func TestFuzzWithCustomAttributes(t *testing.T) {
	// Configure custom attributes for controlled input ranges
	attrs := attributes.NewFTAttributes()
	attrs.IntegerAttr = attributes.IntegerAttributesImpl[int]{
		Min:           0,
		Max:           100,
		AllowNegative: false,
		AllowZero:     true,
	}

	// Test division function with controlled inputs
	ft := &ftesting.FTesting{}
	ft.WithFunction(func(a, b int) (int, error) {
		return Divide(a, b)
	}).WithAttributes(attrs).WithIterations(50)

	// Execute and track any errors
	errorCount := 0
	for i := uint(0); i < 50; i++ {
		_, err := ft.ApplyFunction()
		if err != nil {
			// Expected to encounter division by zero sometimes
			errorCount++
		}
	}

	t.Logf("Encountered %d errors during fuzzing (expected due to random zero divisors)", errorCount)
}

// TestFuzzStringFunction demonstrates fuzzing a string processing function
func TestFuzzStringFunction(t *testing.T) {
	attrs := attributes.NewFTAttributes()
	attrs.StringAttr = attributes.StringAttributes{
		MinLen: 1,
		MaxLen: 20,
	}

	ft := &ftesting.FTesting{}
	ft.WithFunction(func(s string) int {
		return len(s)
	}).WithAttributes(attrs).WithIterations(100)

	for i := uint(0); i < 100; i++ {
		ok, err := ft.ApplyFunction()
		if err != nil {
			t.Fatalf("String function failed: %v", err)
		}
		if !ok {
			t.Errorf("String function execution returned false")
		}
	}
}

// TestFuzzMultipleTypes demonstrates fuzzing with multiple different parameter types
func TestFuzzMultipleTypes(t *testing.T) {
	attrs := attributes.NewFTAttributes()
	attrs.IntegerAttr = attributes.IntegerAttributesImpl[int]{
		Min: 1,
		Max: 1000,
	}
	attrs.StringAttr = attributes.StringAttributes{
		MinLen: 3,
		MaxLen: 10,
	}
	attrs.BoolAttr = attributes.BoolAttributes{}

	ft := &ftesting.FTesting{}
	ft.WithFunction(func(id int, name string, active bool) string {
		if active {
			return name
		}
		return ""
	}).WithAttributes(attrs).WithIterations(50)

	for i := uint(0); i < 50; i++ {
		ok, err := ft.ApplyFunction()
		if err != nil {
			t.Fatalf("Multi-type function failed: %v", err)
		}
		if !ok {
			t.Errorf("Multi-type function execution returned false")
		}
	}
}

// TestFuzzFloatOperations demonstrates fuzzing floating-point operations
func TestFuzzFloatOperations(t *testing.T) {
	attrs := attributes.NewFTAttributes()
	attrs.FloatAttr = attributes.FloatAttributesImpl[float64]{
		Min:        -100.0,
		Max:        100.0,
		FiniteOnly: true,
		NonZero:    false,
	}

	ft := &ftesting.FTesting{}
	ft.WithFunction(func(a, b float64) float64 {
		// Test various float operations
		result := a + b
		result = result * 2.0
		if result > 0 {
			result = result / 2.0
		}
		return result
	}).WithAttributes(attrs).WithIterations(100)

	successCount := 0
	for i := uint(0); i < 100; i++ {
		ok, err := ft.ApplyFunction()
		if err != nil {
			t.Logf("Float operation error at iteration %d: %v", i, err)
		} else if ok {
			successCount++
		}
	}

	t.Logf("Successfully executed %d/%d float operations", successCount, 100)
}

// TestGenerateInputsDirectly demonstrates direct input generation without execution
func TestGenerateInputsDirectly(t *testing.T) {
	attrs := attributes.NewFTAttributes()
	attrs.IntegerAttr = attributes.IntegerAttributesImpl[int]{
		Min: 10,
		Max: 20,
	}

	ft := &ftesting.FTesting{}
	ft.WithFunction(func(x, y int) int {
		return x + y
	}).WithAttributes(attrs)

	// Generate inputs without executing
	inputs, err := ft.GenerateInputs()
	if err != nil {
		t.Fatalf("Failed to generate inputs: %v", err)
	}

	if len(inputs) != 2 {
		t.Errorf("Expected 2 inputs, got %d", len(inputs))
	}

	// Verify inputs are in expected range
	for i, input := range inputs {
		if val, ok := input.(int); ok {
			if val < 10 || val > 20 {
				t.Errorf("Input %d value %d is outside expected range [10, 20]", i, val)
			}
			t.Logf("Generated input %d: %v", i, val)
		}
	}
}
