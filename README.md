# goTestUtils

[![Go Reference](https://pkg.go.dev/badge/github.com/laiambryant/gotestutils.svg)](https://pkg.go.dev/github.com/laiambryant/gotestutils)
[![Go Report Card](https://goreportcard.com/badge/github.com/laiambryant/gotestutils)](https://goreportcard.com/report/github.com/laiambryant/gotestutils)
[![GitHub license](https://img.shields.io/github/license/laiambryant/gotestutils.svg)](https://github.com/laiambryant/gotestutils/blob/main/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/laiambryant/gotestutils.svg)](https://github.com/laiambryant/gotestutils/issues)
[![GitHub stars](https://img.shields.io/github/stars/laiambryant/gotestutils.svg)](https://github.com/laiambryant/gotestutils/stargazers)
[![Coverage Status](https://coveralls.io/repos/github/laiambryant/goTestUtils/badge.svg?branch=main)](https://coveralls.io/github/laiambryant/goTestUtils?branch=main)

My Favourite testing utility methods for go. Includes utilities for:

- Characterization testing
- Stress testing
- More to come...

## Installation

```bash
go get github.com/laiambryant/gotestutils
```

## Examples

All code examples referenced in this README are available as complete, runnable examples in the [`examples/`](examples/) folder. Each example file contains working tests that demonstrate the usage patterns described below.

## Characterization Testing

The `ctesting` package provides a powerful framework for characterization testing in Go. Characterization tests capture and verify the current behavior of existing code by comparing actual outputs and errors against expected values. This approach is particularly valuable when working with legacy code, refactoring existing systems, or documenting the behavior of complex functions.

### Core Components

#### CharacterizationTest Structure

The `CharacterizationTest[t comparable]` struct is the foundation of the characterization testing framework:

```go
type CharacterizationTest[t comparable] struct {
    err            error           // Actual error (populated during execution)
    ExpectedErr    error           // Expected error
    output         t               // Actual output (populated during execution)
    ExpectedOutput t               // Expected output value
    F              gtu.TestFunc[t] // The test function to execute
}
```

#### Creating Characterization Tests

Use `NewCharacterizationTest` to create test instances that capture the current behavior of your code. Each test consists of an expected output value, an expected error (or nil if no error is expected), and a function to execute. The test function should match the `TestFunc[t]` signature, returning a value of type `t` and an error. This approach allows you to document and verify the exact behavior of functions, making it easier to detect unintended changes during refactoring or maintenance.

### Execution Methods

#### Basic Test Execution

`VerifyCharacterizationTests` executes a suite of characterization tests and returns detailed results. It runs each test function in the provided slice, captures the actual outputs and errors, then compares them against the expected values. The function returns a boolean slice indicating which tests passed or failed, along with an updated test suite containing the actual results. This separation of execution and reporting allows for custom result processing when needed.

#### Integrated Test Execution and Reporting

`VerifyCharacterizationTestsAndResults` combines test execution and result reporting into a single convenient function call. It automatically executes all tests in the suite, compares results against expected values, and reports outcomes directly to the provided `testing.T` instance. This streamlined approach reduces boilerplate code in test functions while providing comprehensive logging of both successful and failed test cases.

#### Error Checking Modes

The framework supports two error comparison strategies to accommodate different testing needs:

**Deep Error Checking (recommended):** Performs exact string comparison of error messages, ensuring precise matching of error content. This mode is ideal when you need to verify specific error text or when working with custom error types that have meaningful messages.

**Shallow Error Checking:** Uses Go's `errors.Is()` function for error comparison, which handles error wrapping and allows for more flexible error matching. This mode is useful when testing with wrapped errors or when the exact error message is less important than the error type or cause.

### Characterization Testing Examples

Complete examples for characterization testing:

- [`examples/basic_characterization_test.go`](examples/basic_characterization_test.go) - Creating characterization tests that capture expected behavior
- [`examples/execution_examples_test.go`](examples/execution_examples_test.go) - Test execution patterns, error checking modes, and integrated reporting
- [`examples/string_processing_test.go`](examples/string_processing_test.go) - Testing string processing functions with various inputs including edge cases
- [`examples/complex_data_structures_test.go`](examples/complex_data_structures_test.go) - Testing functions that work with structs and complex data types
- [`examples/http_response_test.go`](examples/http_response_test.go) - Characterization testing with HTTP responses and mock servers
- [`examples/error_handling_test.go`](examples/error_handling_test.go) - Testing various error scenarios including division by zero and validation errors
- [`examples/manual_result_processing_test.go`](examples/manual_result_processing_test.go) - Advanced result processing with custom failure handling and error analysis

### Test Result Analysis

The framework provides detailed reporting for both successful and failed tests:

#### Success Reporting

```text
test number 1: SUCCESS [ERRORS] got error {<nil>}, expected {<nil>}, [VALUES] got {42} expected {42}
```

#### Failure Reporting

```text
test number 2: ERROR [ERRORS] got error {<nil>}, expected {division by zero}, [VALUES] got {0} expected {0}
```

## Stress Testing Framework

The `stesting` package provides a comprehensive framework for stress testing Go functions to evaluate their performance, reliability, and behavior under load. Stress tests execute a function repeatedly for a specified number of iterations to identify potential issues, memory leaks, race conditions, or performance degradation.

### Stress Test Components

#### StressTest Structure

The `StressTest[fRetType, testVarType]` struct is the foundation of the stress testing framework. See the [source code](stesting/stesting.go) for the complete struct definition.

```go
type StressTest[fRetType comparable, testVarType comparable] struct {
    iterations uint64 // Number of iterations to run
    testVar    *testVarType // Variable who's value should be tested at the end of the test
    F          gtu.TestFunc[fRetType] // The test function to execute
}
```

#### Creating a Stress Test

See [`examples/basic_stress_test.go`](examples/basic_stress_test.go) for complete examples of creating stress test instances.

### Stress Test Execution Methods

#### Sequential Execution

`RunStressTest` executes the test function repeatedly in a single goroutine, running each iteration sequentially. This approach provides predictable, deterministic execution order and is essential when testing functions that are not thread-safe or when you need to measure sequential performance characteristics. Sequential execution also makes it easier to reproduce and debug issues since there are no concurrency-related variables affecting the test results.

**Use cases:**

- Testing functions that are not thread-safe
- Measuring sequential performance
- Simple reliability testing

#### Parallel Execution

`RunParallelStressTest` distributes the stress test workload across multiple goroutines, allowing concurrent execution of the test function. You specify the maximum number of worker goroutines, and the framework distributes iterations among them using a work queue pattern. This approach is valuable for testing concurrent safety, identifying race conditions, simulating realistic load scenarios, and evaluating performance under parallel execution. The function stops immediately upon encountering the first error and properly synchronizes all workers before returning.

**Use cases:**

- Testing concurrent safety and race conditions
- Load testing with realistic concurrency
- Performance testing under parallel execution
- Identifying deadlocks or synchronization issues

#### File Output Testing

The framework provides functions to save stress test results to files for detailed analysis. `RunStressTestWithFilePathOut` creates a file at the specified path and writes each iteration's output, while `RunStressTestWithFileOut` uses an existing file handle. This capability is particularly useful for analyzing output patterns across many iterations, investigating intermittent issues that only appear under sustained load, creating audit trails for compliance testing, and performing post-execution analysis of performance trends or data patterns.

**Use cases:**

- Analyzing output patterns over many iterations
- Debugging intermittent issues
- Performance profiling and trend analysis
- Compliance testing with audit trails

### Stress Testing Examples

Complete examples for stress testing:

- [`examples/basic_stress_test.go`](examples/basic_stress_test.go) - Basic stress testing patterns with sequential and parallel execution
- [`examples/web_service_stress_test.go`](examples/web_service_stress_test.go) - Stress testing web service calls with concurrent workers
- [`examples/memory_stress_test.go`](examples/memory_stress_test.go) - Testing memory-intensive functions with file output for analysis
- [`examples/stress_error_handling_test.go`](examples/stress_error_handling_test.go) - Error handling in stress tests with detailed iteration information

### Stress Testing Best Practices

1. **Choose appropriate iteration counts**: Start with smaller numbers and increase based on your needs
2. **Monitor resource usage**: Use tools like `go test -race` and `go test -memprofile` with stress tests
3. **Use parallel tests judiciously**: Not all functions are safe for concurrent testing
4. **Combine with benchmarks**: Use `go test -bench` alongside stress tests for comprehensive performance analysis
5. **File output for debugging**: Use file output when investigating intermittent failures
6. **Gradual load increase**: Start with fewer workers/iterations and gradually increase to find breaking points

### Integration with Go Testing

Stress tests integrate seamlessly with Go's testing framework. See [`examples/basic_stress_test.go`](examples/basic_stress_test.go) for examples of integration patterns.
