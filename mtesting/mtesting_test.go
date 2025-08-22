package mtesting

import (
	"testing"

	"github.com/laiambryant/gotestutils/mtesting/attributes"
)

var (
	sumFunc = func(a int, b int) int {
		return a + b
	}
	mta = attributes.MTAttributes{
		IA: attributes.IntegerAttributes{
			Min: 10,
			Max: 100,
		},
	}
)

func TestMTesting(t *testing.T) {
	mt := MTesting{}
	mt = *mt.WithFunction(sumFunc).WithIterations(1000).WithAttributes(mta)
	in, err := mt.GenerateInputs()
	t.Logf("inputs: %v, error: %s", in, err)
}

func TestMTestingEmptyF(t *testing.T) {
	mt := MTesting{}
	mt = *mt.WithFunction(nil)
	mt.GenerateInputs()
}

func TestMTestingFNotFunc(t *testing.T) {
	mt := MTesting{}
	mt = *mt.WithFunction(1)
	if _, err := mt.GenerateInputs(); err == nil {
		t.Errorf("should return an error")
	}
}
