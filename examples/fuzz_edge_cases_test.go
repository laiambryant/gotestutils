package examples

import (
	"fmt"
	"strings"
	"testing"

	"github.com/laiambryant/gotestutils/ftesting"
	"github.com/laiambryant/gotestutils/ftesting/attributes"
)

// TestFuzzErrorDetection demonstrates using fuzzing to find edge cases that cause errors
func TestFuzzErrorDetection(t *testing.T) {
	// Function that has a hidden bug with certain inputs
	buggyFunction := func(s string) (int, error) {
		if strings.Contains(s, "ERROR") {
			return 0, fmt.Errorf("invalid input")
		}
		// Bug: crashes on very long strings
		if len(s) > 50 {
			panic("string too long")
		}
		return len(s), nil
	}

	attrs := attributes.NewFTAttributes()
	attrs.StringAttr = attributes.StringAttributes{
		MinLen: 1,
		MaxLen: 100, // Will occasionally generate strings > 50
	}

	ft := &ftesting.FTesting{}
	ft.WithFunction(buggyFunction).WithAttributes(attrs).WithIterations(100)

	errorCount := 0
	panicCount := 0

	for i := uint(0); i < 100; i++ {
		// Catch panics
		func() {
			defer func() {
				if r := recover(); r != nil {
					panicCount++
					t.Logf("Caught panic at iteration %d: %v", i, r)
				}
			}()

			ok, err := ft.ApplyFunction()
			if err != nil {
				errorCount++
			}
			_ = ok
		}()
	}

	t.Logf("Found %d errors and %d panics during fuzzing", errorCount, panicCount)
	if panicCount > 0 {
		t.Logf("Successfully discovered the string length bug!")
	}
}

// TestFuzzBoundaryConditions demonstrates finding boundary condition bugs
func TestFuzzBoundaryConditions(t *testing.T) {
	// Function with off-by-one error
	indexFunction := func(arr []int, idx int) (int, error) {
		if idx < 0 || idx >= len(arr) {
			return 0, fmt.Errorf("index out of bounds")
		}
		// Bug: should check idx < len(arr) not <=
		if idx <= len(arr) {
			return arr[idx], nil
		}
		return 0, fmt.Errorf("unreachable")
	}

	attrs := attributes.NewFTAttributes()
	attrs.SliceAttr = attributes.SliceAttributes{
		MinLen: 1,
		MaxLen: 10,
		ElementAttrs: attributes.IntegerAttributesImpl[int]{
			Min: 0,
			Max: 100,
		},
	}
	attrs.IntegerAttr = attributes.IntegerAttributesImpl[int]{
		Min:           -5,
		Max:           15,
		AllowNegative: true,
	}

	ft := &ftesting.FTesting{}
	ft.WithFunction(indexFunction).WithAttributes(attrs).WithIterations(200)

	panicCount := 0
	errorCount := 0
	successCount := 0

	for i := uint(0); i < 200; i++ {
		func() {
			defer func() {
				if r := recover(); r != nil {
					panicCount++
				}
			}()

			ok, err := ft.ApplyFunction()
			if err != nil {
				errorCount++
			} else if ok {
				successCount++
			}
		}()
	}

	t.Logf("Results: %d successes, %d errors, %d panics", successCount, errorCount, panicCount)
	if panicCount > 0 {
		t.Logf("Successfully discovered the boundary condition bug!")
	}
}

// TestFuzzTypeConversion demonstrates finding type conversion issues
func TestFuzzTypeConversion(t *testing.T) {
	// Function that performs risky type conversions
	conversionFunc := func(val int64) int {
		// Potential overflow when converting to int on 32-bit systems
		return int(val)
	}

	attrs := attributes.NewFTAttributes()
	attrs.IntegerAttr = attributes.IntegerAttributesImpl[int64]{
		Min: -1000000000,
		Max: 1000000000,
	}

	ft := &ftesting.FTesting{}
	ft.WithFunction(conversionFunc).WithAttributes(attrs).WithIterations(100)

	for i := uint(0); i < 100; i++ {
		inputs, err := ft.GenerateInputs()
		if err != nil {
			t.Fatalf("Failed to generate inputs: %v", err)
		}

		val := inputs[0].(int64)

		ok, err := ft.ApplyFunction()
		if err != nil {
			t.Errorf("Conversion function failed with input %d: %v", val, err)
		}
		if !ok {
			t.Errorf("Conversion function returned false")
		}
	}
}

// TestFuzzConcurrencySafety demonstrates testing for race conditions
func TestFuzzConcurrencySafety(t *testing.T) {
	// Shared state that might have race conditions
	var counter int

	unsafeIncrement := func(n int) int {
		// This is intentionally unsafe for concurrent access
		counter += n
		return counter
	}

	attrs := attributes.NewFTAttributes()
	attrs.IntegerAttr = attributes.IntegerAttributesImpl[int]{
		Min: 1,
		Max: 10,
	}

	ft := &ftesting.FTesting{}
	ft.WithFunction(unsafeIncrement).WithAttributes(attrs)

	// Run sequentially first
	counter = 0
	for i := uint(0); i < 20; i++ {
		_, err := ft.ApplyFunction()
		if err != nil {
			t.Fatalf("Function failed: %v", err)
		}
	}
	sequentialResult := counter

	t.Logf("Sequential execution result: %d", sequentialResult)
	// Note: For actual race detection, use `go test -race`
}

