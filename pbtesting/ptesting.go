package pbtesting

import (
	"reflect"

	s "github.com/laiambryant/gotestutils/pbtesting/properties"
)

type PBTest struct {
	Func                       func(...any) []any
	PreconditionValidationFunc func(...any) bool
	Predicates                 []s.Predicate
	iterations                 uint
}

func (pbt PBTest) Run() {
	for i := uint(0); i < pbt.iterations; i++ {
		inputs, _ := pbt.generateInputs()
		if !pbt.validatePreconditions(inputs...) {
			continue
		}
		outs, _ := pbt.applyFunction(inputs...)
		// Optional: apply predicates on outputs here if desired in Run
		_ = outs
	}
}

func extractFArgTypes(f any) (inputTypes []reflect.Type, ouputTypes []reflect.Type) {
	types := reflect.TypeOf(f)
	for i := 0; i < types.NumIn(); i++ {
		inputTypes = append(inputTypes, types.In(i))
	}
	for i := 0; i < types.NumOut(); i++ {
		ouputTypes = append(ouputTypes, types.Out(i))
	}
	return
}

func createInstances(types []reflect.Type, isZero bool) []any {
	instances := make([]any, len(types))
	for i, t := range types {
		newValue := reflect.New(t).Elem()
		if isZero {
			instances[i] = newValue.Interface()
			continue
		}
		getRandomValue(newValue)
		instances[i] = newValue.Interface()
	}
	return instances
}

func (pbt PBTest) validatePreconditions(args ...any) bool {
	if pbt.PreconditionValidationFunc == nil {
		return true
	}
	return pbt.PreconditionValidationFunc(args...)
}

func (pbt PBTest) applyFunction(args ...any) ([]any, error) {
	if pbt.Func == nil {
		return nil, nil
	}
	return pbt.Func(args...), nil
}

func (pbt PBTest) generateInputs() ([]any, error) {
	if pbt.Func == nil {
		return nil, nil
	}
	inTypes, _ := extractFArgTypes(pbt.Func)
	args := make([]any, len(inTypes))
	for i, t := range inTypes {
		v := reflect.New(t).Elem()
		getRandomValue(v)
		args[i] = v.Interface()
	}
	return args, nil
}

func (pbt PBTest) satisfyAll(val any) bool {
	if len(pbt.Predicates) == 0 {
		return true
	}
	for _, prop := range pbt.Predicates {
		if !prop.Verify(val) {
			return false
		}
	}
	return true
}
