package stesting

import (
	"fmt"
	"os"
	"sync"

	gtu "github.com/laiambryant/gotestutils/testing"
)

// StressTest represents a parameterized stress testing framework that executes
// a test function repeatedly with a specified number of iterations.
//
// Type parameters:
//   - fRetType: The comparable return type of the test function
//   - testVarType: The comparable type of the test variable being used
//
// Fields:
//   - iterations: The number of times the test function will be executed
//   - testVar: A pointer to the test variable used during stress testing
//   - F: The test function to be executed, must conform to gtu.TestFunc[fRetType]
//
// This struct is designed to facilitate performance and reliability testing
// by running the same test function multiple times and collecting results.
type StressTest[fRetType comparable, testVarType comparable] struct {
	iterations uint32
	testVar    *testVarType
	F          gtu.TestFunc[fRetType]
}

// NewStressTest creates a new StressTest instance for running stress tests on a function.
// It takes the number of iterations to run, a test function that returns a comparable type,
// and a pointer to test variables of a comparable type.
//
// Type parameters:
//   - fRetType: the return type of the function being tested, must be comparable
//   - testVarType: the type of test variables, must be comparable
//
// Parameters:
//   - iterations: the number of times to execute the test function
//   - f: the test function to be executed repeatedly
//   - testVar: pointer to the test variables used by the function
//
// Returns:
//   - stressTest: a configured StressTest instance ready for execution
func NewStressTest[fRetType comparable, testVarType comparable](
	iterations uint32,
	f gtu.TestFunc[fRetType],
	testVar *testVarType,
) (stressTest StressTest[fRetType, testVarType]) {
	return StressTest[fRetType, testVarType]{
		iterations: iterations,
		testVar:    testVar,
		F:          f,
	}
}

// RunStressTest executes a stress test by running the specified function F for the given number of iterations.
// It takes a StressTest struct containing the function to test and iteration count.
// The function returns true and nil if all iterations complete successfully.
// If any iteration fails, it returns false and a StressTestingError containing the failing iteration index and the original error.
//
// Type parameters:
//   - fRetType: the return type of the function being tested (must be comparable)
//   - testVarType: the type of test variables used (must be comparable)
//
// Parameters:
//   - stressTest: pointer to StressTest struct containing the test configuration
//
// Returns:
//   - success: true if all iterations passed, false if any iteration failed
//   - err: nil on success, StressTestingError on failure containing iteration details
func RunStressTest[fRetType comparable, testVarType comparable](
	stressTest *StressTest[fRetType, testVarType],
) (success bool, err error) {
	for range stressTest.iterations {
		_, err = stressTest.F()
		if err != nil {
			return false, StressTestingError{Err: err}
		}
	}
	return true, nil
}

// RunParallelStressTest executes a stress test function concurrently across multiple workers.
// It runs the provided stress test's function for the specified number of iterations,
// distributing the work among up to maxWorkers goroutines.
//
// Type parameters:
//   - fRetType: the return type of the stress test function (must be comparable)
//   - testVarType: the type of test variables used in the stress test (must be comparable)
//
// Parameters:
//   - stressTest: a pointer to the StressTest instance containing the function to test
//     and the number of iterations to perform
//   - maxWorkers: the maximum number of concurrent goroutines to use for executing the test
//
// Returns:
//   - success: true if all iterations completed without errors, false if any iteration failed
//   - r_err: nil on success, or a StressTestingError containing details about the first
//     error encountered during execution
//
// The function stops execution and returns immediately upon encountering the first error.
// All workers are properly synchronized and cleaned up before returning.
func RunParallelStressTest[fRetType comparable, testVarType comparable](
	stressTest *StressTest[fRetType, testVarType],
	maxWorkers uint32,
) (success bool, rErr error) {
	errchan, jobs := make(chan error, stressTest.iterations), make(chan uint32)
	var wg sync.WaitGroup
	wg.Add(int(maxWorkers))
	for range maxWorkers {
		go func() {
			defer wg.Done()
			workerFunc(jobs, stressTest, errchan)
		}()
	}
	go func() {
		for i := uint32(0); i < stressTest.iterations; i++ {
			jobs <- i
		}
		close(jobs)
	}()
	for range stressTest.iterations {
		if rErr = <-errchan; rErr != nil {
			wg.Wait()
			close(errchan)
			if ste, ok := rErr.(StressTestingError); ok {
				return false, ste
			}
		}
	}
	wg.Wait()
	close(errchan)
	return true, nil
}

