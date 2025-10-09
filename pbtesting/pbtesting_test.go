package pbtesting

import (
	"fmt"
	"reflect"
	"testing"

	p "github.com/laiambryant/gotestutils/pbtesting/properties/predicates"
)

var f1 func(a int) int = func(a int) int {
	return a
}

var f2 func(a int, b int) int = func(a int, b int) int {
	return a + b
}

var funcAnyToAny func(any) any = func(input any) any {
	if val, ok := input.(int); ok {
		return val * 2
	}
	return input
}

var funcVariadicAnyToAny func(...any) any = func(args ...any) any {
	if len(args) == 0 {
		return nil
	}
	if len(args) == 1 {
		return args[0]
	}
	sum := 0
	for _, arg := range args {
		if val, ok := arg.(int); ok {
			sum += val
		}
	}
	return sum
}

type mockPredicate struct {
	shouldPass bool
	name       string
}

func (m mockPredicate) Verify(val any) bool {
	return m.shouldPass
}

func (m mockPredicate) String() string {
	return m.name
}

func TestNewPBTest(t *testing.T) {
	pbt := NewPBTest(f1)
	if pbt.f == nil {
		t.Error("Expected function to be set")
	}
	if pbt.iterations != 1 {
		t.Errorf("Expected default iterations to be 1, got %d", pbt.iterations)
	}
}

func TestWithIterations(t *testing.T) {
	pbt := NewPBTest(f1)
	result := pbt.WithIterations(10)
	if result.iterations != 10 {
		t.Errorf("Expected iterations to be 10, got %d", result.iterations)
	}
	if result != pbt {
		t.Error("Expected WithIterations to return the same instance for chaining")
	}
}

func TestWithPredicates(t *testing.T) {
	pbt := NewPBTest(f1)
	pred1 := mockPredicate{shouldPass: true, name: "pred1"}
	pred2 := mockPredicate{shouldPass: false, name: "pred2"}
	result := pbt.WithPredicates(pred1, pred2)
	if len(result.predicates) != 2 {
		t.Errorf("Expected 2 predicates, got %d", len(result.predicates))
	}
	if result != pbt {
		t.Error("Expected WithPredicates to return the same instance for chaining")
	}
}

func TestWithArgAttributes(t *testing.T) {
	pbt := NewPBTest(f1)
	attrs := []any{"attr1", 42, true}
	result := pbt.WithArgAttributes(attrs...)
	if len(result.argAttrs) != 3 {
		t.Errorf("Expected 3 attributes, got %d", len(result.argAttrs))
	}
	if result != pbt {
		t.Error("Expected WithArgAttributes to return the same instance for chaining")
	}
}

func TestWithT(t *testing.T) {
	pbt := NewPBTest(f1)
	result := pbt.WithT(t)
	if result.t != t {
		t.Error("Expected testing.T to be set")
	}
	if result != pbt {
		t.Error("Expected WithT to return the same instance for chaining")
	}
}

func TestWithF(t *testing.T) {
	pbt := NewPBTest(f1)
	result := pbt.WithF(f2)
	if result.f == nil {
		t.Error("Expected function to be set")
	}
	if result != pbt {
		t.Error("Expected WithF to return the same instance for chaining")
	}
}

func TestApplyFunction_AnyToAny(t *testing.T) {
	pbt := NewPBTest(funcAnyToAny)
	result, err := pbt.applyFunction(5)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expected := []any{5}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected result to be %v, got %v", expected, result)
	}
}

func TestApplyFunction_VariadicAnyToAny(t *testing.T) {
	pbt := NewPBTest(funcVariadicAnyToAny)
	result, err := pbt.applyFunction(1, 2, 3)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 6 {
		t.Errorf("Expected result to be 6, got %v", result)
	}
}

func TestApplyFunction_NilFunction(t *testing.T) {
	pbt := NewPBTest(nil)
	result, err := pbt.applyFunction(1, 2, 3)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
}

func TestApplyFunction_InvalidFunction(t *testing.T) {
	pbt := NewPBTest("not a function")
	result, err := pbt.applyFunction(1, 2)
	if err == nil {
		t.Error("Expected error for invalid function type")
	}
	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
	if _, ok := err.(*InvalidFunctionProvidedError); !ok {
		t.Errorf("Expected InvalidFunctionProvidedError, got %T", err)
	}
}

