package pbtesting

import (
	"reflect"
	"testing"

	"github.com/laiambryant/gotestutils/ftesting"
	p "github.com/laiambryant/gotestutils/pbtesting/properties/predicates"
	"github.com/laiambryant/gotestutils/utils"
)

type PBTest struct {
	t          *testing.T
	f          any
	predicates []p.Predicate
	iterations uint
	argAttrs   []any
}

type PBTestOut struct {
	Output     any
	Predicates []p.Predicate
	ok         bool
}

type returnTypes interface {
	any | []any
}

func NewPBTest(f any) *PBTest { return &PBTest{f: f, iterations: 1} }

func (pbt *PBTest) WithIterations(n uint) *PBTest { pbt.iterations = n; return pbt }

func (pbt *PBTest) WithPredicates(preds ...p.Predicate) *PBTest { pbt.predicates = preds; return pbt }

func (pbt *PBTest) WithArgAttributes(attrs ...any) *PBTest { pbt.argAttrs = attrs; return pbt }

func (pbt *PBTest) WithT(t *testing.T) *PBTest { pbt.t = t; return pbt }

func (pbt *PBTest) WithF(f any) *PBTest { pbt.f = f; return pbt }

func (pbt *PBTest) Run() (retOut []PBTestOut, err error) {
	if pbt.f == nil {
		return []PBTestOut{}, nil
	}
	for i := uint(0); i < pbt.iterations; i++ {
		fuzzTest := (&ftesting.FTesting{}).WithFunction(pbt.f)
		inputs, err := fuzzTest.GenerateInputs()
		if err != nil {
			return nil, err
		}
		outs, _ := pbt.applyFunction(inputs...)
		if pbt.haspredicates() {
			switch ret := outs.(type) {
			case []any:
				for _, out := range ret {
					retOut = pbt.validatePredicates(retOut, out)
				}
			case any:
				retOut = pbt.validatePredicates(retOut, ret)
			}
		}
	}
	return retOut, nil
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

func (pbt *PBTest) applyFunction(args ...any) (returnTypes, error) {
	if pbt.f == nil {
		return nil, nil
	}
	switch fn := pbt.f.(type) {
	case func(any) any:
		return fn(args), nil
	case func(...any) any:
		return fn(args...), nil
	}
	fValue := reflect.ValueOf(pbt.f)
	fType := fValue.Type()

	if fType.Kind() != reflect.Func {
		return nil, &InvalidFunctionProvidedError{pbt.f}
	}
	reflectArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		argValue := reflect.ValueOf(arg)
		expectedType := fType.In(i)
		if argValue.Type() != expectedType {
			if argValue.Type().ConvertibleTo(expectedType) {
				argValue = argValue.Convert(expectedType)
			} else {
				return nil, &InvalidFunctionProvidedError{pbt.f}
			}
		}
		reflectArgs[i] = argValue
	}
	results := fValue.Call(reflectArgs)
	if len(results) == 0 {
		return nil, nil
	} else if len(results) == 1 {
		return results[0].Interface(), nil
	} else {
		retSlice := make([]any, len(results))
		for i, result := range results {
			retSlice[i] = result.Interface()
		}
		return retSlice, nil
	}
}

func (pbt *PBTest) satisfyAll(val any) (ok bool, failedpredicates []p.Predicate) {
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

func (pbt *PBTest) haspredicates() bool {
	return pbt.predicates != nil
}

func FilterPBTTestOut(in []PBTestOut) []PBTestOut {
	return utils.Filter(in, func(po PBTestOut) bool {
		return !po.ok
	})
}
