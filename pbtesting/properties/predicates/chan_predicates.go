package predicates

import "reflect"

type ChanBufferMin struct{ Min int }
type ChanBufferMax struct{ Max int }
type ChanBufferRange struct{ Min, Max int }

func (p ChanBufferMin) Verify(v any) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Chan {
		return true
	}
	return rv.Cap() >= p.Min
}
func (p ChanBufferMax) Verify(v any) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Chan {
		return true
	}
	return rv.Cap() <= p.Max
}
func (p ChanBufferRange) Verify(v any) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Chan {
		return true
	}
	c := rv.Cap()
	return c >= p.Min && c <= p.Max
}