func TestApplyFunction_ConcreteFunction(t *testing.T) {
	validFunc := func(a int, b int) int { return a + b }
	pbt := NewPBTest(validFunc)
	result, err := pbt.applyFunction(1, 2)
	if err != nil {
		t.Errorf("Expected no error for valid function, got %v", err)
	}
	if result != 3 {
		t.Errorf("Expected result 3, got %v", result)
	}
}

func TestApplyFunction_TypeConversion(t *testing.T) {
	funcInt64 := func(a int64) int64 { return a * 2 }
	pbt := NewPBTest(funcInt64)
	var input int32 = 42
	result, err := pbt.applyFunction(input)
	if err != nil {
		t.Errorf("Expected no error for convertible types, got %v", err)
	}
	expected := int64(84)
	if result != expected {
		t.Errorf("Expected result %v, got %v", expected, result)
	}
}

func TestApplyFunction_TypeConversionFloat(t *testing.T) {
	funcFloat64 := func(a float64) float64 { return a + 0.5 }
	pbt := NewPBTest(funcFloat64)
	var input float32 = 3.5
	result, err := pbt.applyFunction(input)
	if err != nil {
		t.Errorf("Expected no error for convertible types, got %v", err)
	}
	expected := float64(4.0)
	if result != expected {
		t.Errorf("Expected result %v, got %v", expected, result)
	}
}

func TestApplyFunction_NonConvertibleTypes(t *testing.T) {
	funcInt := func(a int) int { return a * 2 }
	pbt := NewPBTest(funcInt)
	result, err := pbt.applyFunction("not a number")
	if err == nil {
		t.Error("Expected error for non-convertible types")
	}
	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
	if _, ok := err.(*InvalidFunctionProvidedError); !ok {
		t.Errorf("Expected InvalidFunctionProvidedError, got %T", err)
	}
}

func TestApplyFunction_MultipleReturnValues(t *testing.T) {
	funcMultiReturn := func(a, b int) (int, int, int) {
		return a + b, a - b, a * b
	}
	pbt := NewPBTest(funcMultiReturn)
	result, err := pbt.applyFunction(5, 3)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	resultSlice, ok := result.([]any)
	if !ok {
		t.Errorf("Expected []any result, got %T", result)
	}
	if len(resultSlice) != 3 {
		t.Errorf("Expected 3 return values, got %d", len(resultSlice))
	}
	if resultSlice[0] != 8 {
		t.Errorf("Expected first return value 8, got %v", resultSlice[0])
	}
	if resultSlice[1] != 2 {
		t.Errorf("Expected second return value 2, got %v", resultSlice[1])
	}
	if resultSlice[2] != 15 {
		t.Errorf("Expected third return value 15, got %v", resultSlice[2])
	}
}

func TestApplyFunction_ZeroReturnValues(t *testing.T) {
	called := false
	funcNoReturn := func(a int) {
		called = true
	}
	pbt := NewPBTest(funcNoReturn)
	result, err := pbt.applyFunction(42)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil result for function with no return, got %v", result)
	}
	if !called {
		t.Error("Expected function to be called")
	}
}

func TestApplyFunction_StructTypeConversion(t *testing.T) {
	type CustomInt int
	funcCustom := func(a CustomInt) CustomInt { return a * 2 }
	pbt := NewPBTest(funcCustom)
	result, err := pbt.applyFunction(21)
	if err != nil {
		t.Errorf("Expected no error for convertible custom types, got %v", err)
	}
	expected := CustomInt(42)
	if result != expected {
		t.Errorf("Expected result %v, got %v", expected, result)
	}
}

func TestSatisfyAll_NoPredicates(t *testing.T) {
	pbt := NewPBTest(f1)
	ok, failed := pbt.satisfyAll(42)
	if !ok {
		t.Error("Expected satisfyAll to return true when no predicates")
	}
	if failed != nil {
		t.Error("Expected no failed predicates")
	}
}

func TestSatisfyAll_PassingPredicates(t *testing.T) {
	pred1 := mockPredicate{shouldPass: true, name: "pred1"}
	pred2 := mockPredicate{shouldPass: true, name: "pred2"}
	pbt := NewPBTest(f1).WithPredicates(pred1, pred2)
	ok, failed := pbt.satisfyAll(42)
	if !ok {
		t.Error("Expected satisfyAll to return true when all predicates pass")
	}
	if len(failed) != 0 {
		t.Errorf("Expected no failed predicates, got %d", len(failed))
	}
}

