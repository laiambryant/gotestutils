package generation

import (
	"fmt"
	"reflect"

	a "github.com/laiambryant/gotestutils/mtesting/attributes"
)

func GenerateValueForTypeWithAttr(attr a.Attributes, t reflect.Type) (any, error) {
	switch attr := attr.(type) {
	case a.IntegerAttributes:
		return generateIntegerValue(attr)
	case a.FloatAttributes:
		return generateFloatValue(attr)
	case a.ComplexAttributes:
		return generateComplexValue(attr)
	case a.StringAttributes:
		return generateStringValue(attr)
	case a.BoolAttributes:
		return generateBoolValue(attr)
	case a.SliceAttributes:
		return generateSliceValue(attr)
	case a.MapAttributes:
		return generateMapValue(attr)
	case a.PointerAttributes:
		return generatePointerValue(attr)
	case a.StructAttributes:
		return generateStructValue(attr)
	case a.ArrayAttributes:
		return generateArrayValue(attr)
	case a.ChanAttributes:
		return generateChanValue(attr)
	}
	return nil, UnknownTypeError{reflect.TypeOf(attr)}
}

func generateIntegerValue(ia a.IntegerAttributes) (any, error) {
	return nil, nil
}

func generateFloatValue(fa a.FloatAttributes) (any, error) {
	return nil, nil
}

func generateComplexValue(ca a.ComplexAttributes) (any, error) {
	return nil, nil
}

func generateStringValue(sa a.StringAttributes) (any, error) {
	return nil, nil
}

func generateBoolValue(ba a.BoolAttributes) (any, error) {
	return nil, nil
}

func generateSliceValue(sa a.SliceAttributes) (any, error) {
	return nil, nil
}

func generateSliceValueWithDepth(a a.SliceAttributes, depth int) (any, error) {
	return nil, nil
}

func generateMapValue(ma a.MapAttributes) (any, error) {
	return nil, nil
}

func generateMapValueWithDepth(a a.MapAttributes, depth int) (any, error) { return generateMapValue(a) }

func generatePointerValue(pa a.PointerAttributes) (any, error) {
	return nil, nil
}

func generatePointerValueWithDepth(a a.PointerAttributes, depth int) (any, error) { return nil, nil }

func generateStructValue(sa a.StructAttributes) (any, error) {
	return nil, nil
}

func generateStructValueWithDepth(a a.StructAttributes, depth int) (any, error) {
	return nil, nil
}

func generateArrayValue(aa a.ArrayAttributes) (any, error) {
	return nil, nil
}

func generateArrayValueWithDepth(a a.ArrayAttributes, depth int) (any, error) {
	return nil, nil
}

func generateChanValue(ca a.ChanAttributes) (any, error) {
	return nil, nil
}

func genSignedInteger(a a.IntegerAttributes) (any, error) { return nil, nil }

func genUnsignedInteger(a a.IntegerAttributes) (any, error) { return nil, nil }

func enforceSignedZero(val, min, max int64, a a.IntegerAttributes) (any, error) { return nil, nil }

func applyParity(val int64, a a.IntegerAttributes) (any, error) { return nil, nil }

func applyMultipleSigned(val, min, max int64, a a.IntegerAttributes) (any, error) { return nil, nil }

func chooseInSetSigned(current int64, a a.IntegerAttributes) (any, error) { return nil, nil }

func applyExcludeSigned(val int64, a a.IntegerAttributes) (any, error) { return nil, nil }

func isIntKind(k reflect.Kind) (any, error) { return nil, nil }

func isUintKind(k reflect.Kind) (any, error) { return nil, nil }

func randIntWithin(min, max int64) (any, error) { return nil, nil }

func randUintWithin(min, max uint64) (any, error) { return nil, nil }

func alignIntMultiple(val, k, min, max int64) (any, error) { return nil, nil }

func alignUintMultiple(val, k, min, max uint64) (any, error) { return nil, nil }

func inIntExcludeSet(val int64, set []int64) (any, error) { return nil, nil }

func inUintExcludeSet(val uint64, set []int64) (any, error) { return nil, nil }

func float64Pow10(n int) (any, error) { return nil, nil }

type GenerationError struct {
	Op  string
	Msg string
}

func (e GenerationError) Error() string { return fmt.Sprintf("generation %s: %s", e.Op, e.Msg) }
