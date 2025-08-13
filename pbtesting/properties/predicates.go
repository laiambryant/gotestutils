package properties

import (
	"math"
	"reflect"
	"regexp"
	"slices"
)

type Predicate interface {
	Verify(any) bool
}

type BoolMustBeTrue struct{}

type BoolMustBeFalse struct{}

func (p BoolMustBeTrue) Verify(v any) bool  { b, ok := v.(bool); return !ok || b }
func (p BoolMustBeFalse) Verify(v any) bool { b, ok := v.(bool); return !ok || !b }

// Signed integers

type IntMin struct{ Min int64 }
type IntMax struct{ Max int64 }
type IntRange struct{ Min, Max int64 }
type IntNonZero struct{ Required bool }
type IntEvenOnly struct{ Enabled bool }
type IntOddOnly struct{ Enabled bool }
type IntMultipleOf struct{ K int64 }
type IntInSet struct{ Values []int64 }
type IntNotInSet struct{ Values []int64 }
type IntSigned struct{ AllowNegative bool }
type IntCanIncludeZero struct{ Allowed bool }

func (p IntMin) Verify(v any) bool   { n, ok := asInt64(v); return !ok || n >= p.Min }
func (p IntMax) Verify(v any) bool   { n, ok := asInt64(v); return !ok || n <= p.Max }
func (p IntRange) Verify(v any) bool { n, ok := asInt64(v); return !ok || (n >= p.Min && n <= p.Max) }
func (p IntNonZero) Verify(v any) bool {
	if !p.Required {
		return true
	}
	n, ok := asInt64(v)
	return !ok || n != 0
}
func (p IntEvenOnly) Verify(v any) bool {
	if !p.Enabled {
		return true
	}
	n, ok := asInt64(v)
	return !ok || n%2 == 0
}
func (p IntOddOnly) Verify(v any) bool {
	if !p.Enabled {
		return true
	}
	n, ok := asInt64(v)
	return !ok || n%2 != 0
}
func (p IntMultipleOf) Verify(v any) bool {
	if p.K == 0 {
		return true
	}
	n, ok := asInt64(v)
	return !ok || n%p.K == 0
}
func (p IntInSet) Verify(v any) bool {
	n, ok := asInt64(v)
	if !ok {
		return false
	}
	return slices.Contains(p.Values, n)
}
func (p IntNotInSet) Verify(v any) bool {
	n, ok := asInt64(v)
	if !ok {
		return true
	}
	return !slices.Contains(p.Values, n)
}
func (p IntSigned) Verify(v any) bool {
	if p.AllowNegative {
		return true
	}
	n, ok := asInt64(v)
	return !ok || n >= 0
}
func (p IntCanIncludeZero) Verify(v any) bool {
	if p.Allowed {
		return true
	}
	n, ok := asInt64(v)
	return !ok || n != 0
}

// Unsigned integers

type UintMin struct{ Min uint64 }
type UintMax struct{ Max uint64 }
type UintRange struct{ Min, Max uint64 }
type UintNonZero struct{ Required bool }
type UintMultipleOf struct{ K uint64 }
type UintInSet struct{ Values []uint64 }
type UintNotInSet struct{ Values []uint64 }
type UintCanIncludeZero struct{ Allowed bool }

func (p UintMin) Verify(v any) bool { n, ok := asUint64(v); return !ok || n >= p.Min }
func (p UintMax) Verify(v any) bool { n, ok := asUint64(v); return !ok || n <= p.Max }
func (p UintRange) Verify(v any) bool {
	n, ok := asUint64(v)
	return !ok || (n >= p.Min && n <= p.Max)
}
func (p UintNonZero) Verify(v any) bool {
	if !p.Required {
		return true
	}
	n, ok := asUint64(v)
	return !ok || n != 0
}
func (p UintMultipleOf) Verify(v any) bool {
	if p.K != 0 {
		n, ok := asUint64(v)
		return !ok || n%p.K == 0
	} else {
		return true
	}
}
func (p UintInSet) Verify(v any) bool {
	n, ok := asUint64(v)
	if !ok {
		return false
	}
	return slices.Contains(p.Values, n)
}
func (p UintNotInSet) Verify(v any) bool {
	n, ok := asUint64(v)
	if !ok {
		return true
	}
	return !slices.Contains(p.Values, n)
}
func (p UintCanIncludeZero) Verify(v any) bool {
	if p.Allowed {
		return true
	}
	n, ok := asUint64(v)
	return !ok || n != 0
}

