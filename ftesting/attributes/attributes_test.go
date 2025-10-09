package attributes

import (
	"reflect"
	"testing"

	ctesting "github.com/laiambryant/gotestutils/ctesting"
	p "github.com/laiambryant/gotestutils/pbtesting/properties/predicates"
)

func TestGetAttributesMethods(t *testing.T) {
	type testCase struct {
		name string
		in   Attributes
		want any
	}
	cases := []testCase{
		{"IntegerAttributesImpl", IntegerAttributesImpl[int64]{AllowNegative: true, AllowZero: true, Max: 10, Min: -5}, IntegerAttributesImpl[int64]{AllowNegative: true, AllowZero: true, Max: 10, Min: -5}},
		{"UnsignedIntegerAttributesImpl", UnsignedIntegerAttributesImpl[uint64]{AllowNegative: false, AllowZero: true, Max: 100, Min: 0}, UnsignedIntegerAttributesImpl[uint64]{AllowNegative: false, AllowZero: true, Max: 100, Min: 0}},
		{"FloatAttributesImpl", FloatAttributesImpl[float64]{Min: 1.1, Max: 2.2, NonZero: true, FiniteOnly: true, AllowNaN: true, AllowInf: true, Precision: 3}, FloatAttributesImpl[float64]{Min: 1.1, Max: 2.2, NonZero: true, FiniteOnly: true, AllowNaN: true, AllowInf: true, Precision: 3}},
		{"ComplexAttributesImpl", ComplexAttributesImpl[complex128]{RealMin: -1, RealMax: 1, ImagMin: -2, ImagMax: 2, MagnitudeMin: 0.5, MagnitudeMax: 10, AllowNaN: true, AllowInf: true}, ComplexAttributesImpl[complex128]{RealMin: -1, RealMax: 1, ImagMin: -2, ImagMax: 2, MagnitudeMin: 0.5, MagnitudeMax: 10, AllowNaN: true, AllowInf: true}},
		{"StringAttributes", StringAttributes{MinLen: 1, MaxLen: 5, Prefix: "pre", Suffix: "suf", Contains: "mid", UniqueChars: true}, StringAttributes{MinLen: 1, MaxLen: 5, Prefix: "pre", Suffix: "suf", Contains: "mid", UniqueChars: true}},
		{"SliceAttributes", SliceAttributes{MinLen: 1, MaxLen: 3, Unique: true, Sorted: true, ElementPreds: []p.Predicate{}, ElementAttrs: IntegerAttributesImpl[int64]{}}, SliceAttributes{MinLen: 1, MaxLen: 3, Unique: true, Sorted: true, ElementPreds: []p.Predicate{}, ElementAttrs: IntegerAttributesImpl[int64]{}}},
		{"BoolAttributes", BoolAttributes{ForceTrue: true}, BoolAttributes{ForceTrue: true}},
		{"MapAttributes", MapAttributes{MinSize: 1, MaxSize: 3, KeyPreds: []p.Predicate{}, ValuePreds: []p.Predicate{}, KeyAttrs: StringAttributes{}, ValueAttrs: IntegerAttributesImpl[int64]{}}, MapAttributes{MinSize: 1, MaxSize: 3, KeyPreds: []p.Predicate{}, ValuePreds: []p.Predicate{}, KeyAttrs: StringAttributes{}, ValueAttrs: IntegerAttributesImpl[int64]{}}},
		{"PointerAttributes", PointerAttributes{AllowNil: true, Depth: 2, Inner: IntegerAttributesImpl[int64]{}}, PointerAttributes{AllowNil: true, Depth: 2, Inner: IntegerAttributesImpl[int64]{}}},
		{"StructAttributes", StructAttributes{FieldAttrs: map[string]any{"A": IntegerAttributesImpl[int64]{}, "B": FloatAttributesImpl[float64]{}}}, StructAttributes{FieldAttrs: map[string]any{"A": IntegerAttributesImpl[int64]{}, "B": FloatAttributesImpl[float64]{}}}},
		{"ArrayAttributes", ArrayAttributes{Length: 3, Sorted: true, ElementAttrs: IntegerAttributesImpl[int64]{}}, ArrayAttributes{Length: 3, Sorted: true, ElementAttrs: IntegerAttributesImpl[int64]{}}},
	}
	var suite []ctesting.CharacterizationTest[bool]
	for _, tc := range cases {
		toExec := tc
		suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
			got := toExec.in.GetAttributes()
			if reflect.TypeOf(got) != reflect.TypeOf(toExec.want) {
				return false, nil
			}
			if !reflect.DeepEqual(got, toExec.want) {
				return false, nil
			}
			return true, nil
		}))
	}
	results, _ := ctesting.VerifyCharacterizationTestsAndResults(t, suite, true)
	for i, passed := range results {
		if !passed {
			t.Fatalf("characterization test %d failed", i+1)
		}
	}
}

