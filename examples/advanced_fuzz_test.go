package examples

import (
	"testing"

	"github.com/laiambryant/gotestutils/ftesting"
	"github.com/laiambryant/gotestutils/ftesting/attributes"
)

// TestFuzzComplexNumbers demonstrates fuzzing with complex number types
func TestFuzzComplexNumbers(t *testing.T) {
	attrs := attributes.NewFTAttributes()
	attrs.ComplexAttr = attributes.ComplexAttributesImpl[complex128]{
		RealMin: -10.0,
		RealMax: 10.0,
		ImagMin: -10.0,
		ImagMax: 10.0,
	}

	ft := &ftesting.FTesting{}
	ft.WithFunction(func(c complex128) float64 {
		// Calculate magnitude of complex number
		return real(c)*real(c) + imag(c)*imag(c)
	}).WithAttributes(attrs).WithIterations(50)

	for i := uint(0); i < 50; i++ {
		ok, err := ft.ApplyFunction()
		if err != nil {
			t.Fatalf("Complex number function failed: %v", err)
		}
		if !ok {
			t.Errorf("Complex number function returned false")
		}
	}
}

// TestFuzzSlices demonstrates fuzzing with slice inputs
func TestFuzzSlices(t *testing.T) {
	attrs := attributes.NewFTAttributes()
	attrs.SliceAttr = attributes.SliceAttributes{
		MinLen: 1,
		MaxLen: 10,
		ElementAttrs: attributes.IntegerAttributesImpl[int]{
			Min: 0,
			Max: 100,
		},
	}

	ft := &ftesting.FTesting{}
	ft.WithFunction(func(nums []int) int {
		sum := 0
		for _, n := range nums {
			sum += n
		}
		return sum
	}).WithAttributes(attrs).WithIterations(50)

	for i := uint(0); i < 50; i++ {
		ok, err := ft.ApplyFunction()
		if err != nil {
			t.Fatalf("Slice function failed: %v", err)
		}
		if !ok {
			t.Errorf("Slice function returned false")
		}
	}
}

// TestFuzzMaps demonstrates fuzzing with map inputs
func TestFuzzMaps(t *testing.T) {
	attrs := attributes.NewFTAttributes()
	attrs.MapAttr = attributes.MapAttributes{
		MinSize: 1,
		MaxSize: 5,
		KeyAttrs: attributes.StringAttributes{
			MinLen: 2,
			MaxLen: 5,
		},
		ValueAttrs: attributes.IntegerAttributesImpl[int]{
			Min: 0,
			Max: 100,
		},
	}

	ft := &ftesting.FTesting{}
	ft.WithFunction(func(data map[string]int) int {
		count := 0
		for _, v := range data {
			count += v
		}
		return count
	}).WithAttributes(attrs).WithIterations(30)

	for i := uint(0); i < 30; i++ {
		ok, err := ft.ApplyFunction()
		if err != nil {
			t.Fatalf("Map function failed: %v", err)
		}
		if !ok {
			t.Errorf("Map function returned false")
		}
	}
}

// TestFuzzPointers demonstrates fuzzing with pointer inputs
func TestFuzzPointers(t *testing.T) {
	attrs := attributes.NewFTAttributes()
	attrs.PointerAttr = attributes.PointerAttributes{
		AllowNil: true,
		Depth:    1,
		Inner: attributes.IntegerAttributesImpl[int]{
			Min: 1,
			Max: 100,
		},
	}

	ft := &ftesting.FTesting{}
	ft.WithFunction(func(ptr *int) int {
		if ptr == nil {
			return 0
		}
		return *ptr
	}).WithAttributes(attrs).WithIterations(50)

	nilCount := 0
	nonNilCount := 0

	for i := uint(0); i < 50; i++ {
		inputs, err := ft.GenerateInputs()
		if err != nil {
			t.Fatalf("Failed to generate inputs: %v", err)
		}

		if ptr, ok := inputs[0].(*int); ok {
			if ptr == nil {
				nilCount++
			} else {
				nonNilCount++
			}
		}

		ok, err := ft.ApplyFunction()
		if err != nil {
			t.Fatalf("Pointer function failed: %v", err)
		}
		if !ok {
			t.Errorf("Pointer function returned false")
		}
	}

	t.Logf("Generated %d nil pointers and %d non-nil pointers", nilCount, nonNilCount)
}

