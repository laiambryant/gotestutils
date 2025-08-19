package predicates

type Predicate interface{ Verify(any) bool }
