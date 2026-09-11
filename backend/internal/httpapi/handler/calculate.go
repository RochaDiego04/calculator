package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/RochaDiego04/calculator/backend/internal/calc"
	"github.com/RochaDiego04/calculator/backend/internal/httpapi/respond"
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
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, h.maxBodyBytes))
	dec.DisallowUnknownFields()

	var req calculateRequest
	if err := dec.Decode(&req); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			respond.Errors(w, http.StatusRequestEntityTooLarge,
				respond.Error{Code: "REQUEST_TOO_LARGE", Message: "request body is too large"})
			return
		}
		respond.Errors(w, http.StatusBadRequest,
			respond.Error{Code: "INVALID_JSON", Message: "request body is not valid JSON"})
		return
	}

	if req.A == nil {
		respond.Errors(w, http.StatusUnprocessableEntity,
			respond.Error{Code: "MISSING_OPERAND", Message: "a is required", Field: "a"})
		return
	}

	result, err := calc.Apply(req.Operation, *req.A, req.B)
	if err != nil {
		writeCalcError(w, err)
		return
	}

	respond.JSON(w, http.StatusOK, calculateResponse{
		Operation: req.Operation,
		A:         *req.A,
		B:         req.B,
		Result:    result,
	})
}

func writeCalcError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, calc.ErrMissingOperand):
		respond.Errors(w, http.StatusUnprocessableEntity,
			respond.Error{Code: "MISSING_OPERAND", Message: "b is required", Field: "b"})
	case errors.Is(err, calc.ErrUnknownOperation):
		respond.Errors(w, http.StatusUnprocessableEntity,
			respond.Error{Code: "UNKNOWN_OPERATION", Message: "unknown operation", Field: "operation"})
	case errors.Is(err, calc.ErrDivisionByZero):
		respond.Errors(w, http.StatusUnprocessableEntity,
			respond.Error{Code: "DIVISION_BY_ZERO", Message: "cannot divide by zero", Field: "b"})
	case errors.Is(err, calc.ErrNegativeSquareRoot):
		respond.Errors(w, http.StatusUnprocessableEntity,
			respond.Error{Code: "NEGATIVE_SQUARE_ROOT", Message: "cannot take the square root of a negative number", Field: "a"})
	case errors.Is(err, calc.ErrResultNotRepresentable):
		respond.Errors(w, http.StatusUnprocessableEntity,
			respond.Error{Code: "RESULT_NOT_REPRESENTABLE", Message: "result is not representable as a JSON number"})
	default:
		respond.Errors(w, http.StatusInternalServerError,
			respond.Error{Code: "INTERNAL", Message: "internal server error"})
	}
}
