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
		{"IntegerAttributes", IntegerAttributes[int64]{AllowNegative: true, AllowZero: true, Max: 10, Min: -5, InSet: []int64{1, 2}, NotInSet: []int64{3}}, IntegerAttributes[int64]{AllowNegative: true, AllowZero: true, Max: 10, Min: -5, InSet: []int64{1, 2}, NotInSet: []int64{3}}},
		{"UnsignedIntegerAttributes", UnsignedIntegerAttributes[uint64]{AllowNegative: false, AllowZero: true, Max: 100, Min: 0, InSet: []uint64{10, 20}, NotInSet: []uint64{30}}, UnsignedIntegerAttributes[uint64]{AllowNegative: false, AllowZero: true, Max: 100, Min: 0, InSet: []uint64{10, 20}, NotInSet: []uint64{30}}},
		{"FloatAttributes", FloatAttributes[float64]{Min: 1.1, Max: 2.2, NonZero: true, FiniteOnly: true, AllowNaN: true, AllowInf: true, Precision: 3}, FloatAttributes[float64]{Min: 1.1, Max: 2.2, NonZero: true, FiniteOnly: true, AllowNaN: true, AllowInf: true, Precision: 3}},
		{"ComplexAttributes", ComplexAttributes[complex128]{RealMin: -1, RealMax: 1, ImagMin: -2, ImagMax: 2, MagnitudeMin: 0.5, MagnitudeMax: 10, AllowNaN: true, AllowInf: true}, ComplexAttributes[complex128]{RealMin: -1, RealMax: 1, ImagMin: -2, ImagMax: 2, MagnitudeMin: 0.5, MagnitudeMax: 10, AllowNaN: true, AllowInf: true}},
		{"StringAttributes", StringAttributes{MinLen: 1, MaxLen: 5, Prefix: "pre", Suffix: "suf", Contains: "mid", UniqueChars: true}, StringAttributes{MinLen: 1, MaxLen: 5, Prefix: "pre", Suffix: "suf", Contains: "mid", UniqueChars: true}},
		{"SliceAttributes", SliceAttributes{MinLen: 1, MaxLen: 3, Unique: true, Sorted: true, ElementPreds: []p.Predicate{}, ElementAttrs: IntegerAttributes[int64]{}}, SliceAttributes{MinLen: 1, MaxLen: 3, Unique: true, Sorted: true, ElementPreds: []p.Predicate{}, ElementAttrs: IntegerAttributes[int64]{}}},
		{"BoolAttributes", BoolAttributes{ForceTrue: true}, BoolAttributes{ForceTrue: true}},
		{"MapAttributes", MapAttributes{MinSize: 1, MaxSize: 3, KeyPreds: []p.Predicate{}, ValuePreds: []p.Predicate{}, KeyAttrs: StringAttributes{}, ValueAttrs: IntegerAttributes[int64]{}}, MapAttributes{MinSize: 1, MaxSize: 3, KeyPreds: []p.Predicate{}, ValuePreds: []p.Predicate{}, KeyAttrs: StringAttributes{}, ValueAttrs: IntegerAttributes[int64]{}}},
		{"ChanAttributes", ChanAttributes{MinBuffer: 0, MaxBuffer: 2, ElemAttrs: FloatAttributes[float64]{}}, ChanAttributes{MinBuffer: 0, MaxBuffer: 2, ElemAttrs: FloatAttributes[float64]{}}},
		{"FuncAttributes", FuncAttributes{Deterministic: true, PanicProbability: 0.1, ReturnZeroValues: true}, FuncAttributes{Deterministic: true, PanicProbability: 0.1, ReturnZeroValues: true}},
		{"InterfaceAttributes", InterfaceAttributes{AllowedConcrete: []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}}, InterfaceAttributes{AllowedConcrete: []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}}},
		{"PointerAttributes", PointerAttributes{AllowNil: true, Depth: 2, Inner: IntegerAttributes[int64]{}}, PointerAttributes{AllowNil: true, Depth: 2, Inner: IntegerAttributes[int64]{}}},
		{"StructAttributes", StructAttributes{FieldAttrs: map[string]any{"A": IntegerAttributes[int64]{}, "B": FloatAttributes[float64]{}}}, StructAttributes{FieldAttrs: map[string]any{"A": IntegerAttributes[int64]{}, "B": FloatAttributes[float64]{}}}},
		{"ArrayAttributes", ArrayAttributes{Length: 3, Sorted: true, ElementAttrs: IntegerAttributes[int64]{}}, ArrayAttributes{Length: 3, Sorted: true, ElementAttrs: IntegerAttributes[int64]{}}},
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
