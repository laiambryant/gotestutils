package properties

import (
	"reflect"
	"testing"
)

func TestGetAttributesMethods(t *testing.T) {
	cases := []struct {
		name string
		in   Attributes
		want any
	}{
		{"IntegerAttributes", IntegerAttributes{Signed: true, AllowNegative: true, AllowZero: true, Max: 10, Min: -5, EvenOnly: true, MultipleOf: 2, InSet: []int64{1, 2}, NotInSet: []int64{3}}, IntegerAttributes{Signed: true, AllowNegative: true, AllowZero: true, Max: 10, Min: -5, EvenOnly: true, MultipleOf: 2, InSet: []int64{1, 2}, NotInSet: []int64{3}}},
		{"FloatAttributes", FloatAttributes{Min: 1.1, Max: 2.2, NonZero: true, FiniteOnly: true, AllowNaN: true, AllowInf: true, Precision: 3}, FloatAttributes{Min: 1.1, Max: 2.2, NonZero: true, FiniteOnly: true, AllowNaN: true, AllowInf: true, Precision: 3}},
		{"ComplexAttributes", ComplexAttributes{RealMin: -1, RealMax: 1, ImagMin: -2, ImagMax: 2, MagnitudeMin: 0.5, MagnitudeMax: 10, AllowNaN: true, AllowInf: true}, ComplexAttributes{RealMin: -1, RealMax: 1, ImagMin: -2, ImagMax: 2, MagnitudeMin: 0.5, MagnitudeMax: 10, AllowNaN: true, AllowInf: true}},
		{"StringAttributes", StringAttributes{MinLen: 1, MaxLen: 5, Prefix: "pre", Suffix: "suf", Contains: "mid", UniqueChars: true}, StringAttributes{MinLen: 1, MaxLen: 5, Prefix: "pre", Suffix: "suf", Contains: "mid", UniqueChars: true}},
		{"SliceAttributes", SliceAttributes{MinLen: 1, MaxLen: 3, Unique: true, Sorted: true, ElementPreds: []Predicate{}, ElementAttrs: IntegerAttributes{}}, SliceAttributes{MinLen: 1, MaxLen: 3, Unique: true, Sorted: true, ElementPreds: []Predicate{}, ElementAttrs: IntegerAttributes{}}},
		{"BoolAttributes", BoolAttributes{ForceTrue: true}, BoolAttributes{ForceTrue: true}},
		{"MapAttributes", MapAttributes{MinSize: 1, MaxSize: 3, KeyPreds: []Predicate{}, ValuePreds: []Predicate{}, KeyAttrs: StringAttributes{}, ValueAttrs: IntegerAttributes{}}, MapAttributes{MinSize: 1, MaxSize: 3, KeyPreds: []Predicate{}, ValuePreds: []Predicate{}, KeyAttrs: StringAttributes{}, ValueAttrs: IntegerAttributes{}}},
		{"ChanAttributes", ChanAttributes{MinBuffer: 0, MaxBuffer: 2, ElemAttrs: FloatAttributes{}}, ChanAttributes{MinBuffer: 0, MaxBuffer: 2, ElemAttrs: FloatAttributes{}}},
		{"FuncAttributes", FuncAttributes{Deterministic: true, PanicProbability: 0.1, ReturnZeroValues: true}, FuncAttributes{Deterministic: true, PanicProbability: 0.1, ReturnZeroValues: true}},
		{"InterfaceAttributes", InterfaceAttributes{AllowedConcrete: []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}}, InterfaceAttributes{AllowedConcrete: []reflect.Type{reflect.TypeOf(1), reflect.TypeOf("")}}},
		{"PointerAttributes", PointerAttributes{AllowNil: true, Depth: 2, Inner: IntegerAttributes{}}, PointerAttributes{AllowNil: true, Depth: 2, Inner: IntegerAttributes{}}},
		{"StructAttributes", StructAttributes{FieldAttrs: map[string]any{"A": IntegerAttributes{}, "B": FloatAttributes{}}}, StructAttributes{FieldAttrs: map[string]any{"A": IntegerAttributes{}, "B": FloatAttributes{}}}},
		{"ArrayAttributes", ArrayAttributes{Length: 3, Sorted: true, ElementAttrs: IntegerAttributes{}, ElementPreds: []Predicate{}}, ArrayAttributes{Length: 3, Sorted: true, ElementAttrs: IntegerAttributes{}, ElementPreds: []Predicate{}}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.in.GetAttributes()
			if reflect.TypeOf(got) != reflect.TypeOf(tc.want) {
				t.Fatalf("expected type %T got %T", tc.want, got)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("mismatch for %s:\nexpected: %#v\n     got: %#v", tc.name, tc.want, got)
			}
		})
	}
}

// TestAttributesInterfaceConformance ensures all structs satisfy the Attributes interface at runtime (already asserted at compile time in source).
func TestAttributesInterfaceConformance(t *testing.T) {
	_ = []Attributes{
		IntegerAttributes{}, FloatAttributes{}, ComplexAttributes{}, StringAttributes{}, SliceAttributes{}, BoolAttributes{}, MapAttributes{}, ChanAttributes{}, FuncAttributes{}, InterfaceAttributes{}, PointerAttributes{}, StructAttributes{}, ArrayAttributes{},
	}
}
