package ftesting

import (
	"testing"

	"github.com/laiambryant/gotestutils/ftesting/attributes"
)

var (
	sumFunc = func(a int, b int) int {
		return a + b
	}
	mta = attributes.FTAttributes{
		IntegerAttr: attributes.IntegerAttributesImpl[int]{
			Min: 10,
			Max: 100,
		},
	}
)

var mockT = &testing.T{}

func TestFTesting(t *testing.T) {
	mt := FTesting{}
	mt = *mt.WithFunction(sumFunc).WithIterations(1000).WithAttributes(mta)
	in, err := mt.GenerateInputs()
	t.Logf("inputs: %v, error: %s", in, err)
}

func TestFTestingEmptyF(t *testing.T) {
	mt := FTesting{}
	mt = *mt.WithFunction(nil)
	mt.GenerateInputs()
}

func TestFTestingFNotFunc(t *testing.T) {
	mt := FTesting{}
	mt = *mt.WithFunction(1)
	if _, err := mt.GenerateInputs(); err == nil {
		t.Errorf("should return an error")
	}
}

func TestFTestingVerify(t *testing.T) {
	mt := FTesting{}
	mt = *mt.WithFunction(sumFunc).WithIterations(10).WithAttributes(mta)
	mt.Verify()
}
func TestFTestingVerifyWithTesting(t *testing.T) {
	mt := FTesting{t: mockT}
	mt = *mt.WithFunction(sumFunc).WithIterations(10).WithAttributes(mta)
	mt.t = mockT
	mt.Verify()
	mt2 := FTesting{}
	mt2 = *mt2.WithFunction(nil).WithIterations(10).WithAttributes(mta)
	mt2.t = mockT
	mt2.Verify()
}

func TestFTestingVerifyWithNonFunction(t *testing.T) {
	mt := FTesting{}
	mt = *mt.WithFunction("not a function").WithIterations(10).WithAttributes(mta)
	mt.t = mockT
	mt.Verify()
}

func TestFTestingVerifyWithPanicFunction(t *testing.T) {
	panicFunc := func(a int, b int) int {
		panic("test panic")
	}
	mt := FTesting{}
	mt = *mt.WithFunction(panicFunc).WithIterations(10).WithAttributes(mta)
	mt.t = t
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Expected panic caught during Verify: %v", r)
		} else {
			t.Error("Expected a panic but none occurred")
		}
	}()
	mt.Verify()
}