// Floats

type FloatMin struct{ Min float64 }
type FloatMax struct{ Max float64 }
type FloatRange struct{ Min, Max float64 }
type FloatNonZero struct{ Required bool }
type FloatFiniteOnly struct{ Enabled bool }
type FloatAllowNaN struct{ Allowed bool }
type FloatAllowInf struct{ Allowed bool }
type FloatPrecisionMax struct{ Decimals int }

func (p FloatMin) Verify(v any) bool { n, ok := asFloat64(v); return !ok || n >= p.Min }
func (p FloatMax) Verify(v any) bool { n, ok := asFloat64(v); return !ok || n <= p.Max }
func (p FloatRange) Verify(v any) bool {
	n, ok := asFloat64(v)
	return !ok || (n >= p.Min && n <= p.Max)
}
func (p FloatNonZero) Verify(v any) bool {
	if !p.Required {
		return true
	}
	n, ok := asFloat64(v)
	return !ok || n != 0
}
func (p FloatFiniteOnly) Verify(v any) bool {
	if !p.Enabled {
		return true
	}
	n, ok := asFloat64(v)
	return !ok || (!math.IsNaN(n) && !math.IsInf(n, 0))
}
func (p FloatAllowNaN) Verify(v any) bool {
	if p.Allowed {
		return true
	}
	n, ok := asFloat64(v)
	return !ok || !math.IsNaN(n)
}
func (p FloatAllowInf) Verify(v any) bool {
	if p.Allowed {
		return true
	}
	n, ok := asFloat64(v)
	return !ok || !math.IsInf(n, 0)
}
func (p FloatPrecisionMax) Verify(v any) bool {
	n, ok := asFloat64(v)
	if !ok || p.Decimals <= 0 {
		return true
	}
	m := math.Pow10(p.Decimals)
	return math.Abs(n*m-math.Round(n*m)) < 1e-9
}

// Complex

type ComplexRealRange struct{ Min, Max float64 }
type ComplexImagRange struct{ Min, Max float64 }
type ComplexMagnitudeRange struct{ Min, Max float64 }
type ComplexAllowNaN struct{ Allowed bool }
type ComplexAllowInf struct{ Allowed bool }

func (p ComplexRealRange) Verify(v any) bool {
	c, ok := asComplex128(v)
	return !ok || (real(c) >= p.Min && real(c) <= p.Max)
}
func (p ComplexImagRange) Verify(v any) bool {
	c, ok := asComplex128(v)
	return !ok || (imag(c) >= p.Min && imag(c) <= p.Max)
}
func (p ComplexMagnitudeRange) Verify(v any) bool {
	c, ok := asComplex128(v)
	if !ok {
		return true
	}
	r := math.Hypot(real(c), imag(c))
	return r >= p.Min && r <= p.Max
}
func (p ComplexAllowNaN) Verify(v any) bool {
	c, ok := asComplex128(v)
	if !ok || p.Allowed {
		return true
	}
	return !(math.IsNaN(real(c)) || math.IsNaN(imag(c)))
}
func (p ComplexAllowInf) Verify(v any) bool {
	c, ok := asComplex128(v)
	if !ok || p.Allowed {
		return true
	}
	return !(math.IsInf(real(c), 0) || math.IsInf(imag(c), 0))
}

// String

type StringLenMin struct{ Min int }
type StringLenMax struct{ Max int }
type StringLenRange struct{ Min, Max int }
type StringAllowedRunes struct{ Runes []rune }
type StringRegex struct{ Pattern string }
type StringPrefix struct{ Prefix string }
type StringSuffix struct{ Suffix string }
type StringContains struct{ Substr string }

