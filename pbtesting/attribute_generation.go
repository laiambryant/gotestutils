package pbtesting

import (
	"reflect"

	properties "github.com/laiambryant/gotestutils/pbtesting/properties"
)

func generateValueForTypeWithAttr(t reflect.Type, attr any, depth int) (reflect.Value, error) {
	return reflect.Value{}, nil
}

func generateIntegerValue(v reflect.Value, a properties.IntegerAttributes) error { return nil }

func generateFloatValue(v reflect.Value, a properties.FloatAttributes) error { return nil }

func generateComplexValue(v reflect.Value, a properties.ComplexAttributes) {}

func generateStringValue(v reflect.Value, a properties.StringAttributes) {}

func generateBoolValue(v reflect.Value, a properties.BoolAttributes) {}

func generateSliceValue(v reflect.Value, a properties.SliceAttributes, depth int) error { return nil }

func generateMapValue(v reflect.Value, a properties.MapAttributes, depth int) error { return nil }

func generatePointerValue(v reflect.Value, a properties.PointerAttributes, depth int) error {
	return nil
}

func generateStructValue(v reflect.Value, a properties.StructAttributes, depth int) {}

func generateArrayValue(v reflect.Value, a properties.ArrayAttributes, depth int) {}

func generateChanValue(v reflect.Value, a properties.ChanAttributes) {}

func genSignedInteger(a properties.IntegerAttributes) int64 { return 0 }

func genUnsignedInteger(a properties.IntegerAttributes) uint64 { return 0 }

func enforceSignedZero(val, min, max int64, a properties.IntegerAttributes) int64 { return 0 }

func applyParity(val int64, a properties.IntegerAttributes) int64 { return 0 }

func applyMultipleSigned(val, min, max int64, a properties.IntegerAttributes) int64 { return 0 }

func chooseInSetSigned(current int64, a properties.IntegerAttributes) int64 { return 0 }

func applyExcludeSigned(val int64, a properties.IntegerAttributes) int64 { return 0 }

func isIntKind(k reflect.Kind) bool { return false }

func isUintKind(k reflect.Kind) bool { return false }

func randIntWithin(min, max int64) int64 { return 0 }

func randUintWithin(min, max uint64) uint64 { return 0 }

func alignIntMultiple(val, k, min, max int64) int64 { return 0 }

func alignUintMultiple(val, k, min, max uint64) uint64 { return 0 }

func inIntExcludeSet(val int64, set []int64) bool { return false }

func inUintExcludeSet(val uint64, set []int64) bool { return false }

func float64Pow10(n int) float64 { return 0 }
