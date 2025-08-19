package predicates

import "testing"

func TestSliceProperties(t *testing.T) {
    assertProp(t, SliceLenMin{Min: 3}, []int{1, 2}, false)
    assertProp(t, SliceLenMin{Min: 3}, []int{1, 2, 3}, true)
    assertProp(t, SliceLenMin{Min: 3}, "SASA", false)
    assertProp(t, SliceLenMax{Max: 2}, []int{1, 2, 3}, false)
    assertProp(t, SliceLenMax{Max: 3}, []int{1, 2, 3}, true)
    assertProp(t, SliceLenMax{Max: 3}, "SASA", false)
    assertProp(t, SliceLenRange{Min: 2, Max: 3}, []int{1}, false)
    assertProp(t, SliceLenRange{Min: 2, Max: 3}, []int{1, 2}, true)
    assertProp(t, SliceLenRange{Min: 2, Max: 3}, "", false)
    props := []Predicate{IntMin{Min: 2}}
    assertProp(t, SliceElementPredicates{Props: props}, []int{2, 3, 4}, true)
    assertProp(t, SliceElementPredicates{Props: props}, []int{1, 3, 4}, false)
    assertProp(t, SliceElementPredicates{Props: props}, "", false)
}
