package attributes

import (
	"reflect"
	"testing"

	"github.com/laiambryant/gotestutils/ctesting"
)

type constIntAttr struct{}

func (c constIntAttr) GetAttributes() any                   { return c }
func (c constIntAttr) GetReflectType() reflect.Type         { return reflect.TypeOf(int(0)) }
func (c constIntAttr) GetRandomValue() any                  { return 7 }
func (c constIntAttr) GetDefaultImplementation() Attributes { return c }

func TestArrayAttributes(t *testing.T) {
	var suite []ctesting.CharacterizationTest[bool]

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := ArrayAttributes{Length: 5, ElementAttrs: IntegerAttributesImpl[int]{}}
		got := attr.GetAttributes()
		expected := ArrayAttributes{Length: 5, ElementAttrs: IntegerAttributesImpl[int]{}}
		return reflect.DeepEqual(got, expected), nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := ArrayAttributes{}
		got := attr.GetDefaultImplementation()
		return got != nil && reflect.TypeOf(got) == reflect.TypeOf(attr), nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attrs := ArrayAttributes{
			Length:       5,
			ElementAttrs: reflect.TypeOf(int(0)),
		}
		expectedType := reflect.ArrayOf(5, reflect.TypeOf(int(0)))
		reflectType := attrs.GetReflectType()
		return reflectType == expectedType, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attrs := ArrayAttributes{
			Length:       -5,
			ElementAttrs: IntegerAttributesImpl[int]{},
		}
		reflectType := attrs.GetReflectType()
		return reflectType == nil, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attrs := ArrayAttributes{
			Length:       5,
			ElementAttrs: nil,
		}
		reflectType := attrs.GetReflectType()
		return reflectType == nil, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attrs := ArrayAttributes{
			Length:       5,
			ElementAttrs: nilTypeReturningAttribute{},
		}
		reflectType := attrs.GetReflectType()
		return reflectType == nil, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := ArrayAttributes{Length: 0, ElementAttrs: IntegerAttributesImpl[int]{}}
		result := attr.GetRandomValue()
		return result == nil, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := ArrayAttributes{Length: -5, ElementAttrs: IntegerAttributesImpl[int]{}}
		result := attr.GetRandomValue()
		return result == nil, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := ArrayAttributes{Length: 5, ElementAttrs: "not an attribute"}
		result := attr.GetRandomValue()
		return result == nil, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attrs := ArrayAttributes{
			Length:       3,
			ElementAttrs: nilReturningAttribute{},
		}
		expectedArray := [3]int{0, 0, 0}
		result := attrs.GetRandomValue()
		if arr, ok := result.([3]int); ok {
			return reflect.DeepEqual(arr, expectedArray), nil
		}
		return false, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attrs := ArrayAttributes{
			Length:       4,
			ElementAttrs: constIntAttr{},
		}
		expectedArray := [4]int{7, 7, 7, 7}
		result := attrs.GetRandomValue()
		if arr, ok := result.([4]int); ok {
			return reflect.DeepEqual(arr, expectedArray), nil
		}
		return false, nil
	}))

	results, _ := ctesting.VerifyCharacterizationTestsAndResults(t, suite, true)
	for i, passed := range results {
		if !passed {
			t.Fatalf("ArrayAttributes test %d failed", i+1)
		}
	}
}
