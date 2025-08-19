package predicates

import (
    "reflect"
    "testing"
)

func TestPointerAllowNil(t *testing.T) {
    var p *int
    assertProp(t, PointerAllowNil{Allowed: true}, p, true)
    assertProp(t, PointerAllowNil{Allowed: false}, p, false)
    x := 3
    p = &x
    assertProp(t, PointerAllowNil{Allowed: false}, p, true)
}

func TestInterfaceAllowedConcrete(t *testing.T) {
    typ := reflect.TypeOf(0)
    assertProp(t, InterfaceAllowedConcrete{Types: []reflect.Type{typ}}, 3, true)
    assertProp(t, InterfaceAllowedConcrete{Types: []reflect.Type{reflect.TypeOf("")}}, 3, false)
}

func TestChanBufferProperties(t *testing.T) {
    ch := make(chan int, 2)
    assertProp(t, ChanBufferMin{Min: 3}, ch, false)
    assertProp(t, ChanBufferMin{Min: 2}, ch, true)
    assertProp(t, ChanBufferMax{Max: 1}, ch, false)
    assertProp(t, ChanBufferMax{Max: 2}, ch, true)
    assertProp(t, ChanBufferRange{Min: 1, Max: 2}, ch, true)
    assertProp(t, ChanBufferRange{Min: 3, Max: 4}, ch, false)
}
