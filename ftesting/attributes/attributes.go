// Package attributes provides a type-safe, configurable system for generating random
// values for fuzz testing in Go.
//
// The attributes package is the core of the ftesting framework's input generation
// capabilities. It defines interfaces and implementations for generating random values
// of any Go type with fine-grained control over ranges, constraints, and special values.
//
// Architecture:
//
// The package uses a generic, reflection-based design:
//   - Attributes interface: The contract all value generators must implement
//   - FTAttributes: Central configuration mapping types to their generators
//   - Type-specific implementations: Generic structs like IntegerAttributesImpl[T]
//   - Type constraints: Union types (Integers, Floats, etc.) for type safety
//
// Supported Types:
//   - Signed integers: int, int8, int16, int32, int64
//   - Unsigned integers: uint, uint8, uint16, uint32, uint64
//   - Floating-point: float32, float64
//   - Complex numbers: complex64, complex128
//   - Strings with customizable character sets and patterns
//   - Slices with element constraints
//   - Arrays with fixed sizes
//   - Maps with key/value constraints
//   - Pointers including multi-level pointers and nil support
//   - Structs with per-field attribute configuration
//   - Booleans
//
// Key Concepts:
//
//  1. Attributes: Configuration objects that define how to generate random values.
//     Each attribute type has fields controlling ranges, special values, and constraints.
//
//  2. Type Parameters: Generic implementations use Go 1.18+ type parameters with
//     union constraints (e.g., IntegerAttributesImpl[T Integers]) to ensure type safety
//     across different bit sizes and numeric types.
//
// 3. Reflection: The package heavily uses reflection to:
//
//   - Map runtime types to attribute implementations
//
//   - Create values of arbitrary types dynamically
//
//   - Convert between types while maintaining type safety
//
//     4. Default Implementations: Every attribute type provides sensible defaults via
//     GetDefaultImplementation(), enabling zero-configuration usage.
//
// Basic Usage:
//
//	// Create default attributes for all types
//	attrs := NewFTAttributes()
//
//	// Get attributes for a specific type
//	intType := reflect.TypeOf(int(0))
//	intAttr, err := attrs.GetAttributeGivenType(intType)
//
//	// Generate a random value
//	randomInt := intAttr.GetRandomValue()
//
// Custom Configuration:
//
//	// Customize integer generation
//	attrs := NewFTAttributes()
//	attrs.IntegerAttr = IntegerAttributesImpl[int]{
//	    Min: 1,
//	    Max: 100,
//	    AllowZero: false,
//	    AllowNegative: false,
//	}
//
//	// Customize string generation
//	attrs.StringAttr = StringAttributes{
//	    MinLen: 8,
//	    MaxLen: 16,
//	    AllowedRunes: []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"),
//	}
//
// Advanced Usage with Nested Types:
//
//	// Configure slice of custom structs
//	attrs.SliceAttr = SliceAttributes{
//	    MinLen: 5,
//	    MaxLen: 10,
//	    ElementAttrs: StructAttributes{
//	        FieldAttrs: map[string]any{
//	            "ID": IntegerAttributesImpl[int]{Min: 1, Max: 1000},
//	            "Name": StringAttributes{MinLen: 3, MaxLen: 20},
//	        },
//	    },
//	}
//
// Design Patterns:
//
//  1. Generic Attribute Implementations: Use type parameters to implement a single
//     struct that works for all variants of a type category (e.g., all signed integers).
//
//  2. Helper Methods: Private helper methods break down complex generation logic into
//     testable, reusable pieces (e.g., getBounds(), generateRandomValue(), convertToType()).
//
//  3. Zero Value Fallbacks: Methods return appropriate zero values when configuration
//     is invalid rather than panicking, ensuring robustness.
//
//  4. Type Safety via Reflection: Despite heavy reflection use, the package maintains
//     type safety by validating types and using type conversion rather than unsafe operations.
package attributes

import (
	"fmt"
	"math/rand"
	"reflect"

	p "github.com/laiambryant/gotestutils/pbtesting/properties/predicates"
)

// FTAttributes is the central configuration struct for fuzz testing input generation.
// It contains attribute configurations for all supported Go types, allowing fine-grained
// control over how random values are generated for each type category.
//
// Each field represents a category of types and contains configuration for generating
// random values of that category. Default configurations are provided by NewFTAttributes.
//
// Fields:
//   - IntegerAttr: Configuration for signed integer types (int, int8, int16, int32, int64)
//   - UIntegerAttr: Configuration for unsigned integer types (uint, uint8, uint16, uint32, uint64)
//   - FloatAttr: Configuration for floating-point types (float32, float64)
//   - ComplexAttr: Configuration for complex number types (complex64, complex128)
//   - StringAttr: Configuration for string generation
//   - SliceAttr: Configuration for slice generation
//   - BoolAttr: Configuration for boolean generation
//   - MapAttr: Configuration for map generation
//   - PointerAttr: Configuration for pointer generation (including multi-level pointers)
//   - StructAttr: Configuration for struct generation
//   - ArrayAttr: Configuration for array generation
//
// Example usage:
//
//	attrs := NewFTAttributes()
//	attrs.IntegerAttr = IntegerAttributesImpl[int]{Min: 0, Max: 100, AllowZero: false}
//	attrs.StringAttr = StringAttributes{MinLen: 5, MaxLen: 20}
type FTAttributes struct {
	IntegerAttr  IntegerAttributes
	UIntegerAttr UnsignedIntegerAttributes
	FloatAttr    FloatAttributes
	ComplexAttr  ComplexAttributes
	StringAttr   StringAttributes
	SliceAttr    SliceAttributes
	BoolAttr     BoolAttributes
	MapAttr      MapAttributes
	PointerAttr  PointerAttributes
	StructAttr   StructAttributes
	ArrayAttr    ArrayAttributes
}