func TestSatisfyAll_FailingPredicates(t *testing.T) {
	pred1 := mockPredicate{shouldPass: true, name: "pred1"}
	pred2 := mockPredicate{shouldPass: false, name: "pred2"}
	pred3 := mockPredicate{shouldPass: false, name: "pred3"}
	pbt := NewPBTest(f1).WithPredicates(pred1, pred2, pred3)
	ok, failed := pbt.satisfyAll(42)
	if ok {
		t.Error("Expected satisfyAll to return false when predicates fail")
	}
	if len(failed) != 2 {
		t.Errorf("Expected 2 failed predicates, got %d", len(failed))
	}
}

func TestHasPredicates(t *testing.T) {
	pbt1 := NewPBTest(f1)
	if pbt1.haspredicates() {
		t.Error("Expected haspredicates to return false when no predicates set")
	}
	pred := mockPredicate{shouldPass: true, name: "pred"}
	pbt2 := NewPBTest(f1).WithPredicates(pred)
	if !pbt2.haspredicates() {
		t.Error("Expected haspredicates to return true when predicates are set")
	}
}

func TestValidatePredicates_Passing(t *testing.T) {
	pred := mockPredicate{shouldPass: true, name: "pred"}
	pbt := NewPBTest(f1).WithPredicates(pred)
	var retOut []PBTestOut
	result := pbt.validatePredicates(retOut, 42)
	if len(result) != 1 {
		t.Errorf("Expected 1 result, got %d", len(result))
	}
	if !result[0].ok {
		t.Error("Expected result to be ok")
	}
	if result[0].Predicates != nil {
		t.Error("Expected nil predicates for passing case")
	}
	if result[0].Output != 42 {
		t.Errorf("Expected output to be 42, got %v", result[0].Output)
	}
}

func TestValidatePredicates_Failing(t *testing.T) {
	pred := mockPredicate{shouldPass: false, name: "pred"}
	pbt := NewPBTest(f1).WithPredicates(pred)
	var retOut []PBTestOut
	result := pbt.validatePredicates(retOut, 42)
	if len(result) != 1 {
		t.Errorf("Expected 1 result, got %d", len(result))
	}
	if result[0].ok {
		t.Error("Expected result to not be ok")
	}
	if len(result[0].Predicates) != 1 {
		t.Errorf("Expected 1 failed predicate, got %d", len(result[0].Predicates))
	}
	if result[0].Output != 42 {
		t.Errorf("Expected output to be 42, got %v", result[0].Output)
	}
}

func TestRun_ArrayOutput_WithPredicates(t *testing.T) {
	testFunc := func(args ...any) any {
		return []any{1, 2, 3}
	}
	pred := mockPredicate{shouldPass: true, name: "pred"}
	pbt := NewPBTest(testFunc).WithPredicates(pred)
	var retOut []PBTestOut
	arrayOutput := []any{1, 2, 3}
	for _, out := range arrayOutput {
		retOut = pbt.validatePredicates(retOut, out)
	}
	if len(retOut) != 3 {
		t.Errorf("Expected 3 results, got %d", len(retOut))
	}
	for i, result := range retOut {
		if !result.ok {
			t.Errorf("Expected result %d to be ok", i)
		}
	}
}

func TestRun_SingleOutput_WithPredicates(t *testing.T) {
	pred := mockPredicate{shouldPass: false, name: "failing_pred"}
	pbt := NewPBTest(funcVariadicAnyToAny).WithPredicates(pred)
	var retOut []PBTestOut
	singleOutput := 42
	retOut = pbt.validatePredicates(retOut, singleOutput)
	if len(retOut) != 1 {
		t.Errorf("Expected 1 result, got %d", len(retOut))
	}
	if retOut[0].ok {
		t.Error("Expected result to not be ok")
	}
}

func TestRun_NoPredicates(t *testing.T) {
	pbt := NewPBTest(funcVariadicAnyToAny).WithT(t).WithIterations(2)
	if pbt.haspredicates() {
		t.Error("Expected haspredicates to return false when no predicates")
	}
}

func TestRun_MultipleIterations(t *testing.T) {
	pbt := NewPBTest(funcVariadicAnyToAny).WithIterations(3)
	if pbt.iterations != 3 {
		t.Errorf("Expected iterations to be 3, got %d", pbt.iterations)
	}
	pbt2 := NewPBTest(funcVariadicAnyToAny).WithIterations(5)
	if pbt2.iterations != 5 {
		t.Errorf("Expected iterations to be 5, got %d", pbt2.iterations)
	}
}

