# goTestUtils Examples

This folder contains complete, runnable examples for all the functionality demonstrated in the main README.

## Characterization Testing Examples

- [`basic_characterization_test.go`](basic_characterization_test.go) - Basic characterization test creation
- [`execution_examples_test.go`](execution_examples_test.go) - Test execution patterns and error checking modes
- [`string_processing_test.go`](string_processing_test.go) - Testing string processing functions
- [`complex_data_structures_test.go`](complex_data_structures_test.go) - Testing with structs and complex types
- [`http_response_test.go`](http_response_test.go) - HTTP response testing with mock servers
- [`error_handling_test.go`](error_handling_test.go) - Comprehensive error condition testing
- [`manual_result_processing_test.go`](manual_result_processing_test.go) - Custom result processing

## Fuzz Testing Examples

- [`basic_fuzz_test.go`](basic_fuzz_test.go) - Basic fuzz testing with default and custom attributes
  - Simple integer functions
  - String processing with custom attributes
  - Multi-parameter functions with different types
  - Floating-point operations
  - Direct input generation examples

- [`advanced_fuzz_test.go`](advanced_fuzz_test.go) - Advanced fuzz testing with complex types
  - Complex number fuzzing
  - Slice and array fuzzing
  - Map fuzzing with custom key/value attributes
  - Pointer fuzzing (including nil handling)
  - Struct fuzzing with dynamic types
  - Unsigned integer fuzzing
  - Advanced string attributes (prefix, suffix, allowed characters)
  - Boolean input combinations

- [`fuzz_edge_cases_test.go`](fuzz_edge_cases_test.go) - Using fuzzing to discover bugs
  - Error detection and panic recovery
  - Boundary condition bugs
  - Type conversion issues
  - Concurrency safety testing
  - Data validation fuzzing
  - Performance characteristic discovery
  - Nil handling verification
  - Composite type fuzzing (nested maps/slices)

## Stress Testing Examples

- [`basic_stress_test.go`](basic_stress_test.go) - Basic stress testing patterns
- [`web_service_stress_test.go`](web_service_stress_test.go) - Web service stress testing
- [`memory_stress_test.go`](memory_stress_test.go) - Memory allocation stress testing
- [`stress_error_handling_test.go`](stress_error_handling_test.go) - Error handling in stress tests

## Helper Files

- [`helpers.go`](helpers.go) - Shared helper functions used across examples

## Running the Examples

To run all examples:

```bash
cd examples
go test -v
```

To run specific examples:

```bash
go test -v -run TestBasicStressExample
go test -v -run TestStringProcessing
go test -v -run TestBasicFuzzExample
go test -v -run TestFuzzComplexNumbers
```

To run only fuzz tests:

```bash
go test -v -run Fuzz
```

To run with race detection (recommended for concurrency tests):

```bash
go test -race -v
```

## Notes

- All examples are designed to be self-contained and runnable
- Some tests may take longer due to stress testing iterations or extensive fuzzing
- The web service stress test uses an external service (httpbin.org) and may fail if the service is unavailable
- Examples demonstrate both successful test cases and intentional failures for educational purposes
- Fuzz testing examples may discover different edge cases on different runs due to randomness
- Some fuzz edge case examples intentionally trigger panics to demonstrate error discovery capabilities
