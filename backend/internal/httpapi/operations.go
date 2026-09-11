package httpapi

import (
	"net/http"

	"github.com/RochaDiego04/calculator/backend/internal/calc"
)

func (h *Handler) Operations(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, calc.Operations())
}
