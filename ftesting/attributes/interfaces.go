package attributes

import "reflect"

// Attributes is the main interface that all attribute types must implement
type Attributes interface {
	GetAttributes() any
	GetReflectType() reflect.Type
	GetDefaultImplementation() Attributes
	GetRandomValue() any
}

type AttributesStruct interface {
	GetAttributeGivenType(t reflect.Type) (retA Attributes, err error)
}

// Type Interfaces

// Integers defines the constraint for signed integer types
type Integers interface {
	int | int8 | int16 | int32 | int64
}

// UnsignedIntegers defines the constraint for unsigned integer types
type UnsignedIntegers interface {
	uint | uint8 | uint16 | uint32 | uint64
}

// Floats defines the constraint for floating-point types
type Floats interface {
	float32 | float64
}

// Complex defines the constraint for complex number types
type Complex interface {
	complex64 | complex128
}

// Specific Attribute Interfaces for the 4 generic fields

// IntegerAttributeInterface defines the interface for integer attributes
type IntegerAttributes interface {
	GetAttributes() any
	GetReflectType() reflect.Type
	GetDefaultImplementation() Attributes
	GetRandomValue() any
}

// UnsignedIntegerAttributeInterface defines the interface for unsigned integer attributes
type UnsignedIntegerAttributes interface {
	GetAttributes() any
	GetReflectType() reflect.Type
	GetDefaultImplementation() Attributes
	GetRandomValue() any
}

// FloatAttributeInterface defines the interface for float attributes
type FloatAttributes interface {
	GetAttributes() any
	GetReflectType() reflect.Type
	GetDefaultImplementation() Attributes
	GetRandomValue() any
}

// ComplexAttributeInterface defines the interface for complex attributes
type ComplexAttributes interface {
	GetAttributes() any
	GetReflectType() reflect.Type
	GetDefaultImplementation() Attributes
	GetRandomValue() any
}