func TestGetReflectTypeMethods(t *testing.T) {
	type testCase struct {
		name     string
		in       Attributes
		wantType reflect.Type
	}
	cases := []testCase{
		{"IntegerAttributesImpl", IntegerAttributesImpl[int64]{}, reflect.TypeOf(int64(0))},
		{"UnsignedIntegerAttributesImpl", UnsignedIntegerAttributesImpl[uint64]{}, reflect.TypeOf(uint64(0))},
		{"FloatAttributesImpl", FloatAttributesImpl[float64]{}, reflect.TypeOf(float64(0))},
		{"ComplexAttributesImpl", ComplexAttributesImpl[complex128]{}, reflect.TypeOf(complex128(0))},
		{"StringAttributes", StringAttributes{}, reflect.TypeOf("")},
		{"BoolAttributes", BoolAttributes{}, reflect.TypeOf(true)},
	}
	var suite []ctesting.CharacterizationTest[bool]
	for _, tc := range cases {
		toExec := tc
		suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
			got := toExec.in.GetReflectType()
			return got == toExec.wantType, nil
		}))
	}
	results, _ := ctesting.VerifyCharacterizationTestsAndResults(t, suite, true)
	for i, passed := range results {
		if !passed {
			t.Fatalf("GetReflectType test %d failed", i+1)
		}
	}
}

func TestGetDefaultImplementationMethods(t *testing.T) {
	type testCase struct {
		name string
		in   Attributes
	}
	cases := []testCase{
		{"IntegerAttributesImpl", IntegerAttributesImpl[int64]{}},
		{"UnsignedIntegerAttributesImpl", UnsignedIntegerAttributesImpl[uint64]{}},
		{"FloatAttributesImpl", FloatAttributesImpl[float64]{}},
		{"ComplexAttributesImpl", ComplexAttributesImpl[complex128]{}},
		{"StringAttributes", StringAttributes{}},
		{"SliceAttributes", SliceAttributes{}},
		{"BoolAttributes", BoolAttributes{}},
		{"MapAttributes", MapAttributes{}},
		{"PointerAttributes", PointerAttributes{}},
		{"StructAttributes", StructAttributes{}},
		{"ArrayAttributes", ArrayAttributes{}},
	}
	var suite []ctesting.CharacterizationTest[bool]
	for _, tc := range cases {
		toExec := tc
		suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
			got := toExec.in.GetDefaultImplementation()
			if got == nil {
				return false, nil
			}
			return reflect.TypeOf(got) == reflect.TypeOf(toExec.in), nil
		}))
	}
	results, _ := ctesting.VerifyCharacterizationTestsAndResults(t, suite, true)
	for i, passed := range results {
		if !passed {
			t.Fatalf("GetDefaultImplementation test %d failed", i+1)
		}
	}
}

func TestGetRandomValueMethods(t *testing.T) {
	type testCase struct {
		name string
		in   Attributes
	}
	cases := []testCase{
		{"IntegerAttributesImpl", IntegerAttributesImpl[int64]{Min: -10, Max: 10}},
		{"UnsignedIntegerAttributesImpl", UnsignedIntegerAttributesImpl[uint64]{Min: 0, Max: 100}},
		{"FloatAttributesImpl", FloatAttributesImpl[float64]{Min: -1.0, Max: 1.0}},
		{"ComplexAttributesImpl", ComplexAttributesImpl[complex128]{RealMin: -1.0, RealMax: 1.0, ImagMin: -1.0, ImagMax: 1.0}},
		{"StringAttributes", StringAttributes{MinLen: 1, MaxLen: 10}},
		{"SliceAttributes", SliceAttributes{MinLen: 1, MaxLen: 3, ElementAttrs: IntegerAttributesImpl[int]{Min: 0, Max: 10}}},
		{"BoolAttributes", BoolAttributes{}},
		{"MapAttributes", MapAttributes{MinSize: 1, MaxSize: 3, KeyAttrs: StringAttributes{MinLen: 1, MaxLen: 5}, ValueAttrs: IntegerAttributesImpl[int]{Min: 0, Max: 10}}},
		{"PointerAttributes", PointerAttributes{AllowNil: true, Depth: 1, Inner: IntegerAttributesImpl[int]{Min: 0, Max: 10}}},
		{"StructAttributes", StructAttributes{FieldAttrs: map[string]any{"TestField": IntegerAttributesImpl[int]{Min: 0, Max: 10}}}},
		{"ArrayAttributes", ArrayAttributes{Length: 3, ElementAttrs: IntegerAttributesImpl[int]{Min: 0, Max: 10}}},
	}
	var suite []ctesting.CharacterizationTest[bool]
	for _, tc := range cases {
		toExec := tc
		suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
			got := toExec.in.GetRandomValue()
			return got != nil || isNilValidForType(toExec.in), nil
		}))
	}
	results, _ := ctesting.VerifyCharacterizationTestsAndResults(t, suite, true)
	for i, passed := range results {
		if !passed {
			t.Fatalf("GetRandomValue test %d failed", i+1)
		}
	}
}

