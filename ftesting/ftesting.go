package ftesting

import (
	"fmt"
	"reflect"
	"testing"

	a "github.com/laiambryant/gotestutils/ftesting/attributes"
)

type FTesting struct {
	f          any // Must be a function type
	iterations uint
	attributes a.AttributesStruct
	t          *testing.T
}

func (mt *FTesting) WithIterations(n uint) *FTesting { mt.iterations = n; return mt }
func (mt *FTesting) WithFunction(f any) *FTesting {
	mt.f = f
	return mt
}
func (mt *FTesting) WithAttributes(a a.AttributesStruct) *FTesting {
	mt.attributes = a
	return mt
}
func (mt *FTesting) GenerateInputs() ([]any, error) {
	if mt.f == nil {
		return nil, nil
	}
	if reflect.TypeOf(mt.f).Kind() != reflect.Func {
		return nil, fmt.Errorf("f is not a function: %v", mt.f)
	}
	fType := reflect.TypeOf(mt.f)
	args := make([]any, fType.NumIn())
	for i := 0; i < fType.NumIn(); i++ {
		argType := fType.In(i)
		v := mt.attributes.GetAttributeGivenType(argType).GetRandomValue()
		args[i] = v
	}
	return args, nil
}

func (mt *FTesting) ApplyFunction() (bool, error) {
	if mt.f == nil {
		return false, fmt.Errorf("function is nil")
	}

	inputs, err := mt.GenerateInputs()
	if err != nil {
		return false, fmt.Errorf("failed to generate inputs: %w", err)
	}

	args := make([]reflect.Value, len(inputs))
	for i, input := range inputs {
		args[i] = reflect.ValueOf(input)
	}

	fValue := reflect.ValueOf(mt.f)
	_ = fValue.Call(args)

	// if the function does not panic we assume the test to be succesful
	return true, nil
}

func (mt *FTesting) Verify() {
	if mt.t == nil {
		return
	}
	ok, err := mt.ApplyFunction()
	if err != nil {
		mt.t.Errorf("Test Failed with error: [%s]", err.Error())
	}
	if !ok {
		mt.t.Error("Test Failed")
	}
}
