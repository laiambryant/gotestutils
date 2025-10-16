package attributes

import (
	"reflect"
	"testing"

	ctesting "github.com/laiambryant/gotestutils/ctesting"
)

func TestComplexAttributes(t *testing.T) {
	var suite []ctesting.CharacterizationTest[bool]

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := ComplexAttributesImpl[complex128]{RealMin: -1, RealMax: 1, ImagMin: -2, ImagMax: 2, MagnitudeMin: 0.5, MagnitudeMax: 10, AllowNaN: true, AllowInf: true}
		got := attr.GetAttributes()
		expected := ComplexAttributesImpl[complex128]{RealMin: -1, RealMax: 1, ImagMin: -2, ImagMax: 2, MagnitudeMin: 0.5, MagnitudeMax: 10, AllowNaN: true, AllowInf: true}
		return reflect.DeepEqual(got, expected), nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := ComplexAttributesImpl[complex128]{}
		got := attr.GetReflectType()
		expected := reflect.TypeOf(complex128(0))
		return got == expected, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := ComplexAttributesImpl[complex128]{}
		got := attr.GetDefaultImplementation()
		return got != nil && reflect.TypeOf(got) == reflect.TypeOf(attr), nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := ComplexAttributesImpl[complex128]{RealMin: -1.0, RealMax: 1.0, ImagMin: -1.0, ImagMax: 1.0}
		got := attr.GetRandomValue()
		return got != nil, nil
	}))

	// more edge cases
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := ComplexAttributesImpl[complex128]{RealMin: -10.0, RealMax: 10.0, ImagMin: -5.0, ImagMax: 5.0}
		result := attr.GetRandomValue()
		if result == nil {
			return false, nil
		}
		c := result.(complex128)
		real := real(c)
		imag := imag(c)
		return real >= -10.0 && real <= 10.0 && imag >= -5.0 && imag <= 5.0, nil
	}))

	// invalid bounds tests
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := ComplexAttributesImpl[complex128]{RealMin: 10.0, RealMax: 5.0, ImagMin: -5.0, ImagMax: 5.0}
		result := attr.GetRandomValue()
		if result == nil {
			return false, nil
		}
		c := result.(complex128)
		real := real(c)
		return real >= -10.0 && real <= 10.0, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := ComplexAttributesImpl[complex128]{RealMin: -5.0, RealMax: 5.0, ImagMin: 10.0, ImagMax: 5.0}
		result := attr.GetRandomValue()
		if result == nil {
			return false, nil
		}
		c := result.(complex128)
		imag := imag(c)
		return imag >= -10.0 && imag <= 10.0, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := ComplexAttributesImpl[complex128]{RealMin: 10.0, RealMax: 5.0, ImagMin: 10.0, ImagMax: 5.0}
		result := attr.GetRandomValue()
		if result == nil {
			return false, nil
		}
		c := result.(complex128)
		real := real(c)
		imag := imag(c)
		return real >= -10.0 && real <= 10.0 && imag >= -10.0 && imag <= 10.0, nil
	}))

	// Additional tests from individual test functions
	// TestComplexAttributes_Valid
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attrs := ComplexAttributesImpl[complex128]{RealMin: -10.0, RealMax: 10.0, ImagMin: -5.0, ImagMax: 5.0}
		result := attrs.GetRandomValue()
		if result == nil {
			return false, nil
		}
		c := result.(complex128)
		real := real(c)
		imag := imag(c)
		return real >= -10.0 && real <= 10.0 && imag >= -5.0 && imag <= 5.0, nil
	}))

	// TestComplexAttributes_InvalidImagBounds
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attrs := ComplexAttributesImpl[complex128]{
			RealMin: -5.0,
			RealMax: 5.0,
			ImagMin: 10.0,
			ImagMax: 5.0,
		}
		result := attrs.GetRandomValue()
		if result == nil {
			return false, nil
		}
		c := result.(complex128)
		imag := imag(c)
		return imag >= -10.0 && imag <= 10.0, nil
	}))

	// TestComplexAttributes_BothInvalidBounds
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attrs := ComplexAttributesImpl[complex128]{
			RealMin: 10.0,
			RealMax: 5.0,
			ImagMin: 10.0,
			ImagMax: 5.0,
		}
		result := attrs.GetRandomValue()
		if result == nil {
			return false, nil
		}
		c := result.(complex128)
		real := real(c)
		imag := imag(c)
		return real >= -10.0 && real <= 10.0 && imag >= -10.0 && imag <= 10.0, nil
	}))

	// TestComplexAttributes_InvalidRealBounds
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attrs := ComplexAttributesImpl[complex128]{
			RealMin: 10.0,
			RealMax: 5.0,
			ImagMin: -5.0,
			ImagMax: 5.0,
		}
		result := attrs.GetRandomValue()
		if result == nil {
			return false, nil
		}
		c := result.(complex128)
		real := real(c)
		return real >= -10.0 && real <= 10.0, nil
	}))

	results, _ := ctesting.VerifyCharacterizationTestsAndResults(t, suite, true)
	for i, passed := range results {
		if !passed {
			t.Fatalf("ComplexAttributes test %d failed", i+1)
		}
	}
}