func TestFilterPBTTestOut(t *testing.T) {
	testData := []PBTestOut{
		{Output: 1, ok: true, Predicates: nil},
		{Output: 2, ok: false, Predicates: []p.Predicate{mockPredicate{shouldPass: false, name: "pred"}}},
		{Output: 3, ok: true, Predicates: nil},
		{Output: 4, ok: false, Predicates: []p.Predicate{mockPredicate{shouldPass: false, name: "pred"}}},
	}
	filtered := FilterPBTTestOut(testData)
	if len(filtered) != 2 {
		t.Errorf("Expected 2 filtered results, got %d", len(filtered))
	}
	for _, result := range filtered {
		if result.ok {
			t.Error("Expected all filtered results to have ok: false")
		}
	}
}

func TestMethodChaining(t *testing.T) {
	pred := mockPredicate{shouldPass: true, name: "pred"}
	pbt := NewPBTest(funcVariadicAnyToAny).
		WithT(t).
		WithIterations(5).
		WithPredicates(pred).
		WithArgAttributes("attr1", 42).
		WithF(funcAnyToAny)
	if pbt.t != t {
		t.Error("Expected testing.T to be set")
	}
	if pbt.iterations != 5 {
		t.Error("Expected iterations to be 5")
	}
	if len(pbt.predicates) != 1 {
		t.Error("Expected 1 predicate")
	}
	if len(pbt.argAttrs) != 2 {
		t.Error("Expected 2 attributes")
	}
	expectedFunc := reflect.ValueOf(funcAnyToAny)
	actualFunc := reflect.ValueOf(pbt.f)
	if expectedFunc.Pointer() != actualFunc.Pointer() {
		t.Error("Expected function to be funcAnyToAny")
	}
}

func TestRun_SwitchStatementCoverage(t *testing.T) {
	pred := mockPredicate{shouldPass: true, name: "pred"}
	pbt := &PBTest{
		f:          funcVariadicAnyToAny,
		predicates: []p.Predicate{pred},
		iterations: 1,
	}
	var retOut1 []PBTestOut
	arrayOut := []any{1, 2, 3}
	for _, out := range arrayOut {
		retOut1 = pbt.validatePredicates(retOut1, out)
	}
	if len(retOut1) != 3 {
		t.Errorf("Expected 3 results for array case, got %d", len(retOut1))
	}
	var retOut2 []PBTestOut
	singleOut := 42
	retOut2 = pbt.validatePredicates(retOut2, singleOut)
	if len(retOut2) != 1 {
		t.Errorf("Expected 1 result for single case, got %d", len(retOut2))
	}
}

func TestPBTestOut(t *testing.T) {
	pred := mockPredicate{shouldPass: false, name: "test_pred"}
	out1 := PBTestOut{
		Output:     "test_output",
		Predicates: []p.Predicate{pred},
		ok:         false,
	}
	if out1.ok {
		t.Error("Expected PBTestOut.ok to be false")
	}
	if out1.Output != "test_output" {
		t.Errorf("Expected Output to be 'test_output', got %v", out1.Output)
	}
	if len(out1.Predicates) != 1 {
		t.Errorf("Expected 1 predicate, got %d", len(out1.Predicates))
	}
	out2 := PBTestOut{
		Predicates: nil,
		ok:         true,
	}
	if !out2.ok {
		t.Error("Expected PBTestOut.ok to be true")
	}
	if out2.Predicates != nil {
		t.Error("Expected Predicates to be nil for passing case")
	}
}

func TestSatisfyAll_EdgeCases(t *testing.T) {
	pbt := &PBTest{predicates: []p.Predicate{}}
	ok, failed := pbt.satisfyAll(42)
	if !ok {
		t.Error("Expected satisfyAll to return true for empty predicates slice")
	}
	if failed != nil {
		t.Error("Expected no failed predicates for empty slice")
	}
}

func TestReturnTypesInterface(t *testing.T) {
	var rt1 returnTypes = "string"
	var rt2 returnTypes = []any{1, 2, 3}
	var rt3 returnTypes = 42
	var rt4 returnTypes = []int{1, 2, 3}
	_ = rt1
	_ = rt2
	_ = rt3
	_ = rt4
	t.Log("returnTypes interface test passed - all types assignable")
}

