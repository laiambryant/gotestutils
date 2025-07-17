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

## Quick Start

### Characterization Testing

```go
import "github.com/laiambryant/gotestutils/ctesting"

// Example usage
testSuite := []ctesting.CharacterizationTest[int]{
    ctesting.NewCharacterizationTest(3, nil, func() (int, error) { 
        return sum(1, 2), nil 
    }),
}
results, _ := ctesting.VerifyCharacterizationTestsAndResults(t, testSuite)
```

### Stress Testing

```go
import "github.com/laiambryant/gotestutils/stesting"

// Basic stress test - runs a function 1000 times sequentially
stressTest := stesting.NewStressTest[int, any](1000, func() (int, error) {
    return expensiveOperation(), nil
}, nil)
success, err := stesting.RunStressTest(&stressTest)

// Parallel stress test - distributes work across multiple goroutines
success, err = stesting.RunParallelStressTest(&stressTest, 10) // 10 workers

// Stress test with file output - saves results to a file
success, err = stesting.RunStressTestWithFilePathOut(&stressTest, "output.txt")
```

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

Use `NewCharacterizationTest` to create test instances that capture expected behavior:

```go
// Test expecting successful execution with specific output
test := ctesting.NewCharacterizationTest(42, nil, func() (int, error) {
    return calculateResult(6, 7), nil
})

// Test expecting an error condition
test := ctesting.NewCharacterizationTest(0, fmt.Errorf("division by zero"), func() (int, error) {
    return divide(10, 0)
})
```

### Execution Methods

#### Basic Test Execution

`VerifyCharacterizationTests` executes a suite of tests and returns detailed results:

```go
testSuite := []ctesting.CharacterizationTest[int]{
    ctesting.NewCharacterizationTest(3, nil, func() (int, error) { 
        return sum(1, 2), nil 
    }),
    ctesting.NewCharacterizationTest(10, nil, func() (int, error) { 
        return multiply(2, 5), nil 
    }),
}

// Execute tests with deep error checking
results, updatedSuite := ctesting.VerifyCharacterizationTests(testSuite, true)

// Process results manually
for i, passed := range results {
    if passed {
        fmt.Printf("Test %d: PASSED\n", i+1)
    } else {
        fmt.Printf("Test %d: FAILED\n", i+1)
    }
}
```

#### Integrated Test Execution and Reporting

`VerifyCharacterizationTestsAndResults` combines execution and reporting for convenience:

```go
// Automatically execute tests and report results to testing.T
results, updatedSuite := ctesting.VerifyCharacterizationTestsAndResults(t, testSuite, true)
```

#### Error Checking Modes

The framework supports two error checking modes:

**Deep Error Checking (recommended):**

```go
// Performs exact error message comparison
results, _ := ctesting.VerifyCharacterizationTests(testSuite, true)
```

**Shallow Error Checking:**

```go
// Uses errors.Is() for error comparison
results, _ := ctesting.VerifyCharacterizationTests(testSuite, false)
```

### Advanced Examples

#### Testing String Processing Functions

```go
func TestStringProcessingCharacterization(t *testing.T) {
    testSuite := []ctesting.CharacterizationTest[string]{
        // Test normal input
        ctesting.NewCharacterizationTest("HELLO", nil, func() (string, error) {
            return strings.ToUpper("hello"), nil
        }),
        
        // Test empty string
        ctesting.NewCharacterizationTest("", nil, func() (string, error) {
            return strings.ToUpper(""), nil
        }),
        
        // Test with special characters
        ctesting.NewCharacterizationTest("HELLO, WORLD!", nil, func() (string, error) {
            return strings.ToUpper("hello, world!"), nil
        }),
    }
    
    ctesting.VerifyCharacterizationTestsAndResults(t, testSuite, true)
}
```

#### Testing Complex Data Structures

