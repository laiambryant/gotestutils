package pbtesting

import "testing"

func TestExtractFargTypes(t *testing.T) {
	f := func(a int, b int, c int) (retA int, retB int, errA error) {
		return 1, 2, nil
	}
	in, out := extractFArgTypes(f)
	if len(in) != 3 {
		t.Errorf("Expected 3 input types, got %d", len(in))
	}
	if len(out) != 3 {
		t.Errorf("Expected 3 output types, got %d", len(out))
	}
	t.Logf("Input Types: %v\nOutput Types: %v", in, out)
}
