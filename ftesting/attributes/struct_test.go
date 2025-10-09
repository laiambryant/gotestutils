package attributes

import (
	"reflect"
	"testing"
)

func TestStructAttributes_EmptyFieldAttrs(t *testing.T) {
	attr := StructAttributes{FieldAttrs: map[string]any{}}
	result := attr.GetRandomValue()
	if result != nil {
		t.Errorf("Expected nil for empty field attributes, got %v", result)
	}
}

func TestStructAttributes_InvalidFieldType(t *testing.T) {
	attr := StructAttributes{FieldAttrs: map[string]any{"Field1": "not an attribute"}}
	result := attr.GetRandomValue()
	if result != nil {
		t.Errorf("Expected nil for invalid field type, got %v", result)
	}
}

func TestStructAttributes_TypeConversion(t *testing.T) {
	attrs := StructAttributes{
		FieldAttrs: map[string]any{
			"CustomField": IntegerAttributesImpl[int]{
				Max: 50,
			},
		},
	}

	fields := []reflect.StructField{
		{
			Name: "CustomField",
			Type: reflect.TypeOf(CustomInt(0)),
		},
	}
	structType := reflect.StructOf(fields)
	structValue := reflect.New(structType).Elem()

	fieldAttr := attrs.FieldAttrs["CustomField"]
	if intAttrs, ok := fieldAttr.(Attributes); ok {
		randVal := intAttrs.GetRandomValue()
		if randVal != nil {
			fieldValue := reflect.ValueOf(randVal)
			field := structValue.FieldByName("CustomField")

			if fieldValue.Type().ConvertibleTo(field.Type()) {
				field.Set(fieldValue.Convert(field.Type()))

				result := field.Interface().(CustomInt)
				if int(result) < 0 || int(result) > 50 {
					t.Errorf("Expected value in range [0, 50], got %d", result)
				}
			}
		}
	}
}

func TestStructAttributes_NilFieldValue(t *testing.T) {
	attrs := StructAttributes{
		FieldAttrs: map[string]any{
			"Field1": nilReturningAttribute{},
		},
	}

	result := attrs.GetRandomValue()
	if result == nil {
		t.Fatal("Expected non-nil struct")
	}

	structValue := reflect.ValueOf(result)
	field := structValue.FieldByName("Field1")
	if !field.IsValid() {
		t.Error("Expected valid field")
	}
	if field.Int() != 0 {
		t.Errorf("Expected zero value for field, got %v", field.Interface())
	}
}

func TestStructAttributes_GetReflectType_WithReflectType(t *testing.T) {
	attrs := StructAttributes{
		FieldAttrs: map[string]any{
			"Field1": reflect.TypeOf(int(0)),
			"Field2": reflect.TypeOf(""),
		},
	}

	reflectType := attrs.GetReflectType()
	if reflectType == nil {
		t.Fatal("Expected non-nil reflect type for struct")
	}

	if reflectType.Kind() != reflect.Struct {
		t.Errorf("Expected struct kind, got %v", reflectType.Kind())
	}

	if reflectType.NumField() != 2 {
		t.Errorf("Expected 2 fields, got %d", reflectType.NumField())
	}
}

func TestStructAttributes_GetReflectType_WithNilFieldType(t *testing.T) {
	attrs := StructAttributes{
		FieldAttrs: map[string]any{
			"Field1": nilTypeReturningAttribute{},
		},
	}

	reflectType := attrs.GetReflectType()
	if reflectType != nil {
		t.Errorf("Expected nil reflect type when field returns nil type, got %v", reflectType)
	}
}

func TestStructAttributes_GetReflectType_Mixed(t *testing.T) {
	attrs := StructAttributes{
		FieldAttrs: map[string]any{
			"Field1": IntegerAttributesImpl[int]{},
			"Field2": reflect.TypeOf(""),
		},
	}

	reflectType := attrs.GetReflectType()
	if reflectType == nil {
		t.Fatal("Expected non-nil reflect type for struct")
	}

	if reflectType.Kind() != reflect.Struct {
		t.Errorf("Expected struct kind, got %v", reflectType.Kind())
	}
}

func TestStructAttributes_UnsettableField(t *testing.T) {
	attrs := StructAttributes{
		FieldAttrs: map[string]any{
			"ExportedField": IntegerAttributesImpl[int]{Max: 10},
		},
	}

	result := attrs.GetRandomValue()
	if result == nil {
		t.Fatal("Expected non-nil struct")
	}
	structValue := reflect.ValueOf(result)
	if structValue.Kind() != reflect.Struct {
		t.Errorf("Expected struct, got %v", structValue.Kind())
	}
}

func TestStructAttributes_FieldConversion(t *testing.T) {
	attrs := StructAttributes{
		FieldAttrs: map[string]any{
			"Field1": IntegerAttributesImpl[int]{Max: 10},
			"Field2": FloatAttributesImpl[float64]{Max: 10.0},
		},
	}

	result := attrs.GetRandomValue()
	if result == nil {
		t.Fatal("Expected non-nil struct")
	}

	structValue := reflect.ValueOf(result)
	field1 := structValue.FieldByName("Field1")
	field2 := structValue.FieldByName("Field2")

	if !field1.IsValid() || !field2.IsValid() {
		t.Error("Expected valid fields")
	}
}

func TestStructAttributes_NonConvertibleField(t *testing.T) {
	attrs := StructAttributes{
		FieldAttrs: map[string]any{
			"TestField": IntegerAttributesImpl[int]{Max: 10},
		},
	}

	result := attrs.GetRandomValue()
	if result == nil {
		t.Fatal("Expected non-nil struct")
	}

	structValue := reflect.ValueOf(result)
	if structValue.Kind() != reflect.Struct {
		t.Errorf("Expected struct, got %v", structValue.Kind())
	}
}

func TestStructAttributes_GetReflectTypeNilStruct(t *testing.T) {
	attrs := StructAttributes{
		FieldAttrs: nil,
	}

	reflectType := attrs.GetReflectType()
	if reflectType != nil {
		t.Fatal("Expected nil reflect type for nil struct")
	}
}
