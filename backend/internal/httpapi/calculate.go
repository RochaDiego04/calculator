package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/RochaDiego04/calculator/backend/internal/calc"
)

type calculateRequest struct {
	Operation calc.Operation `json:"operation"`
	A         *float64       `json:"a"`
	B         *float64       `json:"b"`
}

type calculateResponse struct {
	Operation calc.Operation `json:"operation"`
	A         float64        `json:"a"`
	B         *float64       `json:"b,omitempty"`
	Result    float64        `json:"result"`
}

func (h *Handler) Calculate(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	dec.DisallowUnknownFields()

	var req calculateRequest
	if err := dec.Decode(&req); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			writeErrors(w, http.StatusRequestEntityTooLarge,
				apiError{Code: "REQUEST_TOO_LARGE", Message: "request body exceeds the 1 MB limit"})
			return
		}
		writeErrors(w, http.StatusBadRequest,
			apiError{Code: "INVALID_JSON", Message: "request body is not valid JSON"})
		return
	}

	if req.A == nil {
		writeErrors(w, http.StatusUnprocessableEntity,
			apiError{Code: "MISSING_OPERAND", Message: "a is required", Field: "a"})
		return
	}

	result, err := calc.Apply(req.Operation, *req.A, req.B)
	if err != nil {
		writeCalcError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, calculateResponse{
		Operation: req.Operation,
		A:         *req.A,
		B:         req.B,
		Result:    result,
	})
}

func writeCalcError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, calc.ErrMissingOperand):
		writeErrors(w, http.StatusUnprocessableEntity,
			apiError{Code: "MISSING_OPERAND", Message: "b is required", Field: "b"})
	case errors.Is(err, calc.ErrUnknownOperation):
		writeErrors(w, http.StatusUnprocessableEntity,
			apiError{Code: "UNKNOWN_OPERATION", Message: "unknown operation", Field: "operation"})
	case errors.Is(err, calc.ErrDivisionByZero):
		writeErrors(w, http.StatusUnprocessableEntity,
			apiError{Code: "DIVISION_BY_ZERO", Message: "cannot divide by zero", Field: "b"})
	case errors.Is(err, calc.ErrNegativeSquareRoot):
		writeErrors(w, http.StatusUnprocessableEntity,
			apiError{Code: "NEGATIVE_SQUARE_ROOT", Message: "cannot take the square root of a negative number", Field: "a"})
	case errors.Is(err, calc.ErrResultNotRepresentable):
		writeErrors(w, http.StatusUnprocessableEntity,
			apiError{Code: "RESULT_NOT_REPRESENTABLE", Message: "result is not representable as a JSON number"})
	default:
		writeErrors(w, http.StatusInternalServerError,
			apiError{Code: "INTERNAL", Message: "internal server error"})
	}
}
