package generation

import (
	"reflect"

	"github.com/laiambryant/gotestutils/mtesting/attributes"
	a "github.com/laiambryant/gotestutils/mtesting/attributes"
)

func GenerateValueForTypeWithAttr(attr attributes.Attributes) (reflect.Value, error) {
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
	return reflect.Value{}, nil
}

func generateIntegerValue(a a.IntegerAttributes) (reflect.Value, error) {
	return reflect.Value{}, nil
}

func generateFloatValue(a a.FloatAttributes) (reflect.Value, error) {
	return reflect.Value{}, nil
}

func generateComplexValue(a a.ComplexAttributes) (reflect.Value, error) {
	return reflect.Value{}, nil
}

func generateStringValue(a a.StringAttributes) (reflect.Value, error) {
	return reflect.Value{}, nil
}

func generateBoolValue(a a.BoolAttributes) (reflect.Value, error) {
	return reflect.Value{}, nil
}

func generateSliceValue(a a.SliceAttributes) (reflect.Value, error) {
	// For now only depth 0
	return reflect.Value{}, nil
}

// TODO: Implement this at a later moment
func generateSliceValueWithDepth(a a.SliceAttributes, depth int) (reflect.Value, error) {
	return reflect.Value{}, nil
}

func generateMapValue(a a.MapAttributes) (reflect.Value, error) {
	return reflect.Value{}, nil
}

// TODO: Implement this at a later moment
func generateMapValueWithDepth(a a.MapAttributes, depth int) (reflect.Value, error) {
	return reflect.Value{}, nil
}

func generatePointerValue(a a.PointerAttributes) (reflect.Value, error) {
	return reflect.Value{}, nil
}

// TODO: Implement this at a later moment
func generatePointerValueWithDepth(a a.PointerAttributes, depth int) (reflect.Value, error) {
	return reflect.Value{}, nil
}

func generateStructValue(a a.StructAttributes) (reflect.Value, error) {
	return reflect.Value{}, nil
}

// TODO: Implement this at a later moment
func generateStructValueWithDepth(a a.StructAttributes, depth int) (reflect.Value, error) {
	return reflect.Value{}, nil
}

func generateArrayValue(a a.ArrayAttributes) (reflect.Value, error) {
	return reflect.Value{}, nil
}

// TODO: Implement this at a later moment
func generateArrayValueWithDepth(a a.ArrayAttributes, depth int) (reflect.Value, error) {
	return reflect.Value{}, nil
}

func generateChanValue(a a.ChanAttributes) (reflect.Value, error) {
	return reflect.Value{}, nil
}

func genSignedInteger(a a.IntegerAttributes) (reflect.Value, error) { return reflect.Value{}, nil }

func genUnsignedInteger(a a.IntegerAttributes) (reflect.Value, error) { return reflect.Value{}, nil }

func enforceSignedZero(val, min, max int64, a a.IntegerAttributes) (reflect.Value, error) {
	return reflect.Value{}, nil
}

func applyParity(val int64, a a.IntegerAttributes) (reflect.Value, error) {
	return reflect.Value{}, nil
}

func applyMultipleSigned(val, min, max int64, a a.IntegerAttributes) (reflect.Value, error) {
	return reflect.Value{}, nil
}

func chooseInSetSigned(current int64, a a.IntegerAttributes) (reflect.Value, error) {
	return reflect.Value{}, nil
}

func applyExcludeSigned(val int64, a a.IntegerAttributes) (reflect.Value, error) {
	return reflect.Value{}, nil
}

func isIntKind(k reflect.Kind) (reflect.Value, error) { return reflect.Value{}, nil }

func isUintKind(k reflect.Kind) (reflect.Value, error) { return reflect.Value{}, nil }

func randIntWithin(min, max int64) (reflect.Value, error) { return reflect.Value{}, nil }

func randUintWithin(min, max uint64) (reflect.Value, error) { return reflect.Value{}, nil }

func alignIntMultiple(val, k, min, max int64) (reflect.Value, error) { return reflect.Value{}, nil }

func alignUintMultiple(val, k, min, max uint64) (reflect.Value, error) { return reflect.Value{}, nil }

func inIntExcludeSet(val int64, set []int64) (reflect.Value, error) { return reflect.Value{}, nil }

func inUintExcludeSet(val uint64, set []int64) (reflect.Value, error) { return reflect.Value{}, nil }

func float64Pow10(n int) (reflect.Value, error) { return reflect.Value{}, nil }
