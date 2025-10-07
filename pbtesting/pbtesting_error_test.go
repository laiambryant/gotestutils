package pbtesting

import (
	"errors"
	"strings"
	"testing"
)

type mockPredicateForError struct {
	name string
}

func (m mockPredicateForError) Verify(val any) bool {
	return false
}

func (m mockPredicateForError) String() string {
	return m.name
}

func TestInvalidPropertyError(t *testing.T) {
	pred := mockPredicateForError{name: "test_predicate"}
	err := InvalidPropertyError{predicate: pred}
	expectedMsg := "invalid property: test_predicate"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestFunctionNotProvidedError(t *testing.T) {
	err := FunctionNotProvidedError{}
	expectedMsg := "a function must be provided for the property-based test suite to work"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestInvalidFunctionProvidedError(t *testing.T) {
	testFunc := func(a, b int) int { return a + b }
	err := InvalidFunctionProvidedError{f: testFunc}
	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "Invalid function provided to pbt") {
		t.Errorf("Expected error message to contain 'Invalid function provided to pbt', got '%s'", errorMsg)
	}
	if !strings.Contains(errorMsg, "function:") {
		t.Errorf("Expected error message to contain 'function:', got '%s'", errorMsg)
	}
}

func TestInvalidFunctionProvidedError_NilFunction(t *testing.T) {
	err := InvalidFunctionProvidedError{f: nil}
	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "Invalid function provided to pbt") {
		t.Errorf("Expected error message to contain 'Invalid function provided to pbt', got '%s'", errorMsg)
	}
	if !strings.Contains(errorMsg, "<nil>") {
		t.Errorf("Expected error message to contain '<nil>', got '%s'", errorMsg)
	}
}

func TestErrorVariables(t *testing.T) {
	if ErrZeroRangeNonZeroRequested == nil {
		t.Error("Expected ErrZeroRangeNonZeroRequested to be defined")
	}
	expectedMsg1 := "zero range but non-zero requested"
	if ErrZeroRangeNonZeroRequested.Error() != expectedMsg1 {
		t.Errorf("Expected ErrZeroRangeNonZeroRequested message '%s', got '%s'",
			expectedMsg1, ErrZeroRangeNonZeroRequested.Error())
	}
	if ErrMinGreaterThanMax == nil {
		t.Error("Expected ErrMinGreaterThanMax to be defined")
	}
	expectedMsg2 := "minimum is greater than maximum"
	if ErrMinGreaterThanMax.Error() != expectedMsg2 {
		t.Errorf("Expected ErrMinGreaterThanMax message '%s', got '%s'",
			expectedMsg2, ErrMinGreaterThanMax.Error())
	}
}

func TestErrorTypeAssertions(t *testing.T) {
	pred := mockPredicateForError{name: "test"}
	err1 := InvalidPropertyError{predicate: pred}
	if _, ok := any(err1).(InvalidPropertyError); !ok {
		t.Error("Expected InvalidPropertyError to be assertable to its own type")
	}
	err2 := FunctionNotProvidedError{}
	if _, ok := any(err2).(FunctionNotProvidedError); !ok {
		t.Error("Expected FunctionNotProvidedError to be assertable to its own type")
	}
	err3 := InvalidFunctionProvidedError{f: "test"}
	if _, ok := any(err3).(InvalidFunctionProvidedError); !ok {
		t.Error("Expected InvalidFunctionProvidedError to be assertable to its own type")
	}
}

func TestPredefinedErrorsAreDifferent(t *testing.T) {
	if ErrZeroRangeNonZeroRequested == ErrMinGreaterThanMax {
		t.Error("Expected predefined errors to be different instances")
	}
	var err1 error = ErrZeroRangeNonZeroRequested
	var err2 error = ErrMinGreaterThanMax
	if errors.Is(err1, err2) {
		t.Error("Expected predefined errors to not be equal")
	}
}