```go
type UserProfile struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func TestUserProfileCreation(t *testing.T) {
    expectedProfile := UserProfile{ID: 1, Name: "John Doe", Age: 30}
    
    testSuite := []ctesting.CharacterizationTest[UserProfile]{
        ctesting.NewCharacterizationTest(expectedProfile, nil, func() (UserProfile, error) {
            return createUserProfile("John Doe", 30), nil
        }),
        
        // Test with edge case - empty name should return error
        ctesting.NewCharacterizationTest(UserProfile{}, 
            fmt.Errorf("name cannot be empty"), func() (UserProfile, error) {
            return createUserProfile("", 30)
        }),
    }
    
    ctesting.VerifyCharacterizationTestsAndResults(t, testSuite, true)
}
```

#### Testing HTTP Response Handling

```go
type APIResponse struct {
    StatusCode int
    Body       string
}

func TestAPIResponseCharacterization(t *testing.T) {
    // Mock HTTP server setup (simplified)
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(200)
        w.Write([]byte(`{"status": "ok"}`))
    }))
    defer server.Close()
    
    expectedResponse := APIResponse{
        StatusCode: 200,
        Body:       `{"status": "ok"}`,
    }
    
    testSuite := []ctesting.CharacterizationTest[APIResponse]{
        ctesting.NewCharacterizationTest(expectedResponse, nil, func() (APIResponse, error) {
            resp, err := http.Get(server.URL)
            if err != nil {
                return APIResponse{}, err
            }
            defer resp.Body.Close()
            
            body, err := io.ReadAll(resp.Body)
            if err != nil {
                return APIResponse{}, err
            }
            
            return APIResponse{
                StatusCode: resp.StatusCode,
                Body:       string(body),
            }, nil
        }),
    }
    
    ctesting.VerifyCharacterizationTestsAndResults(t, testSuite, true)
}
```

#### Testing Error Conditions

```go
func TestErrorHandlingCharacterization(t *testing.T) {
    testSuite := []ctesting.CharacterizationTest[int]{
        // Test division by zero
        ctesting.NewCharacterizationTest(0, 
            fmt.Errorf("division by zero"), func() (int, error) {
            return safeDivide(10, 0)
        }),
        
        // Test negative input validation
        ctesting.NewCharacterizationTest(0, 
            fmt.Errorf("negative numbers not allowed"), func() (int, error) {
            return processPositiveNumber(-5)
        }),
        
        // Test successful operation
        ctesting.NewCharacterizationTest(5, nil, func() (int, error) {
            return safeDivide(25, 5)
        }),
    }
    
    ctesting.VerifyCharacterizationTestsAndResults(t, testSuite, true)
}
```

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

### Best Practices

1. **Start with existing behavior**: Capture what the code currently does before making changes
2. **Test edge cases**: Include boundary conditions, empty inputs, and error scenarios
3. **Use meaningful expected values**: Choose test inputs that exercise different code paths
4. **Group related tests**: Organize tests by functionality or component being characterized
5. **Update tests when behavior changes**: Characterization tests should evolve with intentional changes
6. **Combine with unit tests**: Use characterization tests alongside traditional unit tests for comprehensive coverage
7. **Test both success and failure paths**: Ensure error handling is properly characterized
8. **Use deep error checking**: Enable deep error checking for precise error message validation

### Integration Patterns

#### Legacy Code Documentation

```go
func TestLegacyCalculatorCharacterization(t *testing.T) {
    // Document the current behavior of legacy calculator
    testSuite := []ctesting.CharacterizationTest[float64]{
        // Document floating point precision behavior
        ctesting.NewCharacterizationTest(0.1, nil, func() (float64, error) {
            return legacyCalculator.Add(0.05, 0.05), nil
        }),
        
        // Document overflow behavior
        ctesting.NewCharacterizationTest(math.Inf(1), nil, func() (float64, error) {
            return legacyCalculator.Multiply(math.MaxFloat64, 2), nil
        }),
    }
    
    ctesting.VerifyCharacterizationTestsAndResults(t, testSuite, true)
}
```

