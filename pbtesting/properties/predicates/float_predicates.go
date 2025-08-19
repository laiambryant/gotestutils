package predicates

import "math"

type FloatMin struct{ Min float64 }
type FloatMax struct{ Max float64 }
type FloatRange struct{ Min, Max float64 }
type FloatNonZero struct{ Required bool }
type FloatFiniteOnly struct{ Enabled bool }
type FloatAllowNaN struct{ Allowed bool }
type FloatAllowInf struct{ Allowed bool }
type FloatPrecisionMax struct{ Decimals int }

func (p FloatMin) Verify(v any) bool { n, ok := asFloat64(v); return !ok || n >= p.Min }
func (p FloatMax) Verify(v any) bool { n, ok := asFloat64(v); return !ok || n <= p.Max }
func (p FloatRange) Verify(v any) bool {
	n, ok := asFloat64(v)
	return !ok || (n >= p.Min && n <= p.Max)
}
func (p FloatNonZero) Verify(v any) bool {
	if !p.Required {
		return true
	}
	n, ok := asFloat64(v)
	return !ok || n != 0
}
func (p FloatFiniteOnly) Verify(v any) bool {
	if !p.Enabled {
		return true
	}
	n, ok := asFloat64(v)
	return !ok || (!math.IsNaN(n) && !math.IsInf(n, 0))
}
func (p FloatAllowNaN) Verify(v any) bool {
	if p.Allowed {
		return true
	}
	n, ok := asFloat64(v)
	return !ok || !math.IsNaN(n)
}
func (p FloatAllowInf) Verify(v any) bool {
	if p.Allowed {
		return true
	}
	n, ok := asFloat64(v)
	return !ok || !math.IsInf(n, 0)
}
func (p FloatPrecisionMax) Verify(v any) bool {
	n, ok := asFloat64(v)
	if !ok || p.Decimals <= 0 {
		return true
	}
	m := math.Pow10(p.Decimals)
	return math.Abs(n*m-math.Round(n*m)) < 1e-9
}
