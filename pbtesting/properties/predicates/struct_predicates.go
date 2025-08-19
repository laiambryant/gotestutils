package predicates

import "reflect"

type StructFieldPredicates struct{ Fields map[string][]Predicate }

func (p StructFieldPredicates) Verify(v any) bool {
	if len(p.Fields) == 0 {
		return true
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return true
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return true
	}
	rt := rv.Type()
	for name, props := range p.Fields {
		f, ok := rt.FieldByName(name)
		if !ok {
			continue
		}
		fv := rv.FieldByIndex(f.Index)
		for _, prop := range props {
			if !prop.Verify(fv.Interface()) {
				return false
			}
		}
	}
	return true
}