#### Refactoring Safety Net

```go
func TestRefactoringSafetyNet(t *testing.T) {
    // Before refactoring: capture current behavior
    testSuite := []ctesting.CharacterizationTest[string]{
        ctesting.NewCharacterizationTest("processed: hello", nil, func() (string, error) {
            return oldProcessingFunction("hello"), nil
        }),
    }
    
    // Run before refactoring
    ctesting.VerifyCharacterizationTestsAndResults(t, testSuite, true)
    
    // After refactoring, update the function call but keep expected results
    testSuite[0] = ctesting.NewCharacterizationTest("processed: hello", nil, func() (string, error) {
        return newProcessingFunction("hello"), nil
    })
    
    // Verify refactored code produces same results
    ctesting.VerifyCharacterizationTestsAndResults(t, testSuite, true)
}
```

### Manual Result Processing

For advanced use cases, process results manually:

```go
results, testSuiteRes := ctesting.VerifyCharacterizationTests(testSuite, true)

for i, passed := range results {
    test := testSuiteRes[i]
    if !passed {
        // Custom failure handling
        log.Printf("Test %d failed: expected %v, got %v", 
            i+1, test.ExpectedOutput, test.output)
        
        // Custom error analysis
        if test.err != test.ExpectedErr {
            log.Printf("Error mismatch: expected %v, got %v", 
                test.ExpectedErr, test.err)
        }
    }
}

// Use ctesting.VerifyResults for standard reporting
ctesting.VerifyResults(t, results, testSuiteRes)
```

## Stress Testing Framework

The `stesting` package provides a comprehensive framework for stress testing Go functions to evaluate their performance, reliability, and behavior under load. Stress tests execute a function repeatedly for a specified number of iterations to identify potential issues, memory leaks, race conditions, or performance degradation.

### Stress Test Components

#### StressTest Structure

The `StressTest[fRetType, testVarType]` struct is the foundation of the stress testing framework:

```go
type StressTest[fRetType comparable, testVarType comparable] struct {
    iterations uint64           // Number of times to execute the test function
    testVar    *testVarType    // Pointer to test variables (optional)
    F          gtu.TestFunc[fRetType] // The function to stress test
}
```

#### Creating a Stress Test

Use `NewStressTest` to create a new stress test instance:

```go
// Test a function that returns an integer
stressTest := stesting.NewStressTest[int, any](
    10000,                    // 10,000 iterations
    func() (int, error) {     // Function to test
        return fibonacci(20), nil
    },
    nil,                      // No test variables needed
)
```

### Stress Test Execution Methods

#### Sequential Execution

`RunStressTest` executes the test function sequentially:

```go
success, err := stesting.RunStressTest(&stressTest)
if !success {
    t.Errorf("Stress test failed: %v", err)
}
```

**Use cases:**

- Testing functions that are not thread-safe
- Measuring sequential performance
- Simple reliability testing

#### Parallel Execution

`RunParallelStressTest` distributes work across multiple goroutines for concurrent testing:

```go
maxWorkers := uint32(8) // Use 8 concurrent workers
success, err := stesting.RunParallelStressTest(&stressTest, maxWorkers)
if !success {
    t.Errorf("Parallel stress test failed: %v", err)
}
```

**Use cases:**

- Testing concurrent safety and race conditions
- Load testing with realistic concurrency
- Performance testing under parallel execution
- Identifying deadlocks or synchronization issues

#### File Output Testing

Save stress test results to files for analysis:

```go
// Option 1: Using file path
success, err := stesting.RunStressTestWithFilePathOut(&stressTest, "results.txt")

// Option 2: Using an existing file handle
file, _ := os.Create("detailed_results.txt")
success, err := stesting.RunStressTestWithFileOut(&stressTest, *file)
```

