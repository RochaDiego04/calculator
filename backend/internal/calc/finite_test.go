package calc

import (
	"errors"
	"math"
	"testing"
)

func TestEnsureFinite(t *testing.T) {
	cases := []struct {
		name    string
		v       float64
		want    float64
		wantErr error
	}{
		{"finite", 42.5, 42.5, nil},
		{"zero", 0, 0, nil},
		{"NaN", math.NaN(), 0, ErrResultNotRepresentable},
		{"positive infinity", math.Inf(1), 0, ErrResultNotRepresentable},
		{"negative infinity", math.Inf(-1), 0, ErrResultNotRepresentable},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ensureFinite(tc.v)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("ensureFinite(%v) error = %v, want %v", tc.v, err, tc.wantErr)
			}
			if tc.wantErr == nil && got != tc.want {
				t.Errorf("ensureFinite(%v) = %v, want %v", tc.v, got, tc.want)
			}
		})
	}
}
