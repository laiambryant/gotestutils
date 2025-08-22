package utils

import (
	"reflect"
	"testing"
)

func TestExtractFargTypes(t *testing.T) {
	f := func(a int, b int, c int) (retA int, retB int, errA error) {
		return 1, 2, nil
	}
	in, out := ExtractFArgTypes(f)
	if len(in) != 3 {
		t.Errorf("Expected 3 input types, got %d", len(in))
	}
	if len(out) != 3 {
		t.Errorf("Expected 3 output types, got %d", len(out))
	}
	t.Logf("Input Types: %v\nOutput Types: %v", in, out)
}

func TestFilter(t *testing.T) {
	t.Run("filter integers", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5, 6}
		evens := Filter(numbers, func(n int) bool { return n%2 == 0 })
		expected := []int{2, 4, 6}

		if !reflect.DeepEqual(evens, expected) {
			t.Errorf("expected %v, got %v", expected, evens)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		empty := []int{}
		result := Filter(empty, func(n int) bool { return n > 0 })

		if len(result) != 0 {
			t.Errorf("expected empty slice, got %v", result)
		}
	})

	t.Run("no matches", func(t *testing.T) {
		numbers := []int{1, 3, 5, 7}
		evens := Filter(numbers, func(n int) bool { return n%2 == 0 })

		if len(evens) != 0 {
			t.Errorf("expected empty slice, got %v", evens)
		}
	})

	t.Run("all matches", func(t *testing.T) {
		numbers := []int{2, 4, 6, 8}
		evens := Filter(numbers, func(n int) bool { return n%2 == 0 })

		if !reflect.DeepEqual(evens, numbers) {
			t.Errorf("expected %v, got %v", numbers, evens)
		}
	})
}
