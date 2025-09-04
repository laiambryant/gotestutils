package mtesting

import (
	"fmt"
	"reflect"
	"testing"

	a "github.com/laiambryant/gotestutils/mtesting/attributes"
	gen "github.com/laiambryant/gotestutils/mtesting/generation"
	"github.com/laiambryant/gotestutils/utils"
)

type MTesting[T any] struct {
	f          any
	iterations uint
	attributes a.MTAttributes[T]
	t          *testing.T
}

func (mt *MTesting[T]) WithIterations(n uint) *MTesting[T]              { mt.iterations = n; return mt }
func (mt *MTesting[T]) WithFunction(f any) *MTesting[T]                 { mt.f = f; return mt }
func (mt *MTesting[T]) WithAttributes(a a.MTAttributes[T]) *MTesting[T] { mt.attributes = a; return mt }
func (mt *MTesting[T]) WithT(t *testing.T) *MTesting[T]                 { mt.t = t; return mt }
func (mt *MTesting[T]) GenerateInputs() ([]any, error) {
	if mt.f == nil {
		return nil, nil
	}
	if reflect.TypeOf(mt.f).Kind() != reflect.Func {
		return nil, fmt.Errorf("f is not a function: %v", mt.f)
	}
	inTypes, _ := utils.ExtractFArgTypes(mt.f)
	args := make([]any, len(inTypes))
	for i, t := range inTypes {
		v, err := gen.GenerateValueForTypeWithAttr(mt.attributes.GetAttributeGivenType(t), t)
		if err != nil && mt.t != nil {
			mt.t.Logf("GenerateValueForTypeWithAttr failed for arg %d (%v): %v", i, t, err)
		} else {
			args[i] = v
		}
	}
	return args, nil
}

func (mt *MTesting[T]) ApplyFunction() (bool, error) {
	return true, nil
}
