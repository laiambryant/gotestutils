package mtesting

import (
	"testing"

	a "github.com/laiambryant/gotestutils/mtesting/attributes"
	gen "github.com/laiambryant/gotestutils/mtesting/generation"
	"github.com/laiambryant/gotestutils/utils"
)

type MTesting struct {
	f          func(...any) []any
	iterations uint
	attributes a.MTAttributes
	t          *testing.T
}

func (mt *MTesting) WithIterations(n uint) *MTesting             { mt.iterations = n; return mt }
func (mt *MTesting) WithFunction(f func(...any) []any) *MTesting { mt.f = f; return mt }
func (mt *MTesting) WithAttributes(a a.MTAttributes) *MTesting   { mt.attributes = a; return mt }

func (mt *MTesting) GenerateInputs() ([]any, error) {
	if mt.f == nil {
		return nil, nil
	}
	inTypes, _ := utils.ExtractFArgTypes(mt.f)
	args := make([]any, len(inTypes))
	for i, t := range inTypes {
		v, err := gen.GenerateValueForTypeWithAttr(mt.attributes.GetAttributeGivenType(t))
		if err != nil {
			if mt.t != nil {
				mt.t.Logf("GenerateValueForTypeWithAttr failed for arg %d (%v): %v", i, t, err)
			}
			args[i] = v
		}
	}
	return args, nil
}
