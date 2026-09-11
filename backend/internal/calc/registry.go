package calc

type Operation string

const (
	OpAdd        Operation = "add"
	OpSubtract   Operation = "subtract"
	OpMultiply   Operation = "multiply"
	OpDivide     Operation = "divide"
	OpPower      Operation = "power"
	OpSqrt       Operation = "sqrt"
	OpPercentage Operation = "percentage"
)

type Descriptor struct {
	Name   Operation `json:"name"`
	Symbol string    `json:"symbol"`
	Arity  int       `json:"arity"`
	Label  string    `json:"label"`
}

var descriptors = []Descriptor{
	{OpAdd, "+", 2, "Add"},
	{OpSubtract, "−", 2, "Subtract"},
	{OpMultiply, "×", 2, "Multiply"},
	{OpDivide, "÷", 2, "Divide"},
	{OpPower, "^", 2, "Power"},
	{OpSqrt, "√", 1, "Square Root"},
	{OpPercentage, "%", 2, "Percentage"},
}

func Operations() []Descriptor {
	out := make([]Descriptor, len(descriptors))
	copy(out, descriptors)
	return out
}

func descriptorFor(op Operation) (Descriptor, bool) {
	for _, d := range descriptors {
		if d.Name == op {
			return d, true
		}
	}
	return Descriptor{}, false
}

func Apply(op Operation, a float64, b *float64) (float64, error) {
	d, ok := descriptorFor(op)
	if !ok {
		return 0, ErrUnknownOperation
	}
	if d.Arity == 2 && b == nil {
		return 0, ErrMissingOperand
	}

	var result float64
	var err error
	switch op {
	case OpAdd:
		result = Add(a, *b)
	case OpSubtract:
		result = Subtract(a, *b)
	case OpMultiply:
		result = Multiply(a, *b)
	case OpDivide:
		result, err = Divide(a, *b)
	case OpPower:
		result, err = Power(a, *b)
	case OpSqrt:
		result, err = SquareRoot(a)
	case OpPercentage:
		result = Percentage(a, *b)
	default:
		return 0, ErrUnknownOperation
	}
	if err != nil {
		return 0, err
	}
	return ensureFinite(result)
}
