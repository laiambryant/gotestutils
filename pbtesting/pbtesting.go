package pbtesting

import (
	"testing"

	"github.com/laiambryant/gotestutils/mtesting"
	p "github.com/laiambryant/gotestutils/pbtesting/properties/predicates"
	"github.com/laiambryant/gotestutils/utils"
)

type PBTest struct {
	t          *testing.T
	f          func(...any) []any
	predicates []p.Predicate
	iterations uint
	argAttrs   []any
}

type PBTestOut struct {
	Output     any
	Predicates []p.Predicate
	ok         bool
}

func NewPBTest(f func(...any) []any) *PBTest { return &PBTest{f: f, iterations: 1} }

func (pbt *PBTest) WithIterations(n uint) *PBTest { pbt.iterations = n; return pbt }

func (pbt *PBTest) WithPredicates(preds ...p.Predicate) *PBTest { pbt.predicates = preds; return pbt }

func (pbt *PBTest) WithArgAttributes(attrs ...any) *PBTest { pbt.argAttrs = attrs; return pbt }

func (pbt *PBTest) Run() (retOut []PBTestOut) {
	for i := uint(0); i < pbt.iterations; i++ {
		monkeyTest := (&mtesting.MTesting[int]{}).WithFunction(pbt.f)
		inputs, _ := monkeyTest.GenerateInputs()
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

func (pbt *PBTest) applyFunction(args ...any) ([]any, error) {
	if pbt.f == nil {
		return nil, nil
	}
	return pbt.f(args...), nil
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
