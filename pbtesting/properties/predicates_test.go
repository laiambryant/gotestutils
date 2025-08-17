package properties

import (
	"math"
	"reflect"
	"testing"

	"github.com/laiambryant/gotestutils/ctesting"
)

func assertProp(t *testing.T, p Predicate, val any, expect bool) {
	if p == nil {
		if !expect {
			t.Fatalf("nil predicate expected=%v but treated as true", expect)
		}
		return
	}
	name := reflect.TypeOf(p).Name()
	got := p.Verify(val)
	if got != expect {
		t.Fatalf("%s.Verify(%#v) = %v, want %v", name, val, got, expect)
	}
}

func TestBoolProperties(t *testing.T) {
	assertProp(t, BoolMustBeTrue{}, true, true)
	assertProp(t, BoolMustBeTrue{}, false, false)
	assertProp(t, BoolMustBeFalse{}, false, true)
	assertProp(t, BoolMustBeFalse{}, true, false)
	assertProp(t, BoolMustBeTrue{}, 1, true)
}

func TestIntProperties(t *testing.T) {
	assertProp(t, IntMin{Min: 5}, 4, false)
	assertProp(t, IntMin{Min: 5}, 5, true)
	assertProp(t, IntMax{Max: 5}, 6, false)
	assertProp(t, IntMax{Max: 5}, 5, true)
	assertProp(t, IntRange{Min: 3, Max: 5}, 2, false)
	assertProp(t, IntRange{Min: 3, Max: 5}, 3, true)
	assertProp(t, IntRange{Min: 3, Max: 5}, 5, true)
	assertProp(t, IntRange{Min: 3, Max: 5}, 6, false)
	assertProp(t, IntNonZero{Required: true}, 0, false)
	assertProp(t, IntNonZero{Required: true}, 7, true)
	assertProp(t, IntNonZero{Required: false}, 0, true)
	assertProp(t, IntEvenOnly{Enabled: true}, 2, true)
	assertProp(t, IntEvenOnly{Enabled: true}, 3, false)
	assertProp(t, IntEvenOnly{Enabled: false}, 0, true)
	assertProp(t, IntOddOnly{Enabled: true}, 3, true)
	assertProp(t, IntOddOnly{Enabled: true}, 2, false)
	assertProp(t, IntOddOnly{Enabled: false}, 0, true)
	assertProp(t, IntMultipleOf{K: 3}, 6, true)
	assertProp(t, IntMultipleOf{K: 3}, 7, false)
	assertProp(t, IntMultipleOf{K: 0}, 7, true)
	assertProp(t, IntInSet{Values: []int64{1, 2, 3}}, 2, true)
	assertProp(t, IntInSet{Values: []int64{1, 2, 3}}, "a", false)
	assertProp(t, IntNotInSet{Values: []int64{1, 2, 3}}, 2, false)
	assertProp(t, IntNotInSet{Values: []int64{1, 2, 3}}, int8(2), false)
	assertProp(t, IntNotInSet{Values: []int64{1, 2, 3}}, int16(4), true)
	assertProp(t, IntNotInSet{Values: []int64{1, 2, 3}}, int32(4), true)
	assertProp(t, IntNotInSet{Values: []int64{1, 2, 3}}, int64(4), true)
	assertProp(t, IntNotInSet{Values: []int64{1, 2, 3}}, "a", true)
	assertProp(t, IntSigned{AllowNegative: false}, -1, false)
	assertProp(t, IntSigned{AllowNegative: false}, 0, true)
	assertProp(t, IntSigned{AllowNegative: true}, 0, true)
	assertProp(t, IntCanIncludeZero{Allowed: false}, 0, false)
	assertProp(t, IntCanIncludeZero{Allowed: false}, 1, true)
	assertProp(t, IntCanIncludeZero{Allowed: true}, 1, true)

}

