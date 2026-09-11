package calc

import "errors"

var (
	ErrDivisionByZero         = errors.New("division by zero")
	ErrNegativeSquareRoot     = errors.New("square root of a negative number")
	ErrResultNotRepresentable = errors.New("result is not representable as a JSON number")
	ErrUnknownOperation       = errors.New("unknown operation")
	ErrMissingOperand         = errors.New("missing operand")
)
