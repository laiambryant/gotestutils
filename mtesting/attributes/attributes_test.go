package attributes

import (
	"reflect"
	"testing"

	ctesting "github.com/laiambryant/gotestutils/ctesting"
	p "github.com/laiambryant/gotestutils/pbtesting/properties/predicates"
)

func TestGetAttributesMethods(t *testing.T) {
	type testCase struct {
		name string
		in   Attributes
		want any
	}
	cases := []testCase{
		{"IntegerAttributesImpl", IntegerAttributesImpl[int64]{AllowNegative: true, AllowZero: true, Max: 10, Min: -5}, IntegerAttributesImpl[int64]{AllowNegative: true, AllowZero: true, Max: 10, Min: -5}},
		{"UnsignedIntegerAttributesImpl", UnsignedIntegerAttributesImpl[uint64]{AllowNegative: false, AllowZero: true, Max: 100, Min: 0}, UnsignedIntegerAttributesImpl[uint64]{AllowNegative: false, AllowZero: true, Max: 100, Min: 0}},
		{"FloatAttributesImpl", FloatAttributesImpl[float64]{Min: 1.1, Max: 2.2, NonZero: true, FiniteOnly: true, AllowNaN: true, AllowInf: true, Precision: 3}, FloatAttributesImpl[float64]{Min: 1.1, Max: 2.2, NonZero: true, FiniteOnly: true, AllowNaN: true, AllowInf: true, Precision: 3}},
		{"ComplexAttributesImpl", ComplexAttributesImpl[complex128]{RealMin: -1, RealMax: 1, ImagMin: -2, ImagMax: 2, MagnitudeMin: 0.5, MagnitudeMax: 10, AllowNaN: true, AllowInf: true}, ComplexAttributesImpl[complex128]{RealMin: -1, RealMax: 1, ImagMin: -2, ImagMax: 2, MagnitudeMin: 0.5, MagnitudeMax: 10, AllowNaN: true, AllowInf: true}},
		{"StringAttributes", StringAttributes{MinLen: 1, MaxLen: 5, Prefix: "pre", Suffix: "suf", Contains: "mid", UniqueChars: true}, StringAttributes{MinLen: 1, MaxLen: 5, Prefix: "pre", Suffix: "suf", Contains: "mid", UniqueChars: true}},
		{"SliceAttributes", SliceAttributes{MinLen: 1, MaxLen: 3, Unique: true, Sorted: true, ElementPreds: []p.Predicate{}, ElementAttrs: IntegerAttributesImpl[int64]{}}, SliceAttributes{MinLen: 1, MaxLen: 3, Unique: true, Sorted: true, ElementPreds: []p.Predicate{}, ElementAttrs: IntegerAttributesImpl[int64]{}}},
		{"BoolAttributes", BoolAttributes{ForceTrue: true}, BoolAttributes{ForceTrue: true}},
		{"MapAttributes", MapAttributes{MinSize: 1, MaxSize: 3, KeyPreds: []p.Predicate{}, ValuePreds: []p.Predicate{}, KeyAttrs: StringAttributes{}, ValueAttrs: IntegerAttributesImpl[int64]{}}, MapAttributes{MinSize: 1, MaxSize: 3, KeyPreds: []p.Predicate{}, ValuePreds: []p.Predicate{}, KeyAttrs: StringAttributes{}, ValueAttrs: IntegerAttributesImpl[int64]{}}},
		{"PointerAttributes", PointerAttributes{AllowNil: true, Depth: 2, Inner: IntegerAttributesImpl[int64]{}}, PointerAttributes{AllowNil: true, Depth: 2, Inner: IntegerAttributesImpl[int64]{}}},
		{"StructAttributes", StructAttributes{FieldAttrs: map[string]any{"A": IntegerAttributesImpl[int64]{}, "B": FloatAttributesImpl[float64]{}}}, StructAttributes{FieldAttrs: map[string]any{"A": IntegerAttributesImpl[int64]{}, "B": FloatAttributesImpl[float64]{}}}},
		{"ArrayAttributes", ArrayAttributes{Length: 3, Sorted: true, ElementAttrs: IntegerAttributesImpl[int64]{}}, ArrayAttributes{Length: 3, Sorted: true, ElementAttrs: IntegerAttributesImpl[int64]{}}},
	}
	var suite []ctesting.CharacterizationTest[bool]
	for _, tc := range cases {
		toExec := tc
		suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
			got := toExec.in.GetAttributes()
			if reflect.TypeOf(got) != reflect.TypeOf(toExec.want) {
				return false, nil
			}
			if !reflect.DeepEqual(got, toExec.want) {
				return false, nil
			}
			return true, nil
		}))
	}
	results, _ := ctesting.VerifyCharacterizationTestsAndResults(t, suite, true)
	for i, passed := range results {
		if !passed {
			t.Fatalf("characterization test %d failed", i+1)
		}
	}
}

