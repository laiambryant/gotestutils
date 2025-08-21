package predicates

import "testing"

type TestPredicate struct{}

func (t TestPredicate) Verify(a any) bool {
	return true
}

func TestPredicateImplementations(t *testing.T) {
	_ = []Predicate{}
	for _, p := range []Predicate{TestPredicate{}} {
		_ = p.Verify(nil)
	}
}
