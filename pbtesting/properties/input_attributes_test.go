package properties

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
		{"IntegerAttributes", IntegerAttributes{Signed: true, AllowNegative: true, AllowZero: true, Max: 10, Min: -5, EvenOnly: true, MultipleOf: 2, InSet: []int64{1, 2}, NotInSet: []int64{3}}, IntegerAttributes{Signed: true, AllowNegative: true, AllowZero: true, Max: 10, Min: -5, EvenOnly: true, MultipleOf: 2, InSet: []int64{1, 2}, NotInSet: []int64{3}}},
		{"FloatAttributes", FloatAttributes{Min: 1.1, Max: 2.2, NonZero: true, FiniteOnly: true, AllowNaN: true, AllowInf: true, Precision: 3}, FloatAttributes{Min: 1.1, Max: 2.2, NonZero: true, FiniteOnly: true, AllowNaN: true, AllowInf: true, Precision: 3}},
		{"ComplexAttributes", ComplexAttributes{RealMin: -1, RealMax: 1, ImagMin: -2, ImagMax: 2, MagnitudeMin: 0.5, MagnitudeMax: 10, AllowNaN: true, AllowInf: true}, ComplexAttributes{RealMin: -1, RealMax: 1, ImagMin: -2, ImagMax: 2, MagnitudeMin: 0.5, MagnitudeMax: 10, AllowNaN: true, AllowInf: true}},
		{"StringAttributes", StringAttributes{MinLen: 1, MaxLen: 5, Prefix: "pre", Suffix: "suf", Contains: "mid", UniqueChars: true}, StringAttributes{MinLen: 1, MaxLen: 5, Prefix: "pre", Suffix: "suf", Contains: "mid", UniqueChars: true}},
		{"SliceAttributes", SliceAttributes{MinLen: 1, MaxLen: 3, Unique: true, Sorted: true, ElementPreds: []p.Predicate{}, ElementAttrs: IntegerAttributes{}}, SliceAttributes{MinLen: 1, MaxLen: 3, Unique: true, Sorted: true, ElementPreds: []p.Predicate{}, ElementAttrs: IntegerAttributes{}}},
		{"BoolAttributes", BoolAttributes{ForceTrue: true}, BoolAttributes{ForceTrue: true}},
		{"MapAttributes", MapAttributes{MinSize: 1, MaxSize: 3, KeyPreds: []p.Predicate{}, ValuePreds: []p.Predicate{}, KeyAttrs: StringAttributes{}, ValueAttrs: IntegerAttributes{}}, MapAttributes{MinSize: 1, MaxSize: 3, KeyPreds: []p.Predicate{}, ValuePreds: []p.Predicate{}, KeyAttrs: StringAttributes{}, ValueAttrs: IntegerAttributes{}}},
		{"ChanAttributes", ChanAttributes{MinBuffer: 0, MaxBuffer: 2, ElemAttrs: FloatAttributes{}}, ChanAttributes{MinBuffer: 0, MaxBuffer: 2, ElemAttrs: FloatAttributes{}}},
		{"FuncAttributes", FuncAttributes{Deterministic: true, PanicProbability: 0.1, ReturnZeroValues: true}, FuncAttributes{Deterministic: true, PanicProbability: 0.1, ReturnZeroValues: true}},
		{"InterfaceAttributes", InterfaceAttributes{AllowedConcrete: []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}}, InterfaceAttributes{AllowedConcrete: []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}}},
		{"PointerAttributes", PointerAttributes{AllowNil: true, Depth: 2, Inner: IntegerAttributes{}}, PointerAttributes{AllowNil: true, Depth: 2, Inner: IntegerAttributes{}}},
		{"StructAttributes", StructAttributes{FieldAttrs: map[string]any{"A": IntegerAttributes{}, "B": FloatAttributes{}}}, StructAttributes{FieldAttrs: map[string]any{"A": IntegerAttributes{}, "B": FloatAttributes{}}}},
		{"ArrayAttributes", ArrayAttributes{Length: 3, Sorted: true, ElementAttrs: IntegerAttributes{}, ElementPreds: []p.Predicate{}}, ArrayAttributes{Length: 3, Sorted: true, ElementAttrs: IntegerAttributes{}, ElementPreds: []p.Predicate{}}},
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
