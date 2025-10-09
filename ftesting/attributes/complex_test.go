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

	results, _ := ctesting.VerifyCharacterizationTestsAndResults(t, suite, true)
	for i, passed := range results {
		if !passed {
			t.Fatalf("ComplexAttributes test %d failed", i+1)
		}
	}
}

func TestComplexAttributes_Valid(t *testing.T) {
	attrs := ComplexAttributesImpl[complex128]{RealMin: -10.0, RealMax: 10.0, ImagMin: -5.0, ImagMax: 5.0}
	result := attrs.GetRandomValue()
	if result == nil {
		t.Fatal("Expected non-nil complex number")
	}
	c := result.(complex128)
	real := real(c)
	imag := imag(c)
	if real < -10.0 || real > 10.0 {
		t.Errorf("Expected real part in range [-10, 10], got %f", real)
	}
	if imag < -5.0 || imag > 5.0 {
		t.Errorf("Expected imaginary part in range [-5, 5], got %f", imag)
	}
}

func TestComplexAttributes_InvalidImagBounds(t *testing.T) {
	attrs := ComplexAttributesImpl[complex128]{
		RealMin: -5.0,
		RealMax: 5.0,
		ImagMin: 10.0,
		ImagMax: 5.0,
	}

	result := attrs.GetRandomValue()
	if result == nil {
		t.Fatal("Expected non-nil complex number")
	}

	c := result.(complex128)
	imag := imag(c)
	if imag < -10.0 || imag > 10.0 {
		t.Errorf("Expected imaginary part in default range [-10, 10], got %f", imag)
	}
}

func TestComplexAttributes_BothInvalidBounds(t *testing.T) {
	attrs := ComplexAttributesImpl[complex128]{
		RealMin: 10.0,
		RealMax: 5.0,
		ImagMin: 10.0,
		ImagMax: 5.0,
	}

	result := attrs.GetRandomValue()
	if result == nil {
		t.Fatal("Expected non-nil complex number")
	}

	c := result.(complex128)
	real := real(c)
	imag := imag(c)
	if real < -10.0 || real > 10.0 {
		t.Errorf("Expected real part in default range [-10, 10], got %f", real)
	}
	if imag < -10.0 || imag > 10.0 {
		t.Errorf("Expected imaginary part in default range [-10, 10], got %f", imag)
	}
}

func TestComplexAttributes_InvalidRealBounds(t *testing.T) {
	attrs := ComplexAttributesImpl[complex128]{
		RealMin: 10.0,
		RealMax: 5.0,
		ImagMin: -5.0,
		ImagMax: 5.0,
	}

	result := attrs.GetRandomValue()
	if result == nil {
		t.Fatal("Expected non-nil complex number")
	}

	c := result.(complex128)
	real := real(c)
	if real < -10.0 || real > 10.0 {
		t.Errorf("Expected real part in default range [-10, 10], got %f", real)
	}
}
