package mtesting

import (
	"fmt"
	"reflect"
	"testing"

	a "github.com/laiambryant/gotestutils/mtesting/attributes"
	"github.com/laiambryant/gotestutils/utils"
)

type MTesting struct {
	f          any
	iterations uint
	attributes a.AttributesStruct
	t          *testing.T
}

func (mt *MTesting) WithIterations(n uint) *MTesting { mt.iterations = n; return mt }
func (mt *MTesting) WithFunction(f any) *MTesting    { mt.f = f; return mt }
func (mt *MTesting) WithAttributes(a a.AttributesStruct) *MTesting {
	mt.attributes = a
	return mt
}
func (mt *MTesting) GenerateInputs() ([]any, error) {
	if mt.f == nil {
		return nil, nil
	}
	if reflect.TypeOf(mt.f).Kind() != reflect.Func {
		return nil, fmt.Errorf("f is not a function: %v", mt.f)
	}
	inTypes, _ := utils.ExtractFArgTypes(mt.f)
	args := make([]any, len(inTypes))
	for i, t := range inTypes {
		v := mt.attributes.GetAttributeGivenType(t).GetRandomValue()
		args[i] = v
	}
	return args, nil
}

func (mt *MTesting) ApplyFunction() (bool, error) {
	return true, nil
}

func (mt *MTesting) Verify() {
	if mt.t == nil {
		return
	}
	ok, err := mt.ApplyFunction()
	if err != nil {
		mt.t.Errorf("Test Failed with error: [%s]", err.Error())
	}
	if !ok {
		mt.t.Error("Test Failed:")
	}
}
