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
- Property Based Testing
- Fuzz testing
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

### Characterization Test Execution

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
    testVar    *testVarType // Value to test at the end of the test
    F          gtu.TestFunc[fRetType] // The test function to execute
}
```

#### Creating a Stress Test

See [`examples/basic_stress_test.go`](examples/basic_stress_test.go) for complete examples of creating stress test instances.

### Stress Test Execution Methods

#### Sequential Execution

`RunStressTest` executes the test function repeatedly in a single goroutine, running each iteration sequentially. This approach provides predictable, deterministic execution order and is essential when testing functions that are not thread-safe or when you need to measure sequential performance characteristics. Sequential execution also makes it easier to reproduce and debug issues since there are no concurrency-related variables affecting the test results.

#### Parallel Execution

`RunParallelStressTest` distributes the stress test workload across multiple goroutines, allowing concurrent execution of the test function. You specify the maximum number of worker goroutines, and the framework distributes iterations among them using a work queue pattern. This approach is valuable for testing concurrent safety, identifying race conditions, simulating realistic load scenarios, and evaluating performance under parallel execution. The function stops immediately upon encountering the first error and properly synchronizes all workers before returning.

#### File Output Testing

The framework provides functions to save stress test results to files for detailed analysis. `RunStressTestWithFilePathOut` creates a file at the specified path and writes each iteration's output, while `RunStressTestWithFileOut` uses an existing file handle. This capability is particularly useful for analyzing output patterns across many iterations, investigating intermittent issues that only appear under sustained load, creating audit trails for compliance testing, and performing post-execution analysis of performance trends or data patterns.

### Stress Testing Examples

Complete examples for stress testing:

- [`examples/basic_stress_test.go`](examples/basic_stress_test.go) - Basic stress testing patterns with sequential and parallel execution
- [`examples/web_service_stress_test.go`](examples/web_service_stress_test.go) - Stress testing web service calls with concurrent workers
- [`examples/memory_stress_test.go`](examples/memory_stress_test.go) - Testing memory-intensive functions with file output for analysis
- [`examples/stress_error_handling_test.go`](examples/stress_error_handling_test.go) - Error handling in stress tests with detailed iteration information

### Stress Testing Integration

Stress tests integrate seamlessly with Go's testing framework. See [`examples/basic_stress_test.go`](examples/basic_stress_test.go) for examples of integration patterns.

## Property-Based Testing Framework

The `pbtesting` package provides a comprehensive property-based testing framework that validates functions satisfy mathematical properties and logical invariants across randomly generated inputs. Unlike traditional example-based tests that check specific input-output pairs, property-based tests verify that certain properties hold for all inputs within a domain.

### Core Concepts

Property-based testing validates that functions satisfy predicates (properties) across many random inputs. For example, instead of testing that `abs(-5) == 5`, you verify the property "absolute value is always non-negative" for thousands of random inputs.

#### PBTest Structure

The `PBTest` struct is the foundation of property-based testing:

```go
type PBTest struct {
    t          *testing.T      // Testing instance for reporting
    f          any             // Function under test
    predicates []Predicate     // Properties to validate
    iterations uint            // Number of test iterations
    argAttrs   []any          // Custom input generation attributes
}
```

#### Creating Property-Based Tests

Use `NewPBTest` to create test instances that validate function properties:

```go
test := NewPBTest(myFunc).
    WithIterations(1000).
    WithPredicates(nonNegative, lessThan100).
    WithT(t)
```

### Property Test Execution

#### Run - Basic Execution

`Run()` executes property tests with default random input generation:

```go
results, err := test.Run()
if err != nil {
    t.Fatal(err)
}

failures := FilterPBTTestOut(results)
if len(failures) > 0 {
    t.Errorf("Found %d property violations", len(failures))
}
```

#### RunWithAttributes - Constrained Input Generation

`RunWithAttributes()` provides fine-grained control over random input generation by accepting custom attributes. This allows you to constrain the input space to specific ranges, types, or characteristics:

```go
// Example: Constrain integers to positive range
attrs := attributes.NewFTAttributes()
attrs.IntegerAttr = attributes.IntegerAttributesImpl[int]{
    Min:           1,
    Max:           100,
    AllowNegative: false,
    AllowZero:     false,
}

