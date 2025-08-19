package predicates

import "math"

type ComplexRealRange struct{ Min, Max float64 }
type ComplexImagRange struct{ Min, Max float64 }
type ComplexMagnitudeRange struct{ Min, Max float64 }
type ComplexAllowNaN struct{ Allowed bool }
type ComplexAllowInf struct{ Allowed bool }

func (p ComplexRealRange) Verify(v any) bool {
	c, ok := asComplex128(v)
	return !ok || (real(c) >= p.Min && real(c) <= p.Max)
}
func (p ComplexImagRange) Verify(v any) bool {
	c, ok := asComplex128(v)
	return !ok || (imag(c) >= p.Min && imag(c) <= p.Max)
}
func (p ComplexMagnitudeRange) Verify(v any) bool {
	c, ok := asComplex128(v)
	if !ok {
		return true
	}
	r := math.Hypot(real(c), imag(c))
	return r >= p.Min && r <= p.Max
}
func (p ComplexAllowNaN) Verify(v any) bool {
	c, ok := asComplex128(v)
	if !ok || p.Allowed {
		return true
	}
	return !(math.IsNaN(real(c)) || math.IsNaN(imag(c)))
}
func (p ComplexAllowInf) Verify(v any) bool {
	c, ok := asComplex128(v)
	if !ok || p.Allowed {
		return true
	}
	return !(math.IsInf(real(c), 0) || math.IsInf(imag(c), 0))
}