func TestUintProperties(t *testing.T) {
	assertProp(t, UintMin{Min: 5}, uint(4), false)
	assertProp(t, UintMin{Min: 5}, int(-1), true)
	assertProp(t, UintMin{Min: 5}, int8(-1), true)
	assertProp(t, UintMin{Min: 5}, int16(-1), true)
	assertProp(t, UintMin{Min: 5}, int32(-1), true)
	assertProp(t, UintMin{Min: 5}, int64(-1), true)
	assertProp(t, UintMin{Min: 5}, uint(5), true)
	assertProp(t, UintMax{Max: 5}, uint(6), false)
	assertProp(t, UintMax{Max: 5}, uint(5), true)
	assertProp(t, UintMax{Max: 5}, "ass", true)
	assertProp(t, UintRange{Min: 3, Max: 5}, uint64(2), false)
	assertProp(t, UintRange{Min: 3, Max: 5}, uint(4), true)
	assertProp(t, UintNonZero{Required: true}, uint(0), false)
	assertProp(t, UintNonZero{Required: true}, uint(7), true)
	assertProp(t, UintNonZero{Required: false}, uint(7), true)
	assertProp(t, UintMultipleOf{K: 4}, uint(8), true)
	assertProp(t, UintMultipleOf{K: 0}, uint64(8), true)
	assertProp(t, UintMultipleOf{K: 4}, uint32(9), false)
	assertProp(t, UintInSet{Values: []uint64{1, 2, 3}}, uint8(2), true)
	assertProp(t, UintInSet{Values: []uint64{1, 2, 3}}, uint16(4), false)
	assertProp(t, UintInSet{Values: []uint64{1, 2, 3}}, int(4), false)
	assertProp(t, UintInSet{Values: []uint64{1, 2, 3}}, "johnny", false)
	assertProp(t, UintNotInSet{Values: []uint64{1, 2, 3}}, "joestar", true)
	assertProp(t, UintNotInSet{Values: []uint64{1, 2, 3}}, int8(2), false)
	assertProp(t, UintNotInSet{Values: []uint64{1, 2, 3}}, int16(2), false)
	assertProp(t, UintNotInSet{Values: []uint64{1, 2, 3}}, int32(4), true)
	assertProp(t, UintCanIncludeZero{Allowed: false}, int64(0), false)
	assertProp(t, UintCanIncludeZero{Allowed: true}, uint(1), true)
}

func TestFloatProperties(t *testing.T) {
	assertProp(t, FloatMin{Min: 1.5}, float32(1.4), false)
	assertProp(t, FloatMin{Min: 1.5}, 1.5, true)
	assertProp(t, FloatMax{Max: 2.5}, 2.6, false)
	assertProp(t, FloatMax{Max: 2.5}, 2.5, true)
	assertProp(t, FloatRange{Min: 1, Max: 2}, 0.9, false)
	assertProp(t, FloatRange{Min: 1, Max: 2}, 1.5, true)
	assertProp(t, FloatNonZero{Required: true}, 0.0, false)
	assertProp(t, FloatNonZero{Required: true}, 0.1, true)
	assertProp(t, FloatNonZero{Required: false}, 0.1, true)
	assertProp(t, FloatFiniteOnly{Enabled: true}, math.NaN(), false)
	assertProp(t, FloatFiniteOnly{Enabled: true}, math.Inf(1), false)
	assertProp(t, FloatFiniteOnly{Enabled: false}, math.Inf(1), true)
	assertProp(t, FloatFiniteOnly{Enabled: true}, 3.14, true)
	assertProp(t, FloatAllowNaN{Allowed: false}, math.NaN(), false)
	assertProp(t, FloatAllowNaN{Allowed: true}, math.NaN(), true)
	assertProp(t, FloatAllowInf{Allowed: false}, math.Inf(-1), false)
	assertProp(t, FloatAllowInf{Allowed: true}, math.Inf(-1), true)
	assertProp(t, FloatPrecisionMax{Decimals: 2}, 1.23, true)
	assertProp(t, FloatPrecisionMax{Decimals: 2}, "sw", true)
	assertProp(t, FloatPrecisionMax{Decimals: 2}, 1.234, false)
}

func TestComplexProperties(t *testing.T) {
	assertProp(t, ComplexRealRange{Min: -1, Max: 1}, complex(1.5, 0), false)
	assertProp(t, ComplexRealRange{Min: -1, Max: 1}, complex(0.5, 0), true)
	assertProp(t, ComplexImagRange{Min: -1, Max: 1}, complex64(complex(0, 1.5)), false)
	assertProp(t, ComplexImagRange{Min: -1, Max: 1}, complex128(complex(0, -0.5)), true)
	assertProp(t, ComplexMagnitudeRange{Min: 0, Max: 2}, complex(2, 0), true)
	assertProp(t, ComplexMagnitudeRange{Min: 0, Max: 1}, complex(1, 1), false)
	assertProp(t, ComplexAllowNaN{Allowed: false}, complex(math.NaN(), 0), false)
	assertProp(t, ComplexAllowNaN{Allowed: true}, complex(math.NaN(), 0), true)
	assertProp(t, ComplexAllowInf{Allowed: false}, complex(math.Inf(1), 0), false)
	assertProp(t, ComplexAllowInf{Allowed: true}, complex(math.Inf(1), 0), true)
	assertProp(t, ComplexMagnitudeRange{Min: 0, Max: 2}, "dsada", true)
}