func TestAttributeErrors(t *testing.T) {
	uate := UnsupportedAttributeTypeError{reflect.TypeFor[int]().Kind()}
	naate := NotAnAttributeTypeError{reflect.TypeFor[int]()}
	if uate.k == reflect.Invalid || naate.Type == nil {
		t.Error("error setting type")
	}
}

func isNilValidForType(attr Attributes) bool {
	switch attr.(type) {
	case PointerAttributes:
		return true
	default:
		return false
	}
}

func TestGetAttributeGivenTypeNil(t *testing.T) {
	attributes := NewFTAttributes()
	_, err := attributes.GetAttributeGivenType(nil)
	if err == nil {
		t.Error("expected NilTypeError")
	}
	if _, ok := err.(NilTypeError); !ok {
		t.Error("expected error to be of type NilTypeError")
	}
}

func TestGetAttributeGivenType_KindMapHit(t *testing.T) {
	attributes := NewFTAttributes()
	testCases := []struct {
		name     string
		typ      reflect.Type
		expected Attributes
	}{
		{"int", reflect.TypeOf(int(0)), attributes.IntegerAttr},
		{"int32", reflect.TypeOf(int32(0)), attributes.IntegerAttr},
		{"uint", reflect.TypeOf(uint(0)), attributes.UIntegerAttr},
		{"float64", reflect.TypeOf(float64(0)), attributes.FloatAttr},
		{"complex128", reflect.TypeOf(complex128(0)), attributes.ComplexAttr},
		{"string", reflect.TypeOf(""), attributes.StringAttr},
		{"bool", reflect.TypeOf(true), attributes.BoolAttr},
		{"slice", reflect.TypeOf([]int{}), attributes.SliceAttr},
		{"map", reflect.TypeOf(map[string]int{}), attributes.MapAttr},
		{"pointer", reflect.TypeOf(new(int)), attributes.PointerAttr},
		{"struct", reflect.TypeOf(struct{}{}), attributes.StructAttr},
		{"array", reflect.TypeOf([3]int{}), attributes.ArrayAttr},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := attributes.GetAttributeGivenType(tc.typ)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if reflect.TypeOf(result) != reflect.TypeOf(tc.expected) {
				t.Errorf("expected type %T, got %T", tc.expected, result)
			}
		})
	}
}

func TestGetAttributeGivenType_KindNotInMap(t *testing.T) {
	attributes := NewFTAttributes()
	testCases := []struct {
		name string
		typ  reflect.Type
	}{
		{"chan", reflect.TypeOf(make(chan int))},
		{"func", reflect.TypeOf(func() {})},
		{"interface", reflect.TypeOf((*interface{})(nil)).Elem()},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := attributes.GetAttributeGivenType(tc.typ)
			if err == nil {
				t.Error("expected error for unsupported type")
			}
			if _, ok := err.(UnsupportedAttributeTypeError); !ok {
				t.Errorf("expected UnsupportedAttributeTypeError, got %T: %v", err, err)
			}
		})
	}
}

func TestGetAttributeGivenType_NilAttribute(t *testing.T) {
	attributes := FTAttributes{
		IntegerAttr: nil,
	}
	result, err := attributes.GetAttributeGivenType(reflect.TypeOf(int(0)))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result == nil {
		t.Error("expected default implementation, got nil")
	}
}

type nilAttributeType struct{}

