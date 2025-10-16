package attributes

import (
	"reflect"
	"testing"

	"github.com/laiambryant/gotestutils/ctesting"
)

func TestStructAttributes(t *testing.T) {
	var suite []ctesting.CharacterizationTest[bool]

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StructAttributes{FieldAttrs: map[string]any{"Field1": IntegerAttributesImpl[int]{}}}
		got := attr.GetAttributes()
		expected := StructAttributes{FieldAttrs: map[string]any{"Field1": IntegerAttributesImpl[int]{}}}
		return reflect.DeepEqual(got, expected), nil
	}))
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StructAttributes{}
		got := attr.GetDefaultImplementation()
		return got != nil && reflect.TypeOf(got) == reflect.TypeOf(attr), nil
	}))
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StructAttributes{FieldAttrs: map[string]any{"Field1": IntegerAttributesImpl[int]{}}}
		got := attr.GetRandomValue()
		return got != nil, nil
	}))
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StructAttributes{FieldAttrs: map[string]any{}}
		result := attr.GetRandomValue()
		return result == nil, nil
	}))
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StructAttributes{FieldAttrs: map[string]any{"Field1": "not an attribute"}}
		result := attr.GetRandomValue()
		return result == nil, nil
	}))
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attrs := StructAttributes{
			FieldAttrs: map[string]any{
				"Field1": reflect.TypeOf(int(0)),
				"Field2": reflect.TypeOf(""),
			},
		}
		reflectType := attrs.GetReflectType()
		return reflectType != nil && reflectType.Kind() == reflect.Struct && reflectType.NumField() == 2, nil
	}))
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attrs := StructAttributes{
			FieldAttrs: map[string]any{
				"Field1": nilTypeReturningAttribute{},
			},
		}
		reflectType := attrs.GetReflectType()
		return reflectType == nil, nil
	}))
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attrs := StructAttributes{
			FieldAttrs: nil,
		}
		reflectType := attrs.GetReflectType()
		return reflectType == nil, nil
	}))
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attrs := StructAttributes{
			FieldAttrs: map[string]any{
				"Field1": nilReturningAttribute{},
			},
		}
		result := attrs.GetRandomValue()
		if result == nil {
			return false, nil
		}
		structValue := reflect.ValueOf(result)
		field := structValue.FieldByName("Field1")
		return field.IsValid() && field.Int() == 0, nil
	}))

	results, _ := ctesting.VerifyCharacterizationTestsAndResults(t, suite, true)
	for i, passed := range results {
		if !passed {
			t.Fatalf("StructAttributes test %d failed", i+1)
		}
	}
}

type CustomString string

func TestStructAttributes_SetFieldValueConversion(t *testing.T) {
	type TestStruct struct {
		IntField    int32
		FloatField  float32
		StringField CustomString
	}
	attrs := StructAttributes{
		FieldAttrs: map[string]any{
			"IntField":    IntegerAttributesImpl[int]{Min: 1, Max: 10},
			"FloatField":  FloatAttributesImpl[float64]{Min: 1.0, Max: 10.0},
			"StringField": StringAttributes{MinLen: 3, MaxLen: 5},
		},
	}
	structType := reflect.TypeOf(TestStruct{})
	structValue := reflect.New(structType).Elem()
	intAttr := attrs.FieldAttrs["IntField"].(IntegerAttributesImpl[int])
	intRandVal := intAttr.GetRandomValue()
	intFieldValue := reflect.ValueOf(intRandVal)
	intField := structValue.FieldByName("IntField")
	if intFieldValue.Type().AssignableTo(intField.Type()) {
		t.Error("Expected int not to be directly assignable to int32")
	}
	if !intFieldValue.Type().ConvertibleTo(intField.Type()) {
		t.Error("Expected int to be convertible to int32")
	}
	attrs.setFieldValue(intField, intFieldValue)
	if intField.Interface().(int32) == 0 {
		t.Error("Expected int32 field to be set via conversion")
	}
	floatAttr := attrs.FieldAttrs["FloatField"].(FloatAttributesImpl[float64])
	floatRandVal := floatAttr.GetRandomValue()
	floatFieldValue := reflect.ValueOf(floatRandVal)
	floatField := structValue.FieldByName("FloatField")
	if floatFieldValue.Type().AssignableTo(floatField.Type()) {
		t.Error("Expected float64 not to be directly assignable to float32")
	}
	if !floatFieldValue.Type().ConvertibleTo(floatField.Type()) {
		t.Error("Expected float64 to be convertible to float32")
	}
	attrs.setFieldValue(floatField, floatFieldValue)
	if floatField.Interface().(float32) == 0.0 {
		t.Error("Expected float32 field to be set via conversion")
	}
	result := attrs.GetRandomValue()
	if result == nil {
		t.Fatal("Expected non-nil struct with type conversions")
	}
	resultValue := reflect.ValueOf(result)
	if resultValue.Kind() != reflect.Struct {
		t.Fatal("Expected struct result")
	}
	intFieldResult := resultValue.FieldByName("IntField")
	if !intFieldResult.IsValid() || intFieldResult.Interface().(int) == 0 {
		t.Error("Expected IntField to be set via conversion")
	}
	floatFieldResult := resultValue.FieldByName("FloatField")
	if !floatFieldResult.IsValid() || floatFieldResult.Interface().(float64) == 0.0 {
		t.Error("Expected FloatField to be set via conversion")
	}
	stringFieldResult := resultValue.FieldByName("StringField")
	if !stringFieldResult.IsValid() || stringFieldResult.Interface().(string) == "" {
		t.Error("Expected StringField to be set via conversion")
	}
}
