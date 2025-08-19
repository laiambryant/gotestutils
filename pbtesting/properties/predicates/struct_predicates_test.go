package predicates

import "testing"

func TestStructFieldProperties(t *testing.T) {
    type S struct { A int; B string }
    vGood := S{A: 5, B: "hello"}
    vBad := S{A: 1, B: "x"}
    props := StructFieldPredicates{Fields: map[string][]Predicate{ "A": {IntMin{Min: 3}}, "B": {StringLenMin{Min: 2}}, }}
    assertProp(t, props, vGood, true)
    assertProp(t, props, vBad, false)
}

func TestStructFieldPropertiesEdgeCases(t *testing.T) {
    type S struct { A int; B string }
    vGood := S{A: 5, B: "hello"}
    props := StructFieldPredicates{Fields: map[string][]Predicate{ "A": {IntMin{Min: 3}}, "B": {StringLenMin{Min: 2}}, "C": {SliceElementPredicates{Props: nil}}, }}
    assertProp(t, props, vGood, true)
    assertProp(t, nil, vGood, true)
}
