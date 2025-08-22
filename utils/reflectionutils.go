package utils

import "reflect"

func ExtractFArgTypes(f interface{}) (inputTypes []reflect.Type, ouputTypes []reflect.Type) {
	types := reflect.TypeOf(f)
	for i := 0; i < types.NumIn(); i++ {
		inputTypes = append(inputTypes, types.In(i))
	}
	for i := 0; i < types.NumOut(); i++ {
		ouputTypes = append(ouputTypes, types.Out(i))
	}
	return
}
