package predicates

import (
	"math"
	"testing"
)

func TestComplexProperties(t *testing.T) {
	assertProp(t, ComplexRealRange{Min: -1, Max: 1}, complex(1.5, 0), false)
	assertProp(t, ComplexRealRange{Min: -1, Max: 1}, complex(0.5, 0), true)
	assertProp(t, ComplexImagRange{Min: -1, Max: 1}, complex64(complex(0, 1.5)), false)
	assertProp(t, ComplexImagRange{Min: -1, Max: 1}, complex128(complex(0, -0.5)), true)
	assertProp(t, ComplexMagnitudeRange{Min: 0, Max: 2}, complex(2, 0), true)
	assertProp(t, ComplexMagnitudeRange{Min: 0, Max: 1}, complex(1, 1), false)
	assertProp(t, ComplexAllowNaN{Allowed: false}, complex(math.NaN(), 0), false)
	assertProp(t, ComplexAllowNaN{Allowed: true}, complex(math.NaN(), 0), true)
	assertProp(t, ComplexAllowInf{Allowed: false}, complex(math.Inf(1), 0), false)
	assertProp(t, ComplexAllowInf{Allowed: true}, complex(math.Inf(1), 0), true)
	assertProp(t, ComplexMagnitudeRange{Min: 0, Max: 2}, "dsada", true)
}
