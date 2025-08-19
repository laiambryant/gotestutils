package predicates

import "testing"

func TestMapProperties(t *testing.T) {
    m := map[int]int{1: 2, 3: 4}
    keyProps := []Predicate{IntMin{Min: 1}}
    valProps := []Predicate{IntMax{Max: 4}}
    assertProp(t, MapSizeMin{Min: 1}, m, true)
    assertProp(t, MapSizeMin{Min: 1}, "x", false)
    assertProp(t, MapSizeMax{Max: 5}, m, true)
    assertProp(t, MapSizeMax{Max: 5}, "x", false)
    assertProp(t, MapSizeRange{Max: 5}, m, true)
    assertProp(t, MapSizeRange{Max: 5}, "x", false)
    assertProp(t, MapKeyPredicates{Props: keyProps}, m, true)
    assertProp(t, MapKeyPredicates{Props: keyProps}, "x", false)
    assertProp(t, MapValuePredicates{Props: valProps}, m, true)
    assertProp(t, MapValuePredicates{Props: valProps}, "", false)
    assertProp(t, MapKeyPredicates{Props: []Predicate{IntMin{Min: 2}}}, m, false)
    assertProp(t, MapValuePredicates{Props: []Predicate{IntMax{Max: 3}}}, m, false)
}
