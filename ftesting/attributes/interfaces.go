package attributes

import "reflect"

// Attributes is the main interface that all attribute types must implement.
// It defines the contract for generating random values for any Go type.
//
// Implementations of this interface are responsible for:
// 1. Storing configuration about how to generate values (via GetAttributes)
// 2. Reporting the Go type they generate values for (via GetReflectType)
// 3. Providing a default configuration (via GetDefaultImplementation)
// 4. Generating random values according to their configuration (via GetRandomValue)
//
// Methods:
//   - GetAttributes() any: Returns the configuration struct for this attribute
//   - GetReflectType() reflect.Type: Returns the Go type this attribute generates values for
//   - GetDefaultImplementation() Attributes: Returns a default-configured instance
//   - GetRandomValue() any: Generates and returns a random value
//
// Example implementation usage:
//
//	attrs := IntegerAttributesImpl[int]{Min: 0, Max: 100}
//	value := attrs.GetRandomValue() // Returns a random int between 0 and 100
//	reflectType := attrs.GetReflectType() // Returns reflect.TypeOf(int(0))
type Attributes interface {
	GetAttributes() any
	GetReflectType() reflect.Type
	GetDefaultImplementation() Attributes
	GetRandomValue() any
}

// AttributesStruct is the interface for the top-level attributes configuration.
// It maps Go types to their corresponding Attributes implementations.
//
// This interface is implemented by FTAttributes and provides the central
// type-to-generator mapping used by the fuzz testing framework.
//
// Methods:
//   - GetAttributeGivenType(t reflect.Type) (retA Attributes, err error): Maps a type to attributes
//
// Example usage:
//
//	ftAttrs := NewFTAttributes()
//	intAttrs, err := ftAttrs.GetAttributeGivenType(reflect.TypeOf(int(0)))
//	randomInt := intAttrs.GetRandomValue()
type AttributesStruct interface {
	GetAttributeGivenType(t reflect.Type) (retA Attributes, err error)
}

// Type Interfaces

// Integers defines the constraint for signed integer types.
// This constraint is used as a type parameter for IntegerAttributesImpl
// to ensure it only works with signed integer types.
//
// Supported types: int, int8, int16, int32, int64
type Integers interface {
	int | int8 | int16 | int32 | int64
}

// UnsignedIntegers defines the constraint for unsigned integer types.
// This constraint is used as a type parameter for UnsignedIntegerAttributesImpl
// to ensure it only works with unsigned integer types.
//
// Supported types: uint, uint8, uint16, uint32, uint64
type UnsignedIntegers interface {
	uint | uint8 | uint16 | uint32 | uint64
}

// Floats defines the constraint for floating-point types.
// This constraint is used as a type parameter for FloatAttributesImpl
// to ensure it only works with floating-point types.
//
// Supported types: float32, float64
type Floats interface {
	float32 | float64
}

// Complex defines the constraint for complex number types.
// This constraint is used as a type parameter for ComplexAttributesImpl
// to ensure it only works with complex number types.
//
// Supported types: complex64, complex128
type Complex interface {
	complex64 | complex128
}

// Specific Attribute Interfaces for the 4 generic fields

// IntegerAttributes defines the interface for integer attribute implementations.
// This is satisfied by IntegerAttributesImpl[T] for any signed integer type T.
//
// Methods are inherited from the Attributes interface.
type IntegerAttributes interface {
	GetAttributes() any
	GetReflectType() reflect.Type
	GetDefaultImplementation() Attributes
	GetRandomValue() any
}

// UnsignedIntegerAttributes defines the interface for unsigned integer attribute implementations.
// This is satisfied by UnsignedIntegerAttributesImpl[T] for any unsigned integer type T.
//
// Methods are inherited from the Attributes interface.
type UnsignedIntegerAttributes interface {
	GetAttributes() any
	GetReflectType() reflect.Type
	GetDefaultImplementation() Attributes
	GetRandomValue() any
}

// FloatAttributes defines the interface for float attribute implementations.
// This is satisfied by FloatAttributesImpl[T] for any floating-point type T.
//
// Methods are inherited from the Attributes interface.
type FloatAttributes interface {
	GetAttributes() any
	GetReflectType() reflect.Type
	GetDefaultImplementation() Attributes
	GetRandomValue() any
}

// ComplexAttributes defines the interface for complex number attribute implementations.
// This is satisfied by ComplexAttributesImpl[T] for any complex number type T.
//
// Methods are inherited from the Attributes interface.
type ComplexAttributes interface {
	GetAttributes() any
	GetReflectType() reflect.Type
	GetDefaultImplementation() Attributes
	GetRandomValue() any
}
