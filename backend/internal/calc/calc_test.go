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
