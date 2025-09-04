package attributes

type SupportedType interface {
	int | int8 | int16 | int32 | int64 |
		float32 | float64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		complex64 | complex128 |
		string | bool | []any
}

type Set[T SupportedType] struct {
	upperBound T
	lowerBound T
	excluded   []T
	mandatory  []T
}

func (s Set[T]) WithUpperBound(ub T) Set[T] {
	s.upperBound = ub
	return s
}

func (s Set[T]) WithLowerBound(lb T) Set[T] {
	s.lowerBound = lb
	return s
}

func (s *Set[T]) AddToExcluded(excl T) {
	s.excluded = append(s.excluded, excl)
}

func (s *Set[T]) AddToMandatory(mand T) {
	s.mandatory = append(s.mandatory, mand)
}