// TestFuzzStructs demonstrates fuzzing with struct inputs
func TestFuzzStructs(t *testing.T) {
	// Note: Dynamic struct generation with reflection has limitations.
	// For practical use, define concrete struct types.
	// This test demonstrates the attribute system for struct configuration.

	attrs := attributes.NewFTAttributes()
	attrs.StructAttr = attributes.StructAttributes{
		FieldAttrs: map[string]any{
			"ID": attributes.IntegerAttributesImpl[int]{
				Min: 1,
				Max: 1000,
			},
			"Score": attributes.FloatAttributesImpl[float64]{
				Min: 0.0,
				Max: 100.0,
			},
		},
	}

	// For demonstration, we'll test that struct attributes can be retrieved
	structType := attrs.StructAttr.GetReflectType()
	if structType == nil {
		t.Skip("Struct type generation not available for dynamic types")
	}

	t.Logf("Struct type generated: %v", structType)
	t.Logf("Successfully configured struct attributes with %d fields", len(attrs.StructAttr.FieldAttrs))
}

// TestFuzzArrays demonstrates fuzzing with fixed-size array inputs
func TestFuzzArrays(t *testing.T) {
	attrs := attributes.NewFTAttributes()
	attrs.ArrayAttr = attributes.ArrayAttributes{
		Length: 5,
		ElementAttrs: attributes.IntegerAttributesImpl[int]{
			Min: 0,
			Max: 10,
		},
	}

	ft := &ftesting.FTesting{}
	ft.WithFunction(func(arr [5]int) int {
		sum := 0
		for _, v := range arr {
			sum += v
		}
		return sum
	}).WithAttributes(attrs).WithIterations(30)

	for i := uint(0); i < 30; i++ {
		ok, err := ft.ApplyFunction()
		if err != nil {
			t.Fatalf("Array function failed: %v", err)
		}
		if !ok {
			t.Errorf("Array function returned false")
		}
	}
}

// TestFuzzUnsignedIntegers demonstrates fuzzing with unsigned integer types
func TestFuzzUnsignedIntegers(t *testing.T) {
	attrs := attributes.NewFTAttributes()
	attrs.UIntegerAttr = attributes.UnsignedIntegerAttributesImpl[uint]{
		Signed:        false,
		AllowNegative: false,
		AllowZero:     true,
		Min:           0,
		Max:           255,
	}

	ft := &ftesting.FTesting{}
	ft.WithFunction(func(a, b uint) uint {
		return a + b
	}).WithAttributes(attrs).WithIterations(50)

	for i := uint(0); i < 50; i++ {
		ok, err := ft.ApplyFunction()
		if err != nil {
			t.Fatalf("Unsigned integer function failed: %v", err)
		}
		if !ok {
			t.Errorf("Unsigned integer function returned false")
		}
	}
}

// TestFuzzStringAttributes demonstrates advanced string generation options
func TestFuzzStringAttributes(t *testing.T) {
	attrs := attributes.NewFTAttributes()
	attrs.StringAttr = attributes.StringAttributes{
		MinLen:       5,
		MaxLen:       15,
		AllowedRunes: []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"),
		Prefix:       "USER_",
		Suffix:       "_ID",
	}

	ft := &ftesting.FTesting{}
	ft.WithFunction(func(username string) bool {
		// Validate username format
		return len(username) > 7 // PREFIX + at least 1 char + SUFFIX
	}).WithAttributes(attrs).WithIterations(20)

	for i := uint(0); i < 20; i++ {
		inputs, err := ft.GenerateInputs()
		if err != nil {
			t.Fatalf("Failed to generate inputs: %v", err)
		}

		if str, ok := inputs[0].(string); ok {
			t.Logf("Generated username: %s", str)
		}

		ok, err := ft.ApplyFunction()
		if err != nil {
			t.Fatalf("String validation failed: %v", err)
		}
		if !ok {
			t.Errorf("String validation returned false")
		}
	}
}

// TestFuzzBooleanInputs demonstrates fuzzing with boolean parameters
func TestFuzzBooleanInputs(t *testing.T) {
	attrs := attributes.NewFTAttributes()
	attrs.BoolAttr = attributes.BoolAttributes{
		ForceTrue:  false,
		ForceFalse: false,
	}

	ft := &ftesting.FTesting{}
	ft.WithFunction(func(enabled bool, debug bool) string {
		if enabled && debug {
			return "debug-mode"
		} else if enabled {
			return "normal-mode"
		}
		return "disabled"
	}).WithAttributes(attrs).WithIterations(50)

	modes := make(map[string]int)
	for i := uint(0); i < 50; i++ {
		inputs, err := ft.GenerateInputs()
		if err != nil {
			t.Fatalf("Failed to generate inputs: %v", err)
		}

		// Track the modes we encounter
		enabled := inputs[0].(bool)
		debug := inputs[1].(bool)
		if enabled && debug {
			modes["debug-mode"]++
		} else if enabled {
			modes["normal-mode"]++
		} else {
			modes["disabled"]++
		}

		ok, err := ft.ApplyFunction()
		if err != nil {
			t.Fatalf("Boolean function failed: %v", err)
		}
		if !ok {
			t.Errorf("Boolean function returned false")
		}
	}

	t.Logf("Mode distribution: %+v", modes)
}
