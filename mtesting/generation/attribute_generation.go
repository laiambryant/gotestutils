package generation

import (
	"reflect"

	a "github.com/laiambryant/gotestutils/mtesting/attributes"
)

func GenerateValueForTypeWithAttr(attr a.Attributes, t reflect.Type) (any, error) {
	switch attr := attr.(type) {
	case a.IntegerAttributes[int]:
		return generateIntegerValue(attr)
	case a.FloatAttributes[float32]:
		return generateFloatValue(attr)
	case a.ComplexAttributes[complex64]:
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

func generateIntegerValue[T a.Integers](ia a.IntegerAttributes[T]) (any, error) {
	return nil, nil
}
func generateFloatValue[T a.Floats](fa a.FloatAttributes[T]) (any, error) {
	return nil, nil
}

func generateComplexValue[T a.Complex](ca a.ComplexAttributes[T]) (any, error) {
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