func (n nilAttributeType) GetAttributes() any           { return nil }
func (n nilAttributeType) GetReflectType() reflect.Type { return reflect.TypeOf(int(0)) }
func (n nilAttributeType) GetDefaultImplementation() Attributes {
	return IntegerAttributesImpl[int]{AllowNegative: true, AllowZero: true, Max: 100, Min: -100}
}
func (n nilAttributeType) GetRandomValue() any { return 42 }

func TestGetAttributeGivenType_GetAttributesReturnsNil(t *testing.T) {
	attributes := FTAttributes{
		IntegerAttr: nilAttributeType{},
	}
	result, err := attributes.GetAttributeGivenType(reflect.TypeOf(int(0)))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if _, ok := result.(IntegerAttributesImpl[int]); !ok {
		t.Errorf("expected IntegerAttributesImpl, got %T", result)
	}
}

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

func TestGetAttributeGivenType_TypeOfAttributesIsNil(t *testing.T) {
	attributes := FTAttributes{
		IntegerAttr: nilTypeAttribute{},
	}
	result, err := attributes.GetAttributeGivenType(reflect.TypeOf(int(0)))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result == nil {
		t.Error("expected result, got nil")
	}
}

func TestGetAttributeGivenType_ZeroValueAttribute(t *testing.T) {
	attributes := FTAttributes{
		IntegerAttr: IntegerAttributesImpl[int]{}, // Zero value
	}
	result, err := attributes.GetAttributeGivenType(reflect.TypeOf(int(0)))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result == nil {
		t.Error("expected default implementation, got nil")
	}
	defaultImpl := IntegerAttributesImpl[int]{}.GetDefaultImplementation()
	if !reflect.DeepEqual(result, defaultImpl) {
		t.Errorf("expected default implementation, got different value")
	}
}

func TestGetAttributeGivenType_NonZeroValueAttribute(t *testing.T) {
	customAttr := IntegerAttributesImpl[int]{
		AllowNegative: true,
		AllowZero:     false,
		Max:           50,
		Min:           10,
	}
	attributes := FTAttributes{
		IntegerAttr: customAttr,
	}
	result, err := attributes.GetAttributeGivenType(reflect.TypeOf(int(0)))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result == nil {
		t.Error("expected custom attribute, got nil")
	}
	if !reflect.DeepEqual(result, customAttr) {
		t.Errorf("expected custom attribute to be returned as-is")
	}
}

func TestGetDefaultForKind_IntegerTypes(t *testing.T) {
	attributes := NewFTAttributes()
	intKinds := []reflect.Kind{
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
	}
	for _, kind := range intKinds {
		t.Run(kind.String(), func(t *testing.T) {
			result, err := attributes.getDefaultForKind(kind)
			if err != nil {
				t.Errorf("unexpected error for %s: %v", kind, err)
			}
			if result == nil {
				t.Errorf("expected result for %s, got nil", kind)
			}
			expected := IntegerAttributesImpl[int64]{}.GetDefaultImplementation()
			if !reflect.DeepEqual(result, expected) {
				t.Errorf("expected IntegerAttributesImpl default for %s", kind)
			}
		})
	}
}

func TestGetDefaultForKind_UnsignedIntegerTypes(t *testing.T) {
	attributes := NewFTAttributes()
	uintKinds := []reflect.Kind{
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
	}
	for _, kind := range uintKinds {
		t.Run(kind.String(), func(t *testing.T) {
			result, err := attributes.getDefaultForKind(kind)
			if err != nil {
				t.Errorf("unexpected error for %s: %v", kind, err)
			}
			if result == nil {
				t.Errorf("expected result for %s, got nil", kind)
			}
			expected := UnsignedIntegerAttributesImpl[uint64]{}.GetDefaultImplementation()
			if !reflect.DeepEqual(result, expected) {
				t.Errorf("expected UnsignedIntegerAttributesImpl default for %s", kind)
			}
		})
	}
}

func TestGetDefaultForKind_FloatTypes(t *testing.T) {
	attributes := NewFTAttributes()
	floatKinds := []reflect.Kind{reflect.Float32, reflect.Float64}
	for _, kind := range floatKinds {
		t.Run(kind.String(), func(t *testing.T) {
			result, err := attributes.getDefaultForKind(kind)
			if err != nil {
				t.Errorf("unexpected error for %s: %v", kind, err)
			}
			if result == nil {
				t.Errorf("expected result for %s, got nil", kind)
			}
			expected := FloatAttributesImpl[float64]{}.GetDefaultImplementation()
			if !reflect.DeepEqual(result, expected) {
				t.Errorf("expected FloatAttributesImpl default for %s", kind)
			}
		})
	}
}

