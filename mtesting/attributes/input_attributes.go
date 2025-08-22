package attributes

import (
	"reflect"

	p "github.com/laiambryant/gotestutils/pbtesting/properties/predicates"
)

type MTAttributes struct {
	IA  IntegerAttributes
	FA  FloatAttributes
	CA  ComplexAttributes
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

func (mt MTAttributes) GetAttributeGivenType(t reflect.Type) (retA Attributes) {
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
	if retA == reflect.Zero(reflect.TypeOf(retA)).Interface() {
		retA = retA.GetDefaultImplementation()
	}
	return
}

type Attributes interface {
	GetAttributes() any
	GetReflectType() reflect.Type
	GetDefaultImplementation() Attributes
}

type IntegerAttributes struct {
	Signed        bool
	AllowNegative bool
	AllowZero     bool
	Max           uint64
	Min           int64
	EvenOnly      bool
	OddOnly       bool
	MultipleOf    uint64
	InSet         []int64
	NotInSet      []int64
}

func (a IntegerAttributes) GetAttributes() any { return a }
func (a IntegerAttributes) GetReflectType() reflect.Type {
	if a.Signed || a.AllowNegative {
		return reflect.TypeOf(int64(0))
	}
	return reflect.TypeOf(uint64(0))
}
func (a IntegerAttributes) GetDefaultImplementation() Attributes {
	return IntegerAttributes{
		Signed:        true,
		AllowNegative: true,
		AllowZero:     true,
		Max:           100,
		Min:           -100,
	}
}

type FloatAttributes struct {
	Min        float64
	Max        float64
	NonZero    bool
	FiniteOnly bool
	AllowNaN   bool
	AllowInf   bool
	Precision  uint
}

func (a FloatAttributes) GetAttributes() any           { return a }
func (a FloatAttributes) GetReflectType() reflect.Type { return reflect.TypeOf(float64(0)) }
func (a FloatAttributes) GetDefaultImplementation() Attributes {
	return FloatAttributes{
		Min:        -100.0,
		Max:        100.0,
		NonZero:    true,
		FiniteOnly: true,
	}
}

type ComplexAttributes struct {
	RealMin      float64
	RealMax      float64
	ImagMin      float64
	ImagMax      float64
	MagnitudeMin float64
	MagnitudeMax float64
	AllowNaN     bool
	AllowInf     bool
}

func (a ComplexAttributes) GetAttributes() any           { return a }
func (a ComplexAttributes) GetReflectType() reflect.Type { return reflect.TypeOf(complex128(0)) }
func (a ComplexAttributes) GetDefaultImplementation() Attributes {
	return ComplexAttributes{
		RealMin: -10.0,
		RealMax: 10.0,
		ImagMin: -10.0,
		ImagMax: 10.0,
	}
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
		MinLen: 1,
		MaxLen: 5,
		ElementAttrs: IntegerAttributes{
			Signed: true,
		},
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
		ValueAttrs: IntegerAttributes{
			Signed: true,
		},
	}
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
		ElemAttrs: FloatAttributes{
			Min: -10.0,
			Max: 10.0,
		},
	}
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
		Inner: IntegerAttributes{
			Signed: true,
		},
	}
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
			"Field1": IntegerAttributes{
				Signed: true,
			},
			"Field2": FloatAttributes{
				Min: -10.0,
				Max: 10.0,
			},
		},
	}
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
		Length: 5,
		ElementAttrs: IntegerAttributes{
			Signed: true,
		},
	}
}
