package stesting

import (
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"testing"
)

const (
	SuccessTrueMsg         string = "Expected success to be true, got false"
	ExpectedNoErrorMsg     string = "Expected no error, got %v"
	ExpectedErrorGotNilMsg string = "Expected error, got nil"
	ExpecteduUnsuccessMsg  string = "Expected success to be false, got true"
)

var (
	testFunc = func() (bool, error) {
		return true, nil
	}
	testFuncWithErr = func() (bool, error) {
		return false, errors.New("error")
	}
)

func assertSuccessNoError(t *testing.T, success bool, err error) {
	if !success {
		t.Error(SuccessTrueMsg)
	}
	if err != nil {
		t.Errorf(ExpectedNoErrorMsg, err)
	}
}

func assertNoSuccessError(t *testing.T, success bool, err error) {
	if success {
		t.Error(ExpecteduUnsuccessMsg)
	}
	if err == nil {
		t.Error(ExpectedErrorGotNilMsg)
	}
}

func TestRunParallelStressTestSuccess(t *testing.T) {
	var counter int64
	testFunc := func() (int, error) {
		newVal := atomic.AddInt64(&counter, 1)
		return int(newVal), nil
	}
	stressTest := NewStressTest[int, int](100, testFunc, nil)
	success, err := RunParallelStressTest(&stressTest, 4)
	assertSuccessNoError(t, success, err)
}

func TestRunParallelStressTestWithError(t *testing.T) {
	testError := errors.New("error")
	stressTest := NewStressTest[bool, int](10, testFuncWithErr, nil)
	success, err := RunParallelStressTest(&stressTest, 2)
	assertNoSuccessError(t, success, err)
	if ste, ok := err.(StressTestingError); ok {
		if ste.Err.Error() != testError.Error() {
			t.Errorf("Expected wrapped error to be %v, got %v", testError, ste.Err)
		}
	} else {
		t.Errorf("Expected StressTestingError, got %T", err)
	}
}

func TestRunParallelStressTestSingleWorker(t *testing.T) {
	iterations := uint32(50)
	stressTest := NewStressTest[bool, int](iterations, testFunc, nil)
	success, err := RunParallelStressTest(&stressTest, 1)
	assertSuccessNoError(t, success, err)
}

func TestRunParallelStressMultiWorker(t *testing.T) {
	iterations := uint32(50000)
	stressTest := NewStressTest[bool, int](iterations, testFunc, nil)
	success, err := RunParallelStressTest(&stressTest, uint32(10))
	assertSuccessNoError(t, success, err)
}

func TestRunParallelStressTestZeroIterations(t *testing.T) {
	stressTest := NewStressTest[bool, int](0, testFunc, nil)
	success, err := RunParallelStressTest(&stressTest, 4)
	assertSuccessNoError(t, success, err)
}

func TestTestingError(t *testing.T) {
	var errStr = "Error"
	err := StressTestingError{Err: errors.New(errStr), Index: 21}
	if err.Error() != "Error while running stress test at step "+fmt.Sprint(21)+" of testing: "+errStr {
		t.Error("Error message is incorrect")
	}
}

func TestRunStressTestSuccess(t *testing.T) {
	var counter int64
	testFunc := func() (int, error) {
		newVal := atomic.AddInt64(&counter, 1)
		return int(newVal), nil
	}
	stressTest := NewStressTest[int, int](100, testFunc, nil)
	success, err := RunStressTest(&stressTest)
	assertSuccessNoError(t, success, err)
}

func TestRunStressTestZeroIterations(t *testing.T) {
	stressTest := NewStressTest[bool, int](0, testFunc, nil)
	success, err := RunStressTest(&stressTest)
	assertSuccessNoError(t, success, err)
}

func TestRunStressTestSingleIteration(t *testing.T) {
	stressTest := NewStressTest[bool, int](1, testFunc, nil)
	success, err := RunStressTest(&stressTest)
	assertSuccessNoError(t, success, err)
}

func TestRunStressTestError(t *testing.T) {
	testFunc := func() (bool, error) {
		return true, StressTestingError{Err: errors.New(ExpectedNoErrorMsg)}
	}
	stressTest := NewStressTest[bool, int](2, testFunc, nil)
	_, err := RunStressTest(&stressTest)
	if err == nil {
		t.Error("Error expected")
	}
}
func TestRunStressTestWithFileOut(t *testing.T) {
	tempFile, err := os.CreateTemp("", "stress_test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	counter := 0
	testCounterFunc := func() (int, error) {
		counter++
		return counter, nil
	}
	stressTest := NewStressTest[int, int](5, testCounterFunc, nil)
	success, err := RunStressTestWithFileOut(&stressTest, *tempFile)
	assertSuccessNoError(t, success, err)
	tempFile, err = os.OpenFile(tempFile.Name(), os.O_RDONLY, 0644)
	if err != nil {
		t.Fatalf("Failed to open temp file: %v", err)
	}
	content := make([]byte, 1024)
	n, err := tempFile.Read(content)
	if err != nil && err.Error() != "EOF" {
		t.Fatalf("Failed to read from file: %v", err)
	}
	fileContent := string(content[:n])
	if fileContent == "" {
		t.Error("Expected file to contain output, got empty file")
	}
}

func TestRunStressTestWithFileOutError(t *testing.T) {
	tempFile, err := os.CreateTemp("", "stress_test_error_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	testError := errors.New("error")
	stressTest := NewStressTest[bool, int](3, testFunc, nil)
	success, err := RunStressTestWithFileOut(&stressTest, *tempFile)
	assertSuccessNoError(t, success, err)
	if ste, ok := err.(StressTestingError); ok {
		if ste.Err != testError {
			t.Errorf("Expected wrapped error to be %v, got %v", testError, ste.Err)
		}
	}
}

func TestRunStressTestWithFilePathOut(t *testing.T) {
	tempDir := t.TempDir()
	filePath := fmt.Sprintf("%s/test_output.txt", tempDir)
	stressTest := NewStressTest[bool, int](3, testFunc, nil)
	success, err := RunStressTestWithFilePathOut(&stressTest, filePath)
	assertSuccessNoError(t, success, err)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("Expected file to be created")
	}
}

func TestRunStressTestWithFilePathOutError(t *testing.T) {
	tempDir := t.TempDir()
	filePath := fmt.Sprintf("%s/test_error_output.txt", tempDir)
	stressTest := NewStressTest[bool, int](2, testFuncWithErr, nil)
	success, err := RunStressTestWithFilePathOut(&stressTest, filePath)
	if success {
		t.Error(ExpecteduUnsuccessMsg)
	}
	if err == nil {
		t.Error(ExpectedErrorGotNilMsg)
	}
}

func TestCreateAndOpenFileError(t *testing.T) {
	tempDir := t.TempDir()
	f, err := createAndOpenFile(tempDir)
	if err == nil {
		t.Error(ExpectedErrorGotNilMsg)
	} else {
		f.Close()
	}
}
