package predicates

import "reflect"

type InterfaceAllowedConcrete struct{ Types []reflect.Type }

func (p InterfaceAllowedConcrete) Verify(v any) bool {
	if v == nil || len(p.Types) == 0 {
		return true
	}
	rt := reflect.TypeOf(v)
	for _, t := range p.Types {
		if rt.AssignableTo(t) {
			return true
		}
	}
	return false
}
