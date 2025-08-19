package predicates

import "reflect"

type PointerAllowNil struct{ Allowed bool }

func (p PointerAllowNil) Verify(v any) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer {
		return true
	}
	if rv.IsNil() {
		return p.Allowed
	}
	return true
}
