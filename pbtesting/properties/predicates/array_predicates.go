package predicates

import "reflect"

type ArrayElementPredicates struct{ Props []Predicate }
type ArraySorted struct{ Enabled bool }

func (p ArrayElementPredicates) Verify(v any) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Array {
		return false
	}
	if len(p.Props) == 0 {
		return true
	}
	for i := 0; i < rv.Len(); i++ {
		val := rv.Index(i).Interface()
		for _, prop := range p.Props {
			if !prop.Verify(val) {
				return false
			}
		}
	}
	return true
}
func (p ArraySorted) Verify(v any) bool {
	if !p.Enabled {
		return true
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Array {
		return false
	}
	if rv.Len() < 2 {
		return true
	}
	for i := 1; i < rv.Len(); i++ {
		if less(rv.Index(i).Interface(), rv.Index(i-1).Interface()) {
			return false
		}
	}
	return true
}
