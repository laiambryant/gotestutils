package gotestutils

// TestFunc defines the signature for functions that can be tested.
// The function should return a value of type t and an error.
// This allows testing both the return value and error conditions.
//
// Example:
//
//	func() (int, error) { return sum(1, 2), nil }
//	func() (string, error) { return getName(), fmt.Errorf("not found") }
type TestFunc[t comparable] func() (t, error)

// PTestFunc defines the signature for parametrized functions that can be tested with property-based testing.
// The function accepts variadic arguments of type argT and returns a value of type retT and an error.
// This allows property testing frameworks to generate test inputs and verify both return values and error conditions.
//
// Type parameters:
//   - retT: the return type that must be comparable for assertion purposes
//   - argT: the input argument type that can be any type for flexible test case generation
//
// Example:
//
//	func(a, b int) (int, error) { return a + b, nil }
//	func(names ...string) (string, error) { return strings.Join(names, " "), nil }
type PTestFunc[retT comparable, argT any] func(...argT) (retT, error)