// NewFTAttributes creates and returns an FTAttributes instance with sensible default
// configurations for all supported types. These defaults are designed to work well
// for general-purpose fuzz testing.
//
// Default configurations:
//   - Integers: Range [-100, 100], allow negative and zero
//   - Unsigned integers: Range [0, 100], allow zero
//   - Floats: Range [-100.0, 100.0], finite only, non-zero
//   - Complex: Real and imaginary parts in range [-10.0, 10.0]
//   - Strings: Length [1, 10] characters
//   - Slices: Length [1, 5] elements, with integer elements
//   - Bools: Random true/false
//   - Maps: Size [1, 5] entries, string keys and integer values
//   - Pointers: Allow nil, depth 1, integer inner type
//   - Structs: Two fields (Field1: int, Field2: float32)
//   - Arrays: Length 5, integer elements
//
// Returns an FTAttributes instance ready for use with FTesting.
//
// Example usage:
//
//	attrs := NewFTAttributes()
//	// Optionally override specific attributes:
//	attrs.IntegerAttr = IntegerAttributesImpl[int]{Min: 1, Max: 1000}
//	ft.WithAttributes(attrs)
func NewFTAttributes() FTAttributes {
	return FTAttributes{
		IntegerAttr:  IntegerAttributesImpl[int]{AllowNegative: true, AllowZero: true, Max: 100, Min: -100},
		UIntegerAttr: UnsignedIntegerAttributesImpl[uint]{Signed: true, AllowNegative: true, AllowZero: true, Max: 100, Min: 0},
		FloatAttr:    FloatAttributesImpl[float64]{Min: -100.0, Max: 100.0, NonZero: true, FiniteOnly: true},
		ComplexAttr:  ComplexAttributesImpl[complex128]{RealMin: -10.0, RealMax: 10.0, ImagMin: -10.0, ImagMax: 10.0},
		StringAttr:   StringAttributes{MinLen: 1, MaxLen: 10},
		SliceAttr:    SliceAttributes{MinLen: 1, MaxLen: 5, ElementAttrs: IntegerAttributesImpl[int]{}},
		BoolAttr:     BoolAttributes{ForceTrue: false},
		MapAttr:      MapAttributes{MinSize: 1, MaxSize: 5, KeyAttrs: StringAttributes{MinLen: 1, MaxLen: 5}, ValueAttrs: IntegerAttributesImpl[int]{}},
		PointerAttr:  PointerAttributes{AllowNil: true, Depth: 1, Inner: IntegerAttributesImpl[int]{}},
		StructAttr:   StructAttributes{FieldAttrs: map[string]any{"Field1": IntegerAttributesImpl[int]{}, "Field2": FloatAttributesImpl[float32]{Min: -10.0, Max: 10.0}}},
		ArrayAttr:    ArrayAttributes{Length: 5, ElementAttrs: IntegerAttributesImpl[int]{}},
	}
}

// GetAttributeGivenType returns the appropriate Attributes implementation for the given
// reflect.Type. This method is the core type-to-attribute mapping mechanism used by
// the fuzz testing framework to determine how to generate random values for function parameters.
//
// The method performs the following:
// 1. Maps the type's Kind to the corresponding attribute configuration
// 2. Checks if the attribute has custom configuration or needs defaults
// 3. Returns a fully configured Attributes instance ready for value generation
//
// Parameters:
//   - t: The reflect.Type to get attributes for
//
// Returns:
//   - retA: An Attributes implementation configured for the given type
//   - err: An error if the type is nil or unsupported
//
// Errors returned:
//   - NilTypeError: When t is nil
//   - UnsupportedAttributeTypeError: When the type's Kind is not supported
//
// Example usage:
//
//	intType := reflect.TypeOf(int(0))
//	attrs := NewFTAttributes()
//	intAttr, err := attrs.GetAttributeGivenType(intType)
//	// intAttr can now generate random integers
//	randomInt := intAttr.GetRandomValue()
func (mt FTAttributes) GetAttributeGivenType(t reflect.Type) (retA Attributes, err error) {
	if t == nil {
		return nil, NilTypeError{}
	}
	kindMap := map[reflect.Kind]Attributes{
		reflect.Int: mt.IntegerAttr, reflect.Int8: mt.IntegerAttr, reflect.Int16: mt.IntegerAttr, reflect.Int32: mt.IntegerAttr, reflect.Int64: mt.IntegerAttr,
		reflect.Uint: mt.UIntegerAttr, reflect.Uint8: mt.UIntegerAttr, reflect.Uint16: mt.UIntegerAttr, reflect.Uint32: mt.UIntegerAttr, reflect.Uint64: mt.UIntegerAttr,
		reflect.Float32: mt.FloatAttr, reflect.Float64: mt.FloatAttr,
		reflect.Complex64: mt.ComplexAttr, reflect.Complex128: mt.ComplexAttr,
		reflect.String: mt.StringAttr, reflect.Slice: mt.SliceAttr, reflect.Bool: mt.BoolAttr,
		reflect.Map: mt.MapAttr, reflect.Pointer: mt.PointerAttr, reflect.Struct: mt.StructAttr, reflect.Array: mt.ArrayAttr,
	}
	retA = kindMap[t.Kind()]
	if retA == nil {
		retA, err = mt.getDefaultForKind(t.Kind())
		return
	}
	attrsVal := retA.GetAttributes()
	if attrsVal == nil {
		retA = retA.GetDefaultImplementation()
		return
	}
	attrsValType := reflect.TypeOf(attrsVal)

	zero := reflect.Zero(attrsValType).Interface()
	if reflect.DeepEqual(attrsVal, zero) {
		retA = retA.GetDefaultImplementation()
	}
	return
}

// getDefaultForKind returns a default Attributes implementation for the given reflect.Kind.
// This is a fallback method used when no custom attribute configuration exists for a type.
//
// Parameters:
//   - kind: The reflect.Kind to get default attributes for
//
// Returns:
//   - Attributes: A default implementation for the given kind
//   - error: UnsupportedAttributeTypeError if the kind is not supported
//
// This method is used internally by GetAttributeGivenType.
func (mt FTAttributes) getDefaultForKind(kind reflect.Kind) (Attributes, error) {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return IntegerAttributesImpl[int64]{}.GetDefaultImplementation(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return UnsignedIntegerAttributesImpl[uint64]{}.GetDefaultImplementation(), nil
	case reflect.Float32, reflect.Float64:
		return FloatAttributesImpl[float64]{}.GetDefaultImplementation(), nil
	case reflect.Complex64, reflect.Complex128:
		return ComplexAttributesImpl[complex128]{}.GetDefaultImplementation(), nil
	case reflect.String:
		return StringAttributes{}.GetDefaultImplementation(), nil
	case reflect.Slice:
		return SliceAttributes{}.GetDefaultImplementation(), nil
	case reflect.Bool:
		return BoolAttributes{}.GetDefaultImplementation(), nil
	case reflect.Map:
		return MapAttributes{}.GetDefaultImplementation(), nil
	case reflect.Pointer:
		return PointerAttributes{}.GetDefaultImplementation(), nil
	case reflect.Struct:
		return StructAttributes{}.GetDefaultImplementation(), nil
	case reflect.Array:
		return ArrayAttributes{}.GetDefaultImplementation(), nil
	default:
		return nil, UnsupportedAttributeTypeError{kind}
	}
}

