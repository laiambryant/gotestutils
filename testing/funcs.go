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
