package mtesting

import (
	"fmt"
	"reflect"
	"testing"

	a "github.com/laiambryant/gotestutils/mtesting/attributes"
	gen "github.com/laiambryant/gotestutils/mtesting/generation"
	"github.com/laiambryant/gotestutils/utils"
)

type MTesting struct {
	f          any
	iterations uint
	attributes a.MTAttributes
	t          *testing.T
}

func (mt *MTesting) WithIterations(n uint) *MTesting           { mt.iterations = n; return mt }
func (mt *MTesting) WithFunction(f any) *MTesting              { mt.f = f; return mt }
func (mt *MTesting) WithAttributes(a a.MTAttributes) *MTesting { mt.attributes = a; return mt }
func (mt *MTesting) WithT(t *testing.T) *MTesting              { mt.t = t; return mt }
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
		v, err := gen.GenerateValueForTypeWithAttr(mt.attributes.GetAttributeGivenType(t))
		if err != nil && mt.t != nil {
			mt.t.Logf("GenerateValueForTypeWithAttr failed for arg %d (%v): %v", i, t, err)
		} else {
			args[i] = v
		}
	}
	return args, nil
}