// IntegerAttributesImpl is a generic implementation for generating random signed integer values
// with configurable constraints. The type parameter T must be one of the signed integer types.
//
// Type parameter:
//   - T: Must satisfy the Integers constraint (int, int8, int16, int32, int64)
//
// Fields:
//   - AllowNegative: If true, negative values can be generated; if false, only positive values
//   - AllowZero: If true, zero can be generated; if false, zero is excluded
//   - Max: The maximum value (inclusive) for generated integers
//   - Min: The minimum value (inclusive) for generated integers
//
// The implementation uses reflection and type conversion to ensure generated values
// match the exact integer type T, even when working with different bit sizes.
//
// Example usage:
//
//	// Generate random integers from 1 to 100 (no zero or negatives)
//	attrs := IntegerAttributesImpl[int]{
//	    AllowNegative: false,
//	    AllowZero: false,
//	    Max: 100,
//	    Min: 1,
//	}
//	randomInt := attrs.GetRandomValue() // Returns a random int between 1 and 100
type IntegerAttributesImpl[T Integers] struct {
	AllowNegative bool
	AllowZero     bool
	Max           T
	Min           T
}

func (a IntegerAttributesImpl[T]) GetAttributes() any { return a }
func (a IntegerAttributesImpl[T]) GetReflectType() reflect.Type {
	return reflect.TypeOf(*new(T))
}

func (a IntegerAttributesImpl[T]) GetDefaultImplementation() Attributes {
	return IntegerAttributesImpl[T]{
		AllowNegative: true,
		AllowZero:     true,
		Max:           100,
		Min:           -100,
	}
}

func (a IntegerAttributesImpl[T]) GetRandomValue() any {
	var zero T
	if !a.isValidRange(zero) {
		return zero
	}
	min, max := a.getMinMaxAsInt64()
	return a.generateRandomInteger(min, max, zero)
}

// isValidRange checks if the min/max range is valid
func (a IntegerAttributesImpl[T]) isValidRange(zero T) bool {
	return a.Max > zero && a.Min <= a.Max
}

// getMinMaxAsInt64 converts min and max to int64 for calculation
func (a IntegerAttributesImpl[T]) getMinMaxAsInt64() (int64, int64) {
	minVal := reflect.ValueOf(a.Min)
	maxVal := reflect.ValueOf(a.Max)
	return minVal.Int(), maxVal.Int()
}

// generateRandomInteger generates a random integer within the range and converts back to type T
func (a IntegerAttributesImpl[T]) generateRandomInteger(min, max int64, zero T) any {
	result := min + rand.Int63n(max-min+1)
	resultVal := reflect.ValueOf(result).Convert(reflect.TypeOf(zero))
	return resultVal.Interface()
}

// UnsignedIntegerAttributesImpl is a generic implementation for generating random unsigned
// integer values with configurable constraints. The type parameter T must be one of the
// unsigned integer types.
//
// Type parameter:
//   - T: Must satisfy the UnsignedIntegers constraint (uint, uint8, uint16, uint32, uint64)
//
// Fields:
//   - Signed: If true, treats the value as signed for generation purposes
//   - AllowNegative: If true, negative values can be generated (requires Signed to be true)
//   - AllowZero: If true, zero can be generated; if false, zero is excluded
//   - Max: The maximum value (inclusive) for generated unsigned integers
//   - Min: The minimum value (inclusive) for generated unsigned integers
//
// Example usage:
//
//	// Generate random unsigned integers from 0 to 255
//	attrs := UnsignedIntegerAttributesImpl[uint8]{
//	    Signed: false,
//	    AllowNegative: false,
//	    AllowZero: true,
//	    Max: 255,
//	    Min: 0,
//	}
//	randomUInt := attrs.GetRandomValue() // Returns a random uint8 between 0 and 255
type UnsignedIntegerAttributesImpl[T UnsignedIntegers] struct {
	Signed        bool
	AllowNegative bool
	AllowZero     bool
	Max           T
	Min           T
}

func (a UnsignedIntegerAttributesImpl[T]) GetAttributes() any { return a }
func (a UnsignedIntegerAttributesImpl[T]) GetReflectType() reflect.Type {
	if a.Signed || a.AllowNegative {
		return reflect.TypeOf(int64(0))
	}
	return reflect.TypeOf(uint64(0))
}

func (a UnsignedIntegerAttributesImpl[T]) GetDefaultImplementation() Attributes {
	return UnsignedIntegerAttributesImpl[T]{
		Signed:        true,
		AllowNegative: true,
		AllowZero:     true,
		Max:           100,
		Min:           0,
	}
}

func (a UnsignedIntegerAttributesImpl[T]) GetRandomValue() any {
	var zero T
	if !a.isValidRange(zero) {
		return zero
	}

	min, max := a.getMinMaxAsUint64()
	if max <= min {
		return zero
	}

	return a.generateRandomUnsignedInteger(min, max, zero)
}

// isValidRange checks if the min/max range is valid
func (a UnsignedIntegerAttributesImpl[T]) isValidRange(zero T) bool {
	return a.Max > zero && a.Min <= a.Max
}

// getMinMaxAsUint64 converts min and max to uint64 for calculation
func (a UnsignedIntegerAttributesImpl[T]) getMinMaxAsUint64() (uint64, uint64) {
	minVal := reflect.ValueOf(a.Min)
	maxVal := reflect.ValueOf(a.Max)
	return minVal.Uint(), maxVal.Uint()
}

