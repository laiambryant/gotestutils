package predicates

import "slices"

type IntMin struct{ Min int64 }
type IntMax struct{ Max int64 }
type IntRange struct{ Min, Max int64 }
type IntNonZero struct{ Required bool }
type IntEvenOnly struct{ Enabled bool }
type IntOddOnly struct{ Enabled bool }
type IntMultipleOf struct{ K int64 }
type IntInSet struct{ Values []int64 }
type IntNotInSet struct{ Values []int64 }
type IntSigned struct{ AllowNegative bool }
type IntCanIncludeZero struct{ Allowed bool }

func (p IntMin) Verify(v any) bool   { n, ok := asInt64(v); return !ok || n >= p.Min }
func (p IntMax) Verify(v any) bool   { n, ok := asInt64(v); return !ok || n <= p.Max }
func (p IntRange) Verify(v any) bool { n, ok := asInt64(v); return !ok || (n >= p.Min && n <= p.Max) }
func (p IntNonZero) Verify(v any) bool {
	if !p.Required {
		return true
	}
	n, ok := asInt64(v)
	return !ok || n != 0
}
func (p IntEvenOnly) Verify(v any) bool {
	if !p.Enabled {
		return true
	}
	n, ok := asInt64(v)
	return !ok || n%2 == 0
}
func (p IntOddOnly) Verify(v any) bool {
	if !p.Enabled {
		return true
	}
	n, ok := asInt64(v)
	return !ok || n%2 != 0
}
func (p IntMultipleOf) Verify(v any) bool {
	if p.K == 0 {
		return true
	}
	n, ok := asInt64(v)
	return !ok || n%p.K == 0
}
func (p IntInSet) Verify(v any) bool {
	n, ok := asInt64(v)
	if !ok {
		return false
	}
	return slices.Contains(p.Values, n)
}
func (p IntNotInSet) Verify(v any) bool {
	n, ok := asInt64(v)
	if !ok {
		return true
	}
	return !slices.Contains(p.Values, n)
}
func (p IntSigned) Verify(v any) bool {
	if p.AllowNegative {
		return true
	}
	n, ok := asInt64(v)
	return !ok || n >= 0
}
func (p IntCanIncludeZero) Verify(v any) bool {
	if p.Allowed {
		return true
	}
	n, ok := asInt64(v)
	return !ok || n != 0
}

type UintMin struct{ Min uint64 }
type UintMax struct{ Max uint64 }
type UintRange struct{ Min, Max uint64 }
type UintNonZero struct{ Required bool }
type UintMultipleOf struct{ K uint64 }
type UintInSet struct{ Values []uint64 }
type UintNotInSet struct{ Values []uint64 }
type UintCanIncludeZero struct{ Allowed bool }

func (p UintMin) Verify(v any) bool   { n, ok := asUint64(v); return !ok || n >= p.Min }
func (p UintMax) Verify(v any) bool   { n, ok := asUint64(v); return !ok || n <= p.Max }
func (p UintRange) Verify(v any) bool { n, ok := asUint64(v); return !ok || (n >= p.Min && n <= p.Max) }
func (p UintNonZero) Verify(v any) bool {
	if !p.Required {
		return true
	}
	n, ok := asUint64(v)
	return !ok || n != 0
}
func (p UintMultipleOf) Verify(v any) bool {
	if p.K == 0 {
		return true
	}
	n, ok := asUint64(v)
	return !ok || n%p.K == 0
}
func (p UintInSet) Verify(v any) bool {
	n, ok := asUint64(v)
	if !ok {
		return false
	}
	return slices.Contains(p.Values, n)
}
func (p UintNotInSet) Verify(v any) bool {
	n, ok := asUint64(v)
	if !ok {
		return true
	}
	return !slices.Contains(p.Values, n)
}
func (p UintCanIncludeZero) Verify(v any) bool {
	if p.Allowed {
		return true
	}
	n, ok := asUint64(v)
	return !ok || n != 0
}
