package attributes

import (
	"reflect"
)

// Shared helper types used across attribute tests
type CustomInt int

type nilReturningAttribute struct{}

func (n nilReturningAttribute) GetAttributes() any                   { return n }
func (n nilReturningAttribute) GetReflectType() reflect.Type         { return reflect.TypeOf(0) }
func (n nilReturningAttribute) GetRandomValue() any                  { return nil }
func (n nilReturningAttribute) GetDefaultImplementation() Attributes { return n }

type nilAttributeType struct{}

func (n nilAttributeType) GetAttributes() any           { return nil }
func (n nilAttributeType) GetReflectType() reflect.Type { return reflect.TypeOf(int(0)) }
func (n nilAttributeType) GetDefaultImplementation() Attributes {
	return IntegerAttributesImpl[int]{AllowNegative: true, AllowZero: true, Max: 100, Min: -100}
}
func (n nilAttributeType) GetRandomValue() any { return 42 }

type nilTypeAttribute struct{}

func (n nilTypeAttribute) GetAttributes() any {
	var nilMap map[string]int
	return nilMap
}
func (n nilTypeAttribute) GetReflectType() reflect.Type { return reflect.TypeOf(int(0)) }
func (n nilTypeAttribute) GetDefaultImplementation() Attributes {
	return IntegerAttributesImpl[int]{AllowNegative: true, AllowZero: true, Max: 100, Min: -100}
}
func (n nilTypeAttribute) GetRandomValue() any { return 42 }

type nilTypeReturningAttribute struct{}

func (n nilTypeReturningAttribute) GetAttributes() any { return n }
func (n nilTypeReturningAttribute) GetReflectType() reflect.Type {
	return nil // Returns nil type
}
func (n nilTypeReturningAttribute) GetRandomValue() any                  { return nil }
func (n nilTypeReturningAttribute) GetDefaultImplementation() Attributes { return n }

// Small helper used in some tests
func isNilValidForType(attr Attributes) bool {
	switch attr.(type) {
	case PointerAttributes:
		return true
	default:
		return false
	}
}
