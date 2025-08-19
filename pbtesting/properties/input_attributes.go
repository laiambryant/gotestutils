package properties

import (
	"reflect"

	p "github.com/laiambryant/gotestutils/pbtesting/properties/predicates"
)

type Attributes interface {
	GetAttributes() any
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

type FloatAttributes struct {
	Min        float64
	Max        float64
	NonZero    bool
	FiniteOnly bool
	AllowNaN   bool
	AllowInf   bool
	Precision  uint
}

func (a FloatAttributes) GetAttributes() any { return a }

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

func (a ComplexAttributes) GetAttributes() any { return a }

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

func (a StringAttributes) GetAttributes() any { return a }

type SliceAttributes struct {
	MinLen       int
	MaxLen       int
	Unique       bool
	Sorted       bool
	ElementPreds []p.Predicate
	ElementAttrs any
}

func (a SliceAttributes) GetAttributes() any { return a }

type BoolAttributes struct {
	ForceTrue  bool
	ForceFalse bool
}

func (a BoolAttributes) GetAttributes() any { return a }

type MapAttributes struct {
	MinSize    int
	MaxSize    int
	KeyPreds   []p.Predicate
	ValuePreds []p.Predicate
	KeyAttrs   any
	ValueAttrs any
}

func (a MapAttributes) GetAttributes() any { return a }

type ChanAttributes struct {
	MinBuffer int
	MaxBuffer int
	ElemAttrs any
}

func (a ChanAttributes) GetAttributes() any { return a }

type FuncAttributes struct {
	Deterministic    bool
	PanicProbability float64
	ReturnZeroValues bool
}

func (a FuncAttributes) GetAttributes() any { return a }

type InterfaceAttributes struct {
	AllowedConcrete []reflect.Type
}

func (a InterfaceAttributes) GetAttributes() any { return a }

type PointerAttributes struct {
	AllowNil bool
	Depth    int
	Inner    any
}

func (a PointerAttributes) GetAttributes() any { return a }

type StructAttributes struct {
	FieldAttrs map[string]any
}

func (a StructAttributes) GetAttributes() any { return a }

type ArrayAttributes struct {
	Length       int
	Sorted       bool
	ElementAttrs any
	ElementPreds []p.Predicate
}

func (a ArrayAttributes) GetAttributes() any { return a }