func (p StringLenMin) Verify(v any) bool { s, ok := v.(string); return !ok || len(s) >= p.Min }
func (p StringLenMax) Verify(v any) bool { s, ok := v.(string); return !ok || len(s) <= p.Max }
func (p StringLenRange) Verify(v any) bool {
	s, ok := v.(string)
	return !ok || (len(s) >= p.Min && len(s) <= p.Max)
}
func (p StringAllowedRunes) Verify(v any) bool {
	s, ok := v.(string)
	if !ok || len(p.Runes) == 0 {
		return true
	}
	allowed := make(map[rune]struct{}, len(p.Runes))
	for _, r := range p.Runes {
		allowed[r] = struct{}{}
	}
	for _, r := range s {
		if _, ok := allowed[r]; !ok {
			return false
		}
	}
	return true
}
func (p StringRegex) Verify(v any) bool {
	s, ok := v.(string)
	if !ok || p.Pattern == "" {
		return true
	}
	re, err := regexp.Compile(p.Pattern)
	if err != nil {
		return true
	}
	return re.MatchString(s)
}
func (p StringPrefix) Verify(v any) bool {
	s, ok := v.(string)
	return !ok || (len(p.Prefix) == 0 || (len(s) >= len(p.Prefix) && s[:len(p.Prefix)] == p.Prefix))
}
func (p StringSuffix) Verify(v any) bool {
	s, ok := v.(string)
	return !ok || (len(p.Suffix) == 0 || (len(s) >= len(p.Suffix) && s[len(s)-len(p.Suffix):] == p.Suffix))
}
func (p StringContains) Verify(v any) bool {
	s, ok := v.(string)
	return !ok || (p.Substr == "" || (len(s) >= len(p.Substr) && (regexp.MustCompile(regexp.QuoteMeta(p.Substr)).FindStringIndex(s) != nil)))
}

// Slice

type SliceLenMin struct{ Min int }
type SliceLenMax struct{ Max int }
type SliceLenRange struct{ Min, Max int }
type SliceUnique struct{ Enabled bool }
type SliceElementPredicates struct{ Props []Predicate }

func (p SliceLenMin) Verify(v any) bool {
	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() != reflect.Slice {
		return true
	}
	return reflectValue.Len() >= p.Min
}
func (p SliceLenMax) Verify(v any) bool {
	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() != reflect.Slice {
		return true
	}
	return reflectValue.Len() <= p.Max
}
func (p SliceLenRange) Verify(v any) bool {
	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() != reflect.Slice {
		return true
	}
	return reflectValue.Len() >= p.Min && reflectValue.Len() <= p.Max
}
func (p SliceUnique) Verify(v any) bool {
	if !p.Enabled {
		return true
	}
	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() != reflect.Slice {
		return true
	}
	seen := map[any]struct{}{}
	for i := 0; i < reflectValue.Len(); i++ {
		elem := reflectValue.Index(i).Interface()
		if isHashable(reflectValue.Index(i)) {
			if _, ok := seen[elem]; ok {
				return false
			}
			seen[elem] = struct{}{}
		}
	}
	return true
}
func (p SliceElementPredicates) Verify(v any) bool {
	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() != reflect.Slice || len(p.Props) == 0 {
		return true
	}
	for i := 0; i < reflectValue.Len(); i++ {
		for _, prop := range p.Props {
			if !prop.Verify(reflectValue.Index(i).Interface()) {
				return false
			}
		}
	}
	return true
}

// Array

type ArrayElementPredicates struct{ Props []Predicate }
type ArraySorted struct{ Enabled bool }

