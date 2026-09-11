package calc

import (
	"errors"
	"testing"
)

func TestAdd(t *testing.T) {
	cases := []struct {
		name string
		a, b float64
		want float64
	}{
		{"positive", 2, 3, 5},
		{"negative", -2, -3, -5},
		{"decimal", 1.5, 2.25, 3.75},
		{"zero operands", 0, 0, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Add(tc.a, tc.b); got != tc.want {
				t.Errorf("Add(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	cases := []struct {
		name string
		a, b float64
		want float64
	}{
		{"positive", 5, 3, 2},
		{"negative", -5, -3, -2},
		{"decimal", 3.75, 1.5, 2.25},
		{"zero operands", 0, 0, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Subtract(tc.a, tc.b); got != tc.want {
				t.Errorf("Subtract(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	cases := []struct {
		name string
		a, b float64
		want float64
	}{
		{"positive", 4, 3, 12},
		{"negative", -4, 3, -12},
		{"decimal", 1.5, 2, 3},
		{"zero operands", 5, 0, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Multiply(tc.a, tc.b); got != tc.want {
				t.Errorf("Multiply(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	cases := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"positive", 12, 3, 4, nil},
		{"negative", -12, 3, -4, nil},
		{"decimal", 7.5, 2.5, 3, nil},
		{"by zero", 5, 0, 0, ErrDivisionByZero},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Divide(tc.a, tc.b)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Divide(%v, %v) error = %v, want %v", tc.a, tc.b, err, tc.wantErr)
			}
			if tc.wantErr == nil && got != tc.want {
				t.Errorf("Divide(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestPower(t *testing.T) {
	cases := []struct {
		name string
		a, b float64
		want float64
	}{
		{"positive", 2, 3, 8},
		{"negative base", -2, 3, -8},
		{"decimal exponent", 4, 0.5, 2},
		{"zero exponent", 5, 0, 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Power(tc.a, tc.b)
			if err != nil {
				t.Fatalf("Power(%v, %v) error = %v, want nil", tc.a, tc.b, err)
			}
			if got != tc.want {
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
