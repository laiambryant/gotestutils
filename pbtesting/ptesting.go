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
