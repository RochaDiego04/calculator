package calc

import (
	"errors"
	"testing"
)

func f(v float64) *float64 { return &v }

func TestOperations(t *testing.T) {
	got := Operations()
	if len(got) != 7 {
		t.Fatalf("Operations() returned %d descriptors, want 7", len(got))
	}

	want := map[Operation]int{
		OpAdd: 2, OpSubtract: 2, OpMultiply: 2, OpDivide: 2,
		OpPower: 2, OpSqrt: 1, OpPercentage: 2,
	}
	for _, d := range got {
		arity, ok := want[d.Name]
		if !ok {
			t.Errorf("Operations() included unexpected operation %q", d.Name)
			continue
		}
		if d.Arity != arity {
			t.Errorf("descriptor %q arity = %d, want %d", d.Name, d.Arity, arity)
		}
	}

	got[0].Symbol = "mutated"
	if Operations()[0].Symbol == "mutated" {
		t.Error("Operations() returned a slice that aliases internal state")
	}
}

func TestApply(t *testing.T) {
	cases := []struct {
		name    string
		op      Operation
		a       float64
		b       *float64
		want    float64
		wantErr error
	}{
		{"add", OpAdd, 2, f(3), 5, nil},
		{"subtract", OpSubtract, 5, f(3), 2, nil},
		{"multiply", OpMultiply, 4, f(3), 12, nil},
		{"divide", OpDivide, 12, f(3), 4, nil},
		{"power", OpPower, 2, f(3), 8, nil},
		{"sqrt", OpSqrt, 9, nil, 3, nil},
		{"percentage", OpPercentage, 200, f(50), 100, nil},
		{"unknown operation", Operation("cos"), 1, f(1), 0, ErrUnknownOperation},
		{"missing operand", OpAdd, 1, nil, 0, ErrMissingOperand},
		{"division by zero", OpDivide, 5, f(0), 0, ErrDivisionByZero},
		{"negative square root", OpSqrt, -4, nil, 0, ErrNegativeSquareRoot},
		{"add overflow to +Inf", OpAdd, 1e308, f(1e308), 0, ErrResultNotRepresentable},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Apply(tc.op, tc.a, tc.b)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Apply(%v, %v, %v) error = %v, want %v", tc.op, tc.a, tc.b, err, tc.wantErr)
			}
			if tc.wantErr == nil && got != tc.want {
				t.Errorf("Apply(%v, %v, %v) = %v, want %v", tc.op, tc.a, tc.b, got, tc.want)
			}
		})
	}
}
