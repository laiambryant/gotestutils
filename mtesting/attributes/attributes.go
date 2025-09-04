package attributes

import (
	"fmt"
	"math/rand"
	"reflect"

	p "github.com/laiambryant/gotestutils/pbtesting/properties/predicates"
)

type MTAttributes struct {
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

func (mt MTAttributes) GetAttributeGivenType(t reflect.Type) (retA Attributes) {
	kindMap := map[reflect.Kind]Attributes{
		reflect.Int: mt.IntegerAttr, reflect.Int8: mt.IntegerAttr, reflect.Int16: mt.IntegerAttr, reflect.Int32: mt.IntegerAttr, reflect.Int64: mt.IntegerAttr,
		reflect.Uint: mt.UIntegerAttr, reflect.Uint8: mt.UIntegerAttr, reflect.Uint16: mt.UIntegerAttr, reflect.Uint32: mt.UIntegerAttr, reflect.Uint64: mt.UIntegerAttr,
		reflect.Float32: mt.FloatAttr, reflect.Float64: mt.FloatAttr,
		reflect.Complex64: mt.ComplexAttr, reflect.Complex128: mt.ComplexAttr,
		reflect.String: mt.StringAttr, reflect.Slice: mt.SliceAttr, reflect.Bool: mt.BoolAttr,
		reflect.Map: mt.MapAttr, reflect.Ptr: mt.PointerAttr, reflect.Struct: mt.StructAttr, reflect.Array: mt.ArrayAttr,
	}
	retA = kindMap[t.Kind()]
	if retA != nil {
		attrsVal := retA.GetAttributes()
		if attrsVal == nil {
			retA = retA.GetDefaultImplementation()
			return
		}
		attrsValType := reflect.TypeOf(attrsVal)
		if attrsValType == nil {
			retA = retA.GetDefaultImplementation()
			return
		}
		zero := reflect.Zero(attrsValType).Interface()
		if reflect.DeepEqual(attrsVal, zero) {
			retA = retA.GetDefaultImplementation()
		}
	}
	return
}

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
	if a.Max > zero && a.Min <= a.Max {
		minVal := reflect.ValueOf(a.Min)
		maxVal := reflect.ValueOf(a.Max)

		min := minVal.Int()
		max := maxVal.Int()

		if max > min {
			result := min + rand.Int63n(max-min+1)
			resultVal := reflect.ValueOf(result).Convert(reflect.TypeOf(zero))
			return resultVal.Interface()
		}
	}
	return zero
}

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
	if a.Max > zero && a.Min <= a.Max {
		minVal := reflect.ValueOf(a.Min)
		maxVal := reflect.ValueOf(a.Max)

		min := minVal.Uint()
		max := maxVal.Uint()

		if max > min {
			diff := max - min + 1
			if diff > 0 {
				result := min + uint64(rand.Int63n(int64(diff)))
				resultVal := reflect.ValueOf(result).Convert(reflect.TypeOf(zero))
				return resultVal.Interface()
			}
		}
	}
	return zero
}

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
	if a.Max > a.Min {
		minVal := reflect.ValueOf(a.Min)
		maxVal := reflect.ValueOf(a.Max)

		min := minVal.Float()
		max := maxVal.Float()

		result := min + rand.Float64()*(max-min)

		resultVal := reflect.ValueOf(result).Convert(reflect.TypeOf(zero))
		return resultVal.Interface()
	}
	return zero
}

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

	realMin := a.RealMin
	realMax := a.RealMax
	imagMin := a.ImagMin
	imagMax := a.ImagMax

	if realMax <= realMin {
		realMin = -10.0
		realMax = 10.0
	}
	if imagMax <= imagMin {
		imagMin = -10.0
		imagMax = 10.0
	}

	realPart := realMin + rand.Float64()*(realMax-realMin)
	imagPart := imagMin + rand.Float64()*(imagMax-imagMin)

	complexVal := complex(realPart, imagPart)
	resultVal := reflect.ValueOf(complexVal).Convert(reflect.TypeOf(zero))
	return resultVal.Interface()
}

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
	minLen := a.MinLen
	maxLen := a.MaxLen

	if maxLen <= 0 {
		maxLen = 10
	}
	if minLen < 0 {
		minLen = 0
	}
	if minLen > maxLen {
		minLen = maxLen
	}

	length := minLen
	if maxLen > minLen {
		length = minLen + rand.Intn(maxLen-minLen+1)
	}

	allowedRunes := a.AllowedRunes
	if len(allowedRunes) == 0 {
		for i := 32; i <= 126; i++ {
			allowedRunes = append(allowedRunes, rune(i))
		}
	}

	result := make([]rune, length)
	for i := 0; i < length; i++ {
		result[i] = allowedRunes[rand.Intn(len(allowedRunes))]
	}

	generated := string(result)

	if a.Prefix != "" {
		generated = a.Prefix + generated
	}
	if a.Suffix != "" {
		generated = generated + a.Suffix
	}

	return generated
}

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
		} else {
			elemValue = reflect.Zero(elemType)
		}
		result.Index(i).Set(elemValue)
	}
}

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
	if a.ForceTrue {
		return true
	}
	if a.ForceFalse {
		return false
	}
	// Random boolean value
	return rand.Intn(2) == 1
}

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
	if a.AllowNil && rand.Intn(2) == 0 {
		return reflect.Zero(a.GetReflectType()).Interface()
	}

	var innerValue reflect.Value
	if attrs, ok := a.Inner.(Attributes); ok {
		randVal := attrs.GetRandomValue()
		if randVal != nil {
			innerValue = reflect.ValueOf(randVal)
		} else {
			innerType := attrs.GetReflectType()
			if innerType != nil {
				innerValue = reflect.Zero(innerType)
			} else {
				return nil
			}
		}
	} else {
		return nil
	}

	ptrValue := reflect.New(innerValue.Type())
	ptrValue.Elem().Set(innerValue)

	currentPtr := ptrValue
	for i := 1; i < a.Depth; i++ {
		newPtr := reflect.New(currentPtr.Type())
		newPtr.Elem().Set(currentPtr)
		currentPtr = newPtr
	}

	return currentPtr.Interface()
}

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
	structValue := reflect.New(structType).Elem()
	for fieldName, fieldAttr := range a.FieldAttrs {
		field := structValue.FieldByName(fieldName)
		if !field.IsValid() || !field.CanSet() {
			continue
		}
		var fieldValue reflect.Value
		if attrs, ok := fieldAttr.(Attributes); ok {
			randVal := attrs.GetRandomValue()
			if randVal != nil {
				fieldValue = reflect.ValueOf(randVal)
			} else {
				fieldValue = reflect.Zero(field.Type())
			}
		} else {
			fieldValue = reflect.Zero(field.Type())
		}
		if fieldValue.Type().AssignableTo(field.Type()) {
			field.Set(fieldValue)
		} else if fieldValue.Type().ConvertibleTo(field.Type()) {
			field.Set(fieldValue.Convert(field.Type()))
		}
	}

	return structValue.Interface()
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
	if a.Length <= 0 {
		return nil
	}

	// Get element type
	var elemType reflect.Type
	if attrs, ok := a.ElementAttrs.(Attributes); ok {
		elemType = attrs.GetReflectType()
	}

	if elemType == nil {
		return nil
	}

	arrayType := reflect.ArrayOf(a.Length, elemType)
	arrayValue := reflect.New(arrayType).Elem()

	// Generate random elements
	for i := 0; i < a.Length; i++ {
		var elemValue reflect.Value
		if attrs, ok := a.ElementAttrs.(Attributes); ok {
			randVal := attrs.GetRandomValue()
			if randVal != nil {
				elemValue = reflect.ValueOf(randVal)
			} else {
				elemValue = reflect.Zero(elemType)
			}
		} else {
			elemValue = reflect.Zero(elemType)
		}
		arrayValue.Index(i).Set(elemValue)
	}

	return arrayValue.Interface()
}