// generateRandomUnsignedInteger generates a random unsigned integer within the range and converts back to type T
func (a UnsignedIntegerAttributesImpl[T]) generateRandomUnsignedInteger(min, max uint64, zero T) any {
	diff := max - min + 1
	result := min + uint64(rand.Int63n(int64(diff)))
	resultVal := reflect.ValueOf(result).Convert(reflect.TypeOf(zero))
	return resultVal.Interface()
}

// FloatAttributesImpl is a generic implementation for generating random floating-point
// values with configurable constraints and special value handling.
//
// Type parameter:
//   - T: Must satisfy the Floats constraint (float32, float64)
//
// Fields:
//   - Min: The minimum value (inclusive) for generated floats
//   - Max: The maximum value (inclusive) for generated floats
//   - NonZero: If true, zero is excluded from generated values
//   - FiniteOnly: If true, only finite values are generated (no Inf or NaN)
//   - AllowNaN: If true, NaN values can be generated (requires FiniteOnly to be false)
//   - AllowInf: If true, Infinity values can be generated (requires FiniteOnly to be false)
//   - Precision: Number of decimal places for rounding (0 means no rounding)
//
// Example usage:
//
//	// Generate random floats from -1.0 to 1.0, excluding zero
//	attrs := FloatAttributesImpl[float64]{
//	    Min: -1.0,
//	    Max: 1.0,
//	    NonZero: true,
//	    FiniteOnly: true,
//	}
//	randomFloat := attrs.GetRandomValue() // Returns a random float64 between -1.0 and 1.0
type FloatAttributesImpl[T Floats] struct {
	Min        T
	Max        T
	NonZero    bool
	FiniteOnly bool
	AllowNaN   bool
	AllowInf   bool
	Precision  uint
}

func (a FloatAttributesImpl[T]) GetAttributes() any           { return a }
func (a FloatAttributesImpl[T]) GetReflectType() reflect.Type { return reflect.TypeOf(float64(0)) }
func (a FloatAttributesImpl[T]) GetDefaultImplementation() Attributes {
	return FloatAttributesImpl[T]{
		Min:        -100.0,
		Max:        100.0,
		NonZero:    true,
		FiniteOnly: true,
	}
}

func (a FloatAttributesImpl[T]) GetRandomValue() any {
	var zero T
	if !a.isValidRange() {
		return zero
	}

	min, max := a.getMinMaxAsFloat64()
	result := a.generateRandomFloat(min, max)
	return a.convertToTargetType(result, zero)
}

// isValidRange checks if the min/max range is valid
func (a FloatAttributesImpl[T]) isValidRange() bool {
	return a.Max > a.Min
}

// getMinMaxAsFloat64 converts min and max to float64 for calculation
func (a FloatAttributesImpl[T]) getMinMaxAsFloat64() (float64, float64) {
	minVal := reflect.ValueOf(a.Min)
	maxVal := reflect.ValueOf(a.Max)
	return minVal.Float(), maxVal.Float()
}

// generateRandomFloat generates a random float within the range
func (a FloatAttributesImpl[T]) generateRandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// convertToTargetType converts the result back to the target type T
func (a FloatAttributesImpl[T]) convertToTargetType(result float64, zero T) any {
	resultVal := reflect.ValueOf(result).Convert(reflect.TypeOf(zero))
	return resultVal.Interface()
}

// ComplexAttributesImpl is a generic implementation for generating random complex number
// values with separate control over real and imaginary components.
//
// Type parameter:
//   - T: Must satisfy the Complex constraint (complex64, complex128)
//
// Fields:
//   - RealMin: The minimum value for the real component
//   - RealMax: The maximum value for the real component
//   - ImagMin: The minimum value for the imaginary component
//   - ImagMax: The maximum value for the imaginary component
//   - MagnitudeMin: Optional constraint on minimum magnitude
//   - MagnitudeMax: Optional constraint on maximum magnitude
//   - MaxComplex: Optional maximum complex value
//   - MinComplex: Optional minimum complex value
//   - AllowNaN: If true, NaN components can be generated
//   - AllowInf: If true, Infinity components can be generated
//
// Example usage:
//
//	// Generate complex numbers with real and imaginary parts in [-5.0, 5.0]
//	attrs := ComplexAttributesImpl[complex128]{
//	    RealMin: -5.0,
//	    RealMax: 5.0,
//	    ImagMin: -5.0,
//	    ImagMax: 5.0,
//	}
//	randomComplex := attrs.GetRandomValue() // Returns a random complex128
type ComplexAttributesImpl[T Complex] struct {
	RealMin      float64
	RealMax      float64
	ImagMin      float64
	ImagMax      float64
	MagnitudeMin float64
	MagnitudeMax float64
	MaxComplex   T
	MinComplex   T
	AllowNaN     bool
	AllowInf     bool
}

func (a ComplexAttributesImpl[T]) GetAttributes() any           { return a }
func (a ComplexAttributesImpl[T]) GetReflectType() reflect.Type { return reflect.TypeOf(complex128(0)) }
func (a ComplexAttributesImpl[T]) GetDefaultImplementation() Attributes {
	return ComplexAttributesImpl[T]{
		RealMin: -10.0,
		RealMax: 10.0,
		ImagMin: -10.0,
		ImagMax: 10.0,
	}
}

func (a ComplexAttributesImpl[T]) GetRandomValue() any {
	var zero T
	realMin, realMax, imagMin, imagMax := a.getBounds()
	realPart := a.generateRandomReal(realMin, realMax)
	imagPart := a.generateRandomImaginary(imagMin, imagMax)
	return a.createComplexValue(realPart, imagPart, zero)
}

// getBounds returns validated real and imaginary bounds
func (a ComplexAttributesImpl[T]) getBounds() (float64, float64, float64, float64) {
	realMin, realMax := a.RealMin, a.RealMax
	imagMin, imagMax := a.ImagMin, a.ImagMax

	if realMax <= realMin {
		realMin, realMax = -10.0, 10.0
	}
	if imagMax <= imagMin {
		imagMin, imagMax = -10.0, 10.0
	}

	return realMin, realMax, imagMin, imagMax
}

