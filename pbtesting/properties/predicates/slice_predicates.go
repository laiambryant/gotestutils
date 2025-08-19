package predicates

import "reflect"

type SliceLenMin struct{ Min int }
type SliceLenMax struct{ Max int }
type SliceLenRange struct{ Min, Max int }
type SliceUnique struct{ Enabled bool }
type SliceElementPredicates struct{ Props []Predicate }

func (p SliceLenMin) Verify(v any) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Slice {
		return false
	}
	return rv.Len() >= p.Min
}
func (p SliceLenMax) Verify(v any) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Slice {
		return false
	}
	return rv.Len() <= p.Max
}
func (p SliceLenRange) Verify(v any) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Slice {
		return false
	}
	l := rv.Len()
	return l >= p.Min && l <= p.Max
}
func (p SliceElementPredicates) Verify(v any) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Slice {
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
func (p SliceUnique) Verify(v any) bool {
	if !p.Enabled {
		return true
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Slice {
		return false
	}
	seen := make(map[any]struct{})
	for i := 0; i < rv.Len(); i++ {
		ev := rv.Index(i)
		if !isHashable(ev) {
			continue
		}
		k := ev.Interface()
		if _, ok := seen[k]; ok {
			return false
		}
		seen[k] = struct{}{}
	}
	return true
}