func TestStringProperties(t *testing.T) {
	assertProp(t, StringLenMin{Min: 3}, "ab", false)
	assertProp(t, StringLenMin{Min: 3}, "abc", true)
	assertProp(t, StringLenMax{Max: 3}, "abcd", false)
	assertProp(t, StringLenMax{Max: 3}, "abc", true)
	assertProp(t, StringLenRange{Min: 2, Max: 3}, "a", false)
	assertProp(t, StringLenRange{Min: 2, Max: 3}, "ab", true)
	assertProp(t, StringRegex{Pattern: "^a.+z$"}, "abz", true)
	assertProp(t, StringRegex{Pattern: "^a.+z$"}, "ax", false)
	assertProp(t, StringRegex{Pattern: "((a){10000})"}, 1, false)
	assertProp(t, StringPrefix{Prefix: "pre"}, "prefix", true)
	assertProp(t, StringPrefix{Prefix: "pre"}, "xprefix", false)
	assertProp(t, StringSuffix{Suffix: "suf"}, "endsuf", true)
	assertProp(t, StringSuffix{Suffix: "suf"}, "sufend", false)
	assertProp(t, StringContains{Substr: "mid"}, "amidb", true)
	assertProp(t, StringContains{Substr: "mid"}, "none", false)
}

func TestSliceProperties(t *testing.T) {
	assertProp(t, SliceLenMin{Min: 3}, []int{1, 2}, false)
	assertProp(t, SliceLenMin{Min: 3}, []int{1, 2, 3}, true)
	assertProp(t, SliceLenMin{Min: 3}, "SASA", false)
	assertProp(t, SliceLenMax{Max: 2}, []int{1, 2, 3}, false)
	assertProp(t, SliceLenMax{Max: 3}, []int{1, 2, 3}, true)
	assertProp(t, SliceLenMax{Max: 3}, "SASA", false)
	assertProp(t, SliceLenRange{Min: 2, Max: 3}, []int{1}, false)
	assertProp(t, SliceLenRange{Min: 2, Max: 3}, []int{1, 2}, true)
	assertProp(t, SliceLenRange{Min: 2, Max: 3}, "", false)
	props := []Predicate{IntMin{Min: 2}}
	assertProp(t, SliceElementPredicates{Props: props}, []int{2, 3, 4}, true)
	assertProp(t, SliceElementPredicates{Props: props}, []int{1, 3, 4}, false)
	assertProp(t, SliceElementPredicates{Props: props}, "", false)
}

func TestArrayProperties(t *testing.T) {
	props := []Predicate{IntMin{Min: 2}}
	arrGood := [3]int{2, 3, 4}
	arrBad := [3]int{1, 3, 4}
	assertProp(t, ArrayElementPredicates{Props: props}, arrGood, true)
	assertProp(t, ArrayElementPredicates{Props: props}, arrBad, false)
	assertProp(t, ArrayElementPredicates{Props: props}, "", false)
	assertProp(t, ArraySorted{Enabled: true}, [3]int{1, 2, 3}, true)
	assertProp(t, ArraySorted{Enabled: true}, [3]int{2, 1, 3}, false)
	assertProp(t, ArraySorted{Enabled: false}, "", true)
	assertProp(t, ArraySorted{Enabled: true}, "", false)
}

func TestMapProperties(t *testing.T) {
	m := map[int]int{1: 2, 3: 4}
	keyProps := []Predicate{IntMin{Min: 1}}
	valProps := []Predicate{IntMax{Max: 4}}
	assertProp(t, MapSizeMin{Min: 1}, m, true)
	assertProp(t, MapSizeMin{Min: 1}, "x", false)
	assertProp(t, MapSizeMax{Max: 5}, m, true)
	assertProp(t, MapSizeMax{Max: 5}, "x", false)
	assertProp(t, MapSizeRange{Max: 5}, m, true)
	assertProp(t, MapSizeRange{Max: 5}, "x", false)
	assertProp(t, MapKeyPredicates{Props: keyProps}, m, true)
	assertProp(t, MapKeyPredicates{Props: keyProps}, "x", false)
	assertProp(t, MapValuePredicates{Props: valProps}, m, true)
	assertProp(t, MapValuePredicates{Props: valProps}, "", false)
	assertProp(t, MapKeyPredicates{Props: []Predicate{IntMin{Min: 2}}}, m, false)
	assertProp(t, MapValuePredicates{Props: []Predicate{IntMax{Max: 3}}}, m, false)
}

