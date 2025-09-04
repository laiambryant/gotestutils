package mtesting

import (
	"testing"

	"github.com/laiambryant/gotestutils/mtesting/attributes"
)

var (
	sumFunc = func(a int, b int) int {
		return a + b
	}
	mta = attributes.MTAttributes[int]{
		IA: attributes.IntegerAttributes[int]{
			Min: 10,
			Max: 100,
		},
	}
)

func TestMTesting(t *testing.T) {
	mt := MTesting[int]{}
	mt = *mt.WithFunction(sumFunc).WithIterations(1000).WithAttributes(mta)
	in, err := mt.GenerateInputs()
	t.Logf("inputs: %v, error: %s", in, err)
}

func TestMTestingEmptyF(t *testing.T) {
	mt := MTesting[int]{}
	mt = *mt.WithFunction(nil)
	mt.GenerateInputs()
}

func TestMTestingFNotFunc(t *testing.T) {
	mt := MTesting[int]{}
	mt = *mt.WithFunction(1)
	if _, err := mt.GenerateInputs(); err == nil {
		t.Errorf("should return an error")
	}
}