func TestNewPBTest_VariousFunctionTypes(t *testing.T) {
	pbt1 := NewPBTest(nil)
	if pbt1.f != nil {
		t.Error("Expected function to be nil")
	}
	pbt2 := NewPBTest("not a function")
	if pbt2.f != "not a function" {
		t.Error("Expected function field to store the provided value")
	}
	testFunc := func() {}
	pbt3 := NewPBTest(testFunc)
	if pbt3.f == nil {
		t.Error("Expected function to be set")
	}
}

func TestRun_WithIntFunction(t *testing.T) {
	intFunc := func(x int) int {
		return x * 2
	}
	pbt := NewPBTest(intFunc).WithT(t).WithIterations(1)
	results, err := pbt.Run()
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if len(results) != 0 {
		t.Errorf("Expected 0 results without predicates, got %d", len(results))
	}
}

func TestRun_WithStringFunction(t *testing.T) {
	stringFunc := func(s string) string {
		return "Hello " + s
	}
	pbt := NewPBTest(stringFunc).WithT(t).WithIterations(1)
	results, err := pbt.Run()
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if len(results) != 0 {
		t.Errorf("Expected 0 results without predicates, got %d", len(results))
	}
}

func TestRun_WithPredicatesAndIntFunction(t *testing.T) {
	intFunc := func(x int) int {
		return x * 2
	}
	pred := mockPredicate{shouldPass: true, name: "passing_pred"}
	pbt := NewPBTest(intFunc).WithT(t).WithIterations(1).WithPredicates(pred)
	results, err := pbt.Run()
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if len(results) == 0 {
		t.Error("Expected at least 1 result with predicates")
	}
	for _, result := range results {
		if !result.ok {
			t.Error("Expected result to be ok with passing predicate")
		}
	}
}

func TestRun_WithPredicatesAndStringFunction(t *testing.T) {
	stringFunc := func(s string) string {
		return "Hello " + s
	}
	pred := mockPredicate{shouldPass: false, name: "failing_pred"}
	pbt := NewPBTest(stringFunc).WithT(t).WithIterations(1).WithPredicates(pred)
	results, err := pbt.Run()
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if len(results) == 0 {
		t.Error("Expected at least 1 result with predicates")
	}
	for _, result := range results {
		if result.ok {
			t.Error("Expected result to fail with failing predicate")
		}
		if len(result.Predicates) != 1 {
			t.Errorf("Expected 1 failed predicate, got %d", len(result.Predicates))
		}
	}
}

func TestRun_WithMultipleIterations(t *testing.T) {
	intFunc := func(x int) int {
		return x + 1
	}
	pred := mockPredicate{shouldPass: true, name: "pred"}
	pbt := NewPBTest(intFunc).WithT(t).WithIterations(3).WithPredicates(pred)
	results, err := pbt.Run()
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if len(results) == 0 {
		t.Error("Expected results from multiple iterations")
	}
	expectedMinResults := 3
	if len(results) < expectedMinResults {
		t.Errorf("Expected at least %d results from 3 iterations, got %d", expectedMinResults, len(results))
	}
}

func TestRun_WithArrayReturningFunction(t *testing.T) {
	arrayFunc := func(x int) []any {
		return []any{x, x * 2, x * 3}
	}
	pred := mockPredicate{shouldPass: true, name: "pred"}
	pbt := NewPBTest(arrayFunc).WithT(t).WithIterations(1).WithPredicates(pred)
	results, err := pbt.Run()
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	expectedResults := 3
	if len(results) != expectedResults {
		t.Errorf("Expected %d results from array output, got %d", expectedResults, len(results))
	}
	for i, result := range results {
		if !result.ok {
			t.Errorf("Expected result %d to be ok", i)
		}
	}
}

func TestRun_WithSingleValueReturningFunction(t *testing.T) {
	singleFunc := func(x int) int {
		return x + 42
	}
	pred := mockPredicate{shouldPass: false, name: "failing_pred"}
	pbt := NewPBTest(singleFunc).WithT(t).WithIterations(1).WithPredicates(pred)
	results, err := pbt.Run()
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 result from single output, got %d", len(results))
	}
	if results[0].ok {
		t.Error("Expected result to fail with failing predicate")
	}
}

func TestRun_WithMixedPredicates(t *testing.T) {
	boolFunc := func(b bool) bool {
		return !b
	}
	passingPred := mockPredicate{shouldPass: true, name: "pass"}
	failingPred := mockPredicate{shouldPass: false, name: "fail"}
	pbt := NewPBTest(boolFunc).WithT(t).WithIterations(1).WithPredicates(passingPred, failingPred)
	results, err := pbt.Run()
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if len(results) == 0 {
		t.Error("Expected at least 1 result")
	}
	for _, result := range results {
		if result.ok {
			t.Error("Expected result to fail when any predicate fails")
		}
		if len(result.Predicates) != 1 {
			t.Errorf("Expected 1 failed predicate, got %d", len(result.Predicates))
		}
	}
}

