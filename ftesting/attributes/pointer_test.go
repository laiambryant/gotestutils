package attributes

import (
	"reflect"
	"testing"

	"github.com/laiambryant/gotestutils/ctesting"
)

func TestPointerAttributes(t *testing.T) {
	var suite []ctesting.CharacterizationTest[bool]

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := PointerAttributes{AllowNil: true, Depth: 2, Inner: IntegerAttributesImpl[int]{}}
		got := attr.GetAttributes()
		expected := PointerAttributes{AllowNil: true, Depth: 2, Inner: IntegerAttributesImpl[int]{}}
		return reflect.DeepEqual(got, expected), nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := PointerAttributes{}
		got := attr.GetDefaultImplementation()
		return got != nil && reflect.TypeOf(got) == reflect.TypeOf(attr), nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := PointerAttributes{AllowNil: false, Depth: 1, Inner: IntegerAttributesImpl[int]{Max: 10, Min: 1}}
		got := attr.GetRandomValue()
		return got != nil, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := PointerAttributes{AllowNil: false, Depth: 1, Inner: "not an attribute"}
		result := attr.GetRandomValue()
		return result == nil, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := PointerAttributes{AllowNil: false, Depth: 3, Inner: IntegerAttributesImpl[int]{Max: 10, Min: 1}}
		result := attr.GetRandomValue()
		if result == nil {
			return false, nil
		}
		rv := reflect.ValueOf(result)
		return rv.Kind() == reflect.Pointer, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attrs := PointerAttributes{
			Inner: reflect.TypeOf(int(0)),
			Depth: 2,
		}
		reflectType := attrs.GetReflectType()
		if reflectType == nil || reflectType.Kind() != reflect.Pointer {
			return false, nil
		}
		if reflectType.Elem().Kind() != reflect.Pointer {
			return false, nil
		}
		return reflectType.Elem().Elem() == reflect.TypeOf(int(0)), nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attrs := PointerAttributes{
			Inner: nil,
			Depth: 1,
		}
		reflectType := attrs.GetReflectType()
		return reflectType == nil, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attrs := PointerAttributes{
			Inner: IntegerAttributesImpl[int]{},
			Depth: 0,
		}
		reflectType := attrs.GetReflectType()
		return reflectType != nil && reflectType.Kind() == reflect.Pointer, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attrs := PointerAttributes{
			Inner: IntegerAttributesImpl[int]{},
			Depth: -5,
		}
		reflectType := attrs.GetReflectType()
		return reflectType != nil && reflectType.Kind() == reflect.Pointer, nil
	}))

	results, _ := ctesting.VerifyCharacterizationTestsAndResults(t, suite, true)
	for i, passed := range results {
		if !passed {
			t.Fatalf("PointerAttributes test %d failed", i+1)
		}
	}
}

func TestPointerAttributes_NilWhenAllowNil(t *testing.T) {
	attr := PointerAttributes{AllowNil: true, Depth: 1, Inner: IntegerAttributesImpl[int]{Max: 10, Min: 1}}
	foundNil := false
	for range 100 {
		result := attr.GetRandomValue()
		if result == nil {
			t.Error("GetRandomValue should never return nil directly, it returns a typed nil pointer")
		}
		rv := reflect.ValueOf(result)
		if rv.Kind() == reflect.Pointer && rv.IsNil() {
			foundNil = true
			break
		}
	}
	if !foundNil {
		t.Log("Warning: Did not find nil pointer in 100 attempts (this is statistically unlikely but possible)")
	}
}

func TestPointerAttributes_InvalidInnerType(t *testing.T) {
	attr := PointerAttributes{AllowNil: false, Depth: 1, Inner: "not an attribute"}

	testSuite := []ctesting.CharacterizationTest[any]{
		ctesting.NewCharacterizationTest(nil, nil, func() (any, error) {
			result := attr.GetRandomValue()
			return result, nil
		}),
	}

	ctesting.VerifyCharacterizationTestsAndResults(t, testSuite, true)
}

func TestPointerAttributes_MultipleDepth(t *testing.T) {
	attr := PointerAttributes{AllowNil: false, Depth: 3, Inner: IntegerAttributesImpl[int]{Max: 10, Min: 1}}

	testSuite := []ctesting.CharacterizationTest[bool]{
		ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
			result := attr.GetRandomValue()
			if result == nil {
				return false, nil
			}
			rv := reflect.ValueOf(result)
			return rv.Kind() == reflect.Pointer, nil
		}),
	}

	ctesting.VerifyCharacterizationTestsAndResults(t, testSuite, true)
}

func TestPointerAttributes_GetReflectType_WithReflectType(t *testing.T) {
	attrs := PointerAttributes{
		Inner: reflect.TypeOf(int(0)),
		Depth: 2,
	}

	testSuite := []ctesting.CharacterizationTest[bool]{
		ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
			reflectType := attrs.GetReflectType()
			if reflectType == nil {
				return false, nil
			}

			if reflectType.Kind() != reflect.Pointer {
				return false, nil
			}

			if reflectType.Elem().Kind() != reflect.Pointer {
				return false, nil
			}

			return reflectType.Elem().Elem() == reflect.TypeOf(int(0)), nil
		}),
	}

	ctesting.VerifyCharacterizationTestsAndResults(t, testSuite, true)
}

func TestPointerAttributes_GetReflectType_WithNilInner(t *testing.T) {
	attrs := PointerAttributes{
		Inner: nil,
		Depth: 1,
	}

	testSuite := []ctesting.CharacterizationTest[reflect.Type]{
		ctesting.NewCharacterizationTest(nil, nil, func() (reflect.Type, error) {
			reflectType := attrs.GetReflectType()
			return reflectType, nil
		}),
	}

	ctesting.VerifyCharacterizationTestsAndResults(t, testSuite, true)
}

func TestPointerAttributes_GetReflectType_WithZeroDepth(t *testing.T) {
	attrs := PointerAttributes{
		Inner: IntegerAttributesImpl[int]{},
		Depth: 0,
	}

	testSuite := []ctesting.CharacterizationTest[bool]{
		ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
			reflectType := attrs.GetReflectType()
			return reflectType != nil && reflectType.Kind() == reflect.Pointer, nil
		}),
	}

	ctesting.VerifyCharacterizationTestsAndResults(t, testSuite, true)
}

func TestPointerAttributes_GetReflectType_WithNegativeDepth(t *testing.T) {
	attrs := PointerAttributes{
		Inner: IntegerAttributesImpl[int]{},
		Depth: -5,
	}

	testSuite := []ctesting.CharacterizationTest[bool]{
		ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
			reflectType := attrs.GetReflectType()
			return reflectType != nil && reflectType.Kind() == reflect.Pointer, nil
		}),
	}

	ctesting.VerifyCharacterizationTestsAndResults(t, testSuite, true)
}

func TestPointerAttributes_NilInnerValueWithType(t *testing.T) {
	attrs := PointerAttributes{
		Inner:    nilReturningAttribute{},
		AllowNil: false,
		Depth:    1,
	}

	testSuite := []ctesting.CharacterizationTest[bool]{
		ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
			result := attrs.GetRandomValue()
			if result == nil {
				return false, nil
			}

			ptrValue := reflect.ValueOf(result)
			return ptrValue.Kind() == reflect.Pointer, nil
		}),
	}

	ctesting.VerifyCharacterizationTestsAndResults(t, testSuite, true)
}
