package pbtesting

import (
	"math/rand"
	"reflect"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

const (
	defaultMaxSliceLen  = 8
	defaultMaxMapLen    = 6
	defaultMaxStringLen = 16
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randLen(max int) int {
	if max <= 0 {
		return 0
	}
	return rng.Intn(max + 1)
}

func randBool() bool { return rng.Intn(2) == 1 }

func randInt(bits int) int64 {
	if bits <= 0 || bits > 63 {
		bits = 63
	}
	max := int64(1)<<(bits-1) - 1
	min := -int64(1) << (bits - 1)
	if max <= min {
		return 0
	}
	rangeSize := max - min + 1
	if rangeSize <= 0 {
		return 0
	}
	return min + rng.Int63n(rangeSize)
}

func randUint(bits int) uint64 {
	if bits <= 0 || bits > 64 {
		bits = 64
	}
	if bits == 64 {
		return uint64(rng.Int63())<<1 | uint64(rng.Intn(2))
	}
	return uint64(rng.Int63n(1 << bits))
}

func randFloat(bits int) float64 {
	switch bits {
	case 32:
		return float64(rng.Float32())
	default:
		return rng.Float64()
	}
}

func randComplex(bits int) complex128 {
	re := randFloat(bits)
	im := randFloat(bits)
	return complex(re, im)
}

func randString(n int) string {
	if n < 0 {
		n = 0
	}
	runes := make([]rune, n)
	for i := 0; i < n; i++ {
		runes[i] = letters[rng.Intn(len(letters))]
	}
	return string(runes)
}

func makeZeroValue(t reflect.Type) reflect.Value { return reflect.Zero(t) }

func getRandomValue(v reflect.Value) {
	if !v.IsValid() || !v.CanSet() {
		return
	}

	switch v.Kind() {
	case reflect.Bool:
		setRandomBool(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		setRandomSignedInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		setRandomUnsignedInt(v)
	case reflect.Float32:
		setRandomFloat(v, 32)
	case reflect.Float64:
		setRandomFloat(v, 64)
	case reflect.Complex64:
		setRandomComplex(v, 32)
	case reflect.Complex128:
		setRandomComplex(v, 64)
	case reflect.String:
		setRandomString(v)
	case reflect.Slice:
		setRandomSlice(v)
	case reflect.Array:
		setRandomArray(v)
	case reflect.Map:
		setRandomMap(v)
	case reflect.Struct:
		setRandomStruct(v)
	case reflect.Pointer:
		setRandomPointer(v)
	case reflect.Interface:
		return
	case reflect.Func:
		setRandomFunc(v)
	case reflect.Chan:
		setRandomChan(v)
	case reflect.UnsafePointer:
		return
	default:
		return
	}
}

func setRandomBool(v reflect.Value)              { v.SetBool(randBool()) }
func setRandomSignedInt(v reflect.Value)         { bits := v.Type().Bits(); v.SetInt(randInt(int(bits))) }
func setRandomUnsignedInt(v reflect.Value)       { bits := v.Type().Bits(); v.SetUint(randUint(int(bits))) }
func setRandomFloat(v reflect.Value, bits int)   { v.SetFloat(randFloat(bits)) }
func setRandomComplex(v reflect.Value, bits int) { v.SetComplex(randComplex(bits)) }
func setRandomString(v reflect.Value)            { v.SetString(randString(randLen(defaultMaxStringLen))) }

func setRandomSlice(v reflect.Value) {
	ln := randLen(defaultMaxSliceLen)
	slice := reflect.MakeSlice(v.Type(), ln, ln)
	for i := 0; i < ln; i++ {
		getRandomValue(slice.Index(i))
	}
	v.Set(slice)
}

func setRandomArray(v reflect.Value) {
	ln := v.Len()
	arr := reflect.New(v.Type()).Elem()
	for i := 0; i < ln; i++ {
		getRandomValue(arr.Index(i))
	}
	v.Set(arr)
}

func setRandomMap(v reflect.Value) {
	ln := randLen(defaultMaxMapLen)
	m := reflect.MakeMapWithSize(v.Type(), ln)
	keyT := v.Type().Key()
	valT := v.Type().Elem()
	for i := 0; i < ln; i++ {
		k := reflect.New(keyT).Elem()
		getRandomValue(k)
		val := reflect.New(valT).Elem()
		getRandomValue(val)
		m.SetMapIndex(k, val)
	}
	v.Set(m)
}

func setRandomStruct(v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.CanSet() {
			getRandomValue(f)
		}
	}
}

func setRandomPointer(v reflect.Value) {
	elemT := v.Type().Elem()
	elem := reflect.New(elemT)
	v.Set(elem)
	getRandomValue(v.Elem())
}

func setRandomFunc(v reflect.Value) {
	fnType := v.Type()
	stub := reflect.MakeFunc(fnType, func(args []reflect.Value) []reflect.Value {
		outs := make([]reflect.Value, fnType.NumOut())
		for i := 0; i < fnType.NumOut(); i++ {
			outs[i] = makeZeroValue(fnType.Out(i))
		}
		return outs
	})
	v.Set(stub)
}

func setRandomChan(v reflect.Value) { v.Set(reflect.MakeChan(v.Type(), 1)) }
