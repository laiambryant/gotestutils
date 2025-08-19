package predicates

import "reflect"

func asInt64(v any) (int64, bool) {
	switch x := v.(type) {
	case int:
		return int64(x), true
	case int8:
		return int64(x), true
	case int16:
		return int64(x), true
	case int32:
		return int64(x), true
	case int64:
		return x, true
	default:
		return 0, false
	}
}
func asUint64(v any) (uint64, bool) {
	switch x := v.(type) {
	case uint:
		return uint64(x), true
	case uint8:
		return uint64(x), true
	case uint16:
		return uint64(x), true
	case uint32:
		return uint64(x), true
	case uint64:
		return x, true
	case int:
		if x >= 0 {
			return uint64(x), true
		}
		return 0, false
	case int8:
		if x >= 0 {
			return uint64(x), true
		}
		return 0, false
	case int16:
		if x >= 0 {
			return uint64(x), true
		}
		return 0, false
	case int32:
		if x >= 0 {
			return uint64(x), true
		}
		return 0, false
	case int64:
		if x >= 0 {
			return uint64(x), true
		}
		return 0, false
	default:
		return 0, false
	}
}
func asFloat64(v any) (float64, bool) {
	switch x := v.(type) {
	case float32:
		return float64(x), true
	case float64:
		return x, true
	default:
		return 0, false
	}
}
func asComplex128(v any) (complex128, bool) {
	switch x := v.(type) {
	case complex64:
		return complex128(x), true
	case complex128:
		return x, true
	default:
		return 0, false
	}
}
func less(a, b any) bool {
	ra := reflect.ValueOf(a)
	rb := reflect.ValueOf(b)
	if ra.Kind() != rb.Kind() {
		return false
	}
	switch ra.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ai, _ := asInt64(a)
		bi, _ := asInt64(b)
		return ai < bi
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		au, _ := asUint64(a)
		bu, _ := asUint64(b)
		return au < bu
	case reflect.Float32, reflect.Float64:
		af, _ := asFloat64(a)
		bf, _ := asFloat64(b)
		return af < bf
	case reflect.String:
		return a.(string) < b.(string)
	default:
		return false
	}
}
func isHashable(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Func, reflect.Map, reflect.Slice, reflect.Struct:
		return false
	default:
		return v.IsValid() && v.Type().Comparable()
	}
}