func TestGetReflectTypeMethods(t *testing.T) {
	type testCase struct {
		name     string
		in       Attributes
		wantType reflect.Type
	}
	cases := []testCase{
		{"IntegerAttributesImpl", IntegerAttributesImpl[int64]{}, reflect.TypeOf(int64(0))},
		{"UnsignedIntegerAttributesImpl", UnsignedIntegerAttributesImpl[uint64]{}, reflect.TypeOf(uint64(0))},
		{"FloatAttributesImpl", FloatAttributesImpl[float64]{}, reflect.TypeOf(float64(0))},
		{"ComplexAttributesImpl", ComplexAttributesImpl[complex128]{}, reflect.TypeOf(complex128(0))},
		{"StringAttributes", StringAttributes{}, reflect.TypeOf("")},
		{"BoolAttributes", BoolAttributes{}, reflect.TypeOf(true)},
	}
	var suite []ctesting.CharacterizationTest[bool]
	for _, tc := range cases {
		toExec := tc
		suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
			got := toExec.in.GetReflectType()
			return got == toExec.wantType, nil
		}))
	}
	results, _ := ctesting.VerifyCharacterizationTestsAndResults(t, suite, true)
	for i, passed := range results {
		if !passed {
			t.Fatalf("GetReflectType test %d failed", i+1)
		}
	}
}

func TestGetDefaultImplementationMethods(t *testing.T) {
	type testCase struct {
		name string
		in   Attributes
	}
	cases := []testCase{
		{"IntegerAttributesImpl", IntegerAttributesImpl[int64]{}},
		{"UnsignedIntegerAttributesImpl", UnsignedIntegerAttributesImpl[uint64]{}},
		{"FloatAttributesImpl", FloatAttributesImpl[float64]{}},
		{"ComplexAttributesImpl", ComplexAttributesImpl[complex128]{}},
		{"StringAttributes", StringAttributes{}},
		{"SliceAttributes", SliceAttributes{}},
		{"BoolAttributes", BoolAttributes{}},
		{"MapAttributes", MapAttributes{}},
		{"PointerAttributes", PointerAttributes{}},
		{"StructAttributes", StructAttributes{}},
		{"ArrayAttributes", ArrayAttributes{}},
	}
	var suite []ctesting.CharacterizationTest[bool]
	for _, tc := range cases {
		toExec := tc
		suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
			got := toExec.in.GetDefaultImplementation()
			if got == nil {
				return false, nil
			}
			return reflect.TypeOf(got) == reflect.TypeOf(toExec.in), nil
		}))
	}
	results, _ := ctesting.VerifyCharacterizationTestsAndResults(t, suite, true)
	for i, passed := range results {
		if !passed {
			t.Fatalf("GetDefaultImplementation test %d failed", i+1)
		}
	}
}

func TestGetRandomValueMethods(t *testing.T) {
	type testCase struct {
		name string
		in   Attributes
	}
	cases := []testCase{
		{"IntegerAttributesImpl", IntegerAttributesImpl[int64]{Min: -10, Max: 10}},
		{"UnsignedIntegerAttributesImpl", UnsignedIntegerAttributesImpl[uint64]{Min: 0, Max: 100}},
		{"FloatAttributesImpl", FloatAttributesImpl[float64]{Min: -1.0, Max: 1.0}},
		{"ComplexAttributesImpl", ComplexAttributesImpl[complex128]{RealMin: -1.0, RealMax: 1.0, ImagMin: -1.0, ImagMax: 1.0}},
		{"StringAttributes", StringAttributes{MinLen: 1, MaxLen: 10}},
		{"SliceAttributes", SliceAttributes{MinLen: 1, MaxLen: 3, ElementAttrs: IntegerAttributesImpl[int]{Min: 0, Max: 10}}},
		{"BoolAttributes", BoolAttributes{}},
		{"MapAttributes", MapAttributes{MinSize: 1, MaxSize: 3, KeyAttrs: StringAttributes{MinLen: 1, MaxLen: 5}, ValueAttrs: IntegerAttributesImpl[int]{Min: 0, Max: 10}}},
		{"PointerAttributes", PointerAttributes{AllowNil: true, Depth: 1, Inner: IntegerAttributesImpl[int]{Min: 0, Max: 10}}},
		{"StructAttributes", StructAttributes{FieldAttrs: map[string]any{"TestField": IntegerAttributesImpl[int]{Min: 0, Max: 10}}}},
		{"ArrayAttributes", ArrayAttributes{Length: 3, ElementAttrs: IntegerAttributesImpl[int]{Min: 0, Max: 10}}},
	}
	var suite []ctesting.CharacterizationTest[bool]
	for _, tc := range cases {
		toExec := tc
		suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
			got := toExec.in.GetRandomValue()
			return got != nil || isNilValidForType(toExec.in), nil
		}))
	}
	results, _ := ctesting.VerifyCharacterizationTestsAndResults(t, suite, true)
	for i, passed := range results {
		if !passed {
			t.Fatalf("GetRandomValue test %d failed", i+1)
		}
	}
}

func isNilValidForType(attr Attributes) bool {
	switch attr.(type) {
	case PointerAttributes:
		return true
	default:
		return false
	}
}