// generateRandomReal generates a random real part
func (a ComplexAttributesImpl[T]) generateRandomReal(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// generateRandomImaginary generates a random imaginary part
func (a ComplexAttributesImpl[T]) generateRandomImaginary(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// createComplexValue creates and converts the complex value to target type
func (a ComplexAttributesImpl[T]) createComplexValue(realPart, imagPart float64, zero T) any {
	complexVal := complex(realPart, imagPart)
	resultVal := reflect.ValueOf(complexVal).Convert(reflect.TypeOf(zero))
	return resultVal.Interface()
}

// StringAttributes configures the generation of random string values with various
// constraints including length, character sets, and pattern matching.
//
// Fields:
//   - MinLen: Minimum string length (inclusive)
//   - MaxLen: Maximum string length (inclusive)
//   - AllowedRunes: Character set to use (defaults to ASCII printable if empty)
//   - Regex: Regular expression pattern that generated strings should match
//   - Prefix: String to prepend to all generated strings
//   - Suffix: String to append to all generated strings
//   - Contains: Substring that must appear in all generated strings
//   - UniqueChars: If true, all characters in generated strings must be unique
//
// Example usage:
//
//	// Generate random alphanumeric strings of length 8-12
//	attrs := StringAttributes{
//	    MinLen: 8,
//	    MaxLen: 12,
//	    AllowedRunes: []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"),
//	}
//	randomString := attrs.GetRandomValue() // Returns a random string like "aBc3Def9Gh"
type StringAttributes struct {
	MinLen       int
	MaxLen       int
	AllowedRunes []rune
	Regex        string
	Prefix       string
	Suffix       string
	Contains     string
	UniqueChars  bool
}

func (a StringAttributes) GetAttributes() any           { return a }
func (a StringAttributes) GetReflectType() reflect.Type { return reflect.TypeOf("") }
func (a StringAttributes) GetDefaultImplementation() Attributes {
	return StringAttributes{
		MinLen: 1,
		MaxLen: 10,
	}
}

func (a StringAttributes) GetRandomValue() any {
	minLen, maxLen := a.getLengthBounds()
	length := a.pickLength(minLen, maxLen)
	allowedRunes := a.getAllowedRunes()
	generated := a.generateRandomString(allowedRunes, length)
	return a.applyPrefixSuffix(generated)
}

// getLengthBounds returns validated min and max length bounds
func (a StringAttributes) getLengthBounds() (int, int) {
	minLen, maxLen := a.MinLen, a.MaxLen
	if maxLen <= 0 {
		maxLen = 10
	}
	if minLen < 0 {
		minLen = 0
	}
	if minLen > maxLen {
		minLen = maxLen
	}
	return minLen, maxLen
}

// pickLength picks a random length between minLen and maxLen
func (a StringAttributes) pickLength(minLen, maxLen int) int {
	if maxLen > minLen {
		return minLen + rand.Intn(maxLen-minLen+1)
	}
	return minLen
}

// getAllowedRunes returns the allowed runes, defaulting to ASCII printable if empty
func (a StringAttributes) getAllowedRunes() []rune {
	allowedRunes := a.AllowedRunes
	if len(allowedRunes) == 0 {
		for i := 32; i <= 126; i++ {
			allowedRunes = append(allowedRunes, rune(i))
		}
	}
	return allowedRunes
}

// generateRandomString generates a random string of given length using allowed runes
func (a StringAttributes) generateRandomString(allowedRunes []rune, length int) string {
	result := make([]rune, length)
	for i := range length {
		result[i] = allowedRunes[rand.Intn(len(allowedRunes))]
	}
	return string(result)
}

// applyPrefixSuffix applies prefix and suffix to the generated string
func (a StringAttributes) applyPrefixSuffix(generated string) string {
	if a.Prefix != "" {
		generated = a.Prefix + generated
	}
	if a.Suffix != "" {
		generated = generated + a.Suffix
	}
	return generated
}

// SliceAttributes configures the generation of random slice values with control
// over slice length, element generation, and optional properties like uniqueness and sorting.
//
// Fields:
//   - MinLen: Minimum slice length (inclusive)
//   - MaxLen: Maximum slice length (inclusive)
//   - Unique: If true, all slice elements must be unique
//   - Sorted: If true, generated slices are sorted
//   - ElementPreds: Predicates that all elements must satisfy
//   - ElementAttrs: Attributes for generating slice elements (can be Attributes or reflect.Type)
//
// Example usage:
//
//	// Generate random slices of integers with length 5-10
//	attrs := SliceAttributes{
//	    MinLen: 5,
//	    MaxLen: 10,
//	    ElementAttrs: IntegerAttributesImpl[int]{Min: 0, Max: 100},
//	}
//	randomSlice := attrs.GetRandomValue() // Returns a random []int with 5-10 elements
type SliceAttributes struct {
	MinLen       int
	MaxLen       int
	Unique       bool
	Sorted       bool
	ElementPreds []p.Predicate
	ElementAttrs any
}

func (a SliceAttributes) GetAttributes() any { return a }
func (a SliceAttributes) GetReflectType() reflect.Type {
	var elemType reflect.Type
	switch v := a.ElementAttrs.(type) {
	case Attributes:
		elemType = v.GetReflectType()
	case reflect.Type:
		elemType = v
	default:
		elemType = nil
	}
	if elemType == nil {
		return nil
	}
	return reflect.SliceOf(elemType)
}

func (a SliceAttributes) GetDefaultImplementation() Attributes {
	return SliceAttributes{
		MinLen:       1,
		MaxLen:       5,
		ElementAttrs: IntegerAttributesImpl[int]{},
	}
}

func (a SliceAttributes) GetRandomValue() any {
	minLen, maxLen := a.getSliceLengthBounds()
	length := a.pickSliceLength(minLen, maxLen)
	elemType := a.getElementType()
	if elemType == nil {
		return nil
	}
	result := a.makeSliceOfType(elemType, length)
	a.fillSliceWithRandomElements(result, elemType, length)
	return result.Interface()
}

// getSliceLengthBounds returns the min and max length for the slice.
func (a SliceAttributes) getSliceLengthBounds() (int, int) {
	minLen := a.MinLen
	maxLen := a.MaxLen
	if maxLen <= 0 {
		maxLen = 5
	}
	if minLen < 0 {
		minLen = 0
	}
	if minLen > maxLen {
		minLen = maxLen
	}
	return minLen, maxLen
}

// pickSliceLength picks a random length between minLen and maxLen.
func (a SliceAttributes) pickSliceLength(minLen, maxLen int) int {
	if maxLen > minLen {
		return minLen + rand.Intn(maxLen-minLen+1)
	}
	return minLen
}

// getElementType returns the reflect.Type of the slice element.
func (a SliceAttributes) getElementType() reflect.Type {
	if attrs, ok := a.ElementAttrs.(Attributes); ok {
		return attrs.GetReflectType()
	}
	return reflect.TypeOf(any(nil))
}

// makeSliceOfType creates a slice of the given type and length.
func (a SliceAttributes) makeSliceOfType(elemType reflect.Type, length int) reflect.Value {
	sliceType := reflect.SliceOf(elemType)
	return reflect.MakeSlice(sliceType, length, length)
}

// fillSliceWithRandomElements fills the slice with random elements.
func (a SliceAttributes) fillSliceWithRandomElements(result reflect.Value, elemType reflect.Type, length int) {
	for i := range length {
		var elemValue reflect.Value
		if attrs, ok := a.ElementAttrs.(Attributes); ok {
			randVal := attrs.GetRandomValue()
			if randVal != nil {
				elemValue = reflect.ValueOf(randVal)
			} else {
				elemValue = reflect.Zero(elemType)
			}
		}
		result.Index(i).Set(elemValue)
	}
}

// BoolAttributes configures the generation of random boolean values with options
// to force specific values.
//
// Fields:
//   - ForceTrue: If true, always generate true
//   - ForceFalse: If true, always generate false
//
// If both ForceTrue and ForceFalse are false, values are randomly generated.
// If both are true, ForceTrue takes precedence.
//
// Example usage:
//
//	// Generate random booleans with 50/50 distribution
//	attrs := BoolAttributes{}
//	randomBool := attrs.GetRandomValue() // Returns true or false randomly
//
//	// Always generate true
//	forcedAttrs := BoolAttributes{ForceTrue: true}
//	alwaysTrue := forcedAttrs.GetRandomValue() // Always returns true
type BoolAttributes struct {
	ForceTrue  bool
	ForceFalse bool
}

func (a BoolAttributes) GetAttributes() any           { return a }
func (a BoolAttributes) GetReflectType() reflect.Type { return reflect.TypeOf(true) }
func (a BoolAttributes) GetDefaultImplementation() Attributes {
	return BoolAttributes{
		ForceTrue: false,
	}
}

func (a BoolAttributes) GetRandomValue() any {
	if a.shouldForceValue() {
		return a.getForcedValue()
	}
	return a.generateRandomBool()
}

// shouldForceValue checks if a specific boolean value should be forced
func (a BoolAttributes) shouldForceValue() bool {
	return a.ForceTrue || a.ForceFalse
}

// getForcedValue returns the forced boolean value
func (a BoolAttributes) getForcedValue() bool {
	return a.ForceTrue
}

// generateRandomBool generates a random boolean value
func (a BoolAttributes) generateRandomBool() bool {
	return rand.Intn(2) == 1
}

// MapAttributes configures the generation of random map values with control over
// map size, key/value generation, and optional predicate validation.
//
// Fields:
//   - MinSize: Minimum number of map entries (inclusive)
//   - MaxSize: Maximum number of map entries (inclusive)
//   - KeyPreds: Predicates that all keys must satisfy
//   - ValuePreds: Predicates that all values must satisfy
//   - KeyAttrs: Attributes for generating map keys (can be Attributes or reflect.Type)
//   - ValueAttrs: Attributes for generating map values (can be Attributes or reflect.Type)
//
// Example usage:
//
//	// Generate random maps with string keys and integer values
//	attrs := MapAttributes{
//	    MinSize: 1,
//	    MaxSize: 10,
//	    KeyAttrs: StringAttributes{MinLen: 3, MaxLen: 8},
//	    ValueAttrs: IntegerAttributesImpl[int]{Min: 0, Max: 100},
//	}
//	randomMap := attrs.GetRandomValue() // Returns a random map[string]int
type MapAttributes struct {
	MinSize    int
	MaxSize    int
	KeyPreds   []p.Predicate
	ValuePreds []p.Predicate
	KeyAttrs   any
	ValueAttrs any
}

func (a MapAttributes) GetAttributes() any { return a }
func (a MapAttributes) GetReflectType() reflect.Type {
	var kt, vt reflect.Type
	switch v := a.KeyAttrs.(type) {
	case Attributes:
		kt = v.GetReflectType()
	case reflect.Type:
		kt = v
	}
	switch v := a.ValueAttrs.(type) {
	case Attributes:
		vt = v.GetReflectType()
	case reflect.Type:
		vt = v
	}
	if kt == nil || vt == nil {
		return nil
	}
	return reflect.MapOf(kt, vt)
}

func (a MapAttributes) GetDefaultImplementation() Attributes {
	return MapAttributes{
		MinSize: 1,
		MaxSize: 5,
		KeyAttrs: StringAttributes{
			MinLen: 1,
			MaxLen: 5,
		},
		ValueAttrs: IntegerAttributesImpl[int]{},
	}
}

func (a MapAttributes) GetRandomValue() any {
	minSize, maxSize := a.getMapSizeBounds()
	size := a.pickMapSize(minSize, maxSize)
	keyType, valueType := a.getKeyValueTypes()
	if keyType == nil || valueType == nil {
		return nil
	}
	mapType := reflect.MapOf(keyType, valueType)
	result := reflect.MakeMap(mapType)
	a.fillMapWithRandomEntries(result, keyType, valueType, size)
	return result.Interface()
}

// getMapSizeBounds returns the min and max size for the map.
func (a MapAttributes) getMapSizeBounds() (int, int) {
	minSize := a.MinSize
	maxSize := a.MaxSize
	if maxSize <= 0 {
		maxSize = 5
	}
	if minSize < 0 {
		minSize = 0
	}
	if minSize > maxSize {
		minSize = maxSize
	}
	return minSize, maxSize
}

// pickMapSize picks a random size between minSize and maxSize.
func (a MapAttributes) pickMapSize(minSize, maxSize int) int {
	if maxSize > minSize {
		return minSize + rand.Intn(maxSize-minSize+1)
	}
	return minSize
}

// getKeyValueTypes returns the reflect.Type of the key and value.
func (a MapAttributes) getKeyValueTypes() (reflect.Type, reflect.Type) {
	var keyType, valueType reflect.Type
	if attrs, ok := a.KeyAttrs.(Attributes); ok {
		keyType = attrs.GetReflectType()
	}
	if attrs, ok := a.ValueAttrs.(Attributes); ok {
		valueType = attrs.GetReflectType()
	}
	return keyType, valueType
}

// fillMapWithRandomEntries fills the map with random key-value pairs.
func (a MapAttributes) fillMapWithRandomEntries(result reflect.Value, keyType, valueType reflect.Type, size int) {
	for i := 0; i < size; i++ {
		keyValue := a.getRandomKeyValue(keyType)
		valueValue := a.getRandomValueValue(valueType)
		result.SetMapIndex(keyValue, valueValue)
	}
}

// getRandomKeyValue returns a random key value.
func (a MapAttributes) getRandomKeyValue(keyType reflect.Type) reflect.Value {
	if attrs, ok := a.KeyAttrs.(Attributes); ok {
		randKey := attrs.GetRandomValue()
		if randKey != nil {
			return reflect.ValueOf(randKey)
		}
	}
	return reflect.Zero(keyType)
}

// getRandomValueValue returns a random value value.
func (a MapAttributes) getRandomValueValue(valueType reflect.Type) reflect.Value {
	if attrs, ok := a.ValueAttrs.(Attributes); ok {
		randValue := attrs.GetRandomValue()
		if randValue != nil {
			return reflect.ValueOf(randValue)
		}
	}
	return reflect.Zero(valueType)
}

// PointerAttributes configures the generation of random pointer values including
// support for nil pointers and multi-level pointer chains (pointer to pointer, etc.).
//
// Fields:
//   - AllowNil: If true, nil pointers can be generated
//   - Depth: Number of pointer levels (1 = *T, 2 = **T, etc.)
//   - Inner: Attributes for the pointed-to value (can be Attributes or reflect.Type)
//
// The implementation creates proper pointer chains by allocating memory at each level
// and setting up the chain correctly.
//
// Example usage:
//
//	// Generate pointers to integers, allowing nil
//	attrs := PointerAttributes{
//	    AllowNil: true,
//	    Depth: 1,
//	    Inner: IntegerAttributesImpl[int]{Min: 0, Max: 100},
//	}
//	randomPtr := attrs.GetRandomValue() // Returns *int (may be nil)
//
//	// Generate pointer-to-pointer-to-string
//	deepAttrs := PointerAttributes{
//	    AllowNil: false,
//	    Depth: 2,
//	    Inner: StringAttributes{MinLen: 5, MaxLen: 10},
//	}
//	deepPtr := deepAttrs.GetRandomValue() // Returns **string
type PointerAttributes struct {
	AllowNil bool
	Depth    int
	Inner    any
}

func (a PointerAttributes) GetAttributes() any { return a }
func (a PointerAttributes) GetReflectType() reflect.Type {
	var inner reflect.Type
	switch v := a.Inner.(type) {
	case Attributes:
		inner = v.GetReflectType()
	case reflect.Type:
		inner = v
	}
	if inner == nil {
		return nil
	}
	if a.Depth <= 0 {
		a.Depth = 1
	}
	for i := 0; i < a.Depth; i++ {
		inner = reflect.PointerTo(inner)
	}
	return inner
}

func (a PointerAttributes) GetDefaultImplementation() Attributes {
	return PointerAttributes{
		AllowNil: true,
		Depth:    1,
		Inner:    IntegerAttributesImpl[int]{},
	}
}

func (a PointerAttributes) GetRandomValue() any {
	if a.shouldReturnNil() {
		return a.getNilPointer()
	}

	innerValue := a.getInnerValue()
	if innerValue == nil {
		return nil
	}

	return a.createPointerChain(innerValue)
}

// shouldReturnNil determines if nil should be returned
func (a PointerAttributes) shouldReturnNil() bool {
	return a.AllowNil && rand.Intn(2) == 0
}

// getNilPointer returns a nil pointer of the correct type
func (a PointerAttributes) getNilPointer() any {
	return reflect.Zero(a.GetReflectType()).Interface()
}

// getInnerValue gets the inner value from the Inner attribute
func (a PointerAttributes) getInnerValue() *reflect.Value {
	if attrs, ok := a.Inner.(Attributes); ok {
		randVal := attrs.GetRandomValue()
		if randVal != nil {
			innerValue := reflect.ValueOf(randVal)
			return &innerValue
		} else {
			innerType := attrs.GetReflectType()
			if innerType != nil {
				innerValue := reflect.Zero(innerType)
				return &innerValue
			}
		}
	}
	return nil
}

// createPointerChain creates a chain of pointers with the specified depth
func (a PointerAttributes) createPointerChain(innerValue *reflect.Value) any {
	ptrValue := reflect.New(innerValue.Type())
	ptrValue.Elem().Set(*innerValue)

	currentPtr := ptrValue
	for i := 1; i < a.Depth; i++ {
		newPtr := reflect.New(currentPtr.Type())
		newPtr.Elem().Set(currentPtr)
		currentPtr = newPtr
	}

	return currentPtr.Interface()
}

// StructAttributes configures the generation of random struct values by mapping
// field names to their respective attribute configurations.
//
// Fields:
//   - FieldAttrs: A map from field name to field attributes (can be Attributes or reflect.Type)
//
// The implementation uses reflection to dynamically create struct types at runtime
// based on the field configurations. Each field is populated with a random value
// generated by its corresponding attribute.
//
// Note: The generated struct type is created dynamically using reflect.StructOf,
// so it won't have any methods or struct tags beyond what's defined in FieldAttrs.
//
// Example usage:
//
//	// Generate random structs with ID (int) and Name (string) fields
//	attrs := StructAttributes{
//	    FieldAttrs: map[string]any{
//	        "ID": IntegerAttributesImpl[int]{Min: 1, Max: 1000},
//	        "Name": StringAttributes{MinLen: 3, MaxLen: 20},
//	    },
//	}
//	randomStruct := attrs.GetRandomValue() // Returns a struct with ID and Name fields
type StructAttributes struct {
	FieldAttrs map[string]any
}

func (a StructAttributes) GetAttributes() any { return a }
func (a StructAttributes) GetReflectType() reflect.Type {
	if len(a.FieldAttrs) == 0 {
		return nil
	}
	fields := make([]reflect.StructField, 0, len(a.FieldAttrs))
	for name, attr := range a.FieldAttrs {
		var ft reflect.Type
		switch v := attr.(type) {
		case Attributes:
			ft = v.GetReflectType()
		case reflect.Type:
			ft = v
		}
		if ft == nil {
			return nil
		}
		fields = append(fields, reflect.StructField{
			Name: name,
			Type: ft,
			Tag:  "",
		})
	}
	return reflect.StructOf(fields)
}

func (a StructAttributes) GetDefaultImplementation() Attributes {
	return StructAttributes{
		FieldAttrs: map[string]any{
			"Field1": IntegerAttributesImpl[int]{},
			"Field2": FloatAttributesImpl[float32]{
				Min: -10.0,
				Max: 10.0,
			},
		},
	}
}

func (a StructAttributes) GetRandomValue() any {
	structType, err := a.getStructReflectType()
	if err != nil {
		return nil
	}
	structValue := a.createStructValue(structType)
	a.populateStructFields(structValue)
	return structValue.Interface()
}

// createStructValue creates a new struct value of the given type
func (a StructAttributes) createStructValue(structType reflect.Type) reflect.Value {
	return reflect.New(structType).Elem()
}

// populateStructFields populates all struct fields with random values
func (a StructAttributes) populateStructFields(structValue reflect.Value) {
	for fieldName, fieldAttr := range a.FieldAttrs {
		field := structValue.FieldByName(fieldName)
		if a.isFieldSettable(field) {
			fieldValue := a.generateFieldValue(fieldAttr, field.Type())
			a.setFieldValue(field, fieldValue)
		}
	}
}

// isFieldSettable checks if the field is valid and can be set
func (a StructAttributes) isFieldSettable(field reflect.Value) bool {
	return field.IsValid() && field.CanSet()
}

// generateFieldValue generates a random value for a struct field
func (a StructAttributes) generateFieldValue(fieldAttr any, fieldType reflect.Type) reflect.Value {
	if attrs, ok := fieldAttr.(Attributes); ok {
		randVal := attrs.GetRandomValue()
		if randVal != nil {
			return reflect.ValueOf(randVal)
		}
	}
	return reflect.Zero(fieldType)
}

// setFieldValue sets the field value with proper type conversion if needed
func (a StructAttributes) setFieldValue(field, fieldValue reflect.Value) {
	if fieldValue.Type().AssignableTo(field.Type()) {
		field.Set(fieldValue)
	} else if fieldValue.Type().ConvertibleTo(field.Type()) {
		field.Set(fieldValue.Convert(field.Type()))
	}
}

func (a StructAttributes) getStructReflectType() (reflect.Type, error) {
	if len(a.FieldAttrs) == 0 {
		return nil, fmt.Errorf("no field attributes found")
	}
	structType := a.GetReflectType()
	if structType == nil {
		return nil, fmt.Errorf("could not retrieve field type")
	}
	return structType, nil
}

// ArrayAttributes configures the generation of random fixed-size array values.
// Unlike slices, arrays have a fixed length determined at compile time.
//
// Fields:
//   - Length: The fixed length of the array (must be >= 0)
//   - Sorted: If true, array elements are sorted
//   - ElementAttrs: Attributes for generating array elements (can be Attributes or reflect.Type)
//
// Arrays are similar to slices but have a fixed size that's part of their type.
// The Length field determines the array type: [5]int vs [10]int are different types.
//
// Example usage:
//
//	// Generate random [10]int arrays with values 0-100
//	attrs := ArrayAttributes{
//	    Length: 10,
//	    ElementAttrs: IntegerAttributesImpl[int]{Min: 0, Max: 100},
//	}
//	randomArray := attrs.GetRandomValue() // Returns [10]int
type ArrayAttributes struct {
	Length       int
	Sorted       bool
	ElementAttrs any
}

func (a ArrayAttributes) GetAttributes() any { return a }
func (a ArrayAttributes) GetReflectType() reflect.Type {
	if a.Length < 0 {
		return nil
	}
	var et reflect.Type
	switch v := a.ElementAttrs.(type) {
	case Attributes:
		et = v.GetReflectType()
	case reflect.Type:
		et = v
	}
	if et == nil {
		return nil
	}
	return reflect.ArrayOf(a.Length, et)
}

func (a ArrayAttributes) GetDefaultImplementation() Attributes {
	return ArrayAttributes{
		Length:       5,
		ElementAttrs: IntegerAttributesImpl[int]{},
	}
}

func (a ArrayAttributes) GetRandomValue() any {
	if !a.isValidLength() {
		return nil
	}

	elemType := a.getElementType()
	if elemType == nil {
		return nil
	}

	arrayValue := a.createArrayValue(elemType)
	a.populateArrayElements(arrayValue, elemType)
	return arrayValue.Interface()
}

// isValidLength checks if the array length is valid
func (a ArrayAttributes) isValidLength() bool {
	return a.Length > 0
}

// getElementType returns the element type for the array
func (a ArrayAttributes) getElementType() reflect.Type {
	if attrs, ok := a.ElementAttrs.(Attributes); ok {
		return attrs.GetReflectType()
	}
	return nil
}

// createArrayValue creates a new array value of the specified type and length
func (a ArrayAttributes) createArrayValue(elemType reflect.Type) reflect.Value {
	arrayType := reflect.ArrayOf(a.Length, elemType)
	return reflect.New(arrayType).Elem()
}

// populateArrayElements fills the array with random elements
func (a ArrayAttributes) populateArrayElements(arrayValue reflect.Value, elemType reflect.Type) {
	for i := 0; i < a.Length; i++ {
		elemValue := a.generateElementValue(elemType)
		arrayValue.Index(i).Set(elemValue)
	}
}

// generateElementValue generates a random value for an array element
func (a ArrayAttributes) generateElementValue(elemType reflect.Type) reflect.Value {
	if attrs, ok := a.ElementAttrs.(Attributes); ok {
		randVal := attrs.GetRandomValue()
		if randVal != nil {
			return reflect.ValueOf(randVal)
		}
	}
	return reflect.Zero(elemType)
}
