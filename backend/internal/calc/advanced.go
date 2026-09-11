package calc

import "math"

func Power(a, b float64) (float64, error) {
	return ensureFinite(math.Pow(a, b))
}

func SquareRoot(a float64) (float64, error) {
	if a < 0 {
		return 0, ErrNegativeSquareRoot
	}
	return ensureFinite(math.Sqrt(a))
}

func Percentage(a, b float64) float64 {
	return a * b / 100
}