**Use cases:**

- Analyzing output patterns over many iterations
- Debugging intermittent issues
- Performance profiling and trend analysis
- Compliance testing with audit trails

### Stress Test Examples

#### Testing a Web Service Function

```go
func TestWebServiceStress(t *testing.T) {
    // Simulate web service calls
    webServiceCall := func() (int, error) {
        resp, err := http.Get("http://localhost:8080/api/health")
        if err != nil {
            return 0, err
        }
        defer resp.Body.Close()
        return resp.StatusCode, nil
    }
    
    stressTest := stesting.NewStressTest[int, any](500, webServiceCall, nil)
    
    // Test with 20 concurrent workers to simulate real load
    success, err := stesting.RunParallelStressTest(&stressTest, 20)
    if !success {
        t.Errorf("Web service stress test failed: %v", err)
    }
}
```

#### Testing Database Operations

```go
func TestDatabaseStress(t *testing.T) {
    db := setupDatabase() // Your database setup
    
    databaseOp := func() (bool, error) {
        // Simulate database operations
        _, err := db.Query("SELECT COUNT(*) FROM users")
        return err == nil, err
    }
    
    stressTest := stesting.NewStressTest[bool, any](1000, databaseOp, nil)
    
    // Sequential test to avoid overwhelming the database
    success, err := stesting.RunStressTest(&stressTest)
    if !success {
        t.Errorf("Database stress test failed: %v", err)
    }
}
```

#### Memory Allocation Testing

```go
func TestMemoryAllocationStress(t *testing.T) {
    memoryIntensiveFunc := func() (int, error) {
        // Allocate and process large slices
        data := make([]int, 100000)
        for i := range data {
            data[i] = i * 2
        }
        return len(data), nil
    }
    
    stressTest := stesting.NewStressTest[int, any](100, memoryIntensiveFunc, nil)
    
    // Save results to analyze memory patterns
    success, err := stesting.RunStressTestWithFilePathOut(&stressTest, "memory_test.log")
    if !success {
        t.Errorf("Memory stress test failed: %v", err)
    }
    
    // Also run in parallel to test concurrent memory allocation
    success, err = stesting.RunParallelStressTest(&stressTest, 4)
    if !success {
        t.Errorf("Parallel memory stress test failed: %v", err)
    }
}
```

### Error Handling

Stress tests return detailed error information through the `StressTestingError` type:

```go
type StressTestingError struct {
    Index uint64  // The iteration number where the error occurred
    Err   error   // The underlying error
}
```

This allows you to pinpoint exactly which iteration failed and why:

```go
success, err := stesting.RunStressTest(&stressTest)
if !success {
    if ste, ok := err.(stesting.StressTestingError); ok {
        t.Errorf("Test failed at iteration %d: %v", ste.Index, ste.Err)
    }
}
```

### Stress Testing Best Practices

1. **Choose appropriate iteration counts**: Start with smaller numbers and increase based on your needs
2. **Monitor resource usage**: Use tools like `go test -race` and `go test -memprofile` with stress tests
3. **Use parallel tests judiciously**: Not all functions are safe for concurrent testing
4. **Combine with benchmarks**: Use `go test -bench` alongside stress tests for comprehensive performance analysis
5. **File output for debugging**: Use file output when investigating intermittent failures
6. **Gradual load increase**: Start with fewer workers/iterations and gradually increase to find breaking points

### Integration with Go Testing

Stress tests integrate seamlessly with Go's testing framework:

```go
func TestMyFunctionStress(t *testing.T) {
    stressTest := stesting.NewStressTest[string, any](1000, 
        func() (string, error) {
            return myFunction("test"), nil
        }, nil)
    
    success, err := stesting.RunStressTest(&stressTest)
    if !success {
        t.Fatalf("Stress test failed: %v", err)
    }
    
    t.Logf("Successfully completed %d iterations", 1000)
}
