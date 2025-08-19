package pbtesting

import (
	"reflect"

	p "github.com/laiambryant/gotestutils/pbtesting/properties/predicates"
	"github.com/laiambryant/gotestutils/utils"
)

type PBTest struct {
	f          func(...any) []any
	predicates []p.Predicate
	iterations uint
}

type PBTestOut struct {
	Output     any
	Predicates []p.Predicate
	ok         bool
}

func NewPBTest(f func(...any) []any) *PBTest {
	return &PBTest{}
}

func (pbt PBTest) Run() (retOut []PBTestOut) {
	for i := uint(0); i < pbt.iterations; i++ {
		inputs, _ := pbt.generateInputs()
		outs, _ := pbt.applyFunction(inputs...)
		if pbt.haspredicates() {
			for _, out := range outs {
				retOut = pbt.validatePredicates(retOut, out)
			}
		}
	}
	return retOut
}

func (pbt PBTest) validatePredicates(retOut []PBTestOut, out any) []PBTestOut {
	if ok, failedpredicates := pbt.satisfyAll(out); !ok {
		retOut = append(retOut, PBTestOut{
			Output:     out,
			Predicates: failedpredicates,
			ok:         false,
		})
	} else {
		retOut = append(retOut, PBTestOut{
			Output:     out,
			Predicates: nil,
			ok:         true,
		})
	}
	return retOut
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

func (pbt PBTest) applyFunction(args ...any) ([]any, error) {
	if pbt.f == nil {
		return nil, nil
	}
	return pbt.f(args...), nil
}

func (pbt PBTest) generateInputs() ([]any, error) {
	if pbt.f == nil {
		return nil, nil
	}
	inTypes, _ := extractFArgTypes(pbt.f)
	args := make([]any, len(inTypes))
	for i, t := range inTypes {
		v := reflect.New(t).Elem()
		getRandomValue(v)
		args[i] = v.Interface()
	}
	return args, nil
}

func (pbt PBTest) satisfyAll(val any) (ok bool, failedpredicates []p.Predicate) {
	if len(pbt.predicates) == 0 {
		return true, nil
	}
	for _, predicate := range pbt.predicates {
		if !predicate.Verify(val) {
			failedpredicates = append(failedpredicates, predicate)
		}
	}
	if len(failedpredicates) > 0 {
		return false, failedpredicates
	}
	return true, nil
}

func (pbt PBTest) haspredicates() bool {
	return pbt.predicates != nil
}

func FilterPBTTestOut(in []PBTestOut) []PBTestOut {
	return utils.Filter(in, func(po PBTestOut) bool {
		return !po.ok
	})
}