func (p ArrayElementPredicates) Verify(v any) bool {
	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() != reflect.Array || len(p.Props) == 0 {
		return true
	}
	for i := 0; i < reflectValue.Len(); i++ {
		for _, prop := range p.Props {
			if !prop.Verify(reflectValue.Index(i).Interface()) {
				return false
			}
		}
	}
	return true
}
func (p ArraySorted) Verify(v any) bool {
	if !p.Enabled {
		return true
	}
	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() != reflect.Array || reflectValue.Len() < 2 {
		return true
	}
	for i := 1; i < reflectValue.Len(); i++ {
		if less(reflectValue.Index(i).Interface(), reflectValue.Index(i-1).Interface()) {
			return false
		}
	}
	return true
}

// Map

type MapSizeMin struct{ Min int }
type MapSizeMax struct{ Max int }
type MapSizeRange struct{ Min, Max int }
type MapKeyPredicates struct{ Props []Predicate }
type MapValuePredicates struct{ Props []Predicate }

func (p MapSizeMin) Verify(v any) bool {
	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() != reflect.Map {
		return true
	}
	return reflectValue.Len() >= p.Min
}
func (p MapSizeMax) Verify(v any) bool {
	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() != reflect.Map {
		return true
	}
	return reflectValue.Len() <= p.Max
}
func (p MapSizeRange) Verify(v any) bool {
	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() != reflect.Map {
		return true
	}
	return reflectValue.Len() >= p.Min && reflectValue.Len() <= p.Max
}
func (p MapKeyPredicates) Verify(v any) bool {
	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() != reflect.Map || len(p.Props) == 0 {
		return true
	}
	iter := reflectValue.MapRange()
	for iter.Next() {
		key := iter.Key().Interface()
		for _, prop := range p.Props {
			if !prop.Verify(key) {
				return false
			}
		}
	}
	return true
}
func (p MapValuePredicates) Verify(v any) bool {
	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() != reflect.Map || len(p.Props) == 0 {
		return true
	}
	iter := reflectValue.MapRange()
	for iter.Next() {
		val := iter.Value().Interface()
		for _, prop := range p.Props {
			if !prop.Verify(val) {
				return false
			}
		}
	}
	return true
}

// Struct

type StructFieldPredicates struct{ Fields map[string][]Predicate }

func (p StructFieldPredicates) Verify(v any) bool {
	if len(p.Fields) == 0 {
		return true
	}
	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() == reflect.Pointer {
		reflectValue = reflectValue.Elem()
	}
	if reflectValue.Kind() != reflect.Struct {
		return true
	}
	rt := reflectValue.Type()
	for name, props := range p.Fields {
		f, ok := rt.FieldByName(name)
		if !ok {
			continue
		}
		fv := reflectValue.FieldByIndex(f.Index)
		for _, prop := range props {
			if !prop.Verify(fv.Interface()) {
				return false
			}
		}
	}
	return true
}

// Pointer

type PointerAllowNil struct{ Allowed bool }

func (p PointerAllowNil) Verify(v any) bool {
	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() != reflect.Pointer {
		return true
	}
	if reflectValue.IsNil() {
		return p.Allowed
	}
	return true
}

// Interface

type InterfaceAllowedConcrete struct{ Types []reflect.Type }

func (p InterfaceAllowedConcrete) Verify(v any) bool {
	if v == nil || len(p.Types) == 0 {
		return true
	}
	rt := reflect.TypeOf(v)
	for _, t := range p.Types {
		if rt.AssignableTo(t) {
			return true
		}
	}
	return false
}

// Chan

type ChanBufferMin struct{ Min int }
type ChanBufferMax struct{ Max int }
type ChanBufferRange struct{ Min, Max int }

func (p ChanBufferMin) Verify(v any) bool {
	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() != reflect.Chan {
		return true
	}
	return reflectValue.Cap() >= p.Min
}
func (p ChanBufferMax) Verify(v any) bool {
	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() != reflect.Chan {
		return true
	}
	return reflectValue.Cap() <= p.Max
}
func (p ChanBufferRange) Verify(v any) bool {
	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() != reflect.Chan {
		return true
	}
	return reflectValue.Cap() >= p.Min && reflectValue.Cap() <= p.Max
}

// Helpers

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
		sa := a.(string)
		sb := b.(string)
		return sa < sb
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