results, err := test.RunWithAttributes(attrs)
```

**Key Features:**

- Control integer ranges (Min, Max, AllowNegative, AllowZero)
- Constrain float values (Min, Max, FiniteOnly, NonZero)
- Specify string lengths (MinLen, MaxLen)
- Configure slice/array sizes and element constraints
- Customize struct field generation

**Use Cases:**

- Testing functions with domain-specific constraints
- Avoiding edge cases irrelevant to your use case
- Focusing tests on particular input ranges
- Validating behavior within business rule boundaries

### Predicates

Predicates define the properties that function outputs must satisfy. Implement the `Predicate` interface:

```go
type Predicate interface {
    Verify(val any) bool
}
```

**Common Property Patterns:**

- **Idempotence**: `f(f(x)) == f(x)`
- **Commutativity**: `f(a, b) == f(b, a)`
- **Associativity**: `f(f(a, b), c) == f(a, f(b, c))`
- **Identity**: `f(x, identity) == x`
- **Inverse**: `f(g(x)) == x`

### Property-Based Testing Examples

Complete examples demonstrating property-based testing:

- [`examples/basic_property_test.go`](examples/basic_property_test.go) - Fundamental property patterns (idempotence, commutativity, inverse operations)
- [`examples/advanced_property_test.go`](examples/advanced_property_test.go) - Advanced properties with custom attributes, complex types, and error handling
- [`examples/property_patterns_test.go`](examples/property_patterns_test.go) - Reusable property testing patterns (monotonicity, round-trip, boundary conditions)

### Property Testing Integration

Property-based tests integrate seamlessly with Go's testing framework:

```go
func TestSquareRootProperty(t *testing.T) {
    // Property: sqrt(x)² should equal x for non-negative values
    sqrtProperty := Predicate{
        Verify: func(val any) bool {
            result := val.(float64)
            return result >= 0
        },
    }

    attrs := attributes.NewFTAttributes()
    attrs.FloatAttr = attributes.FloatAttributesImpl[float64]{
        Min:        0.0,
        Max:        10000.0,
        FiniteOnly: true,
    }

    test := NewPBTest(math.Sqrt).
        WithIterations(1000).
        WithPredicates(sqrtProperty).
        WithT(t)

    results, err := test.RunWithAttributes(attrs)
    if err != nil {
        t.Fatal(err)
    }

    failures := FilterPBTTestOut(results)
    if len(failures) > 0 {
        t.Errorf("Found %d violations of sqrt property", len(failures))
        for i, failure := range failures {
            if i < 5 { // Show first 5 failures
                t.Logf("  Failure: %v", failure.Output)
            }
        }
    }
}
```

### Benefits of Property-Based Testing

- **Broader Coverage**: Tests thousands of inputs automatically
- **Edge Case Discovery**: Finds corner cases you didn't think to test
- **Living Documentation**: Properties document function behavior
- **Regression Prevention**: Catches behavioral changes across refactoring
- **Domain Modeling**: Codifies business rules and mathematical properties

## Fuzz Testing Framework

The `ftesting` package provides a comprehensive fuzz testing framework that automatically generates random inputs for functions with arbitrary signatures. Unlike traditional testing approaches that require manual input creation, fuzz testing uses reflection and sophisticated attribute systems to generate type-appropriate random values, helping discover edge cases, boundary conditions, and unexpected behaviors.

### Fuzz Testing Fundamentals

Fuzz testing validates function robustness by executing it with thousands of randomly generated inputs. This approach is particularly effective for finding:

- **Edge Cases**: Inputs that cause unexpected behavior or crashes
- **Boundary Conditions**: Values at the limits of acceptable ranges  
- **Input Validation Issues**: Missing or insufficient input checking
- **Panic Conditions**: Inputs that cause runtime panics
- **Performance Issues**: Inputs that cause performance degradation

#### FTesting Structure

The `FTesting` struct is the foundation of fuzz testing:

```go
type FTesting struct {
    f          any                    // Function to test (any signature)
    iterations uint                   // Number of test iterations
    attributes AttributesStruct       // Custom input generation rules
    t          *testing.T             // Testing instance for reporting
}
```

#### Creating Fuzz Tests

Use method chaining to create and configure fuzz test instances:

```go
ft := &ftesting.FTesting{}
ft.WithFunction(myFunc).
   WithIterations(1000).
   WithAttributes(customAttrs)
