package pbtesting

import (
	"reflect"

	s "github.com/laiambryant/gotestutils/pbtesting/properties"
)

type PBTest struct {
	Func                       func(...any) []any
	PreconditionValidationFunc func(any) bool
	Properties                 []s.Property
	iterations                 uint
}

func (pbt PBTest) Run() {
	for range pbt.iterations {
		continue
	}
}

func extractFArgTypes(f any) (inputTypes []reflect.Type, ouputTypes []reflect.Type) {
	types := reflect.TypeOf(f)
	for i := range types.NumIn() {
		inputTypes = append(inputTypes, types.In(i))
	}
	for i := range types.NumOut() {
		ouputTypes = append(ouputTypes, types.Out(i))
	}
	return
}

func createInstances(types []reflect.Type, isZero bool) []any {
	instances := make([]any, len(types))
	for i, t := range types {
		if isZero {
			newValue := reflect.New(t).Elem()
			instances[i] = newValue.Interface()
		}
		// If not zero fill with random values
	}
	return instances
}

func getRandomValue(v reflect.Value, properties []s.Property) {
	switch v.Kind() {
	case reflect.Bool:
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		break
	case reflect.Float32, reflect.Float64:
		break
	case reflect.Complex64, reflect.Complex128:
		break
	case reflect.Array:
		break
	case reflect.Chan:
		break
	case reflect.Func:
		break
	case reflect.Interface:
		break
	case reflect.Map:
		break
	case reflect.Pointer:
		break
	case reflect.Slice:
		break
	case reflect.String:
		break
	case reflect.Struct:
		break
	case reflect.UnsafePointer:
		break
	default:
		break
	}
}
func (pbt PBTest) validatePreconditions(args ...any) bool {
	return pbt.PreconditionValidationFunc(args)
}

func (pbt PBTest) applyFunction(args ...any) (any, error) {
	return pbt.Func(args), nil
}

func (pbt PBTest) generateInputs() ([]any, error) {
	var inputs []any
	for _, property := range pbt.Properties {
		value := property.Verify(nil)
		inputs = append(inputs, value)
		return inputs, InvalidPropertyError{property}
	}
	return inputs, nil
}
