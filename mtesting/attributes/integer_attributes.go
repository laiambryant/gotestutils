package attributes

import "reflect"

type Integers interface {
	int | int8 | int16 | int32 | int64
}

type IntegerAttributes[T Integers] struct {
	AllowNegative bool
	AllowZero     bool
	Max           T
	Min           T
	InSet         []T
	NotInSet      []T
}

func (a IntegerAttributes[T]) GetAttributes() any { return a }
func (a IntegerAttributes[T]) GetReflectType() reflect.Type {
	return reflect.TypeOf(*new(T))
}

func (a IntegerAttributes[T]) GetDefaultImplementation() Attributes {
	return IntegerAttributes[T]{
		AllowNegative: true,
		AllowZero:     true,
		Max:           100,
		Min:           -100,
	}
}

type UnsignedIntegers interface {
	uint | uint8 | uint16 | uint32 | uint64
}

type UnsignedIntegerAttributes[T UnsignedIntegers] struct {
	Signed        bool
	AllowNegative bool
	AllowZero     bool
	Max           T
	Min           T
	InSet         []T
	NotInSet      []T
}

func (a UnsignedIntegerAttributes[T]) GetAttributes() any { return a }
func (a UnsignedIntegerAttributes[T]) GetReflectType() reflect.Type {
	if a.Signed || a.AllowNegative {
		return reflect.TypeOf(int64(0))
	}
	return reflect.TypeOf(uint64(0))
}

func (a UnsignedIntegerAttributes[T]) GetDefaultImplementation() Attributes {
	return UnsignedIntegerAttributes[T]{
		Signed:        true,
		AllowNegative: true,
		AllowZero:     true,
		Max:           100,
		Min:           0,
	}
}
