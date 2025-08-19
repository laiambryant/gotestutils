package predicates

import "testing"

func TestPredicateImplementations(t *testing.T) {
	_ = []Predicate{
		BoolMustBeTrue{}, BoolMustBeFalse{},
		IntMin{}, IntMax{}, IntRange{}, IntNonZero{}, IntEvenOnly{}, IntOddOnly{}, IntMultipleOf{}, IntInSet{}, IntNotInSet{}, IntSigned{}, IntCanIncludeZero{},
		UintMin{}, UintMax{}, UintRange{}, UintNonZero{}, UintMultipleOf{}, UintInSet{}, UintNotInSet{}, UintCanIncludeZero{},
		FloatMin{}, FloatMax{}, FloatRange{}, FloatNonZero{}, FloatFiniteOnly{}, FloatAllowNaN{}, FloatAllowInf{}, FloatPrecisionMax{},
		ComplexRealRange{}, ComplexImagRange{}, ComplexMagnitudeRange{}, ComplexAllowNaN{}, ComplexAllowInf{},
		StringLenMin{}, StringLenMax{}, StringLenRange{}, StringRegex{}, StringPrefix{}, StringSuffix{}, StringContains{},
		SliceLenMin{}, SliceLenMax{}, SliceLenRange{}, SliceUnique{}, SliceElementPredicates{},
		ArrayElementPredicates{}, ArraySorted{},
		MapSizeMin{}, MapSizeMax{}, MapSizeRange{}, MapKeyPredicates{}, MapValuePredicates{},
		StructFieldPredicates{}, PointerAllowNil{}, InterfaceAllowedConcrete{},
		ChanBufferMin{}, ChanBufferMax{}, ChanBufferRange{},
	}
	for _, p := range []Predicate{BoolMustBeTrue{}, BoolMustBeFalse{}} {
		_ = p.Verify(nil)
	}
}