// TestFuzzDataValidation demonstrates using fuzzing for data validation logic
func TestFuzzDataValidation(t *testing.T) {
	// Email validation function (simplified)
	validateEmail := func(email string) bool {
		return strings.Contains(email, "@") &&
			strings.Contains(email, ".") &&
			len(email) > 5
	}

	attrs := attributes.NewFTAttributes()
	attrs.StringAttr = attributes.StringAttributes{
		MinLen: 1,
		MaxLen: 50,
	}

	ft := &ftesting.FTesting{}
	ft.WithFunction(validateEmail).WithAttributes(attrs).WithIterations(100)

	validCount := 0
	invalidCount := 0

	for i := uint(0); i < 100; i++ {
		inputs, err := ft.GenerateInputs()
		if err != nil {
			t.Fatalf("Failed to generate inputs: %v", err)
		}

		email := inputs[0].(string)
		isValid := validateEmail(email)

		if isValid {
			validCount++
			t.Logf("Valid email generated: %s", email)
		} else {
			invalidCount++
		}

		ok, err := ft.ApplyFunction()
		if err != nil {
			t.Fatalf("Validation function failed: %v", err)
		}
		_ = ok
	}

	t.Logf("Generated %d valid and %d invalid emails", validCount, invalidCount)
}

// TestFuzzPerformanceCharacteristics demonstrates finding performance issues
func TestFuzzPerformanceCharacteristics(t *testing.T) {
	// Function with O(n²) complexity
	inefficientSort := func(nums []int) []int {
		// Bubble sort - inefficient for large arrays
		result := make([]int, len(nums))
		copy(result, nums)

		for i := 0; i < len(result); i++ {
			for j := 0; j < len(result)-1; j++ {
				if result[j] > result[j+1] {
					result[j], result[j+1] = result[j+1], result[j]
				}
			}
		}
		return result
	}

	attrs := attributes.NewFTAttributes()
	attrs.SliceAttr = attributes.SliceAttributes{
		MinLen: 1,
		MaxLen: 100, // Try with larger arrays to see performance impact
		ElementAttrs: attributes.IntegerAttributesImpl[int]{
			Min: 0,
			Max: 1000,
		},
	}

	ft := &ftesting.FTesting{}
	ft.WithFunction(inefficientSort).WithAttributes(attrs).WithIterations(20)

	for i := uint(0); i < 20; i++ {
		inputs, err := ft.GenerateInputs()
		if err != nil {
			t.Fatalf("Failed to generate inputs: %v", err)
		}

		arraySize := len(inputs[0].([]int))

		ok, err := ft.ApplyFunction()
		if err != nil {
			t.Fatalf("Sort function failed: %v", err)
		}
		if !ok {
			t.Errorf("Sort function returned false")
		}

		t.Logf("Sorted array of size %d", arraySize)
	}
}

// TestFuzzNilHandling demonstrates testing nil handling with pointers and slices
func TestFuzzNilHandling(t *testing.T) {
	// Function that should handle nil gracefully
	safeLength := func(ptr *[]int) int {
		if ptr == nil {
			return -1
		}
		if *ptr == nil {
			return 0
		}
		return len(*ptr)
	}

	attrs := attributes.NewFTAttributes()
	attrs.PointerAttr = attributes.PointerAttributes{
		AllowNil: true,
		Depth:    1,
		Inner: attributes.SliceAttributes{
			MinLen: 0,
			MaxLen: 10,
			ElementAttrs: attributes.IntegerAttributesImpl[int]{
				Min: 0,
				Max: 100,
			},
		},
	}

	ft := &ftesting.FTesting{}
	ft.WithFunction(safeLength).WithAttributes(attrs).WithIterations(50)

	nilPtrCount := 0
	nonNilCount := 0

	for i := uint(0); i < 50; i++ {
		inputs, err := ft.GenerateInputs()
		if err != nil {
			t.Fatalf("Failed to generate inputs: %v", err)
		}

		if inputs[0] == nil {
			nilPtrCount++
		} else {
			nonNilCount++
		}

		ok, err := ft.ApplyFunction()
		if err != nil {
			t.Fatalf("Safe length function failed: %v", err)
		}
		if !ok {
			t.Errorf("Safe length function returned false")
		}
	}

	t.Logf("Tested with %d nil pointers and %d non-nil pointers", nilPtrCount, nonNilCount)
}

// TestFuzzCompositeTypes demonstrates fuzzing complex composite data structures
func TestFuzzCompositeTypes(t *testing.T) {
	// Function working with nested maps and slices
	processData := func(data map[string][]int) int {
		total := 0
		for _, values := range data {
			for _, v := range values {
				total += v
			}
		}
		return total
	}

	attrs := attributes.NewFTAttributes()
	attrs.MapAttr = attributes.MapAttributes{
		MinSize: 1,
		MaxSize: 3,
		KeyAttrs: attributes.StringAttributes{
			MinLen: 3,
			MaxLen: 8,
		},
		ValueAttrs: attributes.SliceAttributes{
			MinLen: 1,
			MaxLen: 5,
			ElementAttrs: attributes.IntegerAttributesImpl[int]{
				Min: 1,
				Max: 10,
			},
		},
	}

	ft := &ftesting.FTesting{}
	ft.WithFunction(processData).WithAttributes(attrs).WithIterations(30)

	for i := uint(0); i < 30; i++ {
		inputs, err := ft.GenerateInputs()
		if err != nil {
			t.Fatalf("Failed to generate inputs: %v", err)
		}

		data := inputs[0].(map[string][]int)
		t.Logf("Generated data structure with %d keys", len(data))

		ok, err := ft.ApplyFunction()
		if err != nil {
			t.Fatalf("Process data function failed: %v", err)
		}
		if !ok {
			t.Errorf("Process data function returned false")
		}
	}
}
