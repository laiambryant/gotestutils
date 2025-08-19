package predicates

import (
	"reflect"
	"testing"
)

func assertProp(t *testing.T, p Predicate, val any, expect bool) {
	if p == nil {
		if !expect {
			t.Fatalf("nil predicate expected=%v but treated as true", expect)
		}
		return
	}
	name := reflect.TypeOf(p).Name()
	got := p.Verify(val)
	if got != expect {
		t.Fatalf("%s.Verify(%#v) = %v, want %v", name, val, got, expect)
	}
}
