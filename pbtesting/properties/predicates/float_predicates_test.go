package predicates

import (
	"math"
	"testing"
)

func TestFloatProperties(t *testing.T) {
	assertProp(t, FloatMin{Min: 1.5}, float32(1.4), false)
	assertProp(t, FloatMin{Min: 1.5}, 1.5, true)
	assertProp(t, FloatMax{Max: 2.5}, 2.6, false)
	assertProp(t, FloatMax{Max: 2.5}, 2.5, true)
	assertProp(t, FloatRange{Min: 1, Max: 2}, 0.9, false)
	assertProp(t, FloatRange{Min: 1, Max: 2}, 1.5, true)
	assertProp(t, FloatNonZero{Required: true}, 0.0, false)
	assertProp(t, FloatNonZero{Required: true}, 0.1, true)
	assertProp(t, FloatNonZero{Required: false}, 0.1, true)
	assertProp(t, FloatFiniteOnly{Enabled: true}, math.NaN(), false)
	assertProp(t, FloatFiniteOnly{Enabled: true}, math.Inf(1), false)
	assertProp(t, FloatFiniteOnly{Enabled: false}, math.Inf(1), true)
	assertProp(t, FloatFiniteOnly{Enabled: true}, 3.14, true)
	assertProp(t, FloatAllowNaN{Allowed: false}, math.NaN(), false)
	assertProp(t, FloatAllowNaN{Allowed: true}, math.NaN(), true)
	assertProp(t, FloatAllowInf{Allowed: false}, math.Inf(-1), false)
	assertProp(t, FloatAllowInf{Allowed: true}, math.Inf(-1), true)
	assertProp(t, FloatPrecisionMax{Decimals: 2}, 1.23, true)
	assertProp(t, FloatPrecisionMax{Decimals: 2}, "sw", true)
	assertProp(t, FloatPrecisionMax{Decimals: 2}, 1.234, false)
}
