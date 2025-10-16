package ftesting

import (
	"reflect"
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

func TestFTestingGenerateInputsWithNilAttributes(t *testing.T) {
	mt := FTesting{}
	mt = *mt.WithFunction(sumFunc)
	inputs, err := mt.GenerateInputs()
	if err != nil {
		t.Errorf("GenerateInputs should not fail with default attributes: %v", err)
	}
	if len(inputs) != 2 {
		t.Errorf("Expected 2 inputs for sumFunc, got %d", len(inputs))
	}
	if mt.attributes == nil {
		t.Error("Default attributes should have been assigned")
	}
}

func TestFTestingGenerateInputsWithUnsupportedType(t *testing.T) {
	funcWithUnsupportedType := func(ch chan int) int {
		return <-ch
	}
	mt := FTesting{}
	mt = *mt.WithFunction(funcWithUnsupportedType).WithAttributes(mta)
	inputs, err := mt.GenerateInputs()
	if err == nil {
		t.Error("GenerateInputs should return error for unsupported channel type")
	}
	if inputs != nil {
		t.Error("Inputs should be nil when error occurs")
	}
	if err != nil && len(err.Error()) == 0 {
		t.Error("Error message should not be empty")
	}
}

func TestNoFunctionProvidedError(t *testing.T) {
	err := NoFunctionProvidedError{}
	expectedMessage := "no function was provided to ftesting suite"
	if err.Error() != expectedMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedMessage, err.Error())
	}
	var _ error = err
}

func TestInputsGenerationError(t *testing.T) {
	underlyingErr := &NoFunctionProvidedError{}
	err := InputsGenerationError{err: underlyingErr}
	expectedMessage := "error in input generation: no function was provided to ftesting suite"
	actualMessage := err.Error()
	if actualMessage != expectedMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedMessage, actualMessage)
	}
	var _ error = err
	underlyingErr2 := &NotAFunctionError{k: reflect.String}
	err2 := InputsGenerationError{err: underlyingErr2}
	expectedMessage2 := "error in input generation: f is not a function: string"
	actualMessage2 := err2.Error()
	if actualMessage2 != expectedMessage2 {
		t.Errorf("Expected error message '%s', got '%s'", expectedMessage2, actualMessage2)
	}
}
