package examples

import "fmt"

func Sum(a, b int) int {
	return a + b
}

func Multiply(a, b int) int {
	return a * b
}

func CalculateResult(a, b int) int {
	return a*b + 2
}

func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

func ProcessPositiveNumber(n int) (int, error) {
	if n < 0 {
		return 0, fmt.Errorf("negative numbers not allowed")
	}
	return n, nil
}

func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}
