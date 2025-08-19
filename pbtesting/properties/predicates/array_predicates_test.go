package predicates

import "testing"

func TestArrayProperties(t *testing.T) {
    props := []Predicate{IntMin{Min: 2}}
    arrGood := [3]int{2, 3, 4}
    arrBad := [3]int{1, 3, 4}
    assertProp(t, ArrayElementPredicates{Props: props}, arrGood, true)
    assertProp(t, ArrayElementPredicates{Props: props}, arrBad, false)
    assertProp(t, ArrayElementPredicates{Props: props}, "", false)
    assertProp(t, ArraySorted{Enabled: true}, [3]int{1, 2, 3}, true)
    assertProp(t, ArraySorted{Enabled: true}, [3]int{2, 1, 3}, false)
    assertProp(t, ArraySorted{Enabled: false}, "", true)
    assertProp(t, ArraySorted{Enabled: true}, "", false)
}