```

### Fuzz Test Execution

#### Basic Execution

`ApplyFunction()` generates random inputs and executes the test function:

```go
ft := &ftesting.FTesting{}
ft.WithFunction(func(x int, y string) int {
    return x + len(y)
})

for i := 0; i < 100; i++ {
    success, err := ft.ApplyFunction()
    if err != nil {
        t.Errorf("Iteration %d failed: %v", i, err)
    }
}
```

#### Integrated Testing

`Verify()` provides integrated execution and reporting with Go's testing framework:

```go
func TestMyFunction(t *testing.T) {
    ft := &ftesting.FTesting{t: t}
    ft.WithFunction(myFunction).
       WithIterations(1000).
       Verify()
}
```

#### Input Generation

`GenerateInputs()` creates random inputs without execution for custom testing scenarios:

```go
ft.WithFunction(func(x, y int) int { return x + y })
inputs, err := ft.GenerateInputs()
// inputs might be: []any{42, -17}
```

### Attributes System

The attributes system provides fine-grained control over random value generation:

#### Default Attributes

```go
attrs := attributes.NewFTAttributes() // Uses sensible defaults for all types
ft.WithAttributes(attrs)
```

#### Custom Constraints

```go
attrs := attributes.NewFTAttributes()

// Constrain integers to positive range
attrs.IntegerAttr = attributes.IntegerAttributesImpl[int]{
    Min:           1,
    Max:           1000,
    AllowNegative: false,
    AllowZero:     false,
}

// Control string generation
attrs.StringAttr = attributes.StringAttributes{
    MinLen: 5,
    MaxLen: 50,
}

// Configure float constraints
attrs.FloatAttr = attributes.FloatAttributesImpl[float64]{
    Min:        0.0,
    Max:        100.0,
    FiniteOnly: true,    // Exclude NaN, +Inf, -Inf
    NonZero:    true,    // Exclude zero values
}

ft.WithAttributes(attrs)
```

#### Supported Types and Constraints

- **Integers**: Min/Max ranges, zero/negative value control
- **Floats**: Ranges, finite-only mode, zero exclusion
- **Strings**: Length constraints, character set control
- **Booleans**: Force true/false values or random distribution
- **Slices/Arrays**: Length constraints, element generation rules
- **Structs**: Field-by-field attribute configuration
- **Pointers**: Nil probability, depth control
- **Maps**: Size constraints, key/value generation rules

### Fuzz Testing Examples

Complete examples demonstrating fuzz testing:

- [`examples/basic_fuzz_test.go`](examples/basic_fuzz_test.go) - Basic fuzz testing patterns with different data types and custom attributes
- [`examples/advanced_fuzz_test.go`](examples/advanced_fuzz_test.go) - Advanced fuzzing with complex numbers, slices, maps, and composite types
- [`examples/fuzz_edge_cases_test.go`](examples/fuzz_edge_cases_test.go) - Advanced fuzzing for edge case discovery, boundary condition testing, and panic detection

### Integration with Testing Framework

Fuzz tests integrate seamlessly with Go's testing framework:

```go
func TestStringProcessor(t *testing.T) {
    attrs := attributes.NewFTAttributes()
    attrs.StringAttr = attributes.StringAttributes{
        MinLen: 1,
        MaxLen: 100,
    }

    ft := &ftesting.FTesting{t: t}
    ft.WithFunction(func(s string) (string, error) {
        if len(s) > 50 {
            return "", fmt.Errorf("string too long")
        }
        return strings.ToUpper(s), nil
    }).WithAttributes(attrs).WithIterations(500)

    // Test with panic recovery
    for i := 0; i < 500; i++ {
        func() {
            defer func() {
                if r := recover(); r != nil {
                    t.Errorf("Function panicked with: %v", r)
                }
            }()
            
            _, err := ft.ApplyFunction()
            if err != nil {
                // Expected for long strings - this is fine
            }
        }()
    }
}
```

### Benefits of Fuzz Testing

- **Automated Edge Case Discovery**: Finds problematic inputs without manual specification
- **Comprehensive Coverage**: Tests far more input combinations than manual testing
- **Robustness Validation**: Ensures functions handle unexpected inputs gracefully
- **Regression Prevention**: Catches issues introduced during refactoring
- **Security Testing**: Discovers input validation vulnerabilities
- **Performance Analysis**: Identifies inputs that cause performance issues
