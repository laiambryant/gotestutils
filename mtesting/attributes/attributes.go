package attributes

import (
	"reflect"

	p "github.com/laiambryant/gotestutils/pbtesting/properties/predicates"
)

type MTAttributes[T any] struct {
	IA  IntegerAttributes[int]
	FA  FloatAttributes[float32]
	CA  ComplexAttributes[complex64]
	SA  StringAttributes
	SLA SliceAttributes
	BA  BoolAttributes
	MA  MapAttributes
	CHA ChanAttributes
	FNA FuncAttributes
	INA InterfaceAttributes
	PA  PointerAttributes
	STA StructAttributes
	ARA ArrayAttributes
}

func (mt MTAttributes[T]) GetAttributeGivenType(t reflect.Type) (retA Attributes) {
	kindMap := map[reflect.Kind]Attributes{
		reflect.Int: mt.IA, reflect.Int8: mt.IA, reflect.Int16: mt.IA, reflect.Int32: mt.IA, reflect.Int64: mt.IA,
		reflect.Uint: mt.IA, reflect.Uint8: mt.IA, reflect.Uint16: mt.IA, reflect.Uint32: mt.IA, reflect.Uint64: mt.IA,
		reflect.Float32: mt.FA, reflect.Float64: mt.FA,
		reflect.Complex64: mt.CA, reflect.Complex128: mt.CA,
		reflect.String: mt.SA, reflect.Slice: mt.SLA, reflect.Bool: mt.BA,
		reflect.Map: mt.MA, reflect.Chan: mt.CHA, reflect.Func: mt.FNA,
		reflect.Interface: mt.INA, reflect.Ptr: mt.PA, reflect.Struct: mt.STA, reflect.Array: mt.ARA,
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

type Attributes interface {
	GetAttributes() any
	GetReflectType() reflect.Type
	GetDefaultImplementation() Attributes
	GetRandomValue() any
}

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

func (a IntegerAttributes[T]) GetRandomValue() any {
	// TODO: Implement random value generation
	return nil
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

func (a UnsignedIntegerAttributes[T]) GetRandomValue() any {
	// TODO: Implement random value generation
	return nil
}

type Floats interface {
	float32 | float64
}

type FloatAttributes[T Floats] struct {
	Min        T
	Max        T
	NonZero    bool
	FiniteOnly bool
	AllowNaN   bool
	AllowInf   bool
	Precision  uint
}

func (a FloatAttributes[T]) GetAttributes() any           { return a }
func (a FloatAttributes[T]) GetReflectType() reflect.Type { return reflect.TypeOf(float64(0)) }
func (a FloatAttributes[T]) GetDefaultImplementation() Attributes {
	return FloatAttributes[T]{
		Min:        -100.0,
		Max:        100.0,
		NonZero:    true,
		FiniteOnly: true,
	}
}

func (a FloatAttributes[T]) GetRandomValue() any {
	// TODO: Implement random value generation
	return nil
}

type Complex interface {
	complex64 | complex128
}

type ComplexAttributes[T Complex] struct {
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

func (a ComplexAttributes[T]) GetAttributes() any           { return a }
func (a ComplexAttributes[T]) GetReflectType() reflect.Type { return reflect.TypeOf(complex128(0)) }
func (a ComplexAttributes[T]) GetDefaultImplementation() Attributes {
	return ComplexAttributes[T]{
		RealMin: -10.0,
		RealMax: 10.0,
		ImagMin: -10.0,
		ImagMax: 10.0,
	}
}

func (a ComplexAttributes[T]) GetRandomValue() any {
	// TODO: Implement random value generation
	return nil
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
	// TODO: Implement random value generation
	return nil
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
		ElementAttrs: IntegerAttributes[int]{},
	}
}

func (a SliceAttributes) GetRandomValue() any {
	// TODO: Implement random value generation
	return nil
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
	// TODO: Implement random value generation
	return nil
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
		ValueAttrs: IntegerAttributes[int]{},
	}
}

func (a MapAttributes) GetRandomValue() any {
	// TODO: Implement random value generation
	return nil
}

type ChanAttributes struct {
	MinBuffer int
	MaxBuffer int
	ElemAttrs any
}

func (a ChanAttributes) GetAttributes() any { return a }
func (a ChanAttributes) GetReflectType() reflect.Type {
	var et reflect.Type
	switch v := a.ElemAttrs.(type) {
	case Attributes:
		et = v.GetReflectType()
	case reflect.Type:
		et = v
	}
	if et == nil {
		return nil
	}
	return reflect.ChanOf(reflect.BothDir, et)
}

func (a ChanAttributes) GetDefaultImplementation() Attributes {
	return ChanAttributes{
		MinBuffer: 0,
		MaxBuffer: 10,
		ElemAttrs: FloatAttributes[float32]{},
	}
}

func (a ChanAttributes) GetRandomValue() any {
	// TODO: Implement random value generation
	return nil
}

type FuncAttributes struct {
	Deterministic    bool
	PanicProbability float64
	ReturnZeroValues bool
}

func (a FuncAttributes) GetAttributes() any           { return a }
func (a FuncAttributes) GetReflectType() reflect.Type { return reflect.TypeOf(func() {}) }
func (a FuncAttributes) GetDefaultImplementation() Attributes {
	return FuncAttributes{
		Deterministic: true,
	}
}

func (a FuncAttributes) GetRandomValue() any {
	// TODO: Implement random value generation
	return nil
}

type InterfaceAttributes struct {
	AllowedConcrete []reflect.Type
}

func (a InterfaceAttributes) GetAttributes() any { return a }
func (a InterfaceAttributes) GetReflectType() reflect.Type {
	return reflect.TypeOf((*any)(nil)).Elem()
}

func (a InterfaceAttributes) GetDefaultImplementation() Attributes {
	return InterfaceAttributes{
		AllowedConcrete: []reflect.Type{
			reflect.TypeOf(1),
			reflect.TypeOf(""),
		},
	}
}

func (a InterfaceAttributes) GetRandomValue() any {
	// TODO: Implement random value generation
	return nil
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
	d := a.Depth
	if d <= 0 {
		d = 1
	}
	t := inner
	for i := 0; i < d; i++ {
		t = reflect.PointerTo(t)
	}
	return t
}

func (a PointerAttributes) GetDefaultImplementation() Attributes {
	return PointerAttributes{
		AllowNil: true,
		Depth:    1,
		Inner:    IntegerAttributes[int]{},
	}
}

func (a PointerAttributes) GetRandomValue() any {
	// TODO: Implement random value generation
	return nil
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
			"Field1": IntegerAttributes[int]{},
			"Field2": FloatAttributes[float32]{
				Min: -10.0,
				Max: 10.0,
			},
		},
	}
}

func (a StructAttributes) GetRandomValue() any {
	// TODO: Implement random value generation
	return nil
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
		ElementAttrs: IntegerAttributes[int]{},
	}
}

func (a ArrayAttributes) GetRandomValue() any {
	// TODO: Implement random value generation
	return nil
}
