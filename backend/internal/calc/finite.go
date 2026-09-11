package calc

import "math"

func ensureFinite(v float64) (float64, error) {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, ErrResultNotRepresentable
	}
	return v, nil
}