func TestStructFieldProperties(t *testing.T) {
	type S struct {
		A int
		B string
	}
	vGood := S{A: 5, B: "hello"}
	vBad := S{A: 1, B: "x"}
	props := StructFieldPredicates{Fields: map[string][]Predicate{
		"A": {IntMin{Min: 3}},
		"B": {StringLenMin{Min: 2}},
	}}
	assertProp(t, props, vGood, true)
	assertProp(t, props, vBad, false)
}

func TestStructFieldPropertiesEdgeCases(t *testing.T) {
	type S struct {
		A int
		B string
	}
	vGood := S{A: 5, B: "hello"}

	props := StructFieldPredicates{Fields: map[string][]Predicate{
		"A": {IntMin{Min: 3}},
		"B": {StringLenMin{Min: 2}},
		"C": {SliceElementPredicates{Props: nil}},
	}}
	assertProp(t, props, vGood, true)
	assertProp(t, nil, vGood, true)

}

func TestPointerAllowNil(t *testing.T) {
	var p *int
	assertProp(t, PointerAllowNil{Allowed: true}, p, true)
	assertProp(t, PointerAllowNil{Allowed: false}, p, false)
	x := 3
	p = &x
	assertProp(t, PointerAllowNil{Allowed: false}, p, true)
}

func TestInterfaceAllowedConcrete(t *testing.T) {
	typ := reflect.TypeOf(0)
	assertProp(t, InterfaceAllowedConcrete{Types: []reflect.Type{typ}}, 3, true)
	assertProp(t, InterfaceAllowedConcrete{Types: []reflect.Type{reflect.TypeOf("")}}, 3, false)
}

func TestChanBufferProperties(t *testing.T) {
	ch := make(chan int, 2)
	assertProp(t, ChanBufferMin{Min: 3}, ch, false)
	assertProp(t, ChanBufferMin{Min: 2}, ch, true)
	assertProp(t, ChanBufferMax{Max: 1}, ch, false)
	assertProp(t, ChanBufferMax{Max: 2}, ch, true)
	assertProp(t, ChanBufferRange{Min: 1, Max: 2}, ch, true)
	assertProp(t, ChanBufferRange{Min: 3, Max: 4}, ch, false)
}

func TestLessCharacterization(t *testing.T) {
	suite := []ctesting.CharacterizationTest[bool]{
		ctesting.NewCharacterizationTest(true, nil, func() (bool, error) { return less(1, 2), nil }),
		ctesting.NewCharacterizationTest(false, nil, func() (bool, error) { return less(2, 1), nil }),
		ctesting.NewCharacterizationTest(false, nil, func() (bool, error) { return less(2, 2), nil }),
		ctesting.NewCharacterizationTest(true, nil, func() (bool, error) { return less(int8(1), int8(2)), nil }),
		ctesting.NewCharacterizationTest(false, nil, func() (bool, error) { return less(int8(1), int16(2)), nil }),
		ctesting.NewCharacterizationTest(true, nil, func() (bool, error) { return less(uint8(1), uint8(2)), nil }),
		ctesting.NewCharacterizationTest(false, nil, func() (bool, error) { return less(uint8(1), uint16(2)), nil }),
		ctesting.NewCharacterizationTest(true, nil, func() (bool, error) { return less(-1, 0), nil }),
		ctesting.NewCharacterizationTest(false, nil, func() (bool, error) { return less(uint(1), int(2)), nil }),
		ctesting.NewCharacterizationTest(true, nil, func() (bool, error) { return less("a", "b"), nil }),
		ctesting.NewCharacterizationTest(false, nil, func() (bool, error) { return less("b", "a"), nil }),
		ctesting.NewCharacterizationTest(false, nil, func() (bool, error) { return less("a", "a"), nil }),
		ctesting.NewCharacterizationTest(true, nil, func() (bool, error) { return less(float32(1.1), float32(1.2)), nil }),
		ctesting.NewCharacterizationTest(false, nil, func() (bool, error) { return less(float32(1.1), float64(1.2)), nil }),
		ctesting.NewCharacterizationTest(false, nil, func() (bool, error) { return less([]int{1}, []int{1}), nil }),
		ctesting.NewCharacterizationTest(false, nil, func() (bool, error) { return less(nil, nil), nil }),
		ctesting.NewCharacterizationTest(false, nil, func() (bool, error) { return less(nil, 1), nil }),
		ctesting.NewCharacterizationTest(false, nil, func() (bool, error) { return less(1, nil), nil }),
	}
	ctesting.VerifyCharacterizationTestsAndResults(t, suite, false)
}