func TestGetDefaultForKind_ComplexTypes(t *testing.T) {
	attributes := NewFTAttributes()
	complexKinds := []reflect.Kind{reflect.Complex64, reflect.Complex128}
	for _, kind := range complexKinds {
		t.Run(kind.String(), func(t *testing.T) {
			result, err := attributes.getDefaultForKind(kind)
			if err != nil {
				t.Errorf("unexpected error for %s: %v", kind, err)
			}
			if result == nil {
				t.Errorf("expected result for %s, got nil", kind)
			}
			expected := ComplexAttributesImpl[complex128]{}.GetDefaultImplementation()
			if !reflect.DeepEqual(result, expected) {
				t.Errorf("expected ComplexAttributesImpl default for %s", kind)
			}
		})
	}
}

func TestGetDefaultForKind_OtherSupportedTypes(t *testing.T) {
	attributes := NewFTAttributes()
	testCases := []struct {
		kind     reflect.Kind
		expected Attributes
	}{
		{reflect.String, StringAttributes{}.GetDefaultImplementation()},
		{reflect.Slice, SliceAttributes{}.GetDefaultImplementation()},
		{reflect.Bool, BoolAttributes{}.GetDefaultImplementation()},
		{reflect.Map, MapAttributes{}.GetDefaultImplementation()},
		{reflect.Pointer, PointerAttributes{}.GetDefaultImplementation()},
		{reflect.Struct, StructAttributes{}.GetDefaultImplementation()},
		{reflect.Array, ArrayAttributes{}.GetDefaultImplementation()},
	}
	for _, tc := range testCases {
		t.Run(tc.kind.String(), func(t *testing.T) {
			result, err := attributes.getDefaultForKind(tc.kind)
			if err != nil {
				t.Errorf("unexpected error for %s: %v", tc.kind, err)
			}
			if result == nil {
				t.Errorf("expected result for %s, got nil", tc.kind)
			}
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("expected correct default implementation for %s", tc.kind)
			}
		})
	}
}

func TestGetDefaultForKind_UnsupportedTypes(t *testing.T) {
	attributes := NewFTAttributes()
	unsupportedKinds := []reflect.Kind{
		reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Invalid,
		reflect.Uintptr,
		reflect.UnsafePointer,
	}

	for _, kind := range unsupportedKinds {
		t.Run(kind.String(), func(t *testing.T) {
			result, err := attributes.getDefaultForKind(kind)
			if err == nil {
				t.Errorf("expected error for unsupported kind %s", kind)
			}
			if result != nil {
				t.Errorf("expected nil result for unsupported kind %s, got %v", kind, result)
			}
			if _, ok := err.(UnsupportedAttributeTypeError); !ok {
				t.Errorf("expected UnsupportedAttributeTypeError for %s, got %T: %v", kind, err, err)
			}
		})
	}
}

func TestGetDefaultForKind_AllKindsCovered(t *testing.T) {
	attributes := NewFTAttributes()
	supportedKinds := map[reflect.Kind]bool{
		reflect.Int: true, reflect.Int8: true, reflect.Int16: true, reflect.Int32: true, reflect.Int64: true,
		reflect.Uint: true, reflect.Uint8: true, reflect.Uint16: true, reflect.Uint32: true, reflect.Uint64: true,
		reflect.Float32: true, reflect.Float64: true,
		reflect.Complex64: true, reflect.Complex128: true,
		reflect.String: true, reflect.Slice: true, reflect.Bool: true,
		reflect.Map: true, reflect.Pointer: true, reflect.Struct: true, reflect.Array: true,
	}
	for kind := reflect.Invalid; kind <= reflect.UnsafePointer; kind++ {
		result, err := attributes.getDefaultForKind(kind)
		if supportedKinds[kind] {
			if err != nil {
				t.Errorf("expected no error for supported kind %s, got: %v", kind, err)
			}
			if result == nil {
				t.Errorf("expected non-nil result for supported kind %s", kind)
			}
		} else {
			if err == nil {
				t.Errorf("expected error for unsupported kind %s", kind)
			}
			if result != nil {
				t.Errorf("expected nil result for unsupported kind %s", kind)
			}
		}
	}
}
