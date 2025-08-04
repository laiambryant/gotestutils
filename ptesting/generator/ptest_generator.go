package generator

import (
	"math/rand"
	"time"
)

func init() {
	rand.NewSource(time.Now().UnixNano())
}

func RandomBool() bool {
	return rand.Intn(2) == 1
}

func RandomInt() int {
	return rand.Int()
}

func RandomInt8() int8 {
	return int8(rand.Intn(256) - 128)
}

func RandomInt16() int16 {
	return int16(rand.Intn(65536) - 32768)
}

func RandomInt32() int32 {
	return rand.Int31()
}

func RandomInt64() int64 {
	return rand.Int63()
}

func RandomUint() uint {
	return uint(rand.Uint32())
}

func RandomUint8() uint8 {
	return uint8(rand.Intn(256))
}

func RandomUint16() uint16 {
	return uint16(rand.Intn(65536))
}

func RandomUint32() uint32 {
	return rand.Uint32()
}

func RandomUint64() uint64 {
	return rand.Uint64()
}

func RandomFloat32() float32 {
	return rand.Float32()
}

func RandomFloat64() float64 {
	return rand.Float64()
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func RandomByte() byte {
	return byte(rand.Intn(256))
}

func RandomRune() rune {
	return rune(rand.Intn(1114112))
}

func RandomComplex64() complex64 {
	return complex(rand.Float32(), rand.Float32())
}

func RandomComplex128() complex128 {
	return complex(rand.Float64(), rand.Float64())
}