func TestRun_WithNilFunction(t *testing.T) {
	pred := mockPredicate{shouldPass: true, name: "pred"}
	pbt := NewPBTest(nil).WithT(t).WithIterations(1).WithPredicates(pred)
	results, err := pbt.Run()
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if len(results) != 0 {
		t.Errorf("Expected 0 results with nil function, got %d", len(results))
	}
}

func TestRun_ZeroIterations(t *testing.T) {
	intFunc := func(x int) int {
		return x
	}
	pred := mockPredicate{shouldPass: true, name: "pred"}
	pbt := NewPBTest(intFunc).WithT(t).WithIterations(0).WithPredicates(pred)
	results, err := pbt.Run()
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if len(results) != 0 {
		t.Errorf("Expected 0 results with 0 iterations, got %d", len(results))
	}
}

func TestRun_ComplexArrayOutput(t *testing.T) {
	complexFunc := func(x int) []any {
		return []any{"string", x, true, []int{x, x * 2}, map[string]int{"key": x}}
	}
	pred := mockPredicate{shouldPass: true, name: "pred"}
	pbt := NewPBTest(complexFunc).WithT(t).WithIterations(1).WithPredicates(pred)
	results, err := pbt.Run()
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if len(results) != 5 {
		t.Errorf("Expected 5 results from complex array, got %d", len(results))
	}
	for i, result := range results {
		if !result.ok {
			t.Errorf("Expected result %d to be ok", i)
		}
	}
}

func TestRun_WithTwoParameterFunction(t *testing.T) {
	twoParamFunc := func(x int, y string) string {
		return fmt.Sprintf("%s_%d", y, x)
	}
	pred := mockPredicate{shouldPass: true, name: "pred"}
	pbt := NewPBTest(twoParamFunc).WithT(t).WithIterations(1).WithPredicates(pred)
	results, err := pbt.Run()
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if len(results) == 0 {
		t.Error("Expected at least 1 result")
	}
	for _, result := range results {
		if !result.ok {
			t.Error("Expected result to be ok with passing predicate")
		}
	}
}

func TestRun_WithFloatFunction(t *testing.T) {
	floatFunc := func(f float64) float64 {
		return f * 2.5
	}
	pred := mockPredicate{shouldPass: false, name: "failing_pred"}
	pbt := NewPBTest(floatFunc).WithT(t).WithIterations(1).WithPredicates(pred)
	results, err := pbt.Run()
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if len(results) == 0 {
		t.Error("Expected at least 1 result")
	}
	for _, result := range results {
		if result.ok {
			t.Error("Expected result to fail with failing predicate")
		}
	}
}

func TestRun_GenerateInputsError(t *testing.T) {
	// Test that Run() returns an error when GenerateInputs() fails
	// This happens when the function has an unsupported parameter type
	// (e.g., chan, func, interface)

	// Function with a channel parameter (unsupported by GenerateInputs)
	funcWithChannel := func(ch chan int) int {
		return 42
	}

	pbt := NewPBTest(funcWithChannel).WithIterations(1)
	results, err := pbt.Run()

	if err == nil {
		t.Error("Expected error when GenerateInputs fails with unsupported type")
	}

	if results != nil {
		t.Errorf("Expected nil results when error occurs, got %v", results)
	}
}

func TestRun_GenerateInputsErrorWithFunc(t *testing.T) {
	funcWithFunc := func(f func() int) int {
		return f()
	}

	pbt := NewPBTest(funcWithFunc).WithIterations(2)
	results, err := pbt.Run()

	if err == nil {
		t.Error("Expected error when GenerateInputs fails with function parameter")
	}

	if results != nil {
		t.Errorf("Expected nil results when error occurs, got %v", results)
	}
}

func TestRun_GenerateInputsErrorWithInterface(t *testing.T) {
	funcWithInterface := func(i interface{ DoSomething() }) int {
		return 42
	}

	pbt := NewPBTest(funcWithInterface).WithIterations(1)
	results, err := pbt.Run()

	if err == nil {
		t.Error("Expected error when GenerateInputs fails with interface parameter")
	}

	if results != nil {
		t.Errorf("Expected nil results when error occurs, got %v", results)
	}
}
