package predicates

import "reflect"

type MapSizeMin struct{ Min int }
type MapSizeMax struct{ Max int }
type MapSizeRange struct{ Min, Max int }
type MapKeyPredicates struct{ Props []Predicate }
type MapValuePredicates struct{ Props []Predicate }

func (p MapSizeMin) Verify(v any) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Map {
		return false
	}
	return rv.Len() >= p.Min
}
func (p MapSizeMax) Verify(v any) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Map {
		return false
	}
	return rv.Len() <= p.Max
}
func (p MapSizeRange) Verify(v any) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Map {
		return false
	}
	l := rv.Len()
	return l >= p.Min && l <= p.Max
}
func (p MapKeyPredicates) Verify(v any) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Map {
		return false
	}
	if len(p.Props) == 0 {
		return true
	}
	iter := rv.MapRange()
	for iter.Next() {
		k := iter.Key().Interface()
		for _, prop := range p.Props {
			if !prop.Verify(k) {
				return false
			}
		}
	}
	return true
}
func (p MapValuePredicates) Verify(v any) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Map {
		return false
	}
	if len(p.Props) == 0 {
		return true
	}
	iter := rv.MapRange()
	for iter.Next() {
		val := iter.Value().Interface()
		for _, prop := range p.Props {
			if !prop.Verify(val) {
				return false
			}
		}
	}
	return true
}
