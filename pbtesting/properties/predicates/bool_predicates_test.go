package predicates

import "testing"

func TestBoolProperties(t *testing.T) {
	assertProp(t, BoolMustBeTrue{}, true, true)
	assertProp(t, BoolMustBeTrue{}, false, false)
	assertProp(t, BoolMustBeFalse{}, false, true)
	assertProp(t, BoolMustBeFalse{}, true, false)
	assertProp(t, BoolMustBeTrue{}, 1, true)
}