// workerFunc executes stress test iterations for parallel execution.
// It processes job indices from the jobs channel, runs the stress test function,
// and sends results (either nil for success or StressTestingError for failure) to the error channel.
//
// Type parameters:
//   - fRetType: the return type of the stress test function (must be comparable)
//   - testVarType: the type of test variables used in the stress test (must be comparable)
//
// Parameters:
//   - jobs: receive-only channel containing iteration indices to process
//   - stressTest: pointer to the StressTest instance containing the function to execute
//   - errchan: send-only channel for communicating results back to the coordinator
//
// For each job received, the function executes the stress test and sends either:
//   - nil to errchan if the test iteration succeeds
//   - StressTestingError to errchan if the test iteration fails, containing the index and error
func workerFunc[fRetType comparable, testVarType comparable](jobs <-chan uint32, stressTest *StressTest[fRetType, testVarType], errchan chan<- error) {
	for range jobs {
		_, err := stressTest.F()
		if err != nil {
			errchan <- StressTestingError{Err: err}
		} else {
			errchan <- nil
		}
	}
}

// RunStressTestWithFileOut executes a stress test and writes the output of each iteration to a file.
// It runs the test function for the specified number of iterations, writing each result to the provided file.
// The file is automatically closed when the function returns.
//
// Type parameters:
//   - fRetType: the return type of the test function, must be comparable
//   - testVarType: the type of test variables, must be comparable
//
// Parameters:
//   - stressTest: pointer to a StressTest containing the test function and iteration count
//   - file: os.File to write test results to
//
// Returns:
//   - success: true if all iterations completed without error, false otherwise
//   - err: nil on success, or a StressTestingError containing the iteration index and underlying error
//
// The function will stop execution and return false on the first error encountered.
// Each test result is written to the file using Go's %+#v format specifier.
func RunStressTestWithFileOut[fRetType comparable, testVarType comparable](
	stressTest *StressTest[fRetType, testVarType],
	file os.File,
) (success bool, err error) {
	defer file.Close()
	var out fRetType
	for i := uint32(0); i < stressTest.iterations; i++ {
		out, err = stressTest.F()
		file.WriteString(fmt.Sprintf("%+#v\n", out))
		if err != nil {
			return false, StressTestingError{Index: i, Err: err}
		}
	}
	return true, nil
}

// RunStressTestWithFilePathOut executes a stress test and writes the output to a file specified by filePath.
// It creates or opens the file at the given path, runs the stress test with file output,
// and automatically closes the file when done.
//
// Type parameters:
//   - fRetType: the comparable return type of the function being stress tested
//   - testVarType: the comparable type of the test variables used in the stress test
//
// Parameters:
//   - stressTest: pointer to the StressTest instance to execute
//   - filePath: string path to the output file where test results will be written
//
// Returns:
//   - success: boolean indicating whether the stress test completed successfully
//   - err: error if the stress test execution failed, nil otherwise
//
// Note: This function handles file creation and cleanup automatically. If file creation
// fails, the error from the underlying RunStressTestWithFileOut call will be returned.
func RunStressTestWithFilePathOut[fRetType comparable, testVarType comparable](
	stressTest *StressTest[fRetType, testVarType],
	filePath string,
) (success bool, err error) {
	file, _ := createAndOpenFile(filePath)
	defer file.Close()
	return RunStressTestWithFileOut(stressTest, *file)
}

// createAndOpenFile creates a new file or opens an existing file at the specified path
// for writing in append mode. The file is created with 0644 permissions if it doesn't exist.
// Returns a file pointer and any error encountered during the operation.
func createAndOpenFile(filePath string) (f *os.File, err error) {
	return os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}
