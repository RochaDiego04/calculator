package calc

import (
	"errors"
	"testing"
)

func TestPower(t *testing.T) {
	cases := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"positive", 2, 3, 8, nil},
		{"negative base", -2, 3, -8, nil},
		{"decimal exponent", 4, 0.5, 2, nil},
		{"zero exponent", 5, 0, 1, nil},
		{"overflow to +Inf", 1e308, 2, 0, ErrResultNotRepresentable},
		{"NaN from negative base fractional exponent", -8, 1.0 / 3, 0, ErrResultNotRepresentable},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Power(tc.a, tc.b)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Power(%v, %v) error = %v, want %v", tc.a, tc.b, err, tc.wantErr)
			}
			if tc.wantErr == nil && got != tc.want {
				t.Errorf("Power(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestSquareRoot(t *testing.T) {
	cases := []struct {
		name    string
		a       float64
		want    float64
		wantErr error
	}{
		{"positive", 9, 3, nil},
		{"zero", 0, 0, nil},
		{"decimal", 2.25, 1.5, nil},
		{"negative", -4, 0, ErrNegativeSquareRoot},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := SquareRoot(tc.a)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("SquareRoot(%v) error = %v, want %v", tc.a, err, tc.wantErr)
			}
			if tc.wantErr == nil && got != tc.want {
				t.Errorf("SquareRoot(%v) = %v, want %v", tc.a, got, tc.want)
			}
		})
	}
}

func TestPercentage(t *testing.T) {
	cases := []struct {
		name string
		a, b float64
		want float64
	}{
		{"50 percent of 200", 200, 50, 100},
		{"10 percent of 90", 90, 10, 9},
		{"0 percent", 100, 0, 0},
		{"percent of zero", 0, 50, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Percentage(tc.a, tc.b); got != tc.want {
				t.Errorf("Percentage(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}
